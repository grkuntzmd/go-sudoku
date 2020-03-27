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

func (g *Grid) boxLine() bool {
	return g.boxLineGroup(&col, pointCol, pointRow, func(i, c int) int { return i*3 + c/3 }) ||
		g.boxLineGroup(&row, pointRow, pointCol, func(i, r int) int { return r/3*3 + i })
}

func (g *Grid) boxLineGroup(
	gr *group,
	major func(point) uint8,
	minor func(point) uint8,
	boxSel func(int, int) int,
) (res bool) {
	for ui, u := range gr.unit {
		boxes := [10][3]bool{}

		for _, p := range u {
			cell := *g.pt(p)
			for d := 1; d <= 9; d++ {
				if cell&(1<<d) != 0 {
					boxes[d][minor(p)/3] = true
				}
			}
		}

		for d := 1; d <= 9; d++ {
			var index int
			if boxes[d][0] && !boxes[d][1] && !boxes[d][2] {
				index = 0
			} else if !boxes[d][0] && boxes[d][1] && !boxes[d][2] {
				index = 1
			} else if !boxes[d][0] && !boxes[d][1] && boxes[d][2] {
				index = 2
			} else {
				continue
			}

			for i := 0; i < 9; i++ {
				p := box.unit[boxSel(index, ui)][i]

				if major(p) == major(u[index]) {
					continue
				}

				if g.pt(p).andNot(1 << d) {
					g.cellChange(&res, "boxLine: all %d's in %s %d appear in box %d removing from %s\n", d, gr.name, ui, boxSel(index, ui), p)
				}
			}
		}
	}

	return
}
