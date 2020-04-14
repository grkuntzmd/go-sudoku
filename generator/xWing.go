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

// xWing removes candidates. If in 2 columns, say 0 and 7, all instances of a particular digit, say 4, appear in the same two rows, say 4 and 6, then 1 of the 4's must be in (0, 4) or (0, 6) and the other in (7, 4) or (7, 6). Therefore all of the other 4's in those two rows can be removed. The same logic applies if rows and columns are swapped. It returns true if it changes any cells.
func (g *Grid) xWing(verbose uint) bool {
	return g.xWingGroup(&col, &row, verbose) || g.xWingGroup(&row, &col, verbose)
}

func (g *Grid) xWingGroup(majorGroup, minorGroup *group, verbose uint) (res bool) {
	var digits [9][10]cell
	for ui, u := range majorGroup.unit {
		for pi, p := range u {
			cell := *g.pt(p)
			for d := 1; d <= 9; d++ {
				if cell&(1<<d) != 0 {
					digits[ui][d] |= 1 << pi
				}
			}
		}
	}

	for d := 1; d <= 9; d++ {
		for c1i := 0; c1i < 9; c1i++ {
			for c2i := 0; c2i < 9; c2i++ {
				if c1i == c2i {
					continue
				}

				proto := digits[c1i][d]
				if bitCount[proto] == 2 && proto == digits[c2i][d] {
					for minor := 1; minor <= 9; minor++ {
						if proto&(1<<minor) != 0 {
							for mi, m := range minorGroup.unit[minor] {
								if mi == c1i || mi == c2i {
									continue
								}

								if g.pt(m).andNot(1 << d) {
									g.cellChange(&res, verbose, "xWing: in %ss %d and %d, %d appears only in %s %d and 1 other; "+
										"removing from %s\n", majorGroup.name, c1i, c2i, d, minorGroup.name, minor, m)
								}
							}
						}
					}
				}
			}
		}
	}

	return
}
