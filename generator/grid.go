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
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type (
	Grid struct {
		orig  [rows][cols]bool
		cells [rows][cols]cell
	}

	pointCell struct {
		*point
		cell
	}
)

const (
	cols = 9
	rows = 9

	all = 0b1111111110

	botLeft  = "\u2514"
	botRight = "\u2518"
	botT     = "\u2534"
	horizBar = "\u2500"
	leftT    = "\u251c"
	plus     = "\u253c"
	rightT   = "\u2524"
	topLeft  = "\u250c"
	topRight = "\u2510"
	topT     = "\u252c"
	vertBar  = "\u2502"

	green  = "32"
	yellow = "33"
)

var (
	attempts uint
	color    bool
	verbose  uint
)

func init() {
	rand.Seed(time.Now().Unix())

	flag.UintVar(&attempts, "a", 10, "maximum `attempts` to generate a puzzle")
	flag.BoolVar(&color, "c", false, "colorize the output for ANSI terminals")
	flag.UintVar(&verbose, "v", 0, "`verbosity` level; higher emits more messages")
}

// ParseEncoded parses an input string contains 81 digits and dots ('.') representing an initial puzzle layout.
func ParseEncoded(i string) (*Grid, error) {
	if len(i) != 81 {
		return nil, fmt.Errorf("encoded puzzle must contain 81 characters")
	}
	g := Grid{}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			b := i[r*9+c]
			if b == '.' {
				g.cells[r][c] = all
			} else {
				d, err := strconv.Atoi(string(b))
				if err != nil {
					return nil, fmt.Errorf("Illegal character '%c' in encoded puzzle", b)
				}
				g.orig[r][c] = true
				g.cells[r][c] = 1 << d
			}
		}
	}

	return &g, nil
}

// Randomize generates a random puzzle. There is no guarantee that the puzzle will be solvable or have just one solution.
func Randomize() *Grid {
	g := Grid{}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g.cells[r][c] = all
		}
	}

	d := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Shuffle(len(d), func(i, j int) { d[i], d[j] = d[j], d[i] })
	var group [9][9]point
	switch rand.Intn(3) {
	case 0:
		group = box.unit
	case 1:
		group = col.unit
	case 2:
		group = row.unit
	}

	for i, p := range group[rand.Intn(9)] {
		g.cells[p.r][p.c] = 1 << d[i]
	}

	return &g
}

func (g *Grid) allPoints() (res []pointCell) {
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			res = append(res, pointCell{&point{r, c}, g.cells[r][c]})
		}
	}

	return
}

// cellChange is a convenience function that is called by strategy methods when a cell changes value.
func (g *Grid) cellChange(res *bool, format string, a ...interface{}) {
	*res = true
	if verbose >= 1 {
		log.Printf(format, a...)
	}
	if verbose >= 2 {
		g.Display()
	}
}

// Display emits a grid to stdout in a framed format.
func (g *Grid) Display() {
	width := g.maxWidth() + 2 // Add 2 for margins.
	bars := strings.Repeat(horizBar, width*3)
	line := leftT + strings.Join([]string{bars, bars, bars}, plus) + rightT

	// Top line with column headers.
	fmt.Print("\t   ")
	for d := 0; d < 9; d++ {
		fmt.Printf("%s", colorize(yellow, center(strconv.Itoa(d), width)))
		if d == 2 || d == 5 {
			fmt.Print(" ")
		}
	}
	fmt.Println()

	// First frame line.
	fmt.Printf("\t  %s%s%s%s%s%s%s\n", topLeft, bars, topT, bars, topT, bars, topRight)

	// Grid rows.
	for r := 0; r < rows; r++ {
		fmt.Printf("\t%s %s", colorize(yellow, strconv.Itoa(r)), vertBar)
		for c := 0; c < cols; c++ {
			cell := g.cells[r][c]
			orig := g.orig[r][c]
			s := cell.String()
			if s == "123456789" {
				fmt.Printf("%s", center(".", width))
			} else {
				if orig {
					fmt.Printf("%s", colorize(green, center(s, width)))
				} else {
					fmt.Printf("%s", center(s, width))
				}
			}
			if c == 2 || c == 5 {
				fmt.Printf("%s", vertBar)
			}
		}
		fmt.Printf("%s\n", vertBar)
		if r == 2 || r == 5 {
			fmt.Printf("\t  %s\n", line)
		}
	}

	// Bottom line.
	fmt.Printf("\t  %s%s%s%s%s%s%s\n", botLeft, bars, botT, bars, botT, bars, botRight)
}

