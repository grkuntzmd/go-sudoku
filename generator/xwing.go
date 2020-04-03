/*
 * Copyright © 2020, G.Ralph Kuntz, MD.
 *
 * Licensed under the Apache License, Version 2.0(the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIC
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
