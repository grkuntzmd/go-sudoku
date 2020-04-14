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

func TestXYChains(t *testing.T) {
	g := decodeInts([]int{26, 8, 245, 1, 29, 3, 59, 7, 456, 37, 9, 24, 5, 27, 6, 18, 14, 348, 37, 56,
		1, 4, 79, 8, 359, 2, 356, 5, 7, 8, 2, 4, 1, 6, 3, 9, 1, 4, 3, 6, 5, 9, 7, 8, 2, 9, 2, 6, 8,
		3, 7, 4, 5, 1, 68, 3, 7, 9, 16, 5, 2, 14, 48, 28, 56, 25, 3, 16, 4, 18, 9, 7, 4, 1, 9, 7,
		8, 2, 35, 6, 35})
	assert.True(t, g.xyChains(0))
	assert.Equal(t, []int{26, 8, 4, 1, 29, 3, 59, 7, 456, 37, 9, 24, 5, 27, 6, 18, 14, 38, 37, 56,
		1, 4, 79, 8, 39, 2, 36, 5, 7, 8, 2, 4, 1, 6, 3, 9, 1, 4, 3, 6, 5, 9, 7, 8, 2, 9, 2, 6, 8,
		3, 7, 4, 5, 1, 68, 3, 7, 9, 16, 5, 2, 14, 48, 28, 56, 25, 3, 16, 4, 18, 9, 7, 4, 1, 9, 7,
		8, 2, 35, 6, 35}, g.encodeInts())
}