// digitPlaces returns an array of digits containing values where the bits (1 - 9) are set if the corresponding digit appears in that cell.
func (g *Grid) digitPlaces(points [9]point) (res [10]cell) {
	for pi, p := range points {
		cell := *g.pt(&p)
		for d := 1; d <= 9; d++ {
			if cell&(1<<d) != 0 {
				res[d] |= 1 << pi
			}
		}
	}
	return
}

// digitPoints builds a table of points that contain each digit.
func (g *Grid) digitPoints(ps [9]point) (res [10][]point) {
	for _, p := range ps {
		cell := *g.pt(&p)
		for d := 1; d <= 9; d++ {
			if cell&(1<<d) != 0 {
				res[d] = append(res[d], p)
			}
		}
	}

	return
}

// emptyCell returns true if the grid contains at least one empty cell (no digits set).
func (g *Grid) emptyCell() bool {
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if g.cells[r][c] == 0 {
				return true
			}
		}
	}
	return false
}

// maxWidth calculates the width in characters of the widest cell in the grid (maximum number of candidate digits). If the width is 9, it is changed to 1 because we will display only a dot ('.').
func (g *Grid) maxWidth() int {
	width := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			count := bitCount[g.cells[r][c]]
			if count == 9 {
				count = 1
			}
			if width < count {
				width = count
			}
		}
	}

	return width
}

// minPoint find the non-solved point with the least number of candidates and returns that point and true if found, otherwise it returns false.
func (g *Grid) minPoint() (p point, found bool) {
	min := 10
	minPoints := make([]*point, 0)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			cell := g.cells[r][c]
			count := bitCount[cell]
			if count > 1 {
				p = point{r, c}
				if count < min {
					min = count
					minPoints = minPoints[:0]
					minPoints = append(minPoints, &p)
					found = true
				} else if count == min {
					minPoints = append(minPoints, &p)
					found = true
				}
			}
		}
	}

	if found {
		rand.Shuffle(len(minPoints), func(i, j int) { minPoints[i], minPoints[j] = minPoints[j], minPoints[i] })
		return *minPoints[0], true
	}

	return
}

// pt returns the cell at a given point.
func (g *Grid) pt(p *point) *cell {
	return &g.cells[p.r][p.c]
}

// Reduce eliminates candidates from cells using logical methods. For example if a cell contains a single digit candidate, that digit can be removed from all other cells in the same box, row, and column.
func (g *Grid) Reduce() (Level, bool) {
	maxLevel := Trivial

	if g.emptyCell() {
		return Trivial, false
	}

	for {
		if g.solved() {
			return maxLevel, true
		}

		if g.reduceLevel(&maxLevel, Trivial, []func() bool{
			g.nakedSingle,
			g.hiddenSingle,
			g.nakedPair,
			g.nakedTriple,
			g.nakedQuad,
			g.hiddenPair,
			g.hiddenTriple,
			g.hiddenQuad,
			g.pointingLine,
			g.boxLine,
		}) {
			continue
		}

		if g.reduceLevel(&maxLevel, Tough, []func() bool{
			g.xWing,
		}) {
			continue
		}

		// if g.reduceLevel(&maxLevel, Diabolical, []func() bool{}) {
		// 	continue
		// }

		// if g.reduceLevel(&maxLevel, Extreme, []func() bool{}) {
		// 	continue
		// }

		// if g.reduceLevel(&maxLevel, Insane, []func() bool{}) {
		// 	continue
		// }

		break
	}

	return maxLevel, false
}

func (g *Grid) reduceLevel(maxLevel *Level, level Level, fs []func() bool) bool {
	for _, f := range fs {
		if f() {
			if *maxLevel < level {
				*maxLevel = level
			}
			return true
		}
	}

	return false
}

