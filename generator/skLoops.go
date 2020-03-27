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

func (g *Grid) skLoops() (res bool) {
	solved := make(map[point]bool)
	for r := zero; r < rows; r++ {
		for c := zero; c < cols; c++ {
			if bitCount[g.cells[r][c]] == 1 {
				solved[point{r, c}] = true
			}
		}
	}

	rectangles := make(map[[4]point]bool)
	for p1 := range solved {
		for p2 := range solved {
			if p2.r <= p1.r || p2.c <= p1.c {
				continue
			}

			if solved[point{p1.r, p2.c}] && solved[point{p2.r, p1.c}] {
				rectangles[[4]point{p1, point{p1.r, p2.c}, point{p2.r, p1.c}, p2}] = true
			}
		}
	}

	for r := range rectangles {
		p0 := r[0]
		p1 := r[1]
		p2 := r[3] // These are out of order because we want the points to go clockwise around the loop.
		p3 := r[2]

		var immune [rows][cols]bool

		p0l, p0r, p0t, p0b := lobes(p0, &immune)
		p1l, p1r, p1t, p1b := lobes(p1, &immune)
		p2l, p2r, p2t, p2b := lobes(p2, &immune)
		p3l, p3r, p3t, p3b := lobes(p3, &immune)

		if !g.checkSides(r, p0l, p0r, p0t, p0b) || !g.checkSides(r, p1l, p1r, p1t, p1b) || !g.checkSides(r, p2l, p2r, p2t, p2b) || !g.checkSides(r, p3l, p3r, p3t, p3b) {
			continue
		}

		topF, topC, topS := g.findSKLinks(p0l, p0r, p1l, p1r)
		rightF, rightC, rightS := g.findSKLinks(p1t, p1b, p2t, p2b)
		bottomF, bottomC, bottomS := g.findSKLinks(p2l, p2r, p3l, p3r)
		leftF, leftC, leftS := g.findSKLinks(p3t, p3b, p0t, p0b)

		topFCount := bitCount[topF]
		topCCount := bitCount[topC]
		rightFCount := bitCount[rightF]
		rightCCount := bitCount[rightC]
		bottomFCount := bitCount[bottomF]
		bottomCCount := bitCount[bottomC]
		leftFCount := bitCount[leftF]
		leftCCount := bitCount[leftC]

		if topFCount == 0 || topCCount == 0 || rightFCount == 0 || rightCCount == 0 || bottomFCount == 0 || bottomCCount == 0 || leftFCount == 0 || leftCCount == 0 {
			continue
		}

		count := topFCount + topCCount + rightFCount + rightCCount + bottomFCount + bottomCCount + leftFCount + leftCCount

		if count > 16 || topS != rightF || rightS != bottomF || bottomS != leftF || leftS != topF {
			continue
		}

		g.removeSfromSKLoops(&row, p0.r, topC, &immune, &res)
		g.removeSfromSKLoops(&col, p1.c, rightC, &immune, &res)
		g.removeSfromSKLoops(&row, p2.r, bottomC, &immune, &res)
		g.removeSfromSKLoops(&col, p3.c, leftC, &immune, &res)

		g.removeSfromSKLoops(&box, boxOf(p0.r, p0.c), leftS&topF, &immune, &res)
		g.removeSfromSKLoops(&box, boxOf(p1.r, p1.c), topS&rightF, &immune, &res)
		g.removeSfromSKLoops(&box, boxOf(p2.r, p2.c), rightS&bottomF, &immune, &res)
		g.removeSfromSKLoops(&box, boxOf(p3.r, p3.c), bottomS&leftF, &immune, &res)
	}

	return
}

func (g *Grid) removeSfromSKLoops(gr *group, sel uint8, mask cell, immune *[rows][cols]bool, res *bool) {
	for _, p := range gr.unit[sel] {
		if immune[p.r][p.c] {
			continue
		}

		prev := *g.pt(p)
		if g.pt(p).andNot(mask) {
			g.cellChange(res, "skloops: remove %s from %s\n", prev&mask, p)
		}
	}
}

func (g *Grid) checkSides(rect [4]point, l, r, t, b point) bool {
	cl := *g.pt(l)
	cr := *g.pt(r)

	if bitCount[cl|cr] > 4 {
		return false
	}

	ct := *g.pt(t)
	cb := *g.pt(b)

	if bitCount[ct|cb] > 4 {
		return false
	}

	return true
}

func (g *Grid) colorDigit(p point, c color, mask cell, d int, colors *[rows][cols][10]color) {
	cell := *g.pt(p)
	if cell&mask&(1<<d) != 0 {
		colors[p.r][p.c][d] = c
	}
}

func (g *Grid) findSKLinks(p0a, p0b, p1a, p1b point) (cell, cell, cell) {
	c0a := *g.pt(p0a)
	c0b := *g.pt(p0b)
	c1a := *g.pt(p1a)
	c1b := *g.pt(p1b)

	common := c0a & c0b & c1a & c1b

	if bitCount[c0a] == 1 {
		c0a = 0
		common = c0b & c1a & c1b
	}
	if bitCount[c0b] == 1 {
		c0b = 0
		common = c0a & c1a & c1b
	}
	if bitCount[c1a] == 1 {
		c1a = 0
		common = c0a & c0b & c1b
	}
	if bitCount[c1b] == 1 {
		c1b = 0
		common = c0a & c0b & c1a
	}

	c0 := c0a | c0b
	c1 := c1a | c1b

	first := c0 &^ common
	second := c1 &^ common

	return first, common, second
}

func lobes(p point, immune *[rows][cols]bool) (point, point, point, point) {
	pl := point{p.r, (p.c+2)%3 + p.c/3*3}
	pr := point{p.r, (p.c+1)%3 + p.c/3*3}
	pt := point{(p.r+2)%3 + p.r/3*3, p.c}
	pb := point{(p.r+1)%3 + p.r/3*3, p.c}

	immune[p.r][p.c] = true
	immune[pl.r][pl.c] = true
	immune[pr.r][pr.c] = true
	immune[pt.r][pt.c] = true
	immune[pb.r][pb.c] = true

	return pl, pr, pt, pb
}
