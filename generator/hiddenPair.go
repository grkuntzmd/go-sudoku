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

// hiddenPair removes other digits from a pair of cells in a group (box, column, row) when that pair contains the only occurrances of the digits in the group and returns true if it changes any cells.
func (g *Grid) hiddenPair() bool {
	return g.hiddenPairGroup(&box) || g.hiddenPairGroup(&col) || g.hiddenPairGroup(&row)

}

func (g *Grid) hiddenPairGroup(gr *group) (res bool) {
	for pi, ps := range gr.points {
		points := g.digitPoints(ps)

		for d1 := 1; d1 <= 9; d1++ {
			for d2 := 1; d2 <= 9; d2++ {
				if d1 == d2 || len(points[d1]) != 2 || len(points[d2]) != 2 {
					continue
				}

				if comparePointSlices(points[d1], points[d2]) {
					comb := cell(1<<d1 | 1<<d2)
					for k := 0; k < 2; k++ {
						p := points[d1][k]
						if g.pt(p).and(comb) {
							g.cellChange(&res, "hiddenPair: in %s %d limits %s (pair: %s, %s) to %s\n", gr.name, pi, &p, &points[d1][0], &points[d1][1], comb)
						}
					}
				}
			}
		}
	}

	return
}
