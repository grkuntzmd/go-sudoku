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

// jellyfish finds and removes candidates. A jellyfish is a 4 by 4 sixteen-cell pattern, where is in each column (or row), a candidate is only found in four different rows (or columns). The candidate can be removed from all other columns (or row) that line up with the four rows. jellyfish extends swordfish from 3 units to 4.
func (g *Grid) jellyfish(verbose uint) bool {
	return g.jellyfishGroup(&col, verbose) || g.jellyfishGroup(&row, verbose)
}

func (g *Grid) jellyfishGroup(gr *group, verbose uint) (res bool) {
	for p1i, p1s := range gr.unit {
		digits1 := g.digitPlaces(p1s)

		for p2i, p2s := range gr.unit {
			if p1i == p2i {
				continue
			}

			digits2 := g.digitPlaces(p2s)

			for p3i, p3s := range gr.unit {
				if p1i == p3i || p2i == p3i {
					continue
				}

				digits3 := g.digitPlaces(p3s)

				for p4i, p4s := range gr.unit {
					if p1i == p4i || p2i == p4i || p3i == p4i {
						continue
					}

					digits4 := g.digitPlaces(p4s)

					for d := 1; d <= 9; d++ {
						d1 := digits1[d]
						d2 := digits2[d]
						d3 := digits3[d]
						d4 := digits4[d]
						if bitCount[d1] < 2 || bitCount[d2] < 2 || bitCount[d3] < 2 || bitCount[d4] < 2 || bitCount[d1|d2|d3|d4] != 4 {
							continue
						}

						places := (d1 | d2 | d3 | d4).places()
						for pi, ps := range gr.unit {
							if p1i == pi || p2i == pi || p3i == pi || p4i == pi {
								continue
							}

							for _, p := range places {
								if g.pt(ps[p]).andNot(1 << d) {
									g.cellChange(&res, verbose, "jellyfish: (%d, %d, %d, %d), in %s %d, removing %d from position %d (%s)\n", p1i, p2i, p3i, p4i, gr.name, pi, d, p, ps[p])
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
