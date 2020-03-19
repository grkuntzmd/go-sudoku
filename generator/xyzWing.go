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

func (g *Grid) xyzWing() (res bool) {
	for _, u := range box.unit {
		for _, p := range u { // Traverse all cells, using box units for convenience.
			cell := *g.pt(&p)

			if bitCount[cell] != 3 {
				continue
			}

			n := neighbors(&p)

			candidates := g.findCandidates(&p, 2)
			if len(candidates) < 2 {
				continue
			}

			for c1i, p1 := range candidates {
				cell1 := *g.pt(&p1)
				n1 := neighbors(&p1)

				for c2i, p2 := range candidates {
					if c1i == c2i {
						continue
					}

					cell2 := *g.pt(&p2)

					if bitCount[cell1|cell2] != 3 {
						continue
					}

					n2 := neighbors(&p2)

					var overlap [9][9]bool
					for r := 0; r < rows; r++ {
						for c := 0; c < cols; c++ {
							overlap[r][c] = n[r][c] && n1[r][c] && n2[r][c]
						}
					}

					overlap[p.r][p.c] = false

					for r := 0; r < rows; r++ {
						for c := 0; c < cols; c++ {
							if overlap[r][c] {
								bits := cell1 & cell2 & cell
								if (&g.cells[r][c]).andNot(bits) {
									g.cellChange(&res, "xyzWing: %s, %s, %s causes clearing %s from (%d, %d)\n", &p, &p1, &p2, bits, r, c)
								}
							}
						}
					}
				}
			}
		}
	}

	return
}
