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

package perspectiveeditorgo

import (
	"github.com/AletheiaWareLLC/perspectivego"
	//"github.com/golang/protobuf/proto"
	//"log"
	"sync"
)

const (
	BAD  = -1
	GOOD = 1
)

var (
	directions = []*perspectivego.Location{
		// Left
		&perspectivego.Location{
			X: -1,
		},
		// Right
		&perspectivego.Location{
			X: 1,
		},
		// Down
		&perspectivego.Location{
			Y: -1,
		},
		// Up
		&perspectivego.Location{
			Y: 1,
		},
		// Backwards
		&perspectivego.Location{
			Z: -1,
		},
		// Forewards
		&perspectivego.Location{
			Z: 1,
		},
	}
	testedMutex  = sync.RWMutex{}
	visitedMutex = sync.RWMutex{}
)

func Score(puzzle *perspectivego.Puzzle, size uint32) (int, int) {
	// log.Println("Scoring Puzzle:", puzzle)
	// TODO support multiple spheres
	sphere := puzzle.Sphere[0]
	blocks := make(map[string]bool, len(puzzle.Block))
	for _, b := range puzzle.Block {
		blocks[b.Location.String()] = true
	}
	goals := make(map[string]bool, len(puzzle.Goal))
	for _, g := range puzzle.Goal {
		goals[g.Location.String()] = true
	}
	portals := make(map[string]*perspectivego.Location, len(puzzle.Portal))
	for _, p := range puzzle.Portal {
		portals[p.Location.String()] = p.Link
	}
	tested := make(map[string]bool)
	visited := make(map[string]bool)
	score := ScoreDirections(blocks, goals, portals, size, sphere.Location, tested, visited, false)
	penalty := 0
	// Check all blocks were visited
	for _, b := range puzzle.Block {
		if !visited[b.Location.String()] {
			// log.Println("Unvisited Block: " + b.String())
			penalty++
		}
	}
	// Check all portals were visited
	for _, p := range puzzle.Portal {
		if !visited[p.Location.String()] {
			// log.Println("Unvisited Portal: " + p.String())
			penalty++
		}
	}
	return score, penalty
}

func ScoreDirections(blocks, goals map[string]bool, portals map[string]*perspectivego.Location, size uint32, sphere *perspectivego.Location, tested, visited map[string]bool, portaled bool) int {
	limit := 0
	scores := make(chan int)
	posId := sphere.String()
	for _, dir := range directions {
		id := posId + dir.String()
		testedMutex.Lock()
		if tested[id] {
			testedMutex.Unlock()
		} else {
			tested[id] = true
			testedMutex.Unlock()
			limit++
			go func(bs, gs map[string]bool, ps map[string]*perspectivego.Location, s uint32, d *perspectivego.Location, l *perspectivego.Location, t, v map[string]bool, pd bool) {
				scores <- ScoreDirection(bs, gs, ps, s, d, l, t, v, pd)
			}(blocks, goals, portals, size, dir, &perspectivego.Location{
				X: sphere.X,
				Y: sphere.Y,
				Z: sphere.Z,
			}, tested, visited, portaled)
		}
	}
	sum := 0
	count := 0
	for i := 0; i < limit; i++ {
		s := <-scores
		if s >= 0 {
			sum += s
			count++
		}
	}
	if count > 0 {
		return sum / count
	}
	return BAD
}

func ScoreDirection(blocks, goals map[string]bool, portals map[string]*perspectivego.Location, size uint32, direction *perspectivego.Location, sphere *perspectivego.Location, tested, visited map[string]bool, portaled bool) int {
	// log.Println("Scoring Direction:", direction)
	score := 0
	// Tracks portal usage to prevent infinite portal loops
	usage := make(map[string]int)
	for {
		// log.Println("Sphere:", sphere)
		if Abs(sphere.X) > size || Abs(sphere.Y) > size || Abs(sphere.Z) > size {
			// log.Println("Out of Bounds")
			return BAD
		}
		key := sphere.String()
		if goals[key] {
			// log.Println("Goal")
			return score
		}
		if !portaled {
			link, ok := portals[key]
			if ok {
				uses, ok := usage[key]
				if !ok {
					uses = 0
				}
				if uses < 6 {
					// log.Println("Portal")
					sphere.X = link.X
					sphere.Y = link.Y
					sphere.Z = link.Z
					score += GOOD
					portaled = true
					usage[key] = uses + 1
					visitedMutex.Lock()
					visited[key] = true
					visited[link.String()] = true
					visitedMutex.Unlock()
					continue
				} else {
					// log.Println("Infinite Portal Loop")
					return BAD
				}
			}
		}
		next := &perspectivego.Location{
			X: sphere.X + direction.X,
			Y: sphere.Y + direction.Y,
			Z: sphere.Z + direction.Z,
		}
		if blocks[next.String()] {
			// log.Println("Block")
			visitedMutex.Lock()
			visited[next.String()] = true
			visitedMutex.Unlock()
			s := ScoreDirections(blocks, goals, portals, size, sphere, tested, visited, portaled)
			if s >= 0 {
				return s + score + GOOD
			}
			return s
		}
		sphere.X = next.X
		sphere.Y = next.Y
		sphere.Z = next.Z
		portaled = false
	}
}

func Abs(a int32) uint32 {
	if a < 0 {
		return uint32(-a)
	}
	return uint32(a)
}
