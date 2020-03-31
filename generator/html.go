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
	"html/template"
	"strconv"
	"strings"

	s "github.com/ajstarks/svgo"
	"github.com/pkg/browser"
)

const html = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sudoku</title>
</head>
<body>
    {{.Body}}
</body>
</html>
`

// HTML generates the HTML for a grid. The HTML will contain embedded SVG for the actual grid.
func (g *Grid) HTML(showCandidates bool, colors *[rows][cols][10]color) {
	s := g.SVG(2.0, false, showCandidates, colors)

	var b strings.Builder
	t := template.Must(template.New("html").Parse(html))
	if err := t.Execute(&b, struct{ Body template.HTML }{template.HTML(s)}); err != nil {
		panic(err)
	}

	if err := browser.OpenReader(strings.NewReader(b.String())); err != nil {
		panic(err)
	}
}

// SVG returns the standard vector graphics representation for a grid.
func (g *Grid) SVG(scale float64, invert bool, showCandidates bool, colors *[rows][cols][10]color) string {
	const (
		xoffset    = 25
		yoffset    = 25
		gridWidth  = 450
		gridHeight = 450
	)
	var (
		width  = 500 * scale
		height = 500 * scale
		b      strings.Builder
	)

	canvas := s.New(&b)
	canvas.Start(int(width), int(height))
	canvas.Scale(scale)
	if invert {
		canvas.Gtransform("translate(500, 500) rotate(180)")
	}
	canvas.Rect(0, 0, 500, 500, "fill:white")
	canvas.Grid(xoffset, yoffset, gridWidth, gridHeight, gridWidth/9, "stroke:black")
	canvas.Grid(xoffset, yoffset, gridWidth, gridHeight, gridWidth/3, "stroke:black;stroke-width:5;stroke-linecap:round")

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			cell := *&g.cells[r][c]
			digits := cell.String()
			if len(digits) == 1 {
				var color string
				if g.orig[r][c] {
					color = ";fill:green"
				} else {
					color = ";fill:black"
				}

				canvas.Text(c*50+xoffset+25, r*50+yoffset+35, digits, "font:25px sans-serif;;text-anchor:middle"+color)
			} else if showCandidates {
				for d := 1; d <= 9; d++ {
					if cell&(1<<d) != 0 {
						cr := (d - 1) / 3
						cc := (d - 1) % 3
						var color string
						if colors != nil && colors[r][c][d] != black {
							switch colors[r][c][d] {
							case blue:
								color = ";fill:blue"
							case red:
								color = ";fill:red"
							default:
								color = ";fill:black"
							}
						} else {
							color = ";fill:black"
						}

						canvas.Text(c*50+xoffset+10+cc*15, r*50+yoffset+13+cr*15, strconv.Itoa(d), "font:8px sans-serif;;text-anchor:middle"+color)
					}
				}
			}
		}
	}
	if invert {
		canvas.Gend()
	}
	canvas.Gend()

	canvas.End()

	s := b.String()

	i := strings.Index(s, "\n")
	if i > 0 {
		s = s[i+1:]
	}

	return s
}