// Search uses a brute-force descent to solve the grid and returns a slice of grids that may be empty if no solution was found, may contain a single grid if a unique solution was found, or may contain more than one solution.
func (g *Grid) Search(solutions *[]*Grid) {
	if g.solved() {
		*solutions = append(*solutions, g)
		return
	}

	if g.emptyCell() {
		return
	}

	point, found := g.minPoint()
	if !found {
		return
	}

	digits := g.pt(&point).digits()
	rand.Shuffle(len(digits), func(i, j int) { digits[i], digits[j] = digits[j], digits[i] })

	for _, d := range digits {
		cp := *g
		*cp.pt(&point) = 1 << d
		_, solved := cp.Reduce()

		if solved {
			*solutions = append(*solutions, &cp)
			if len(*solutions) > 1 {
				return
			}
			continue
		}

		cp.Search(solutions)
		if len(*solutions) > 1 {
			return
		}
	}

	return
}

// solved checks that a grid is completely solved (all boxes, rows, and columns have each digit appearing exactly once).
func (g *Grid) solved() bool {
	return g.solvedGroup(&box) && g.solvedGroup(&col) && g.solvedGroup(&row)
}

func (g *Grid) solvedGroup(gr *group) bool {
	for _, ps := range gr.unit {
		cells := [10]int{}
		for _, p := range ps {
			cell := *g.pt(&p)

			if g.orig[p.r][p.c] && bitCount[cell] != 1 {
				log.Panicf("changed original cell (%d, %d) to %#b", p.r, p.c, cell)
			}

			if cell == 0 {
				return false
			}

			for d := 1; d <= 9; d++ {
				if cell&(1<<d) != 0 {
					cells[d]++
				}
			}
		}

		for d := 1; d <= 9; d++ {
			if cells[d] != 1 {
				return false
			}
		}
	}

	return true
}

func Worker(tasks chan Level, results chan *Game) {
outer:
	for level := range tasks {
		maxAttempts := attempts

	inner:
		for {
			grid := Randomize()
			solutions := make([]*Grid, 0, 2)
			grid.Search(&solutions)
			if len(solutions) == 0 {
				maxAttempts--
				if maxAttempts == 0 { // If too many attempts, push a nil and start again with a new level.
					results <- nil
					continue outer
				}
			}

			solution := solutions[0]
			points := solution.allPoints()                                                            // Get all points from the first solution.
			rand.Shuffle(len(points), func(i, j int) { points[i], points[j] = points[j], points[i] }) // Shuffle them.

			// Clear the grid.
			for r := 0; r < rows; r++ {
				for c := 0; c < cols; c++ {
					grid.cells[r][c] = all
				}
			}

			// Fill in the first 17 points from the shuffled set. We use 17 points because it is known that the smallest number of cells required to get a unique solution is 17 (https://en.wikipedia.org/wiki/Mathematics_of_Sudoku).
			for p := 0; p < 17; p++ {
				point := points[p].point
				cell := points[p].cell
				grid.pt(point).setTo(cell)
				grid.orig[point.r][point.c] = true
				points = points[1:]
			}

			// Attempt to solve the grid at the given level. Keep adding points until the level is achieved and there is a unique solution.
			for {
				cp := *grid

				l, solved := cp.Reduce()
				solutions = solutions[:0]
				cp.Search(&solutions)

				if solved && l == level && len(solutions) == 1 {
					solution := solutions[0]
					solution.orig = grid.orig
					results <- &Game{level, 81 - len(points), grid, solution}
					continue outer
				}

				if len(points) == 0 { // If there are no more points to add, abort this attempt.
					maxAttempts--
					if maxAttempts == 0 { // If too many attempts, push a nil and start again with a new level.
						results <- nil
						continue outer
					}

					continue inner // Otherwise just try a new randomized grid at the same level.
				}

				// Add the next point from the point slice to the grid and try solving again.
				point := points[0].point
				cell := points[0].cell
				grid.pt(point).setTo(cell)
				grid.orig[point.r][point.c] = true
				points = points[1:]
				continue
			}

		}
	}
}

// center centers a string in the given width field.
func center(s string, w int) string {
	excess := w - len(s)
	lead := excess / 2
	follow := excess - lead
	return fmt.Sprintf("%*s%*s", lead+len(s), s, follow, " ")
}

// colorize adds ANSI escape sequences to display the string in color.
func colorize(c string, s string) string {
	if color {
		return fmt.Sprintf("\x1b[%sm%s\x1b[0m", c, s)
	}

	return fmt.Sprintf("%s", s)
}
