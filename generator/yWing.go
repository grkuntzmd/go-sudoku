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

// yWing removes candidates. If a cell has two candidates (AB) and in a neighboring unit (box, row, or column) of AB is another cell containing AC and in a second neighboring unit of AB is a cell containing BC, then any cell that can be "seen" by AC and BC (in both neighborhoods of AC and BC) that contain C can have C removed. It returns true if it changes any cells.
func (g *Grid) yWing() (res bool) {
	for _, u := range box.unit {
		for _, p := range u { // Traverse all cells, using box units for convenience.
			cell := *g.pt(&p)

			if bitCount[cell] != 2 {
				continue
			}

			candidates := g.findCandidates(&p)
			length := len(candidates)

			for c1i := 0; c1i < length; c1i++ {
				p1 := candidates[c1i]
				cell1 := *g.pt(&p1)
				n1 := neighbors(&p1)

				for c2i := 0; c2i < length; c2i++ {
					if c1i == c2i {
						continue
					}

					p2 := candidates[c2i]
					cell2 := *g.pt(&p2)

					if bitCount[cell1|cell2] != 3 || cell&cell1|cell&cell2 != cell {
						continue
					}

					n2 := neighbors(&p2)

					var overlap [9][9]bool
					for r := 0; r < rows; r++ {
						for c := 0; c < cols; c++ {
							if n1[r][c] && n2[r][c] {
								overlap[r][c] = true
							}
						}
					}

					overlap[p.r][p.c] = false

					for r := 0; r < rows; r++ {
						for c := 0; c < cols; c++ {
							if overlap[r][c] {
								bits := (cell1 | cell2) &^ cell
								if (&g.cells[r][c]).andNot(bits) {
									g.cellChange(&res, "yWing: %s, %s, %s causes clearing %s from (%d, %d)\n", &p, &p1, &p2, bits, r, c)
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

func (g *Grid) findCandidates(curr *point) (res []point) {
	res = append(res, g.findCandidatesUnit(&box.unit[boxOf(curr.r, curr.c)], curr)...)
	res = append(res, g.findCandidatesUnit(&col.unit[curr.c], curr)...)
	res = append(res, g.findCandidatesUnit(&row.unit[curr.r], curr)...)
	return
}

func (g *Grid) findCandidatesUnit(u *[9]point, curr *point) (res []point) {
	cell := *g.pt(curr)
	for _, p := range u {
		if p == *curr {
			continue
		}

		candidate := *g.pt(&p)
		if bitCount[candidate] != 2 || bitCount[cell&candidate] != 1 {
			continue
		}

		res = append(res, p)
	}

	return
}

func neighbors(curr *point) *[9][9]bool {
	var res [9][9]bool
	for _, u := range []*[9]point{&box.unit[boxOf(curr.r, curr.c)], &col.unit[curr.c], &row.unit[curr.r]} {
		for _, p := range u {
			if p == *curr {
				continue
			}

			res[p.r][p.c] = true
		}
	}

	return &res
}
