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

// nakedSingle removes a solved digit from all other candidates in the same unit (box, row, or column) and returns true if it modifies the grid.
func (g *Grid) nakedSingle(verbose uint) bool {
	return g.nakedSingleGroup(&box, verbose) || g.nakedSingleGroup(&col, verbose) || g.nakedSingleGroup(&row, verbose)
}

func (g *Grid) nakedSingleGroup(gr *group, verbose uint) (res bool) {
	for ui, u := range gr.unit {
		for _, p1 := range u {
			cell := *g.pt(p1)
			if bitCount[cell] != 1 {
				continue
			}

			for _, p2 := range u {
				if p1 == p2 {
					continue
				}

				if g.pt(p2).andNot(cell) {
					g.cellChange(&res, verbose, "nakedSingle: in %s %d cell %s allows only %s, removed from %s\n", gr.name, ui, p1, cell, p2)
				}
			}
		}
	}

	return
}
