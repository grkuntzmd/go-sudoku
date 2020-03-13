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
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
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
	flag.IntVar(&level4Count, "4", 0, "`count` of (almost) impossible games to generate")

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

	r := make([]io.Reader, 0)
	for _, i := range input {

		f, err := os.Open(i)
		if err != nil {
			panic(fmt.Sprintf("cannot open %s for reading", i))
		}
		r = append(r, f)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "buildStamp: %s, gitHash: %s, version: %s\n", buildStamp, gitHash, version)
	flag.PrintDefaults()
}

func (i *inputs) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *inputs) String() string {
	return strings.Join(*i, ",")
}
