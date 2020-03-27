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

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
	"time"

	"dogdaze.org/sudoku/generator"
)

type inputs []string

var (
	buildInfo  string
	buildStamp string
	gitHash    string
	version    string

	level0Count int
	level1Count int
	level2Count int
	level3Count int
	// level4Count int

	input inputs
)

func init() {
	flag.IntVar(&level0Count, "0", 0, "`count` of easy games to generate")
	flag.IntVar(&level1Count, "1", 0, "`count` of standard games to generate")
	flag.IntVar(&level2Count, "2", 0, "`count` of hard games to generate")
	flag.IntVar(&level3Count, "3", 0, "`count` of expert games to generate")
	// flag.IntVar(&level4Count, "4", 0, "`count` of extreme (nearly impossible) games to generate")

	flag.Var(&input, "i", "`file` containing input patterns (may be repeated)")

	if buildInfo != "" {
		parts := strings.Split(buildInfo, "|")
		if len(parts) >= 3 {
			buildStamp = parts[0]
			gitHash = parts[1]
			version = parts[2]
		}
	}
}

func main() {
	flag.CommandLine.Usage = usage
	flag.Parse()

	if len(input) > 0 && (level0Count > 0 || level1Count > 0 || level2Count > 0 || level3Count > 0 /* || level4Count > 0 */) {
		usage()
		os.Exit(1)
	}

	if len(input) > 0 { // Handle -i files.
		for _, i := range input {
			f, err := os.Open(i)
			if err != nil {
				fmt.Printf("cannot open %s for reading; skipping\n", i)
				continue
			}
			defer f.Close()

			all := 0
			sol := 0
			s := bufio.NewScanner(f)
			for s.Scan() {
				all++
				line := s.Text()
				fmt.Printf("Encoded: %s\n", line)

				grid, err := generator.ParseEncoded(line)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				grid.Display()
				strategies := make(map[string]bool)
				start := time.Now()
				maxLevel, solved := grid.Reduce(&strategies)

				var names []string
				for n := range strategies {
					names = append(names, n)
				}
				sort.Slice(names, func(i, j int) bool { return names[i] < names[j] })

				grid.Display()
				if solved {
					sol++
					fmt.Printf("level: %s, solved, (%s) in %d milliseconds\n", maxLevel, strings.Join(names, ", "), time.Since(start).Milliseconds())
				} else {
					fmt.Printf("level: %s, not solved (%s) in %d milliseconds\n", maxLevel, strings.Join(names, ", "), time.Since(start).Milliseconds())
					solutions := make([]*generator.Grid, 0)
					grid.Search(&solutions)
					switch len(solutions) {
					case 0:
						fmt.Printf("still not solved after search, (%s) in %d milliseconds\n", strings.Join(names, ", "), time.Since(start).Milliseconds())
					case 1:
						sol++
						fmt.Printf("single solution found, (%s) in %d milliseconds\n", strings.Join(names, ", "), time.Since(start).Milliseconds())
						solutions[0].Display()
					default:
						fmt.Printf("multiple solutions found, (%s) in %d milliseconds\n", strings.Join(names, ", "), time.Since(start).Milliseconds())
						for _, s := range solutions {
							s.Display()
						}
					}
				}
			}
			fmt.Printf("solved %d of %d\n", sol, all)
		}
	} else { // Generate puzzles of levels given in -0, -1, -2, -3, -4.
		numberOfWorkers := runtime.NumCPU()
		numberOfTasks := level0Count + level1Count + level2Count + level3Count // + level4Count

		tasks := make(chan generator.Level, numberOfTasks)
		results := make(chan *generator.Game, numberOfTasks)

		for w := 0; w < numberOfWorkers; w++ {
			go generator.Worker(tasks, results)
		}

		for t := 0; t < level0Count; t++ {
			tasks <- generator.Easy
		}

		for t := 0; t < level1Count; t++ {
			tasks <- generator.Standard
		}

		for t := 0; t < level2Count; t++ {
			tasks <- generator.Hard
		}

		for t := 0; t < level3Count; t++ {
			tasks <- generator.Expert
		}

		// for t := 0; t < level4Count; t++ {
		// 	tasks <- generator.Insane
		// }

		close(tasks)

		for t := 0; t < numberOfTasks; t++ {
			g := <-results
			if g != nil {
				fmt.Printf("%s (%d) %s\n", g.Level, g.Clues, strings.Join(g.Strategies, ", "))
				g.Puzzle.Display()
				g.Solution.Display()
			}
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", path.Base(os.Args[0]))
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nEither -i or level counts (-0, -1, -2, -3, -4) may be used, but not both.")
	fmt.Fprintf(os.Stderr, "\nbuildStamp: %s, gitHash: %s, version: %s\n", buildStamp, gitHash, version)
}

func (i *inputs) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *inputs) String() string {
	return strings.Join(*i, ",")
}
