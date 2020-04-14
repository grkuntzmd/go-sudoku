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

func TestHiddenQuadBox(t *testing.T) {
	g := decodeInts([]int{9, 37, 1, 5, 28, 28, 37, 4, 6, 4, 2, 5, 367, 9, 367, 37, 8, 1, 8, 6, 37, 347,
		1, 347, 59, 2, 59, 5, 3478, 2, 1346789, 378, 346789, 19, 37, 89, 37, 1, 9, 2378, 23578,
		23578, 4, 6, 58, 6, 3478, 3478, 134789, 3578, 345789, 159, 37, 2, 1, 9, 6, 78, 4, 78, 2,
		5, 3, 2, 345, 34, 39, 6, 359, 8, 1, 7, 37, 3578, 378, 23, 235, 1, 6, 9, 4})
	assert.True(t, g.hiddenQuadGroup(&box, 0))
	assert.Equal(t, []int{9, 37, 1, 5, 28, 28, 37, 4, 6, 4, 2, 5, 367, 9, 367, 37, 8, 1, 8, 6, 37,
		347, 1, 347, 59, 2, 59, 5, 3478, 2, 1469, 378, 469, 19, 37, 89, 37, 1, 9, 2378, 23578,
		23578, 4, 6, 58, 6, 3478, 3478, 149, 3578, 49, 159, 37, 2, 1, 9, 6, 78, 4, 78, 2, 5, 3, 2,
		345, 34, 39, 6, 359, 8, 1, 7, 37, 3578, 378, 23, 235, 1, 6, 9, 4}, g.encodeInts())
}

func TestHiddenQuadCol(t *testing.T) {
	g := decodeInts([]int{5679, 3, 279, 24, 45, 27, 2678, 1, 2456789, 567, 1256, 8, 234, 9, 127, 2367,
		356, 234567, 4, 125, 1279, 6, 135, 8, 237, 35, 23579, 138, 128, 123, 5, 7, 6, 9, 4, 138,
		67, 146, 147, 9, 8, 3, 5, 2, 167, 5679, 568, 379, 1, 2, 4, 3678, 368, 3678, 2, 7, 6, 348,
		34, 5, 1, 9, 38, 138, 148, 134, 7, 136, 9, 2368, 3568, 23568, 138, 9, 5, 238, 136, 12, 4,
		7, 368})
	assert.True(t, g.hiddenQuadGroup(&col, 0))
	assert.Equal(t, []int{5679, 3, 279, 24, 45, 27, 2678, 1, 2459, 567, 1256, 8, 234, 9, 127, 2367,
		356, 245, 4, 125, 1279, 6, 135, 8, 237, 35, 259, 138, 128, 123, 5, 7, 6, 9, 4, 138, 67,
		146, 147, 9, 8, 3, 5, 2, 167, 5679, 568, 379, 1, 2, 4, 3678, 368, 3678, 2, 7, 6, 348, 34,
		5, 1, 9, 38, 138, 148, 134, 7, 136, 9, 2368, 3568, 25, 138, 9, 5, 238, 136, 12, 4, 7,
		368}, g.encodeInts())
}

func TestHiddenQuadRow(t *testing.T) {
	g := decodeInts([]int{19, 5689, 1689, 3, 7, 4, 2, 1569, 18, 1379, 35679, 1369, 1569, 8, 2, 1356, 4,
		137, 123479, 2356789, 1234689, 1569, 16, 159, 1356, 13569, 1378, 1479, 79, 149, 157, 3,
		157, 8, 2, 6, 6, 237, 123, 12578, 9, 1578, 135, 135, 4, 8, 23, 5, 12, 4, 6, 9, 7, 13, 5, 4,
		7, 168, 2, 138, 136, 1368, 9, 239, 23689, 23689, 16789, 16, 13789, 4, 1368, 5, 39, 1,
		3689, 4, 5, 389, 7, 368, 2})
	assert.True(t, g.hiddenQuadGroup(&row, 0))
	assert.Equal(t, []int{19, 5689, 1689, 3, 7, 4, 2, 1569, 18, 1379, 35679, 1369, 1569, 8, 2,
		1356, 4, 137, 247, 278, 248, 1569, 16, 159, 1356, 13569, 78, 1479, 79, 149, 157, 3, 157, 8,
		2, 6, 6, 237, 123, 12578, 9, 1578, 135, 135, 4, 8, 23, 5, 12, 4, 6, 9, 7, 13, 5, 4, 7, 168,
		2, 138, 136, 1368, 9, 239, 23689, 23689, 16789, 16, 13789, 4, 1368, 5, 39, 1, 3689, 4, 5,
		389, 7, 368, 2}, g.encodeInts())
}
