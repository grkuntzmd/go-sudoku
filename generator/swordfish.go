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

// swordfish finds and removes candidates. A swordfish is a 3 by 3 nine-cell pattern, where is in each column (or row), a candidate is only found in three different rows (or columns). The candidate can be removed from all other columns (or row) that line up with the three rows.
func (g *Grid) swordfish(verbose uint) bool {
	return g.swordfishGroup(&col, verbose) || g.swordfishGroup(&row, verbose)
}

func (g *Grid) swordfishGroup(gr *group, verbose uint) (res bool) {
	for p1i, p1s := range gr.unit {
		digits1 := g.digitPlaces(p1s)

		for p2i, p2s := range gr.unit {
			if p1i == p2i {
				continue
			}

			digits2 := g.digitPlaces(p2s)

			for p3i, p3s := range gr.unit {
				if p1i == p3i || p2i == p3i {
					continue
				}

				digits3 := g.digitPlaces(p3s)

				for d := 1; d <= 9; d++ {
					d1 := digits1[d]
					d2 := digits2[d]
					d3 := digits3[d]
					if bitCount[d1] < 2 || bitCount[d2] < 2 || bitCount[d3] < 2 || bitCount[d1|d2|d3] != 3 {
						continue
					}

					places := (d1 | d2 | d3).places()
					for pi, ps := range gr.unit {
						if p1i == pi || p2i == pi || p3i == pi {
							continue
						}

						for _, p := range places {
							if g.pt(ps[p]).andNot(1 << d) {
								g.cellChange(&res, verbose, "swordfish: (%d, %d, %d), in %s %d, removing %d from position %d (%s)\n", p1i, p2i, p3i, gr.name, pi, d, p, ps[p])
							}
						}
					}
				}
			}
		}
	}

	return
}
