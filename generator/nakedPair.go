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

// nakedPair checks a group for 2 cells containing only the same pair of values. If present, those values can be removed from all other cells in the group. It returns true if it changes any cells.
func (g *Grid) nakedPair() bool {
	return g.nakedPairGroup(&box) || g.nakedPairGroup(&col) || g.nakedPairGroup(&row)
}

func (g *Grid) nakedPairGroup(gr *group) (res bool) {
	for ui, u := range gr.unit {
	outer:
		for _, p1 := range u {
			cell1 := *g.pt(&p1)
			if bitCount[cell1] != 2 {
				continue
			}

			for _, p2 := range u {
				if p1 == p2 {
					continue
				}

				cell2 := *g.pt(&p2)
				if cell1 != cell2 {
					continue
				}

				for _, p3 := range u {
					if p1 == p3 || p2 == p3 {
						continue
					}

					if g.pt(&p3).andNot(cell1) {
						g.cellChange(&res, "nakedPair: in %s %d removed %s from %s", gr.name, ui, cell1, &p3)
					}
				}
				continue outer
			}
		}
	}

	return
}
