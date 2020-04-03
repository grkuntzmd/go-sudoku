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
	"fmt"
)

type loopKind = int

const (
	niceLoop loopKind = iota
	strongLoop
	weakLoop
)

func (g *Grid) xCycles(verbose uint) (res bool) {
	// Find all strong links. A pair of points form a strong link if they contain the only two instances of a digit within a unit (box, column, or row).
	var strongLinks [10]map[unitLink]bool
	g.findStrongLinks(&box, &strongLinks)
	g.findStrongLinks(&col, &strongLinks)
	g.findStrongLinks(&row, &strongLinks)

	// Find all weak links. A pair of points form a weak link if they contain the two instances of a digit within a unit (box, column, or row). There can be other instances of the digit in the same unit.
	var weakLinks [10]map[unitLink]bool
	g.findXCycleWeakLinks(&box, &weakLinks)
	g.findXCycleWeakLinks(&col, &weakLinks)
	g.findXCycleWeakLinks(&row, &weakLinks)

	for d := 1; d <= 9; d++ { // Process nice chains.
		niceChain := findCycle(d, niceLoop, strongLinks[d], weakLinks[d])

		var overlap [rows][cols]bool
		for _, c := range niceChain {
			if c.strong {
				continue
			} // Only consider weak links.

			nl := neighbors(c.left)
			nr := neighbors(c.right)
			for r := zero; r < rows; r++ {
				for c := zero; c < cols; c++ {
					overlap[r][c] = nl[r][c] && nr[r][c]
				}
			}
			overlap[c.left.r][c.left.c] = false // Members of the chain cannot be removed.
			overlap[c.right.r][c.right.c] = false

			// In each cell that is seen by both ends of a weak chain link, the digit can be removed.
			for r := zero; r < rows; r++ {
				for c := zero; c < cols; c++ {
					if overlap[r][c] {
						if g.pt(point{r, c}).andNot(1 << d) {
							g.cellChange(&res, verbose, "xCycles: nice chain removes %d from %s\n", d, point{r, c})
						}
					}
				}
			}
		}
	}

	if res {
		return
	}

	for d := 1; d <= 9; d++ { // Process strong chains.
		strongChain := findCycle(d, strongLoop, strongLinks[d], weakLinks[d])

		// Find the strong discontinuity (two strong links in a row) and fix the digit at the intersection to the current digit.
		length := len(strongChain)
		if length >= 3 {
			first := strongChain[0]
			last := strongChain[length-1]
			if first.strong && last.strong { // If the first and last links are strong, the discontinuity is the last point in the chain (last.right).
				if g.pt(last.right).setTo(1 << d) {
					g.cellChange(&res, verbose, "xCycles: strong chain sets %s to %d\n", last.right, d)
				}
			} else { // Search the chain for the discontinuity.
				for i := 0; i < length-1; i++ {
					if strongChain[i].strong && strongChain[i+1].strong {
						if g.pt(strongChain[i].right).setTo(1 << d) {
							g.cellChange(&res, verbose, "xCycles: strong chain sets %s to %d\n", strongChain[i].right, d)
						}
						break // Once we find the discontinuity, we can stop looking because there can be only one ("Highlander").
					}
				}
			}
		}
	}

	if res {
		return
	}

	for d := 1; d <= 9; d++ { // Process weak chains.
		weakChain := findCycle(d, weakLoop, strongLinks[d], weakLinks[d])

		// Find the weak discontinuity (two weak links in a row) and remove the current digit at the intersection as a candidate.
		length := len(weakChain)
		if length >= 3 {
			first := weakChain[0]
			last := weakChain[length-1]
			if !first.strong && !last.strong { // If the first and last links are weak, the discontinuity is the last point in the chain (last.right).
				if g.pt(last.right).andNot(1 << d) {
					g.cellChange(&res, verbose, "xCycles: weak chain removes %d from %s\n", d, last.right)
				}
			} else { // Search the chain for the discontinuity.
				for i := 0; i < length-1; i++ {
					if !weakChain[i].strong && !weakChain[i+1].strong {
						if g.pt(weakChain[i].right).andNot(1 << d) {
							g.cellChange(&res, verbose, "xCycles: weak chain removes %d from %s\n", d, weakChain[i].right)
						}
						break // Once we find the discontinuity, we can stop looking because there can be only one ("Highlander").
					}
				}
			}
		}
	}

	return
}

func (g *Grid) checkUnit(d int, p1, p2 point) bool {
	b1 := boxOfPoint(p1)
	b2 := boxOfPoint(p2)
	if b1 == b2 {
		count := 0
		for _, p := range box.unit[b1] {
			if *g.pt(p)&cell(1<<d) != 0 {
				count++
			}
			if count > 2 {
				return true
			}
		}
	}

	if p1.c == p2.c {
		count := 0
		for _, p := range col.unit[b1] {
			if *g.pt(p)&cell(1<<d) != 0 {
				count++
			}
			if count > 2 {
				return true
			}
		}
	}

	if p1.r == p2.r {
		count := 0
		for _, p := range row.unit[b1] {
			if *g.pt(p)&cell(1<<d) != 0 {
				count++
			}
			if count > 2 {
				return true
			}
		}
	}

	return false
}

