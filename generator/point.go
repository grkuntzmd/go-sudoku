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

import "fmt"

type point struct {
	r, c uint8
}

func (p point) String() string {
	return fmt.Sprintf("(%d, %d)", p.r, p.c)
}

func comparePointSlices(p1, p2 []point) bool {
	if len(p1) != len(p2) {
		return false
	}

	for i, v := range p1 {
		if v != p2[i] {
			return false
		}
	}

	return true
}
