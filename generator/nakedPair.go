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

// nakedPair checks a group for 2 cells containing only the same pair of values. If present, those values can be removed from all other cells in the group. It returns true if it changes any cells.
func (g *Grid) nakedPair(verbose uint) bool {
	return g.nakedPairGroup(&box, verbose) || g.nakedPairGroup(&col, verbose) || g.nakedPairGroup(&row, verbose)
}

func (g *Grid) nakedPairGroup(gr *group, verbose uint) (res bool) {
	for ui, u := range gr.unit {
	outer:
		for _, p1 := range u {
			cell1 := *g.pt(p1)
			if bitCount[cell1] != 2 {
				continue
			}

			for _, p2 := range u {
				if p1 == p2 {
					continue
				}

				cell2 := *g.pt(p2)
				if cell1 != cell2 {
					continue
				}

				for _, p3 := range u {
					if p1 == p3 || p2 == p3 {
						continue
					}

					if g.pt(p3).andNot(cell1) {
						g.cellChange(&res, verbose, "nakedPair: in %s %d removed %s from %s\n", gr.name, ui, cell1, p3)
					}
				}
				continue outer
			}
		}
	}

	return
}
