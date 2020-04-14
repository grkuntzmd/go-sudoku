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

func TestSinglesChains_twiceInAUnit(t *testing.T) {
	g := decodeInts([]int{2, 8, 9, 146, 46, 14, 3, 7, 5, 3, 6, 4, 57, 9, 57, 8, 1, 2, 5, 1, 7, 2, 8,
		3, 9, 6, 4, 8, 9, 3, 457, 2, 457, 6, 45, 1, 1, 4, 5, 8, 3, 6, 7, 2, 9, 7, 2, 6, 19,
		45, 19, 45, 8, 3, 4, 5, 1, 3, 7, 8, 2, 9, 6, 69, 7, 2, 4569, 1, 459, 45, 3, 8, 69, 3, 8, 4569, 456, 2, 1, 45, 7})
	assert.True(t, g.singlesChains(0))
	assert.Equal(t, []int{2, 8, 9, 146, 46, 14, 3, 7, 5, 3, 6, 4, 57, 9, 57, 8, 1, 2, 5, 1, 7,
		2, 8, 3, 9, 6, 4, 8, 9, 3, 457, 2, 457, 6, 45, 1, 1, 4, 5, 8, 3, 6, 7, 2, 9, 7, 2, 6,
		19, 45, 19, 4, 8, 3, 4, 5, 1, 3, 7, 8, 2, 9, 6, 69, 7, 2, 469, 1, 49, 45, 3, 8, 69, 3,
		8, 4569, 6, 2, 1, 4, 7}, g.encodeInts())
}

func TestSinglesChains_twoColorsElsewhere(t *testing.T) {
	g := decodeInts([]int{1, 2, 8, 4, 5, 37, 37, 9, 6, 37, 4, 6, 37, 9, 1, 2, 8, 5, 9, 37, 5, 8, 2,
		6, 4, 1, 37, 678, 67, 3, 5, 678, 2, 1, 4, 9, 678, 9, 1, 367, 4, 37, 68, 5, 2, 4, 5, 2,
		1, 68, 9, 68, 37, 37, 36, 36, 4, 27, 1, 5, 9, 27, 8, 2, 8, 7, 9, 3, 4, 5, 6, 1, 5, 1,
		9, 267, 67, 8, 37, 237, 4})
	assert.True(t, g.singlesChains(0))
	assert.Equal(t, []int{1, 2, 8, 4, 5, 37, 37, 9, 6, 37, 4, 6, 37, 9, 1, 2, 8, 5, 9, 37, 5,
		8, 2, 6, 4, 1, 37, 678, 67, 3, 5, 68, 2, 1, 4, 9, 68, 9, 1, 367, 4, 37, 68, 5, 2, 4,
		5, 2, 1, 68, 9, 68, 37, 37, 36, 36, 4, 27, 1, 5, 9, 27, 8, 2, 8, 7, 9, 3, 4, 5, 6, 1,
		5, 1, 9, 26, 67, 8, 37, 237, 4}, g.encodeInts())
}
