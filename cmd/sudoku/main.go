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
	"html/template"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"

	"dogdaze.org/sudoku/generator"
	"github.com/pkg/browser"
)

type (
	inputs []string

	puzzle struct {
		Num int
		generator.Level
		Break bool
		template.HTML
	}

	solution struct {
		Num int
		template.HTML
	}
)

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

	input   inputs
	verbose uint
)

func init() {
	flag.IntVar(&level0Count, "0", 0, "`count` of easy games to generate")
	flag.IntVar(&level1Count, "1", 0, "`count` of standard games to generate")
	flag.IntVar(&level2Count, "2", 0, "`count` of hard games to generate")
	flag.IntVar(&level3Count, "3", 0, "`count` of expert games to generate")
	// flag.IntVar(&level4Count, "4", 0, "`count` of extreme (nearly impossible) games to generate")

	flag.Var(&input, "i", "`file` containing input patterns (may be repeated)")
	flag.UintVar(&verbose, "v", 0, "`verbosity` level; higher emits more messages")

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

				if !grid.Valid() {
					fmt.Fprintln(os.Stderr, "grid is invalid")
					continue
				}

				strategies := make(map[string]bool)
				maxLevel, solved := grid.Reduce(&strategies, verbose)

				var names []string
				for n := range strategies {
					names = append(names, n)
				}
				sort.Slice(names, func(i, j int) bool { return names[i] < names[j] })

				grid.Display()
				if solved {
					sol++
					fmt.Printf("level: %s, solved, (%s)\n", maxLevel, strings.Join(names, ", "))
				} else {
					fmt.Printf("level: %s, not solved (%s)\n", maxLevel, strings.Join(names, ", "))
					solutions := make([]*generator.Grid, 0)
					grid.Search(&solutions)
					switch len(solutions) {
					case 0:
						fmt.Printf("still not solved after search, (%s)\n", strings.Join(names, ", "))
					case 1:
						sol++
						fmt.Printf("single solution found, (%s)\n", strings.Join(names, ", "))
						solutions[0].Display()
					default:
						fmt.Printf("multiple solutions found, (%s)\n", strings.Join(names, ", "))
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
		// 	tasks <- generator.Extreme
		// }

		close(tasks)

		// games := make([]*generator.Game, 0, numberOfTasks)

		for t := 0; t < numberOfTasks; t++ {
			g := <-results
			if g != nil {
				fmt.Printf("%s (%d) %s\n", g.Level, g.Clues, strings.Join(g.Strategies, ", "))
				// g.Puzzle.Display()
				// g.Solution.Display()
				// games = append(games, g)
			}
		}

		// sort.Slice(games, func(i, j int) bool {
		// 	return games[i].Level < games[j].Level
		// })

		// html(games)
	}
}

func (i *inputs) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *inputs) String() string {
	return strings.Join(*i, ",")
}

func html(games []*generator.Game) {
	puzzles := make([]puzzle, 0, len(games))
	solutions := make([]solution, 0, len(games))

	for i, g := range games {
		puzzles = append(puzzles, puzzle{i + 1, g.Level, i%2 == 1, template.HTML(g.Puzzle.SVG(0.8, false, false, nil))})
		solutions = append(solutions, solution{i + 1, template.HTML(g.Solution.SVG(0.3, true, false, nil))})
	}

	t := template.Must(template.New("html").Parse(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Sudoku</title>

			<style>
				.break { page-break-after: always; }
				.solutions {
					display: flex;
					flex-direction: row;
					flex-wrap: wrap;
					justify-content: space-between;
					align-items: flex-start;
				}
			</style>
		</head>
		<body>
			{{ range .Puzzles }}
				<div {{ if .Break }}class="break"{{ end }} style="page-break-inside: avoid;">
					<h2>{{ .Num }} {{ .Level }}</h2>
					<p>{{ .HTML }}</p>
				</div>
			{{ end }}
			<p class="break"></p>
			<div class="solutions">
				{{ range .Solutions }}
					<div style="page-break-inside: avoid;">
						<h4>{{ .Num }}</h4>
						<p>{{ .HTML }}</p>
					</div>
				{{ end }}
			</div>
		</body>
		</html>
	`))

	var b strings.Builder
	if err := t.Execute(&b, struct {
		Puzzles   []puzzle
		Solutions []solution
	}{puzzles, solutions}); err != nil {
		panic(err)
	}

	if err := browser.OpenReader(strings.NewReader(b.String())); err != nil {
		panic(err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", path.Base(os.Args[0]))
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nEither -i or level counts (-0, -1, -2, -3, -4) may be used, but not both.")
	fmt.Fprintf(os.Stderr, "\nbuildStamp: %s, gitHash: %s, version: %s\n", buildStamp, gitHash, version)
}
