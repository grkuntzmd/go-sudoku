/*
 * MIT LICENSE
 *
 * Copyright © 2020, G.Ralph Kuntz, MD.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package generator

var pairIndexes = [][2]uint8{
	{0, 1},
	{0, 2},
	{0, 3},
	{0, 6},
	{1, 2},
	{1, 4},
	{1, 7},
	{2, 5},
	{2, 8},
	{3, 4},
	{3, 5},
	{3, 6},
	{4, 5},
	{4, 7},
	{5, 8},
	{6, 7},
	{6, 8},
	{7, 8},
}

// exocet removes candidates. When 2 of the 3 cells in a box-line intersection together contain 3 or 4 candidates, then in each of the two boxes in the same band but in different lines, if there are cells with the same 3 or 4 candidates, any others can be removed. See https://www.sudokuwiki.org/Exocet for explanation/discussion.
func (g *Grid) exocet(verbose uint) (res bool) {
	for _, b := range box.unit {
		for _, pi := range pairIndexes {
			// Find base cells.
			p1 := b[pi[0]]
			p2 := b[pi[1]]

			cell1 := *g.pt(p1)
			cell2 := *g.pt(p2)
			common := cell1 | cell2

			if bitCount[cell1] < 2 || bitCount[cell2] < 2 {
				continue
			}

			if bc := bitCount[common]; bc < 3 || bc > 4 {
				continue
			}

			b1 := boxOfPoint(p1)

			// Pattern rule 1 satisfied. Now find target cells.
			var targets []point
			if p1.c == p2.c {
				c1 := (p1.c+1)%3 + p1.c/3*3
				c2 := (p1.c+2)%3 + p1.c/3*3
				findTargets(b1, c1, c2, &col, &targets)
			} else { // p1.r == p2.r
				r1 := (p1.r+1)%3 + p1.r/3*3
				r2 := (p1.r+2)%3 + p1.r/3*3
				findTargets(b1, r1, r2, &row, &targets)
			}

			// Pattern rule 2.
			targetPairs := make(map[pair]bool)
			for _, t1 := range targets {
				for _, t2 := range targets {
					if boxOfPoint(t1) == boxOfPoint(t2) || t1.c == t2.c || t1.r == t2.r {
						continue
					}

					t1cell := *g.pt(t1)
					t2cell := *g.pt(t2)

					if t1cell&common != common || t2cell&common != common {
						continue
					}

					// Pattern rule 3: the companion cells of the target cell must not contain the base candidates.
					var c1, c2 point
					if p1.c == p2.c {
						c1 = point{t1.r, t2.c}
						c2 = point{t2.r, t1.c}
					} else { // p1.r == p2.r
						c1 = point{t2.r, t1.c}
						c2 = point{t1.r, t2.c}
					}
					if common&*g.pt(c1) != 0 || common&*g.pt(c2) != 0 {
						continue
					}

					if !targetPairs[pair{t2, t1}] {
						targetPairs[pair{t1, t2}] = true
					}
				}
			}

			if len(targetPairs) == 0 {
				continue
			}

		outer:
			for pair := range targetPairs {
				t1 := pair.left
				t2 := pair.right

				// Find the cross-lines.
				var crossLines [3][]point
				if p1.c == p2.c {
					crossLines = findCrossLines(p1, p2, t1, t2, &row, func(p point) uint8 { return p.r })
				} else { // p1.r == p2.r
					crossLines = findCrossLines(p1, p2, t1, t2, &col, func(p point) uint8 { return p.c })
				}

				// Pattern rule 4. Cross lines must not contain more than 2 instances of each base candidate. Same for cover lines (lines that run perpendiculr to cross lines).
				for i := 0; i < 6; i++ {
					var digits [10]uint8
					for c := 0; c < 3; c++ {
						cell := *g.pt(crossLines[c][i]) & common
						for d := 1; d <= 9; d++ {
							if cell&(1<<d) != 0 {
								digits[d]++
							}
						}
					}

					for d := 1; d <= 9; d++ {
						if digits[d] > 2 {
							continue outer
						}
					}
				}

				for i := 0; i < 3; i++ {
					var digits [10]uint8

					for _, p := range crossLines[i] {
						cell := *g.pt(p) & common
						for d := 1; d <= 9; d++ {
							if cell&(1<<d) != 0 {
								digits[d]++
							}
						}
					}

					for d := 1; d <= 9; d++ {
						if digits[d] > 2 {
							continue outer
						}
					}
				}

				// Elimination rule 1: candidates in the target cells that are not in the base cells can be removed.
				e1 := *g.pt(t1) &^ common
				e2 := *g.pt(t2) &^ common
				if g.pt(t1).andNot(e1) {
					g.cellChange(&res, verbose, "exocet: in %s, remove %s\n", t1, e1)
				}
				if g.pt(t2).andNot(e2) {
					g.cellChange(&res, verbose, "exocet: in %s, remove %s\n", t2, e2)
				}
			}
		}
	}

	return
}

func findCrossLines(p1, p2, t1, t2 point, gr *group, sel func(point) uint8) (crossLines [3][]point) {
	b1 := boxOfPoint(p1)
	bt1 := boxOfPoint(t1)
	bt2 := boxOfPoint(t2)
	var indexes [3]uint8
	if sel(p1)%3 == 0 && sel(p2)%3 == 1 {
		indexes[0] = sel(p1)/3*3 + 2
	} else if sel(p1)%3 == 0 && sel(p2)%3 == 2 {
		indexes[0] = sel(p1)/3*3 + 1
	} else { // sel(p1)%3 == 1 && sel(p2)%3 == 2 {
		indexes[0] = sel(p1) / 3 * 3
	}
	indexes[1] = sel(t1)
	indexes[2] = sel(t2)

	for i, ind := range indexes {
		for _, p := range gr.unit[ind] {
			b := boxOfPoint(p)
			if b1 == b || bt1 == b || bt2 == b {
				continue
			}

			crossLines[i] = append(crossLines[i], p)
		}
	}

	return
}

func findTargets(excludeBox, line1, line2 uint8, gr *group, targets *[]point) {
	for _, p := range gr.unit[line1] {
		if excludeBox == boxOfPoint(p) {
			continue
		}
		*targets = append(*targets, p)
	}
	for _, p := range gr.unit[line2] {
		if excludeBox == boxOfPoint(p) {
			continue
		}
		*targets = append(*targets, p)
	}
}
