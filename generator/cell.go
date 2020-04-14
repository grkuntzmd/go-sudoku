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

func (c cell) lowestSetBit() int {
	for d := 1; d <= 9; d++ {
		if c&(1<<d) != 0 {
			return d
		}
	}

	return 0
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