func (g *Grid) findXCycleWeakLinks(gr *group, weakLinks *[10]map[unitLink]bool) {
	for pi, ps := range gr.unit {
		points := g.digitPoints(ps)

		for d := 1; d <= 9; d++ {
			p := points[d]

			if len(p) < 3 {
				continue
			}

			w := &(*weakLinks)[d]
			if *w == nil {
				*w = make(map[unitLink]bool)
			}

			for _, p1 := range p {
				for _, p2 := range p {
					if p1 == p2 {
						continue
					}

					(*w)[sortLink(unitLink{pair{p1, p2}, d, fmt.Sprintf("%s %d", gr.name, pi), false})] = true
				}
			}
		}
	}
}

func chainLoops(chain []unitLink) bool {
	if len(chain) == 0 {
		return false
	}

	return chain[0].left == chain[len(chain)-1].right
}

func chainValid(checkComplete bool, kind loopKind, chain []unitLink) bool {
	doubleStrong := false
	doubleWeak := false

	switch len(chain) {
	case 0:
		return false
	case 1:
		return true
	default:
		if checkComplete && len(chain) > 2 { // If the chain wraps around and the first and last links match in strength, we need to account for that.
			first := chain[0]
			last := chain[len(chain)-1]
			if first.left == last.right {
				if last.strong == first.strong {
					if first.strong {
						doubleStrong = true
					} else {
						doubleWeak = true
					}
				}
			}
		}

		for i := 0; i < len(chain)-1; i++ {
			if chain[i].strong && chain[i+1].strong {
				if doubleStrong {
					return false
				}
				doubleStrong = true
			} else if !chain[i].strong && !chain[i+1].strong {
				if doubleWeak {
					return false
				}
				doubleWeak = true
			}
		}
	}

	if doubleStrong && doubleWeak {
		return false
	}

	switch kind {
	case niceLoop:
		return !doubleStrong && !doubleWeak
	case strongLoop:
		if checkComplete {
			return doubleStrong && !doubleWeak
		}
		return true
	case weakLoop:
		if checkComplete {
			return !doubleStrong && doubleWeak
		}
		return true
	}

	return false
}

func findCycle(digit int, kind loopKind, strongLinks, weakLinks map[unitLink]bool) (res []unitLink) {
	// Try each strong link as the start of the chain and keep the longest chain we can form.
	for s := range strongLinks {
		chain := []unitLink{s}
		best := []unitLink{}
		findCycleRecursive(digit, kind, chain, &best, strongLinks, weakLinks)

		if len(best) > len(res) {
			res = res[:0]
			for _, c := range best {
				res = append(res, c)
			}
		}
	}

	return
}

func findCycleRecursive(digit int, kind loopKind, chain []unitLink, best *[]unitLink, strongLinks, weakLinks map[unitLink]bool) {
	// If the right side of the last item in the chain links back to the head, we are done. TODO: keep searching for a longer chain.
	if chainValid(true, kind, chain) && chainLoops(chain) {
		if len(chain) > len(*best) {
			*best = (*best)[:0]
			for _, c := range chain {
				*best = append(*best, c)
			}
		}
	}

	last := chain[len(chain)-1]
	var candidates []unitLink

	strongAllowed := false
	weakAllowed := false

	switch kind {
	case niceLoop:
		strongAllowed = !last.strong
		weakAllowed = last.strong
	case strongLoop:
		strongAllowed = true
		weakAllowed = last.strong
	case weakLoop:
		strongAllowed = !last.strong
		weakAllowed = true
	}

	if strongAllowed {
	strongOuter:
		for s := range strongLinks {
			reversed := unitLink{pair{s.right, s.left}, s.digit, s.unit, s.strong}

			for _, c := range chain {
				if s.pair == c.pair || reversed.pair == c.pair { // Already in the chain, skip.
					continue strongOuter
				}
			}

			if last.unit == s.unit {
				continue
			}

			if last.right == s.left {
				candidates = append(candidates, s)
			}
			if last.right == reversed.left {
				candidates = append(candidates, reversed)
			}
		}
	}

	if weakAllowed {
	weakOuter:
		for w := range weakLinks {
			reversed := unitLink{pair{w.right, w.left}, w.digit, w.unit, w.strong}

			for _, c := range chain {
				if w.pair == c.pair || reversed.pair == c.pair { // Already in the chain, skip.
					continue weakOuter
				}
			}

			if last.unit == w.unit {
				continue
			}

			if last.right == w.left {
				candidates = append(candidates, w)
			}
			if last.right == reversed.left {
				candidates = append(candidates, reversed)
			}
		}
	}

	for _, c := range candidates {
		chain = append(chain, c)

		if chainValid(false, kind, chain) {
			findCycleRecursive(digit, kind, chain, best, strongLinks, weakLinks)
		}
		chain = chain[:len(chain)-1]
	}

	return
}
