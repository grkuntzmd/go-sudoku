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

func TestWXYZWing(t *testing.T) {
	g := decodeInts([]int{1689, 169, 189, 1589, 2, 4, 7, 3, 158, 5, 4, 189, 3, 7, 89, 2, 6, 18, 2, 3,
		7, 1568, 15, 568, 159, 189, 4, 7, 12569, 1259, 59, 3, 259, 8, 4, 156, 69, 2569, 3, 4, 8, 1,
		59, 279, 567, 19, 8, 4, 579, 6, 2579, 159, 12, 3, 3, 12, 128, 1678, 14, 678, 46, 5, 9, 148,
		7, 158, 568, 9, 3, 46, 18, 2, 1489, 159, 6, 2, 145, 58, 3, 178, 178})
	assert.True(t, g.wxyzWing(0))
	assert.Equal(t, []int{1689, 169, 189, 1589, 2, 4, 7, 3, 158, 5, 4, 189, 3, 7, 89, 2, 6, 18, 2,
		3, 7, 1568, 15, 568, 159, 189, 4, 7, 1256, 1259, 59, 3, 259, 8, 4, 156, 69, 2569, 3, 4, 8,
		1, 59, 279, 567, 19, 8, 4, 579, 6, 2579, 159, 12, 3, 3, 12, 128, 1678, 14, 678, 46, 5, 9,
		148, 7, 158, 568, 9, 3, 46, 18, 2, 1489, 159, 6, 2, 145, 58, 3, 178, 178}, g.encodeInts())
}
