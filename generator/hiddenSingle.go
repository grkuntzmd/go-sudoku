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

// hiddenSingle solves a cell if it contains the only instance of a digit within its group (box, column, row) and returns true if it changes any cells.
func (g *Grid) hiddenSingle() bool {
	return g.hiddenSingleGroup(&box) || g.hiddenSingleGroup(&col) || g.hiddenSingleGroup(&row)
}

func (g *Grid) hiddenSingleGroup(gr *group) (res bool) {
	for ui, u := range gr.unit {
		points := g.digitPoints(u)

		for d := 1; d <= 9; d++ {
			if len(points[d]) == 1 {
				p := points[d][0]
				if g.pt(p).setTo(1 << d) {
					g.cellChange(&res, "hiddenSingle: in %s %d set %s to %d\n", gr.name, ui, p, d)
				}
			}
		}
	}

	return
}
