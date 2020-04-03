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

// nakedQuad checks a group for 4 cells with the same quad of values. If present, those values can be removed from all other cells in the group. It returns true if it changes any cells.
func (g *Grid) nakedQuad(verbose uint) bool {
	return g.nakedQuadGroup(&box, verbose) || g.nakedQuadGroup(&col, verbose) || g.nakedQuadGroup(&row, verbose)

}

func (g *Grid) nakedQuadGroup(gr *group, verbose uint) (res bool) {
	for ui, u := range gr.unit {
		for _, p1 := range u {
			cell1 := *g.pt(p1)
			count := bitCount[cell1]
			if count == 1 || count > 4 {
				continue
			}

			for _, p2 := range u {
				if p1 == p2 {
					continue
				}

				cell2 := *g.pt(p2)
				count := bitCount[cell2]
				if count == 1 || count > 4 {
					continue
				}

				if bitCount[cell1|cell2] > 4 {
					continue
				}

				for _, p3 := range u {
					if p1 == p3 || p2 == p3 {
						continue
					}

					cell3 := *g.pt(p3)
					count := bitCount[cell3]
					if count == 1 || count > 4 {
						continue
					}

					if bitCount[cell1|cell2|cell3] > 4 {
						continue
					}

					for _, p4 := range u {
						if p1 == p4 || p2 == p4 || p3 == p4 {
							continue
						}

						cell4 := *g.pt(p4)
						count := bitCount[cell4]
						if count == 1 || count > 4 {
							continue
						}

						comb := cell1 | cell2 | cell3 | cell4
						if bitCount[comb] > 4 {
							continue
						}

						for _, p := range u {
							if p1 == p || p2 == p || p3 == p || p4 == p {
								continue
							}

							if g.pt(p).andNot(comb) {
								g.cellChange(&res, verbose, "nakedQuad: in %s %d (%s, %s, %s, %s) removing %s from %s\n", gr.name, ui, p1, p2, p3, p4, comb, p)
							}
						}
					}
				}
			}
		}
	}

	return
}
