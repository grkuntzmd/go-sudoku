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

// xyzWing removes candidates in a manor similar to yWing. The pivot point contains XYZ. If it can see two other cells XZ and YZ, then any cells that can be seen by all three can be cleared of Z.
func (g *Grid) xyzWing(verbose uint) (res bool) {
	for _, u := range box.unit {
		for _, p := range u { // Traverse all cells, using box units for convenience.
			cell := *g.pt(p)

			if bitCount[cell] != 3 {
				continue
			}

			n := neighbors(p)

			candidates := g.findYWingCandidates(p, 2)
			if len(candidates) < 2 {
				continue
			}

			for c1i, p1 := range candidates {
				cell1 := *g.pt(p1)
				n1 := neighbors(p1)

				for c2i, p2 := range candidates {
					if c1i == c2i {
						continue
					}

					cell2 := *g.pt(p2)

					if bitCount[cell1|cell2] != 3 {
						continue
					}

					n2 := neighbors(p2)

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
									g.cellChange(&res, verbose, "xyzWing: %s, %s, %s causes clearing %s from (%d, %d)\n", p, p1, p2, bits, r, c)
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
