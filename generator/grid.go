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
	"math/rand"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type (
	// Grid is the primary data structure for the generator. It contains the candidates for each cell in the 9 x 9 puzzle.
	Grid struct {
		orig  [rows][cols]bool
		cells [rows][cols]cell
	}

	pointCell struct {
		point
		cell
	}
)

const (
	cols = 9
	rows = 9
	zero = uint8(0)

	all = 0b1111111110
)

var (
	attempts  uint
	colorized bool
	encodings bool
)

func init() {
	rand.Seed(time.Now().Unix())

	flag.UintVar(&attempts, "a", 100, "maximum `attempts` to generate a puzzle")
	flag.BoolVar(&colorized, "c", false, "colorize the output for ANSI terminals")
	flag.BoolVar(&encodings, "e", false, "Add encodings to the each grid display (used to write test cases)")
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

	indexes := []int{0, 1, 2}
	rand.Shuffle(len(indexes), func(i, j int) { indexes[i], indexes[j] = indexes[j], indexes[i] })
	for i, index := range indexes {
		u := i*3 + index
		d := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		rand.Shuffle(len(d), func(i, j int) { d[i], d[j] = d[j], d[i] })
		for pi, p := range box.unit[u] {
			*g.pt(p) = 1 << d[pi]
		}
	}

	return &g
}

func (g *Grid) allPoints() (res []pointCell) {
	for r := zero; r < rows; r++ {
		for c := zero; c < cols; c++ {
			res = append(res, pointCell{point{r, c}, g.cells[r][c]})
		}
	}

	return
}

// cellChange is a convenience function that is called by strategy methods when a cell changes value.
func (g *Grid) cellChange(res *bool, verbose uint, format string, a ...interface{}) {
	*res = true
	if verbose >= 1 {
		fmt.Printf(format, a...)
	}
	if verbose >= 2 {
		g.Display()
	}
}

