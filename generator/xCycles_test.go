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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXCyclesNiceLoop(t *testing.T) {
	g := decodeInts([]int{59, 2, 4, 1, 35, 58, 6, 7, 389, 59, 6, 38, 238, 7, 258, 4, 1, 389, 7, 18,
		138, 9, 6, 4, 58, 2, 358, 2, 4, 6, 5, 9, 1, 3, 8, 7, 1, 3, 5, 4, 8, 7, 2, 9, 6, 8, 7, 9,
		6, 2, 3, 1, 5, 4, 4, 18, 128, 38, 35, 9, 7, 6, 258, 3, 5, 28, 7, 1, 6, 9, 4, 28, 6, 9,
		7, 28, 4, 258, 58, 3, 1})
	assert.True(t, g.xCycles(0))
	assert.Equal(t, []int{59, 2, 4, 1, 35, 58, 6, 7, 389, 59, 6, 38, 238, 7, 258, 4, 1, 389, 7,
		18, 13, 9, 6, 4, 58, 2, 35, 2, 4, 6, 5, 9, 1, 3, 8, 7, 1, 3, 5, 4, 8, 7, 2, 9, 6, 8,
		7, 9, 6, 2, 3, 1, 5, 4, 4, 18, 12, 38, 35, 9, 7, 6, 25, 3, 5, 28, 7, 1, 6, 9, 4, 28,
		6, 9, 7, 28, 4, 258, 58, 3, 1}, g.encodeInts())
}

func TestXCyclesStrongLoop(t *testing.T) {
	g := decodeInts([]int{8, 19, 4, 5, 3, 7, 169, 126, 12, 79, 2, 3, 6, 1, 4, 79, 8, 5, 6, 17, 5, 9,
		8, 2, 17, 3, 4, 349, 346, 269, 1, 469, 5, 8, 7, 29, 5, 49, 12, 7, 49, 8, 3, 12, 6,
		179, 8, 1679, 2, 69, 3, 4, 5, 19, 2, 467, 167, 8, 5, 9, 16, 146, 3, 49, 5, 69, 3, 7,
		1, 2, 469, 8, 139, 39, 8, 4, 2, 6, 5, 19, 7})
	assert.True(t, g.xCycles(0))
	assert.Equal(t, []int{8, 19, 4, 5, 3, 7, 169, 126, 12, 79, 2, 3, 6, 1, 4, 79, 8, 5, 6, 17,
		5, 9, 8, 2, 17, 3, 4, 349, 346, 269, 1, 469, 5, 8, 7, 29, 5, 49, 12, 7, 49, 8, 3, 12,
		6, 179, 8, 1679, 2, 69, 3, 4, 5, 19, 2, 467, 167, 8, 5, 9, 16, 146, 3, 49, 5, 69, 3,
		7, 1, 2, 469, 8, 1, 39, 8, 4, 2, 6, 5, 19, 7}, g.encodeInts())
}

func TestXCyclesWeakLoop(t *testing.T) {
	g := decodeInts([]int{2478, 23, 247, 357, 1, 357, 9, 6, 28, 127, 1239, 1279, 6, 8, 37, 4, 5, 12,
		18, 5, 6, 9, 4, 2, 3, 18, 7, 1247, 126, 12457, 157, 36, 8, 17, 137, 9, 3, 8, 17, 17,
		9, 4, 6, 2, 5, 9, 16, 157, 2, 36, 157, 178, 1378, 4, 6, 7, 3, 18, 2, 9, 5, 4, 18, 5,
		129, 8, 4, 7, 6, 12, 19, 3, 12, 4, 129, 138, 5, 13, 1278, 1789, 6})
	assert.True(t, g.xCycles(0))
	assert.Equal(t, []int{2478, 23, 247, 357, 1, 357, 9, 6, 28, 127, 1239, 1279, 6, 8, 37, 4, 5,
		12, 18, 5, 6, 9, 4, 2, 3, 18, 7, 1247, 126, 12457, 157, 36, 8, 17, 137, 9, 3, 8, 17,
		17, 9, 4, 6, 2, 5, 9, 16, 157, 2, 36, 157, 178, 1378, 4, 6, 7, 3, 18, 2, 9, 5, 4, 18,
		5, 129, 8, 4, 7, 6, 12, 19, 3, 12, 4, 29, 138, 5, 13, 1278, 1789, 6}, g.encodeInts())
}

func BenchmarkXCycles(b *testing.B) {
	for n := 0; n < b.N; n++ {
		g := decodeInts([]int{59, 2, 4, 1, 35, 58, 6, 7, 389, 59, 6, 38, 238, 7, 258, 4, 1, 389, 7, 18,
			138, 9, 6, 4, 58, 2, 358, 2, 4, 6, 5, 9, 1, 3, 8, 7, 1, 3, 5, 4, 8, 7, 2, 9, 6, 8, 7,
			9, 6, 2, 3, 1, 5, 4, 4, 18, 128, 38, 35, 9, 7, 6, 258, 3, 5, 28, 7, 1, 6, 9, 4, 28, 6,
			9, 7, 28, 4, 258, 58, 3, 1})
		g.xCycles(0)
	}
	for n := 0; n < b.N; n++ {
		g := decodeInts([]int{8, 19, 4, 5, 3, 7, 169, 126, 12, 79, 2, 3, 6, 1, 4, 79, 8, 5, 6, 17, 5,
			9, 8, 2, 17, 3, 4, 349, 346, 269, 1, 469, 5, 8, 7, 29, 5, 49, 12, 7, 49, 8, 3, 12, 6,
			179, 8, 1679, 2, 69, 3, 4, 5, 19, 2, 467, 167, 8, 5, 9, 16, 146, 3, 49, 5, 69, 3, 7,
			1, 2, 469, 8, 139, 39, 8, 4, 2, 6, 5, 19, 7})
		g.xCycles(0)
	}
	for n := 0; n < b.N; n++ {
		g := decodeInts([]int{2478, 23, 247, 357, 1, 357, 9, 6, 28, 127, 1239, 1279, 6, 8, 37, 4, 5,
			12, 18, 5, 6, 9, 4, 2, 3, 18, 7, 1247, 126, 12457, 157, 36, 8, 17, 137, 9, 3, 8, 17,
			17, 9, 4, 6, 2, 5, 9, 16, 157, 2, 36, 157, 178, 1378, 4, 6, 7, 3, 18, 2, 9, 5, 4, 18,
			5, 129, 8, 4, 7, 6, 12, 19, 3, 12, 4, 129, 138, 5, 13, 1278, 1789, 6})
		g.xCycles(0)
	}
}
