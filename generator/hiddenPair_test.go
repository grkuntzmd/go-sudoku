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

func TestHiddenPairBox(t *testing.T) {
	g := decodeInts([]int{2689, 268, 4, 5, 7, 1689, 18, 3, 1689, 1, 3568, 35689, 4, 69, 689, 578, 2,
		6789, 7, 568, 5689, 169, 2, 3, 4, 15689, 1689, 2356, 12356, 356, 12369, 8, 125679, 127,
		179, 4, 2568, 12568, 7, 1269, 1569, 4, 128, 189, 3, 4, 9, 38, 123, 13, 127, 6, 178, 5,
		5689, 4, 2, 169, 1569, 1569, 3, 15678, 1678, 3568, 3568, 3568, 7, 13456, 1256, 9, 14568,
		1268, 3569, 7, 1, 8, 34569, 2569, 25, 456, 26})
	assert.True(t, g.hiddenPairGroup(&box, 0))
	assert.Equal(t, []int{2689, 268, 4, 5, 7, 1689, 18, 3, 1689, 1, 3568, 35689, 4, 69, 689,
		578, 2, 6789, 7, 568, 5689, 169, 2, 3, 4, 15689, 1689, 2356, 12356, 356, 12369, 8, 125679,
		127, 179, 4, 2568, 12568, 7, 1269, 1569, 4, 128, 189, 3, 4, 9, 38, 123, 13, 127, 6, 178,
		5, 5689, 4, 2, 169, 1569, 1569, 3, 15678, 1678, 3568, 3568, 3568, 7, 34, 1256, 9, 14568,
		1268, 3569, 7, 1, 8, 34, 2569, 25, 456, 26}, g.encodeInts())
}

func TestHiddenPairCol(t *testing.T) {
	g := decodeInts([]int{6, 24, 5, 3, 2489, 478, 18, 49, 178, 8, 34, 1, 79, 5, 6, 49, 2, 37, 23, 7,
		9, 48, 1, 248, 5, 6, 38, 4, 28, 6, 19, 89, 5, 3, 7, 12, 59, 1, 28, 48, 7, 3, 249, 459, 6,
		7, 59, 3, 2, 6, 14, 149, 8, 145, 23, 6, 7, 5, 2348, 248, 28, 1, 9, 59, 59, 4, 17, 28, 17,
		6, 3, 28, 1, 238, 28, 6, 23, 9, 7, 45, 45})
	assert.True(t, g.hiddenPairGroup(&col, 0))
	assert.Equal(t, []int{6, 24, 5, 3, 2489, 478, 18, 49, 178, 8, 34, 1, 79, 5, 6, 49, 2, 37, 23,
		7, 9, 48, 1, 248, 5, 6, 38, 4, 28, 6, 19, 89, 5, 3, 7, 12, 59, 1, 28, 48, 7, 3, 249, 459,
		6, 7, 59, 3, 2, 6, 14, 149, 8, 45, 23, 6, 7, 5, 2348, 248, 28, 1, 9, 59, 59, 4, 17, 28,
		17, 6, 3, 28, 1, 238, 28, 6, 23, 9, 7, 45, 45}, g.encodeInts())
}

func TestHiddenPairRow(t *testing.T) {
	g := decodeInts([]int{269, 24, 5, 3, 248, 478, 18, 469, 178, 8, 34, 1, 9, 5, 467, 46, 2, 37, 2369,
		7, 29, 46, 1, 248, 5, 469, 38, 4, 28, 6, 18, 9, 5, 3, 7, 12, 259, 1, 289, 48, 7, 3, 249,
		45, 6, 7, 59, 3, 2, 46, 146, 149, 8, 145, 23, 6, 7, 5, 2348, 248, 248, 1, 9, 159, 59, 4,
		7, 268, 168, 268, 3, 258, 15, 238, 28, 146, 23, 9, 7, 456, 45})
	assert.True(t, g.hiddenPairGroup(&row, 0))
	assert.Equal(t, []int{69, 24, 5, 3, 248, 478, 18, 69, 178, 8, 34, 1, 9, 5, 467, 46, 2, 37,
		2369, 7, 29, 46, 1, 248, 5, 469, 38, 4, 28, 6, 18, 9, 5, 3, 7, 12, 259, 1, 289, 48, 7,
		3, 249, 45, 6, 7, 59, 3, 2, 46, 146, 149, 8, 145, 23, 6, 7, 5, 2348, 248, 248, 1, 9, 159,
		59, 4, 7, 268, 168, 268, 3, 258, 15, 238, 28, 146, 23, 9, 7, 456, 45}, g.encodeInts())
}
