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

type group struct {
	name string
	unit [rows][cols]point
}

var (
	box = group{name: "box"} // These are all of the coordinates in a box (first dimension).
	col = group{name: "col"} // These are all of the coordinates in a column (first dimension).
	row = group{name: "row"} // These are all of the coordinates in a row (first dimension).
)

func init() {
	for r := zero; r < rows; r++ {
		for c := zero; c < cols; c++ {
			p := point{r, c}
			box.unit[boxOf(r, c)][r%3*3+c%3] = p
			col.unit[c][r] = p
			row.unit[r][c] = p
		}
	}
}

func boxOf(r, c uint8) uint8 {
	return r/3*3 + c/3
}
