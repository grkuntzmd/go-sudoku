/*
 * MIT LICENSE
 *
 * Copyright © 2020, G.Ralph Kuntz, MD.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package generator

import "fmt"

const (
	black color = iota
	blue
	red
)

type (
	color uint8

	pair struct {
		left, right point
	}

	link struct {
		pair
		digit int
	}

	unitLink struct {
		link
		unit   string
		strong bool
	}
)

func (g *Grid) findBivalueLinks(gr *group, links *[10]map[link]bool) {
	for _, ps := range gr.unit {
		points := g.digitPoints(ps)

		for d := 1; d <= 9; d++ {
			dp := points[d]

			for p1i, p1 := range dp {
				if bitCount[*g.pt(p1)] != 2 {
					continue
				}

				for p2i, p2 := range dp {
					if p1i >= p2i {
						continue
					}

					if bitCount[*g.pt(p2)] != 2 {
						continue
					}

					s := &(*links)[d]
					if *s == nil {
						*s = make(map[link]bool)
					}

					(*s)[link{pair{p1, p2}, d}] = true
				}
			}
		}
	}
}

func (g *Grid) findStrongLinks(gr *group, strongLinks *[10]map[unitLink]bool) {
	for pi, ps := range gr.unit {
		points := g.digitPoints(ps)

		for d := 1; d <= 9; d++ {
			p := points[d]

			if len(p) != 2 {
				continue
			}

			s := &(*strongLinks)[d]
			if *s == nil {
				*s = make(map[unitLink]bool)
			}

			(*s)[sortLink(unitLink{link{pair{p[0], p[1]}, d}, fmt.Sprintf("%s %d", gr.name, pi), true})] = true
		}
	}
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

func (c color) String() string {
	switch c {
	case blue:
		return "blue"
	case red:
		return "red"
	default:
		return "black"
	}
}

func (l *link) reverse() link {
	return link{l.pair.reverse(), l.digit}
}

func (p pair) reverse() pair {
	return pair{p.right, p.left}
}

func coloredNeighbors(d int, curr point, influence *[rows][cols][10]bool) {
	for _, u := range []*[9]point{&box.unit[boxOfPoint(curr)], &col.unit[curr.c], &row.unit[curr.r]} {
		for _, p := range u {
			if p == curr {
				continue
			}

			(*influence)[p.r][p.c][d] = true
		}
	}
}

func flipColor(c color) color {
	switch c {
	case blue:
		return red
	case red:
		return blue
	default:
		return black
	}
}

func neighbors(curr point) *[9][9]bool {
	var res [9][9]bool
	for _, u := range []*[9]point{&box.unit[boxOfPoint(curr)], &col.unit[curr.c], &row.unit[curr.r]} {
		for _, p := range u {
			if p == curr {
				continue
			}

			res[p.r][p.c] = true
		}
	}

	return &res
}

func pointCol(p point) uint8 {
	return p.c
}

func pointRow(p point) uint8 {
	return p.r
}

func sortLink(p unitLink) unitLink {
	if p.left.r < p.right.r || p.left.r == p.right.r && p.left.c < p.right.c {
		return p
	}

	p.left, p.right = p.right, p.left
	return p
}
