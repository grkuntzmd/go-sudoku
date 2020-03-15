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

// nakedSingle removes a solved digit from all other candidates in the same unit (box, row, or column) and returns true if it modifies the grid.
func (g *Grid) nakedSingle() bool {
	return g.nakedSingleGroup(&box) || g.nakedSingleGroup(&col) || g.nakedSingleGroup(&row)
}

func (g *Grid) nakedSingleGroup(gr *group) (res bool) {
	for pi, ps := range gr.points {
		for _, p1 := range ps {
			cell := *g.pt(p1)
			if bitCount[cell] != 1 {
				continue
			}

			for _, p2 := range ps {
				if p1 == p2 {
					continue
				}

				if g.pt(p2).andNot(cell) {
					g.cellChange(&res, "nakedSingle: in %s %d cell %s allows only %s, removed from %s", gr.name, pi, &p1, cell, &p2)
				}
			}
		}
	}

	return
}
