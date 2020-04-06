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

func TestExocet(t *testing.T) {
	g := decode([]int{158, 158, 7, 1568, 2, 1569, 1358, 1389, 4, 9, 3, 12458, 1458, 478, 1457, 6,
		1278, 257, 6, 1458, 12458, 3, 4789, 14579, 12578, 12789, 2579, 13478, 14789, 1489, 248,
		34678, 247, 127, 5, 69, 2, 4579, 459, 456, 1, 34567, 37, 69, 8, 13578, 1578, 6, 9, 378,
		257, 4, 127, 237, 1458, 14568, 3, 7, 46, 12, 9, 248, 256, 478, 2, 489, 46, 5, 39, 378,
		34678, 1, 1457, 145679, 1459, 12, 39, 8, 257, 247, 23567})
	assert.True(t, g.exocet(0))
	assert.Equal(t, []int{158, 158, 7, 1568, 2, 1569, 1358, 1389, 4, 9, 3, 12458, 158, 478, 1457,
		6, 1278, 257, 6, 1458, 12458, 3, 4789, 14579, 158, 12789, 2579, 13478, 14789, 1489, 248,
		34678, 247, 127, 5, 69, 2, 4579, 459, 456, 1, 34567, 37, 69, 8, 13578, 1578, 6, 9, 378,
		257, 4, 127, 237, 1458, 14568, 3, 7, 46, 12, 9, 248, 256, 478, 2, 489, 46, 5, 39, 378,
		34678, 1, 1457, 145679, 1459, 12, 39, 8, 257, 247, 23567}, g.encode())
}
