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

func (g *Grid) medusa() (res bool) {
	var pairMaps [10]map[pair]bool
	g.unitPairs(&pairMaps)

	strongLinks := make(map[pair]cell)
	for d := 1; d <= 9; d++ {
		for p := range pairMaps[d] {
			strongLinks[p] |= 1 << d
		}
	}

	linkEnds := make(map[point][]pair)
	for p := range strongLinks {
		linkEnds[p.left] = append(linkEnds[p.left], p)
		linkEnds[p.right] = append(linkEnds[p.right], p)
	}

	used := make(map[point]bool)

	for p, c := range strongLinks {
		if used[p.left] || used[p.right] {
			continue
		}

		used[p.left] = true

		digit := c.lowestSetBit()
		var colors [rows][cols][10]color
		colors[p.left.r][p.left.c][digit] = blue

		g.colorGrid(digit, p.left, &colors, linkEnds, strongLinks, &used)

		// Twice in a cell. If the same color appears twice in a cell, that color can be removed from the whole puzzle.
		blueMoreThanOnce := false
		redMoreThanOnce := false
		for r := zero; r < rows; r++ {
			for c := zero; c < cols; c++ {
				blues := 0
				reds := 0
				for _, c := range colors[r][c] {
					switch c {
					case blue:
						blues++
					case red:
						reds++
					}
				}
				blueMoreThanOnce = blueMoreThanOnce || blues > 1
				redMoreThanOnce = redMoreThanOnce || reds > 1
			}
		}
		if blueMoreThanOnce {
			g.removeColor(blue, &colors, "twice in a cell", &res)
		} else if redMoreThanOnce {
			g.removeColor(red, &colors, "twice in a cell", &res)
		}

		if res {
			return
		}

		// Twice in a unit. If the same color appears twice in a unit (box, column, or row) for the same digit, that color can be removed from the whole puzzle.
		blueMoreThanOnce = false
		redMoreThanOnce = false

		b, r := g.groupColors(&box, &colors)
		blueMoreThanOnce = blueMoreThanOnce || b
		redMoreThanOnce = redMoreThanOnce || r

		b, r = g.groupColors(&col, &colors)
		blueMoreThanOnce = blueMoreThanOnce || b
		redMoreThanOnce = redMoreThanOnce || r

		b, r = g.groupColors(&row, &colors)
		blueMoreThanOnce = blueMoreThanOnce || b
		redMoreThanOnce = redMoreThanOnce || r

		if blueMoreThanOnce {
			g.removeColor(blue, &colors, "twice in a unit", &res)
		} else if redMoreThanOnce {
			g.removeColor(red, &colors, "twice in a unit", &res)
		}

		if res {
			return
		}

		// Two colors in a cell. If a cell contains digits that are colored both blue and red, any non-colored digits can be removed.
		for r := zero; r < rows; r++ {
			for c := zero; c < cols; c++ {
				blueFound := 0
				redFound := 0
				for d := 1; d <= 9; d++ {
					switch colors[r][c][d] {
					case blue:
						blueFound |= 1 << d
					case red:
						redFound |= 1 << d
					}
				}

				if blueFound != 0 && redFound != 0 {
					for d := 1; d <= 9; d++ {
						if blueFound&(1<<d) != 0 || redFound&(1<<d) != 0 {
							continue
						}

						if g.pt(point{r, c}).andNot(1 << d) {
							g.cellChange(&res, "3dMedusa (two colors in a cell): in %s, remove %d\n", point{r, c}, d)
						}
					}
				}
			}
		}

		if res {
			return
		}

		// Two colors elsewhere. In all cells C containing a digit X, if that cell can see a blue X and a red X, then X can be removed from the cell C.
		// Mark the cells that see each blue and red digit.
		var blueInfluence, redInfluence, immune [rows][cols][10]bool
		for r := zero; r < rows; r++ {
			for c := zero; c < cols; c++ {
				for d := 1; d <= 9; d++ {
					switch colors[r][c][d] {
					case blue:
						immune[r][c][d] = true // Cells that are part of the 3d medusa are not eligible for removal.
						coloredNeighbors(d, point{r, c}, &blueInfluence)
					case red:
						immune[r][c][d] = true // Cells that are part of the 3d medusa are not eligible for removal.
						coloredNeighbors(d, point{r, c}, &redInfluence)
					}
				}
			}
		}
		for r := zero; r < rows; r++ {
			for c := zero; c < cols; c++ {
				for d := 1; d <= 9; d++ {
					if !immune[r][c][d] {
						if blueInfluence[r][c][d] && redInfluence[r][c][d] {
							if g.pt(point{r, c}).andNot(1 << d) {
								g.cellChange(&res, "3dMedusa (two colors elsewhere): in %s, remove %d\n", point{r, c}, d)
							}
						}
					}
				}
			}
		}

		if res {
			return
		}

		// Two colors unit and cell. If a cell C containing a digit X can see another cell containing a colored X and in C there is a candidate with the opposite color, X can be removed from C.
		for r := zero; r < rows; r++ {
			for c := zero; c < cols; c++ {
				blueFound := 0
				redFound := 0
				var immune [10]bool
				for d := 1; d <= 9; d++ {
					switch colors[r][c][d] {
					case blue:
						blueFound |= 1 << d
						immune[d] = true
					case red:
						redFound |= 1 << d
						immune[d] = true
					}
				}
				for d := 1; d <= 9; d++ {
					if !immune[d] && blueFound != 0 && canSeeColor(d, point{r, c}, red, &colors) {
						if g.pt(point{r, c}).andNot(1 << d) {
							g.cellChange(&res, "3dMedusa (two colors unit and cell): in %s, remove %d\n", point{r, c}, d)
						}
					} else if !immune[d] && redFound != 0 && canSeeColor(d, point{r, c}, blue, &colors) {
						if g.pt(point{r, c}).andNot(1 << d) {
							g.cellChange(&res, "3dMedusa (two colors unit and cell): in %s, remove %d\n", point{r, c}, d)
						}
					}

				}
			}
		}

		if res {
			return
		}

		// Cell emptied by color.
	outer:
		for r := zero; r < rows; r++ {
		inner:
			for c := zero; c < cols; c++ {
				if bitCount[g.cells[r][c]] == 1 {
					continue
				}

				seeBlue := true
				seeRed := true
				for d := 1; d <= 9; d++ {
					if colors[r][c][d] != black {
						continue inner
					}

					if !blueInfluence[r][c][d] {
						seeBlue = false
					}

					if !redInfluence[r][c][d] {
						seeRed = false
					}
				}
				if seeBlue {
					g.removeColor(blue, &colors, "cell emptied by color", &res)
					break outer
				}
				if seeRed {
					g.removeColor(red, &colors, "cell emptied by color", &res)
					break outer
				}
			}
		}

		if res {
			return
		}
	}

	return
}

