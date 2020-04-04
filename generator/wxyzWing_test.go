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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWXYZWing(t *testing.T) {
	g := decode([]int{1689, 169, 189, 1589, 2, 4, 7, 3, 158, 5, 4, 189, 3, 7, 89, 2, 6, 18, 2, 3,
		7, 168, 15, 568, 159, 189, 4, 7, 12569, 1259, 59, 3, 259, 8, 4, 156, 69, 2569, 3, 4, 8, 1,
		59, 279, 567, 19, 8, 4, 579, 6, 2579, 159, 12, 3, 3, 12, 128, 1678, 14, 678, 46, 5, 9, 148,
		7, 158, 568, 9, 3, 46, 18, 2, 1489, 159, 6, 2, 145, 58, 3, 178, 178})
	assert.True(t, g.wxyzWing(0))
	assert.Equal(t, []int{1689, 169, 189, 1589, 2, 4, 7, 3, 158, 5, 4, 189, 3, 7, 89, 2, 6, 18, 2,
		3, 7, 168, 15, 568, 159, 189, 4, 7, 12569, 125, 59, 3, 259, 8, 4, 156, 69, 2569, 3, 4, 8,
		1, 59, 279, 567, 19, 8, 4, 579, 6, 2579, 159, 12, 3, 3, 12, 128, 1678, 14, 678, 46, 5, 9,
		148, 7, 58, 568, 9, 3, 46, 18, 2, 1489, 159, 6, 2, 145, 58, 3, 178, 178}, g.encode())
}
