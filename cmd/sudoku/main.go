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
	"log"
	"os"
	"path"
	"strings"

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
	level4Count int

	input inputs
)

func init() {
	flag.IntVar(&level0Count, "0", 0, "`count` of easy games to generate")
	flag.IntVar(&level1Count, "1", 0, "`count` of medium games to generate")
	flag.IntVar(&level2Count, "2", 0, "`count` of hard games to generate")
	flag.IntVar(&level3Count, "3", 0, "`count` of ridiculous games to generate")
	flag.IntVar(&level4Count, "4", 0, "`count` of insane (nearly impossible) games to generate")

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

	if len(input) > 0 && (level0Count > 0 || level1Count > 0 || level2Count > 0 || level3Count > 0 || level4Count > 0) {
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
				log.Printf("Encoded: %s", line)

				grid, err := generator.ParseEncoded(line)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				grid.Display()
				maxLevel, solved := grid.Reduce()
				grid.Display()
				if solved {
					sol++
					log.Printf("level: %s, solved", maxLevel)
				} else {
					log.Printf("level: %s, not solved", maxLevel)
					solutions := make([]*generator.Grid, 0)
					grid.Search(&solutions)
					switch len(solutions) {
					case 0:
						log.Println("still not solved after search")
					case 1:
						sol++
						log.Println("single solution found")
						solutions[0].Display()
					default:
						log.Println("multiple solutions found")
						for _, s := range solutions {
							s.Display()
						}
					}
				}
			}
			log.Printf("solved %d of %d", sol, all)
		}
	} else { // Generate puzzles of levels given in -0, -1, -2, -3, -4.
		grid := generator.Randomize()
		grid.Display()
		maxLevel, solved := grid.Reduce()
		grid.Display()
		if solved {
			log.Printf("level: %s, solved", maxLevel)
		} else {
			log.Printf("level: %s, not solved", maxLevel)
			solutions := make([]*generator.Grid, 0)
			grid.Search(&solutions)
			switch len(solutions) {
			case 0:
				log.Println("still not solved after search")
			case 1:
				log.Println("single solution found")
				solutions[0].Display()
			default:
				log.Println("multiple solutions found")
				for _, s := range solutions {
					s.Display()
				}
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