func (g *Grid) colorGrid(digit int, p point, colors *[rows][cols][10]color, linkEnds map[point][]pair, strongLinks map[pair]cell, used *map[point]bool) {
	(*used)[p] = true

	currColor := colors[p.r][p.c][digit]

	// If the point is a bivalue (only two candidate digits in the cell), then color the other one the opposite color.
	cell := *g.pt(p)
	if bitCount[cell] == 2 {
		o := (cell &^ (1 << digit)).lowestSetBit()
		if colors[p.r][p.c][o] == black {
			colors[p.r][p.c][o] = flipColor(currColor)
			g.colorGrid(o, p, colors, linkEnds, strongLinks, used)
		}
	}

	for _, l := range linkEnds[p] {
		if strongLinks[l]&(1<<digit) == 0 {
			continue
		}

		if p == l.left { // If we are at the left end of the link, process the right end.
			if colors[l.right.r][l.right.c][digit] == black {
				colors[l.right.r][l.right.c][digit] = flipColor(currColor)
				g.colorGrid(digit, l.right, colors, linkEnds, strongLinks, used)
			}
		} else { // Process the left end.
			if colors[l.left.r][l.left.c][digit] == black {
				colors[l.left.r][l.left.c][digit] = flipColor(currColor)
				g.colorGrid(digit, l.left, colors, linkEnds, strongLinks, used)
			}
		}
	}
}

func (g *Grid) groupColors(gr *group, colors *[rows][cols][10]color) (bool, bool) {
	blueMoreThanOnce := false
	redMoreThanOnce := false
	for _, ps := range gr.unit {
		var blues, reds [10]int
		for _, p := range ps {
			for d := 1; d <= 9; d++ {
				if *g.pt(p)&(1<<d) != 0 {
					switch colors[p.r][p.c][d] {
					case blue:
						blues[d]++
					case red:
						reds[d]++
					}
				}
			}
		}
		for d := 1; d <= 9; d++ {
			if blues[d] > 2 {
				blueMoreThanOnce = true
			}
			if reds[d] > 2 {
				redMoreThanOnce = true
			}
		}
	}

	return blueMoreThanOnce, redMoreThanOnce
}

func (g *Grid) removeColor(cl color, colors *[rows][cols][10]color, message string, res *bool) {
	for r := zero; r < rows; r++ {
		for c := zero; c < cols; c++ {
			for ci, color := range colors[r][c] {
				if color == cl {
					if g.pt(point{r, c}).andNot(1 << ci) {
						g.cellChange(res, "3dMedusa (%s): in %s, remove %d\n", message, point{r, c}, ci)
					}
				}
			}
		}
	}
}

func canSeeColor(d int, curr point, c color, colors *[rows][cols][10]color) bool {
	for _, u := range []*[9]point{&box.unit[boxOf(curr.r, curr.c)], &col.unit[curr.c], &row.unit[curr.r]} {
		for _, p := range u {
			if p == curr {
				continue
			}

			if colors[p.r][p.c][d] == c {
				return true
			}
		}
	}

	return false
}
