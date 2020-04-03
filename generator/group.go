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

func (u *uint128) unset(p point) {
	bit := p.r*9 + p.c
	if bit < 64 {
		u.ls &^= 1 << bit
	} else {
		u.ms &^= 1 << (bit - 64)
	}
}

func boxOf(r, c uint8) uint8 {
	return r/3*3 + c/3
}

func boxOfPoint(p point) uint8 {
	return boxOf(p.r, p.c)
}

func processInfluence(overlap uint128, f func(uint8, uint8)) {
	for r := zero; r < rows; r++ {
		for c := zero; c < cols; c++ {
			bit := r*9 + c
			if bit < 64 {
				if overlap.ls&(1<<bit) != 0 {
					f(r, c)
				}
			} else {
				if overlap.ms&(1<<(bit-64)) != 0 {
					f(r, c)
				}
			}
		}
	}
}
