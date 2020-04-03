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

func TestXYChains(t *testing.T) {
	g := decode([]int{26, 8, 245, 1, 29, 3, 59, 7, 456, 37, 9, 24, 5, 27, 6, 18, 14, 348, 37, 56,
		1, 4, 79, 8, 359, 2, 356, 5, 7, 8, 2, 4, 1, 6, 3, 9, 1, 4, 3, 6, 5, 9, 7, 8, 2, 9, 2, 6, 8,
		3, 7, 4, 5, 1, 68, 3, 7, 9, 16, 5, 2, 14, 48, 28, 56, 25, 3, 16, 4, 18, 9, 7, 4, 1, 9, 7,
		8, 2, 35, 6, 35})
	assert.True(t, g.xyChains(0))
	assert.Equal(t, []int{26, 8, 4, 1, 29, 3, 59, 7, 456, 37, 9, 24, 5, 27, 6, 18, 14, 38, 37, 56,
		1, 4, 79, 8, 39, 2, 36, 5, 7, 8, 2, 4, 1, 6, 3, 9, 1, 4, 3, 6, 5, 9, 7, 8, 2, 9, 2, 6, 8,
		3, 7, 4, 5, 1, 68, 3, 7, 9, 16, 5, 2, 14, 48, 28, 56, 25, 3, 16, 4, 18, 9, 7, 4, 1, 9, 7,
		8, 2, 35, 6, 35}, g.encode())
}
