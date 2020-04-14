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

// nakedTriple checks a group for 3 cells with the same triple of values. If present, those values can be removed from all other cells in the group. It returns true if it changes any cells.
func (g *Grid) nakedTriple(verbose uint) bool {
	return g.nakedTripleGroup(&box, verbose) || g.nakedTripleGroup(&col, verbose) || g.nakedTripleGroup(&row, verbose)
}

func (g *Grid) nakedTripleGroup(gr *group, verbose uint) (res bool) {
	for ui, u := range gr.unit {
		for _, p1 := range u {
			cell1 := *g.pt(p1)
			count := bitCount[cell1]
			if count == 1 || count > 3 {
				continue
			}

			for _, p2 := range u {
				if p1 == p2 {
					continue
				}

				cell2 := *g.pt(p2)
				count := bitCount[cell2]
				if count == 1 || count > 3 {
					continue
				}

				if bitCount[cell1|cell2] > 3 {
					continue
				}

				for _, p3 := range u {
					if p1 == p3 || p2 == p3 {
						continue
					}

					cell3 := *g.pt(p3)
					count := bitCount[cell3]
					if count == 1 || count > 3 {
						continue
					}

					comb := cell1 | cell2 | cell3
					if bitCount[comb] > 3 {
						continue
					}

					for _, p := range u {
						if p1 == p || p2 == p || p3 == p {
							continue
						}

						if g.pt(p).andNot(comb) {
							g.cellChange(&res, verbose, "nakedTriple: in %s %d (%s, %s, %s) removing %s from %s\n", gr.name, ui, p1, p2, p3, comb, p)
						}
					}
				}
			}
		}
	}

	return
}
