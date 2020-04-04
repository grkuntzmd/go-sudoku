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

import (
	"fmt"
)

type wxyzCandidate struct {
	p1, p2 point
	common cell
}

func (w wxyzCandidate) String() string {
	return fmt.Sprintf("{%s, %s, %s}", w.p1, w.p2, w.common)
}

func (g *Grid) wxyzWing(verbose uint) (res bool) {
	for r := zero; r < rows; r++ {
		for c := zero; c < cols; c++ {
			p := point{r, c}
			cell := g.cells[r][c]
			if bitCount[cell] < 2 {
				continue
			}

			var points []point
			points = append(points, box.unit[boxOf(r, c)][:]...)
			points = append(points, col.unit[c][:]...)
			points = append(points, row.unit[r][:]...)
			for _, p1 := range points {
				if p == p1 {
					continue
				}

				cell1 := *g.pt(p1)

				for _, p2 := range points {
					if p == p2 || p1 == p2 ||
						boxOfPoint(p1) != boxOfPoint(p2) && p1.c != p2.c && p1.r != p2.c {
						continue
					}

					cell2 := *g.pt(p2)

					for _, p3 := range points {
						if p == p3 || p1 == p3 || p2 == p3 {
							continue
						}

						if boxOfPoint(p1) == boxOfPoint(p3) || boxOfPoint(p2) == boxOfPoint(p3) ||
							p1.c == p3.c || p2.c == p3.c ||
							p1.r == p3.r || p2.r == p3.r {
							continue
						}

						cell3 := *g.pt(p3)
						unrestricted := (cell1 | cell2) & cell3

						if bitCount[cell&cell1] < 2 ||
							bitCount[cell&cell2] < 2 ||
							bitCount[cell&cell3] < 2 ||
							bitCount[unrestricted] != 1 {
							continue
						}

						overlap := influence[p.r][p.c].and(influence[p1.r][p1.c]).and(influence[p2.r][p2.r]).and(influence[p3.r][p3.c])
						overlap.unset(p).unset(p1).unset(p2).unset(p3)
						overlap.process(func(r, c uint8) {
							if g.pt(point{r, c}).andNot(unrestricted) {
								g.cellChange(&res, verbose, "wxyzWing: removing %s from (%d, %d) because of %s, %s, %s, %s\n", unrestricted, r, c, p, p1, p2, p3)
							}
						})
					}
				}
			}
		}
	}

	return
}
