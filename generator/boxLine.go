/*
 * MIT LICENSE
 *
 * Copyright © 2020, G.Ralph Kuntz, MD.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package generator

func (g *Grid) boxLine(verbose uint) bool {
	return g.boxLineGroup(&col, verbose, pointCol, pointRow, func(i, c int) int { return i*3 + c/3 }) ||
		g.boxLineGroup(&row, verbose, pointRow, pointCol, func(i, r int) int { return r/3*3 + i })
}

func (g *Grid) boxLineGroup(
	gr *group,
	verbose uint,
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
					g.cellChange(&res, verbose, "boxLine: all %d's in %s %d appear in box %d removing from %s\n", d, gr.name, ui, boxSel(index, ui), p)
				}
			}
		}
	}

	return
}
