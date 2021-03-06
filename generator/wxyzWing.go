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

// wxyzWing removes candidates. A group consists of one "pivot" cell and 3 "wing" cells. The pivot must be able to see all of the wing cells. The group includes 4 digits, exactly one of which must be "unrestricted". A digit is restricted if every occurrance of the digit in the group can see every other occurrance.
func (g *Grid) wxyzWing(verbose uint) (res bool) {
	for r := zero; r < rows; r++ {
		for c := zero; c < cols; c++ {
			p := point{r, c}
			cell := g.cells[r][c]
			if bitCount[cell] < 2 {
				continue
			}

			var points []point
			points = append(points, box.unit[boxOf(r, c)][:]...)
			points = append(points, col.unit[c][:]...)
			points = append(points, row.unit[r][:]...)
			for p1i, p1 := range points {
				if p == p1 {
					continue
				}

				cell1 := *g.pt(p1)
				b1 := boxOfPoint(p1)

				for p2i, p2 := range points {
					if p == p2 || p1 == p2 || p1i >= p2i {
						continue
					}

					b2 := boxOfPoint(p2)

					if b1 != b2 && p1.c != p2.c && p1.r != p2.c {
						continue
					}

					cell2 := *g.pt(p2)

					// p3 is the disjoint wing cell. It cannot see the other two wings cells.
					for _, p3 := range points {
						if p == p3 || p1 == p3 || p2 == p3 {
							continue
						}

						b3 := boxOfPoint(p3)

						// At least one of the wing cells must not be able to see the other wing cells.
						if b1 == b3 || b2 == b3 ||
							p1.c == p3.c || p2.c == p3.c ||
							p1.r == p3.r || p2.r == p3.r {
							continue
						}

						cell3 := *g.pt(p3)

						// There must be a total of 4 digits in the group.
						if bitCount[cell|cell1|cell2|cell3] != 4 {
							continue
						}

						c1 := cell & cell1
						c2 := cell & cell2
						c3 := cell & cell3
						// group := c1 | c2 | c3
						unrestricted := (c1 | c2) & c3

						// The pivot must have at least two digits in common with each wing cell and there must be exactly one unrestricted digit in the group.
						if bitCount[c1] < 2 ||
							bitCount[c2] < 2 ||
							bitCount[c3] < 2 ||
							bitCount[unrestricted] != 1 {
							continue
						}

						overlap := influence[p.r][p.c].and(influence[p1.r][p1.c]).and(influence[p2.r][p2.r]).and(influence[p3.r][p3.c])
						overlap.unset(p).unset(p1).unset(p2).unset(p3) // Remove the group members from the overlap so that they do not get digits cleared.

						overlap.process(func(r, c uint8) {
							if g.pt(point{r, c}).andNot(unrestricted) {
								// fmt.Printf("p: %s (%s), p1: %s (%s), p2: %s (%s), p3: %s (%s), c1: %s, c2: %s, c3: %s, group: %s, unrestricted: %s\n", p, cell, p1, cell1, p2, cell2, p3, cell3, c1, c2, c3, group, unrestricted)
								g.cellChange(&res, verbose, "wxyzWing: removing %s from (%d, %d) because of %s, %s, %s, %s\n", unrestricted, r, c, p, p1, p2, p3)
							}
						})
					}
				}
			}
		}
	}

	return
}
