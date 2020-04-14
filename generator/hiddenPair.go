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

// hiddenPair removes other digits from a pair of cells in a group (box, column, row) when that pair contains the only occurrances of the digits in the group and returns true if it changes any cells.
func (g *Grid) hiddenPair(verbose uint) bool {
	return g.hiddenPairGroup(&box, verbose) || g.hiddenPairGroup(&col, verbose) || g.hiddenPairGroup(&row, verbose)

}

func (g *Grid) hiddenPairGroup(gr *group, verbose uint) (res bool) {
	for ui, u := range gr.unit {
		points := g.digitPoints(u)

		for d1 := 1; d1 <= 9; d1++ {
			for d2 := 1; d2 <= 9; d2++ {
				if d1 == d2 || len(points[d1]) != 2 || len(points[d2]) != 2 {
					continue
				}

				if comparePointSlices(points[d1], points[d2]) {
					comb := cell(1<<d1 | 1<<d2)
					for _, p := range points[d1] {
						if g.pt(p).and(comb) {
							g.cellChange(&res, verbose, "hiddenPair: in %s %d limits %s (pair: %s, %s) to %s\n", gr.name, ui, p, points[d1][0], points[d1][1], comb)
						}
					}
				}
			}
		}
	}

	return
}
