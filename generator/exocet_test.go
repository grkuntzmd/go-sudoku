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

func TestExocet(t *testing.T) {
	g := decodeInts([]int{158, 158, 7, 1568, 2, 1569, 1358, 1389, 4, 9, 3, 12458, 1458, 478, 1457, 6,
		1278, 257, 6, 1458, 12458, 3, 4789, 14579, 12578, 12789, 2579, 13478, 14789, 1489, 248,
		34678, 247, 127, 5, 69, 2, 4579, 459, 456, 1, 34567, 37, 69, 8, 13578, 1578, 6, 9, 378,
		257, 4, 127, 237, 1458, 14568, 3, 7, 46, 12, 9, 248, 256, 478, 2, 489, 46, 5, 39, 378,
		34678, 1, 1457, 145679, 1459, 12, 39, 8, 257, 247, 23567})
	assert.True(t, g.exocet(0))
	assert.Equal(t, []int{158, 158, 7, 1568, 2, 1569, 1358, 1389, 4, 9, 3, 12458, 158, 478, 1457,
		6, 1278, 257, 6, 1458, 12458, 3, 4789, 14579, 158, 12789, 2579, 13478, 14789, 1489, 248,
		34678, 247, 127, 5, 69, 2, 4579, 459, 456, 1, 34567, 37, 69, 8, 13578, 1578, 6, 9, 378,
		257, 4, 127, 237, 1458, 14568, 3, 7, 46, 12, 9, 248, 256, 478, 2, 489, 46, 5, 39, 378,
		34678, 1, 1457, 145679, 1459, 12, 39, 8, 257, 247, 23567}, g.encodeInts())
}
