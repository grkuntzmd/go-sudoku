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

func (g *Grid) xyChains(verbose uint) (res bool) {
	var links [10]map[link]bool
	g.findBivalueLinks(&box, &links)
	g.findBivalueLinks(&col, &links)
	g.findBivalueLinks(&row, &links)

	linkEnds := make(map[point][]link)
	for d := 1; d <= 9; d++ {
		for l := range links[d] {
			linkEnds[l.left] = append(linkEnds[l.left], l)
			linkEnds[l.right] = append(linkEnds[l.right], l)
		}
	}

	for d := 1; d <= 9; d++ {
		for l := range links[d] {
			if l.left.r != 0 || l.left.c != 0 {
				continue
			}

			g.followChain(&res, verbose, []link{l}, links, linkEnds)
			g.followChain(&res, verbose, []link{l.reverse()}, links, linkEnds)
		}
	}

	verbose = 0
	return
}

func (g *Grid) followChain(res *bool, verbose uint, chain []link, links [10]map[link]bool, linkEnds map[point][]link) {
	firstLink := chain[0]
	lastLink := chain[len(chain)-1]

	if len(chain)%2 == 1 &&
		boxOfPoint(firstLink.left) != boxOfPoint(lastLink.right) &&
		firstLink.left.c != lastLink.right.c &&
		firstLink.left.r != lastLink.left.r {
		front := *g.pt(firstLink.left) &^ (1 << firstLink.digit)
		back := *g.pt(lastLink.right) &^ (1 << lastLink.digit)
		if front == back {

			// Find the overlap between the left of front and the right of back.
			firstInfluence := influence[firstLink.left.r][firstLink.left.c]
			lastInfluence := influence[lastLink.right.r][lastLink.right.c]
			overlap := uint128{
				firstInfluence.ms & lastInfluence.ms,
				firstInfluence.ls & lastInfluence.ls,
			}

			// Remove all points in the chain from the overlap.
			for _, l := range chain {
				overlap.unset(l.left)
				overlap.unset(l.right)
			}

			processInfluence(overlap, func(r, c uint8) {
				p := point{r, c}
				if g.pt(p).andNot(front) {
					g.cellChange(res, verbose, "xyChains: remove %s from (%d, %d) because it is seen by %s and %s (chain: %v)\n", front, r, c, firstLink.left, lastLink.right, chain)

					// Once a candidate digit is removed, that point can no longer be a part of any chain since it will not be bivalued.
					for _, l := range linkEnds[p] {
						delete(links[l.digit], l)
					}
					delete(linkEnds, p)
				}
			})
		}
	}

outer:
	for _, l := range linkEnds[lastLink.right] {
		if lastLink.digit == l.digit {
			continue
		}

		for _, c := range chain {
			if l == c || l.reverse() == c {
				continue outer
			}
			if c != lastLink && (l.left == c.left || l.left == c.right || l.right == c.left || l.right == c.right) {
				continue outer
			}
		}

		var newLink link
		if l.left == lastLink.right {
			newLink = l
			chain = append(chain, newLink)
		} else {
			newLink = l.reverse()
			chain = append(chain, newLink)
		}

		g.followChain(res, verbose, chain, links, linkEnds)

		chain = chain[:len(chain)-1]
	}
}
