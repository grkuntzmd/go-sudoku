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

func TestSinglesChains_twiceInAUnit(t *testing.T) {
	g := decode([]int{2, 8, 9, 146, 46, 14, 3, 7, 5, 3, 6, 4, 57, 9, 57, 8, 1, 2, 5, 1, 7, 2, 8,
		3, 9, 6, 4, 8, 9, 3, 457, 2, 457, 6, 45, 1, 1, 4, 5, 8, 3, 6, 7, 2, 9, 7, 2, 6, 19,
		45, 19, 45, 8, 3, 4, 5, 1, 3, 7, 8, 2, 9, 6, 69, 7, 2, 4569, 1, 459, 45, 3, 8, 69, 3, 8, 4569, 456, 2, 1, 45, 7})
	assert.True(t, g.singlesChains(0))
	assert.Equal(t, []int{2, 8, 9, 146, 46, 14, 3, 7, 5, 3, 6, 4, 57, 9, 57, 8, 1, 2, 5, 1, 7,
		2, 8, 3, 9, 6, 4, 8, 9, 3, 457, 2, 457, 6, 45, 1, 1, 4, 5, 8, 3, 6, 7, 2, 9, 7, 2, 6,
		19, 45, 19, 4, 8, 3, 4, 5, 1, 3, 7, 8, 2, 9, 6, 69, 7, 2, 469, 1, 49, 45, 3, 8, 69, 3,
		8, 4569, 6, 2, 1, 4, 7}, g.encode())
}

func TestSinglesChains_twoColorsElsewhere(t *testing.T) {
	g := decode([]int{1, 2, 8, 4, 5, 37, 37, 9, 6, 37, 4, 6, 37, 9, 1, 2, 8, 5, 9, 37, 5, 8, 2,
		6, 4, 1, 37, 678, 67, 3, 5, 678, 2, 1, 4, 9, 678, 9, 1, 367, 4, 37, 68, 5, 2, 4, 5, 2,
		1, 68, 9, 68, 37, 37, 36, 36, 4, 27, 1, 5, 9, 27, 8, 2, 8, 7, 9, 3, 4, 5, 6, 1, 5, 1,
		9, 267, 67, 8, 37, 237, 4})
	assert.True(t, g.singlesChains(0))
	assert.Equal(t, []int{1, 2, 8, 4, 5, 37, 37, 9, 6, 37, 4, 6, 37, 9, 1, 2, 8, 5, 9, 37, 5,
		8, 2, 6, 4, 1, 37, 678, 67, 3, 5, 68, 2, 1, 4, 9, 68, 9, 1, 367, 4, 37, 68, 5, 2, 4,
		5, 2, 1, 68, 9, 68, 37, 37, 36, 36, 4, 27, 1, 5, 9, 27, 8, 2, 8, 7, 9, 3, 4, 5, 6, 1,
		5, 1, 9, 26, 67, 8, 37, 237, 4}, g.encode())
}
