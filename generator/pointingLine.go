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

// pointingLine removes candidates. When a candidate within a box appears only in a single column or row, that candidate can be removed from all cells in the column or row outside of the box. It returns true if it changes any cells.
func (g *Grid) pointingLine() bool {
	return g.pointingLineGroup("col", func(p point) *[9]point {
		return &col.unit[p.c]
	}, func(p point) uint8 {
		return p.c
	}) || g.pointingLineGroup("row", func(p point) *[9]point {
		return &row.unit[p.r]
	}, func(p point) uint8 {
		return p.r
	})
}

func (g *Grid) pointingLineGroup(
	gr string,
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
					g.cellChange(&res, "pointingLine: in box %d removing %d from %s along %s %d\n", ui, d, p, gr, a)
				}
			}
		}
	}

	return
}
