/*
 * ElasticoFS
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package cmd parses the parameters and runs ElasticoFS
package cmd

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	elasticofs "github.com/dutchcoders/elasticofs/fs"
	"github.com/fatih/color"
	"github.com/op/go-logging"
)

var options = flag.String("o", "", "mount options")

func usage() {
	fmt.Fprintf(os.Stderr, "ElasticoFS mount your elasticsearch server.\n\n")

	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s http://127.0.0.1:9200/ MOUNTPOINT\n\n", os.Args[0])
	flag.PrintDefaults()
}

// Main is the actual run function
func Main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 2 {
		usage()
		os.Exit(2)
	}

	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend1Leveled)

	//	if context.GlobalBool("debug") {
	// backend1Leveled.SetLevel(logging.DEBUG, "")
	//	}

	opts := []func(*elasticofs.Config){}

	for _, option := range strings.Split(*options, ",") {
		vals := strings.Split(option, "=")

		switch vals[0] {
		case "uid":
			if len(vals) == 1 {
				fmt.Fprint(os.Stderr, color.RedString("Uid has no value\n"))
				os.Exit(2)
			} else if val, err := strconv.Atoi(vals[1]); err != nil {
				fmt.Fprint(os.Stderr, color.RedString("Uid is not a valid value: %s\n", vals[1]))
				os.Exit(2)
			} else {
				opts = append(opts, elasticofs.Uid(uint32(val)))
			}
		case "gid":
			if len(vals) == 1 {
				fmt.Fprint(os.Stderr, color.RedString("Uid has no value\n"))
				os.Exit(2)
			} else if val, err := strconv.Atoi(vals[1]); err != nil {
				fmt.Fprint(os.Stderr, color.RedString("Gid is not a valid value: %s\n", vals[1]))
				os.Exit(2)
			} else {
				opts = append(opts, elasticofs.Gid(uint32(val)))
			}
		case "cache":
			if len(vals) == 1 {
				fmt.Fprint(os.Stderr, color.RedString("Cache has no value\n"))
				os.Exit(2)
			} else {
				opts = append(opts, elasticofs.CacheDir(vals[1]))
			}
		}
	}

	target := flag.Arg(0)
	mountpoint := flag.Arg(1)

	if fs, err := elasticofs.New(
		append(opts,
			elasticofs.Mountpoint(mountpoint),
			elasticofs.Target(target),
			elasticofs.Debug(),
		)...,
	); err != nil {
		panic(err)
	} else if err := fs.Serve(); err != nil {
		panic(err)
	}
}
