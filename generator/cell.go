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

type cell uint16

var bitCount [1024]int

func init() {
	// Brian Kernighan's algorithm to count bits set to 1.
	for i := 0; i < 1024; i++ {
		n := i
		c := 0
		for n != 0 {
			n &= n - 1
			c++
		}
		bitCount[i] = c
	}
}

func (c *cell) and(o cell) bool {
	prev := *c
	*c &= o
	return *c != prev
}

func (c *cell) andNot(o cell) bool {
	prev := *c
	*c &^= o
	return *c != prev
}

// digits returns a slice containing all the candidate digits in a cell as individual ints.
func (c cell) digits() []int {
	ds := make([]int, 0, 9)
	for d := 1; d <= 9; d++ {
		if c&(1<<d) != 0 {
			ds = append(ds, d)
		}
	}

	return ds
}

func (c *cell) setTo(o cell) bool {
	prev := *c
	*c = o
	return *c != prev
}

func (c cell) String() string {
	var b strings.Builder
	for d := 1; d <= 9; d++ {
		if c&(1<<d) != 0 {
			fmt.Fprintf(&b, "%d", d)
		}
	}
	return b.String()
}
