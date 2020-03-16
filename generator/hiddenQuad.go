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

// hiddenQuad removes other digits from a quad of cells in a group (box, column, row) when that quad contains the only occurrances of the digits in the group. It returns true if it changes any cells.
func (g *Grid) hiddenQuad() bool {
	return g.hiddenQuadGroup(&box) || g.hiddenQuadGroup(&col) || g.hiddenQuadGroup(&row)
}

func (g *Grid) hiddenQuadGroup(gr *group) (res bool) {
	for ui, u := range gr.unit {
		places := g.digitPlaces(u)

		for d1 := 1; d1 <= 9; d1++ {
			p1 := places[d1]
			count := bitCount[p1]
			if count == 1 || count > 4 {
				continue
			}

			for d2 := 1; d2 <= 9; d2++ {
				if d1 == d2 {
					continue
				}

				p2 := places[d2]
				count := bitCount[p2]
				if count == 1 || count > 4 || bitCount[p1|p2] > 4 {
					continue
				}

				for d3 := 1; d3 <= 9; d3++ {
					if d1 == d3 || d2 == d3 {
						continue
					}

					p3 := places[d3]
					count := bitCount[p3]
					if count == 1 || count > 4 || bitCount[p1|p2|p3] != 4 {
						continue
					}

					for d4 := 1; d4 <= 9; d4++ {
						if d1 == d4 || d2 == d4 || d3 == d4 {
							continue
						}

						p4 := places[d4]
						count := bitCount[d4]
						comb := p1 | p2 | p3 | p4
						if count == 1 || count > 4 || bitCount[comb] != 4 {
							continue
						}

						points := [4]point{}
						index := 0
						for pi, p := range u {
							if comb&(1<<pi) != 0 {
								points[index] = p
								index++
							}
						}

						bits := cell(1<<d1 | 1<<d2 | 1<<d3 | 1<<d4)
						for _, p := range points {
							if g.pt(&p).and(bits) {
								g.cellChange(&res, "hiddenQuad: in %s %d limits %s (quad: %s, %s, %s, %s) to %s\n", gr.name, ui, &p, &points[0], &points[1], &points[2], &points[3], bits)
							}
						}
					}
				}
			}
		}
	}

	return
}
