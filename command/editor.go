/*
 * Copyright 2019 Aletheia Ware LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"github.com/AletheiaWareLLC/joygo"
	"github.com/AletheiaWareLLC/perspectiveeditorgo"
	"github.com/AletheiaWareLLC/perspectivego"
	"github.com/golang/protobuf/proto"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "create-world":
			if len(os.Args) > 4 {
				name := os.Args[2]
				size, err := strconv.Atoi(os.Args[3])
				if err != nil {
					log.Fatal(err)
				}
				if size < 0 {
					log.Fatal("World size must be postive")
				}
				if size%2 == 0 {
					log.Fatal("World size must be odd")
				}
				foreground := os.Args[4]
				background := os.Args[5]
				world := &perspectivego.World{
					Name:             name,
					Size:             uint32(size),
					ForegroundColour: foreground,
					BackgroundColour: background,
				}
				writer := os.Stdout
				if len(os.Args) > 6 {
					log.Println("Writing:", os.Args[6])
					file, err := os.OpenFile(os.Args[6], os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
					if err != nil {
						log.Fatal(err)
					}
					defer file.Close()
					writer = file
				}
				if err := perspectivego.WriteWorld(writer, world); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Println("create-world <name> <size> <foreground-colour> <background-colour> (write to stdout)")
				log.Println("create-world <name> <size> <foreground-colour> <background-colour> <output>")
			}
		case "show-world":
			if len(os.Args) > 2 {
				path := os.Args[2]
				world, err := perspectivego.ReadWorldFile(path)
				if err != nil {
					log.Fatal(err)
				}
				log.Println(proto.MarshalTextString(world))
			} else {
				log.Println("show-world <world>")
			}
		case "add-shader":
			if len(os.Args) > 7 {
				path := os.Args[2]
				world, err := perspectivego.ReadWorldFile(path)
				if err != nil {
					log.Fatal(err)
				}
				name := os.Args[3]
				attributes := strings.Split(os.Args[4], ",")
				uniforms := strings.Split(os.Args[5], ",")
				vertex, err := ioutil.ReadFile(os.Args[6])
				if err != nil {
					log.Fatal(err)
				}
				fragment, err := ioutil.ReadFile(os.Args[7])
				if err != nil {
					log.Fatal(err)
				}
				if world.Shader == nil {
					world.Shader = make(map[string]*joygo.Shader)
				}
				world.Shader[name] = &joygo.Shader{
					Name:           name,
					VertexSource:   string(vertex),
					FragmentSource: string(fragment),
					Attributes:     attributes,
					Uniforms:       uniforms,
				}
				if err := perspectivego.WriteWorldFile(path, world); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Println("add-shader <world> <name> <attributes> <uniforms> <vertex-source-file> <fragment-source-file>")
			}
		case "add-puzzle":
			if len(os.Args) > 2 {
				path := os.Args[2]
				world, err := perspectivego.ReadWorldFile(path)
				if err != nil {
					log.Fatal(err)
				}
				reader := os.Stdin
				if len(os.Args) > 3 {
					file, err := os.Open(os.Args[3])
					if err != nil {
						log.Fatal(err)
					}
					reader = file
				}
				puzzle, err := perspectivego.ReadPuzzle(reader)
				world.Puzzle = append(world.Puzzle, puzzle)
				if err := perspectivego.WriteWorldFile(path, world); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Println("add-puzzle <world> (read from stdin)")
				log.Println("add-puzzle <world> <file>")
			}
		case "generate-puzzle":
			if len(os.Args) > 18 {
				size, err := strconv.Atoi(os.Args[2])
				if err != nil {
					log.Fatal(err)
				}
				if size < 0 {
					log.Fatal("World size must be postive")
				}
				if size%2 == 0 {
					log.Fatal("World size must be odd")
				}
				score, err := strconv.Atoi(os.Args[3])
				if err != nil {
					log.Fatal(err)
				}
				description := os.Args[4]
				outlineMesh := os.Args[5]
				outlineColour := os.Args[6]
				goalCount, err := strconv.Atoi(os.Args[7])
				if err != nil {
					log.Fatal("Goal count error:", err)
				}
				goalMesh := strings.Split(os.Args[8], ",")
				goalColour := strings.Split(os.Args[9], ",")
				sphereCount, err := strconv.Atoi(os.Args[10])
				if err != nil {
					log.Fatal("Sphere count error:", err)
				}
				sphereMesh := strings.Split(os.Args[11], ",")
				sphereColour := strings.Split(os.Args[12], ",")
				blockCount, err := strconv.Atoi(os.Args[13])
				if err != nil {
					log.Fatal("Block count error:", err)
				}
				blockMesh := strings.Split(os.Args[14], ",")
				blockColour := strings.Split(os.Args[15], ",")
				portalCount, err := strconv.Atoi(os.Args[16])
				if err != nil {
					log.Fatal("Portal count error:", err)
				}
				if portalCount%2 != 0 {
					log.Fatal("Portal count must be even")
				}
				portalMesh := strings.Split(os.Args[17], ",")
				portalColour := strings.Split(os.Args[18], ",")

				var outline *perspectivego.Outline
				if outlineMesh != "" && outlineColour != "" {
					outline = &perspectivego.Outline{
						Mesh:   outlineMesh,
						Colour: outlineColour,
					}
				}
				puzzle := &perspectivego.Puzzle{}
				if description != "" {
					puzzle.Description = description
				}
				if outline != nil {
					puzzle.Outline = outline
				}
				start := time.Now()
				iteration := 0
				max := 0
				for x := 0; x <= 100; iteration++ {
					if iteration == (x * x * x * x) {
						log.Println(x, "^ 4 =", iteration)
						x++
					}
					perspectiveeditorgo.Generate(puzzle, uint32(size), goalCount, goalMesh, goalColour, sphereCount, sphereMesh, sphereColour, blockCount, blockMesh, blockColour, portalCount, portalMesh, portalColour)
					s, p := perspectiveeditorgo.Score(puzzle, uint32(size))
					puzzle.Target = uint32(s)
					if s > max {
						max = s
						log.Println("Score:", s, "/", score)
						log.Println("Penalty:", p)
						log.Println("Iteration:", iteration)
						log.Println("Elapsed:", time.Since(start))
						log.Println("Puzzle:", puzzle)
						if s > score {
							writer := os.Stdout
							if len(os.Args) > 19 {
								log.Println("Writing:", os.Args[19])
								file, err := os.OpenFile(os.Args[19], os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
								if err != nil {
									log.Fatal(err)
								}
								defer file.Close()
								writer = file
							}
							if err := perspectivego.WritePuzzle(writer, puzzle); err != nil {
								log.Fatal(err)
							}
							break
						}
					}
				}
			} else {
				log.Println("generate-puzzle <size> <score> <description> <outline-mesh> <outline-colour> <goal-count> <goal-mesh...> <goal-colour...> <sphere-count> <sphere-mesh...> <sphere-colour...> <block-count> <block-mesh...> <block-colour...> <portal-count> <portal-mesh...> <portal-colour...>")
			}
		case "generate-world":
			if len(os.Args) > 18 {
				size, err := strconv.Atoi(os.Args[2])
				if err != nil {
					log.Fatal(err)
				}
				if size < 0 {
					log.Fatal("World size must be postive")
				}
				if size%2 == 0 {
					log.Fatal("World size must be odd")
				}
				description := os.Args[3]
				outlineMesh := os.Args[4]
				outlineColour := os.Args[5]
				goalCount, err := strconv.Atoi(os.Args[6])
				if err != nil {
					log.Fatal("Goal count error:", err)
				}
				goalMesh := strings.Split(os.Args[7], ",")
				goalColour := strings.Split(os.Args[8], ",")
				sphereCount, err := strconv.Atoi(os.Args[9])
				if err != nil {
					log.Fatal("Sphere count error:", err)
				}
				sphereMesh := strings.Split(os.Args[10], ",")
				sphereColour := strings.Split(os.Args[11], ",")
				blockCount, err := strconv.Atoi(os.Args[12])
				if err != nil {
					log.Fatal("Block count error:", err)
				}
				blockMesh := strings.Split(os.Args[13], ",")
				blockColour := strings.Split(os.Args[14], ",")
				portalCount, err := strconv.Atoi(os.Args[15])
				if err != nil {
					log.Fatal("Portal count error:", err)
				}
				if portalCount%2 != 0 {
					log.Fatal("Portal count must be even")
				}
				portalMesh := strings.Split(os.Args[16], ",")
				portalColour := strings.Split(os.Args[17], ",")

				var outline *perspectivego.Outline
				if outlineMesh != "" && outlineColour != "" {
					outline = &perspectivego.Outline{
						Mesh:   outlineMesh,
						Colour: outlineColour,
					}
				}
				puzzle := &perspectivego.Puzzle{}
				if description != "" {
					puzzle.Description = description
				}
				if outline != nil {
					puzzle.Outline = outline
				}
				penalties := make(map[int]int)
				start := time.Now()
				iteration := 0
				if len(os.Args) > 18 {
					files, err := ioutil.ReadDir(os.Args[18])
					if err != nil {
						log.Fatal(err)
					}

					for _, file := range files {
						filename := path.Join(os.Args[18], file.Name())
						log.Println("File:", filename)
						file, err := os.Open(filename)
						if err != nil {
							log.Fatal(err)
						}
						defer file.Close()
						puzzle, err := perspectivego.ReadPuzzle(file)
						if err != nil {
							log.Fatal(err)
						}
						s, p := perspectiveeditorgo.Score(puzzle, uint32(size))
						penalties[s] = p
						log.Println("Score:", s)
						log.Println("Penalty:", p)
					}
				}
				for x := 0; x <= 100; iteration++ {
					if iteration == (x * x * x * x) {
						log.Println(x, "^ 4 =", iteration)
						x++
					}
					perspectiveeditorgo.Generate(puzzle, uint32(size), goalCount, goalMesh, goalColour, sphereCount, sphereMesh, sphereColour, blockCount, blockMesh, blockColour, portalCount, portalMesh, portalColour)
					s, p := perspectiveeditorgo.Score(puzzle, uint32(size))
					if s > 0 {
						puzzle.Target = uint32(s)
						penalty, ok := penalties[s]
						if !ok || p < penalty {
							penalties[s] = p
							log.Println("Score:", s)
							log.Println("Penalty:", p)
							log.Println("Iteration:", iteration)
							log.Println("Elapsed:", time.Since(start))
							log.Println("Puzzle:", puzzle)
							writer := os.Stdout
							if len(os.Args) > 18 {
								filename := path.Join(os.Args[18], "/puzzle"+strconv.Itoa(s)+".txt")
								log.Println("Writing:", filename)
								file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
								if err != nil {
									log.Fatal(err)
								}
								defer file.Close()
								writer = file
							}
							if err := perspectivego.WritePuzzle(writer, puzzle); err != nil {
								log.Fatal(err)
							}
						}
					}
				}
			} else {
				log.Println("generate-world <size> <description> <outline-mesh> <outline-colour> <goal-count> <goal-mesh...> <goal-colour...> <sphere-count> <sphere-mesh...> <sphere-colour...> <block-count> <block-mesh...> <block-colour...> <portal-count> <portal-mesh...> <portal-colour...>")
			}
		case "score-puzzle":
			if len(os.Args) > 3 {
				size, err := strconv.Atoi(os.Args[2])
				if err != nil {
					log.Fatal(err)
				}
				if size < 0 {
					log.Fatal("World size must be postive")
				}
				if size%2 == 0 {
					log.Fatal("World size must be odd")
				}
				file, err := os.Open(os.Args[3])
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()
				puzzle, err := perspectivego.ReadPuzzle(file)
				if err != nil {
					log.Fatal(err)
				}
				s, p := perspectiveeditorgo.Score(puzzle, uint32(size))
				log.Println("Score:", s)
				log.Println("Penalty:", p)
			} else {
				log.Println("score-puzzle <size> <path>")
			}
		case "score-world":
			if len(os.Args) > 3 {
				size, err := strconv.Atoi(os.Args[2])
				if err != nil {
					log.Fatal(err)
				}
				if size < 0 {
					log.Fatal("World size must be postive")
				}
				if size%2 == 0 {
					log.Fatal("World size must be odd")
				}
				files, err := ioutil.ReadDir(os.Args[3])
				if err != nil {
					log.Fatal(err)
				}

				for _, file := range files {
					log.Println("File:", file.Name())
					file, err := os.Open(path.Join(os.Args[3], file.Name()))
					if err != nil {
						log.Fatal(err)
					}
					defer file.Close()
					puzzle, err := perspectivego.ReadPuzzle(file)
					if err != nil {
						log.Fatal(err)
					}
					s, p := perspectiveeditorgo.Score(puzzle, uint32(size))
					log.Println("Score:", s)
					log.Println("Penalty:", p)
				}
			} else {
				log.Println("score-world <size> <path>")
			}
		default:
			log.Println("Cannot handle", os.Args[1])
		}
	} else {
		PrintUsage(os.Stdout)
	}
}

func Exists(filename string) bool {
	log.Println("Checking:", filename)
	_, err := os.Stat(filename)
	if err != nil {
		log.Println(err)
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func PrintUsage(output io.Writer) {
	fmt.Fprintln(output, "Perspective Editor Usage:")
	fmt.Fprintln(output, "\tperspective-editor - display usage")
	fmt.Fprintln(output, "\tperspective-editor create-world [name] [size] [foreground-colour] [background-colour] - creates a new world with the given name, size and colour scheme")
	fmt.Fprintln(output, "\tperspective-editor show-world [world] - shows the given world")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "\tperspective-editor add-shader [world] [name] [attributes] [uniforms] [vertex-source-file] [fragment-source-file] - adds a shader with the given name to the world")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "\tperspective-editor add-puzzle [world] - adds a puzzle to the world")
	fmt.Fprintln(output, "\tperspective-editor generate-puzzle [size] [score] [description] [outline-mesh] [outline-colour] [goal-count] [goal-mesh...] [goal-colour...] [sphere-count] [sphere-mesh...] [sphere-colour...] [block-count] [block-mesh...] [block-colour...] [portal-count] [portal-mesh...] [portal-colour...] - generates a new puzzle with the given attributes")
	fmt.Fprintln(output, "\tperspective-editor generate-world [size] [description] [outline-mesh] [outline-colour] [goal-count] [goal-mesh...] [goal-colour...] [sphere-count] [sphere-mesh...] [sphere-colour...] [block-count] [block-mesh...] [block-colour...] [portal-count] [portal-mesh...] [portal-colour...] - generates a new pool of puzzle with the given attributes")
	fmt.Fprintln(output, "\tperspective-editor score-puzzle [size] [path] - scores the puzzle under the given path")
	fmt.Fprintln(output, "\tperspective-editor score-world [size] [path] - scores all puzzles under the given path")
}