// Display emits a grid to stdout in a framed format.
func (g *Grid) Display() {
	const (
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

	if encodings {
		fmt.Printf("encoded: %#v\n", g.encode())
	}
}

// digitPlaces returns an array of digits containing values where the bits (1 - 9) are set if the corresponding digit appears in that cell.
func (g *Grid) digitPlaces(points [9]point) (res [10]positions) {
	for pi, p := range points {
		cell := *g.pt(p)
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
		cell := *g.pt(p)
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

func (g *Grid) encode() []int {
	var encoded []int
	for r := zero; r < rows; r++ {
		for c := zero; c < cols; c++ {
			v, _ := strconv.Atoi(g.cells[r][c].String())
			encoded = append(encoded, v)
		}
	}
	return encoded
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
	minPoints := make([]point, 0)
	for r := zero; r < rows; r++ {
		for c := zero; c < cols; c++ {
			cell := g.cells[r][c]
			count := bitCount[cell]
			if count > 1 {
				p = point{r, c}
				if count < min {
					min = count
					minPoints = minPoints[:0]
					minPoints = append(minPoints, p)
					found = true
				} else if count == min {
					minPoints = append(minPoints, p)
					found = true
				}
			}
		}
	}

	if found {
		rand.Shuffle(len(minPoints), func(i, j int) { minPoints[i], minPoints[j] = minPoints[j], minPoints[i] })
		return minPoints[0], true
	}

	return
}

// pt returns the cell at a given point.
func (g *Grid) pt(p point) *cell {
	return &g.cells[p.r][p.c]
}

// Reduce eliminates candidates from cells using logical methods. For example if a cell contains a single digit candidate, that digit can be removed from all other cells in the same box, row, and column.
func (g *Grid) Reduce(strategies *map[string]bool, verbose uint) (Level, bool) {
	maxLevel := Easy

	if g.emptyCell() {
		return Easy, false
	}

	for {
		if g.solved() {
			return maxLevel, true
		}

		if g.reduceLevel(&maxLevel, Easy, verbose, strategies, []func(uint) bool{
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

		if g.reduceLevel(&maxLevel, Standard, verbose, strategies, []func(uint) bool{
			g.xWing,
			g.yWing,
			g.singlesChains,
			g.swordfish,
			g.xyzWing,
		}) {
			continue
		}

		if g.reduceLevel(&maxLevel, Hard, verbose, strategies, []func(uint) bool{
			g.xCycles,
			g.xyChains,
			g.medusa,
			g.jellyfish,
			g.wxyzWing,
		}) {
			continue
		}

		if g.reduceLevel(&maxLevel, Expert, verbose, strategies, []func(uint) bool{
			g.skLoops,
			g.exocet,
		}) {
			continue
		}

		if g.reduceLevel(&maxLevel, Extreme, verbose, strategies, []func(uint) bool{}) {
			continue
		}

		break
	}

	return maxLevel, false
}

func (g *Grid) reduceLevel(maxLevel *Level, level Level, verbose uint, strategies *map[string]bool, fs []func(uint) bool) bool {
	for _, f := range fs {
		if f(verbose) {
			if strategies != nil {
				name := nameOfFunc(f)
				(*strategies)[name] = true
			}
			if *maxLevel < level {
				*maxLevel = level
			}
			return true
		}
	}

	return false
}

func nameOfFunc(f func(uint) bool) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	i := strings.LastIndex(name, ".")
	if i > 0 {
		name = name[i+1:]
	}
	i = strings.LastIndex(name, "-fm")
	if i > 0 {
		name = name[:i]
	}

	return name
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

	digits := g.pt(point).digits()
	rand.Shuffle(len(digits), func(i, j int) { digits[i], digits[j] = digits[j], digits[i] })

	for _, d := range digits {
		cp := *g
		*cp.pt(point) = 1 << d
		_, solved := cp.Reduce(nil, 0)

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
	for r := zero; r < rows; r++ {
		for c := zero; c < cols; c++ {
			if bitCount[g.cells[r][c]] != 1 {
				return false
			}
		}
	}
	return g.solvedGroup(&box) && g.solvedGroup(&col) && g.solvedGroup(&row)
}

func (g *Grid) solvedGroup(gr *group) bool {
	for _, ps := range gr.unit {
		cells := [10]int{}
		for _, p := range ps {
			cell := *g.pt(p)

			if g.orig[p.r][p.c] && bitCount[cell] != 1 {
				panic(fmt.Sprintf("changed original cell (%d, %d) to %#b", p.r, p.c, cell))
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

// Valid returns true if the grid contains at most one occurance of each digit in each unit.
func (g *Grid) Valid() bool {
	return g.validGroup(&box) && g.validGroup(&col) && g.validGroup(&row)
}

func (g *Grid) validGroup(gr *group) bool {
	for _, u := range gr.unit {
		var seen [10]bool
		for _, p := range u {
			if !g.orig[p.r][p.c] {
				continue
			}
			digit := g.pt(p).lowestSetBit()
			if seen[digit] {
				return false
			}
			seen[digit] = true
		}
	}

	return true
}

// Worker generates puzzles. It removes a requested puzzle level from the tasks channel and attempts to generate a puzzle at the level. If it succeeds, it pushes the puzzle to the results channel. If it cannot generate a puzzle, it pushes nil.
func Worker(tasks chan Level, results chan *Game) {
outer:
	for level := range tasks {
		maxAttempts := attempts

	inner:
		for {
			grid := Randomize()
			solutions := make([]*Grid, 0, 2)
			grid.Search(&solutions)
			if len(solutions) == 0 { // The grid has no solution.
				maxAttempts--
				if maxAttempts == 0 { // If too many attempts, push a nil and start again with a new level.
					results <- nil
					continue outer
				}

				continue inner
			}

			// From https://stackoverflow.com/a/7280517/96233.

			*grid = *solutions[0]                                                                     // Copy the first solution
			points := grid.allPoints()                                                                // Get all points from the first solution.
			rand.Shuffle(len(points), func(i, j int) { points[i], points[j] = points[j], points[i] }) // Shuffle them.

			for len(points) > 0 {
				curr := points[0]
				points = points[1:]
				*grid.pt(curr.point) = all // Clear the cell.

				solutions = solutions[:0]
				grid.Search(&solutions)

				if len(solutions) > 1 { // No longer unique.
					*grid.pt(curr.point) = curr.cell // Put the value back.
				}
			}

			// At this point, grid contains the smallest solution that is unique. Now we test the level.
			cp := *grid
			strategies := make(map[string]bool)
			l, solved := cp.Reduce(&strategies, 0)
			solutions = solutions[:0]
			cp.Search(&solutions)
			if solved && l == level && len(solutions) == 1 {
				solution := solutions[0]
				var clues uint
				for r := 0; r < rows; r++ {
					for c := 0; c < cols; c++ {
						if bitCount[grid.cells[r][c]] == 1 {
							solution.orig[r][c] = true
							clues++
						}
					}
				}

				var s []string
				for n := range strategies {
					s = append(s, n)
				}
				sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })

				results <- &Game{level, clues, s, grid, solution}
				continue outer
			}

			// If we could not find a unique solution, try again.
			maxAttempts--
			if maxAttempts == 0 { // If too many attempts, push a nil and start again with a new level.
				results <- nil
				continue outer
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
	if colorized {
		return fmt.Sprintf("\x1b[%sm%s\x1b[0m", c, s)
	}

	return fmt.Sprintf("%s", s)
}

func decode(encoded []int) *Grid {
	if len(encoded) != 81 {
		panic(fmt.Sprintf("encoding has bad length: %d (should be 81)", len(encoded)))
	}

	g := Grid{}
	for i, e := range encoded {
		s := strconv.Itoa(e)
		c := cell(0)
		for ci := 0; ci < len(s); ci++ {
			v, err := strconv.Atoi(s[ci : ci+1])
			if err != nil {
				panic(fmt.Sprintf("encoding values must be digits -- found %s", s[ci:ci+1]))
			}
			c |= 1 << v
		}
		g.cells[i/9][i%9] = c
	}

	return &g
}
