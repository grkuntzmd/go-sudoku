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

// hiddenTriple removes other digits from a triple of cells in a group (box, column, row) when that triple contains the only occurrances of the digits in the group. It returns true if it changes any cells.
func (g *Grid) hiddenTriple(verbose uint) bool {
	return g.hiddenTripleGroup(&box, verbose) || g.hiddenTripleGroup(&col, verbose) || g.hiddenTripleGroup(&row, verbose)
}

func (g *Grid) hiddenTripleGroup(gr *group, verbose uint) (res bool) {
	for ui, u := range gr.unit {
		places := g.digitPlaces(u)

		for d1 := 1; d1 <= 9; d1++ {
			p1 := places[d1]
			count := bitCount[p1]
			if count == 1 || count > 3 {
				continue
			}

			for d2 := 1; d2 <= 9; d2++ {
				if d1 == d2 {
					continue
				}

				p2 := places[d2]
				count := bitCount[p2]
				if count == 1 || count > 3 || bitCount[p1|p2] > 3 {
					continue
				}

				for d3 := 1; d3 <= 9; d3++ {
					if d1 == d3 || d2 == d3 {
						continue
					}

					p3 := places[d3]
					count := bitCount[p3]
					comb := p1 | p2 | p3
					if count == 1 || count > 3 || bitCount[comb] != 3 {
						continue
					}

					points := [3]point{}
					index := 0
					for pi, p := range u {
						if comb&(1<<pi) != 0 {
							points[index] = p
							index++
						}
					}

					bits := cell(1<<d1 | 1<<d2 | 1<<d3)
					for _, p := range points {
						if g.pt(p).and(bits) {
							g.cellChange(&res, verbose, "hiddenTriple: in %s %d limits %s (triple: %s, %s, %s) to %s\n", gr.name, ui, p, points[0], points[1], points[2], bits)
						}
					}
				}
			}
		}
	}

	return
}
