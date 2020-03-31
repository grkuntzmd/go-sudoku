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

// Level is a type wrapper for the difficulty levels of puzzles.
type Level int

const (
	// Easy puzzle use only basic strategies for solving.
	Easy Level = iota
	// Standard puzzles use more complex strategies.
	Standard
	// Hard puzzles use very challenging strategies.
	Hard
	// Expert puzzles use ridiculously difficult stretegies.
	Expert
	// Extreme puzzles use nearly impossible strategies.
	Extreme
)

func (l Level) String() string {
	switch l {
	case Easy:
		return "Easy"
	case Standard:
		return "Standard"
	case Hard:
		return "Hard"
	case Expert:
		return "Expert"
	case Extreme:
		return "Extreme"
	}

	return ""
}
