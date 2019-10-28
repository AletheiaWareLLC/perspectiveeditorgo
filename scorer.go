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
)

const (
	BAD  = -1
	GOOD = 1
)

var (
	left = &perspectivego.Location{
		X: -1,
	}
	right = &perspectivego.Location{
		X: 1,
	}
	down = &perspectivego.Location{
		Y: -1,
	}
	up = &perspectivego.Location{
		Y: 1,
	}
	backward = &perspectivego.Location{
		Z: -1,
	}
	foreward = &perspectivego.Location{
		Z: 1,
	}
	directions = []*perspectivego.Location{
		left,
		right,
		down,
		up,
		backward,
		foreward,
	}
)

// Score: number of rotations needed to navigate to goal
// Penalty: number of unvisitable elements
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
	tested := make(map[string]int)
	visited := make(map[string]bool)
	rotations, direction := ScoreDirections(blocks, goals, portals, size, sphere.Location, tested, visited, false)
	if direction != down {
		// Add initial rotation
		rotations += 1
	}
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
	return rotations, penalty
}

func ScoreDirections(blocks, goals map[string]bool, portals map[string]*perspectivego.Location, size uint32, sphere *perspectivego.Location, tested map[string]int, visited map[string]bool, portaled bool) (int, *perspectivego.Location) {
	min := BAD
	dir := down
	posId := sphere.String()
	for _, d := range directions {
		id := posId + d.String()
		rotations, ok := tested[id]
		if !ok {
			tested[id] = BAD // Set now, update later to avoid loops
			rotations = ScoreDirection(blocks, goals, portals, size, d, &perspectivego.Location{
				X: sphere.X,
				Y: sphere.Y,
				Z: sphere.Z,
			}, tested, visited, portaled)
			tested[id] = rotations
		}
		if rotations >= 0 && (rotations < min || min == BAD) {
			min = rotations
			dir = d
		}
	}
	return min, dir
}

func ScoreDirection(blocks, goals map[string]bool, portals map[string]*perspectivego.Location, size uint32, direction *perspectivego.Location, sphere *perspectivego.Location, tested map[string]int, visited map[string]bool, portaled bool) int {
	// log.Println("Scoring Direction:", direction)
	rotations := 0
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
			return rotations
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
					portaled = true
					usage[key] = uses + 1
					visited[key] = true
					visited[link.String()] = true
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
			visited[next.String()] = true
			r, _ := ScoreDirections(blocks, goals, portals, size, sphere, tested, visited, portaled)
			if r >= 0 {
				return r + rotations + GOOD
			}
			return r
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
