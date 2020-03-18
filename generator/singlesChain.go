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

type (
	color uint8

	pair struct {
		left, right point
	}
)

const (
	black color = iota
	red
	blue
)

// singlesChain removes candidates by two methods. Prior to removing any candidates, chains are created between cells that contain the only two occurances of a digit in a unit (box, row, or column). The chains connect the units together through the doubly occurring digits. Starting at an arbitrary location in the chain, the cells are alternately colored with two different colors. "Twice in a unit": if the same color occurs twice in a single unit, all cells marked with that color anywhere in the puzzle can be removed. "Two colors elsewhere": if a non-chain cell containing the digit can "see" two cells colored with opposite colors, the digit can be removing from the non-chain cell.
func (g *Grid) singlesChain() (res bool) {
	// Create a pairs set containing cells where the cells contain the only two occurrances of a digit in the unit. We use a set so that the pairs are unique.
	var pairMaps [10]map[pair]bool
	g.unitPairs(&pairMaps)

	// Color the points in the chains.
	for d := 1; d <= 9; d++ {
		pairMap := pairMaps[d]

		for len(pairMap) != 0 {
			colors := make(map[point]color)

			setBoth := true
			for {
				changed := false
				for p := range pairMap {
					if setColors(p, &colors, &pairMap, setBoth) {
						changed = true
						setBoth = false
						delete(pairMap, p)
					}
				}
				if !changed {
					break
				}
			}

			// Separate the colors into two slices.
			var blues, reds []point
			for p, c := range colors {
				switch c {
				case blue:
					blues = append(blues, p)
				case red:
					reds = append(reds, p)
				}
			}

			// Search for "Twice in a unit".
			if g.twiceInAUnit(blues) {
				for _, p := range blues {
					if g.pt(&p).andNot(1 << d) {
						g.cellChange(&res, "singlesChain: in %s, removing %d for twice in a unit", &p, d)
					}
				}
			} else if g.twiceInAUnit(reds) {
				for _, p := range reds {
					if g.pt(&p).andNot(1 << d) {
						g.cellChange(&res, "singlesChain: in %s, removing %d for twice in a unit", &p, d)
					}
				}
			}

			// Search for "Two colors elsewhere".
			for r := 0; r < rows; r++ {
				for c := 0; c < cols; c++ {
					p := point{r, c}

					if *g.pt(&p)&(1<<d) == 0 {
						continue
					}

					if _, ok := colors[p]; ok { // Skip if part of chain.
						continue
					}

					b := boxOf(r, c)
					var seesBlue *point
					for _, blue := range blues {
						if b == boxOf(blue.r, blue.c) || c == blue.c || r == blue.r {
							pt := blue
							seesBlue = &pt
						}
					}
					var seesRed *point
					for _, red := range reds {
						if b == boxOf(red.r, red.c) || c == red.c || r == red.r {
							pt := red
							seesRed = &pt
						}
					}

					if seesBlue != nil && seesRed != nil {
						if g.pt(&p).andNot(1 << d) {
							g.cellChange(&res, "singlesChain: in %s, removing %d for two colors elsewhere (%s, %s)", &p, d, seesBlue, seesRed)
						}
					}
				}
			}
		}
	}

	return
}

func (g *Grid) twiceInAUnit(colors []point) bool {
	for _, p1 := range colors {
		for _, p2 := range colors {
			if p1 == p2 {
				continue
			}

			if boxOf(p1.r, p1.c) == boxOf(p2.r, p2.c) || p1.c == p2.c || p1.r == p2.r {
				return true
			}
		}
	}

	return false
}

func (g *Grid) unitPairs(pairMaps *[10]map[pair]bool) {
	g.unitPairsGroup(&box, pairMaps)
	g.unitPairsGroup(&col, pairMaps)
	g.unitPairsGroup(&row, pairMaps)
}

func (g *Grid) unitPairsGroup(gr *group, pairMaps *[10]map[pair]bool) {
	for _, ps := range gr.unit {
		digits := g.digitPoints(ps)

		for d := uint16(1); d <= 9; d++ {
			points := digits[d]
			if len(points) != 2 {
				continue
			}

			if g.orig[points[0].r][points[0].c] || g.orig[points[1].r][points[1].c] {
				continue
			}

			if (*pairMaps)[d] == nil {
				(*pairMaps)[d] = make(map[pair]bool)
			}
			(*pairMaps)[d][pair{points[0], points[1]}] = true
		}
	}
}

func setColors(p pair, colors *map[point]color, pairMap *map[pair]bool, colorBoth bool) bool {
	colorLeft := (*colors)[p.left]
	colorRight := (*colors)[p.right]

	if colorLeft == black && colorRight == black {
		if colorBoth {
			(*colors)[p.left] = red
			(*colors)[p.right] = blue
			return true
		}
		return false
	}

	if colorLeft == red && colorRight == black {
		(*colors)[p.right] = blue
		return true
	}

	if colorLeft == blue && colorRight == black {
		(*colors)[p.right] = red
		return true
	}

	if colorRight == red && colorLeft == black {
		(*colors)[p.left] = blue
		return true
	}

	if colorRight == blue && colorLeft == black {
		(*colors)[p.left] = red
		return true
	}

	delete(*pairMap, p)
	return false
}
