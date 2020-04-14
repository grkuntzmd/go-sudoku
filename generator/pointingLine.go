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

// pointingLine removes candidates. When a candidate within a box appears only in a single column or row, that candidate can be removed from all cells in the column or row outside of the box. It returns true if it changes any cells.
func (g *Grid) pointingLine(verbose uint) bool {
	return g.pointingLineGroup("col", verbose, func(p point) *[9]point {
		return &col.unit[p.c]
	}, func(p point) uint8 {
		return p.c
	}) || g.pointingLineGroup("row", verbose, func(p point) *[9]point {
		return &row.unit[p.r]
	}, func(p point) uint8 {
		return p.r
	})
}

func (g *Grid) pointingLineGroup(
	gr string,
	verbose uint,
	sel func(point) *[9]point,
	axis func(point) uint8,
) (res bool) {
	for ui, u := range box.unit {
		points := g.digitPoints(u)

		// Loop through the digits and determine if all of them are on the same line (col or row). If so, then all other cells in that line that are not in the current box can have those digits removed.
	outer:
		for d := 1; d <= 9; d++ {
			if len(points[d]) == 0 {
				return false
			}

			a := axis(points[d][0])
			for _, p := range points[d][1:] {
				if axis(p) != a { // All points for each digit must be on the same axis (row or column).
					continue outer
				}
			}

			for _, p := range sel(points[d][0]) {
				if int(p.r/3*3+p.c/3) == ui { // If the point is on the axis (row or column) of interest, skip it.
					continue
				}

				if g.pt(p).andNot(1 << d) {
					g.cellChange(&res, verbose, "pointingLine: in box %d removing %d from %s along %s %d\n", ui, d, p, gr, a)
				}
			}
		}
	}

	return
}
