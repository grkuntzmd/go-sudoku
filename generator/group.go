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

import (
	"fmt"
	"strings"
)

type (
	group struct {
		name string
		unit
	}

	unit [rows][cols]point

	uint128 struct {
		ms, ls uint64
	}
)

var (
	box = group{name: "box"} // These are all of the coordinates in a box (first dimension).
	col = group{name: "col"} // These are all of the coordinates in a column (first dimension).
	row = group{name: "row"} // These are all of the coordinates in a row (first dimension).

	influence [rows][cols]uint128 // These bit masks contain a 1 for locations that the point "shadows" (can see).
)

func init() {
	for r := zero; r < rows; r++ {
		for c := zero; c < cols; c++ {
			p := point{r, c}
			box.unit[boxOf(r, c)][r%3*3+c%3] = p
			col.unit[c][r] = p
			row.unit[r][c] = p
		}
	}

	for r := zero; r < rows; r++ {
		for c := zero; c < cols; c++ {
			for _, p := range box.unit[boxOf(r, c)] {
				if p.r == r && p.c == c {
					continue
				}

				bit := p.r*9 + p.c
				if bit < 64 {
					influence[r][c].ls |= (1 << bit)
				} else {
					influence[r][c].ms |= (1 << (bit - 64))
				}
			}
			for _, p := range col.unit[c] {
				if p.r == r && p.c == c {
					continue
				}

				bit := p.r*9 + p.c
				if bit < 64 {
					influence[r][c].ls |= (1 << bit)
				} else {
					influence[r][c].ms |= (1 << (bit - 64))
				}
			}
			for _, p := range row.unit[r] {
				if p.r == r && p.c == c {
					continue
				}

				bit := p.r*9 + p.c
				if bit < 64 {
					influence[r][c].ls |= (1 << bit)
				} else {
					influence[r][c].ms |= (1 << (bit - 64))
				}
			}
		}
	}
}

func (u uint128) and(other uint128) uint128 {
	return uint128{u.ms & other.ms, u.ls & other.ls}
}

func (u uint128) process(f func(uint8, uint8)) {
	for r := zero; r < rows; r++ {
		for c := zero; c < cols; c++ {
			bit := r*9 + c
			if bit < 64 {
				if u.ls&(1<<bit) != 0 {
					f(r, c)
				}
			} else {
				if u.ms&(1<<(bit-64)) != 0 {
					f(r, c)
				}
			}
		}
	}
}

func (u uint128) String() string {
	var b strings.Builder
	for r := zero; r < rows; r++ {
		for c := zero; c < cols; c++ {
			bit := r*9 + c
			if bit < 64 {
				if u.ls&(1<<bit) != 0 {
					fmt.Fprint(&b, "1")
				} else {
					fmt.Fprint(&b, "0")
				}
			} else {
				if u.ms&(1<<(bit-64)) != 0 {
					fmt.Fprint(&b, "1")
				} else {
					fmt.Fprint(&b, "0")
				}
			}
		}
		fmt.Fprintln(&b)
	}
	return b.String() //fmt.Sprintf("%64.64b%64.64b", u.ms, u.ls)
}

func (u *uint128) unset(p point) *uint128 {
	bit := p.r*9 + p.c
	if bit < 64 {
		u.ls &^= 1 << bit
	} else {
		u.ms &^= 1 << (bit - 64)
	}

	return u
}

func boxOf(r, c uint8) uint8 {
	return r/3*3 + c/3
}

func boxOfPoint(p point) uint8 {
	return boxOf(p.r, p.c)
}
