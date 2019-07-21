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
	mutex = sync.RWMutex{}
)

func Score(puzzle *perspectivego.Puzzle, size uint32) int {
	// log.Println("Scoring Puzzle:", puzzle)
	// TODO support multiple spheres
	sphere := puzzle.Sphere[0]
	return ScoreDirections(puzzle, size, sphere.Location, make(map[string]bool))
}

func ScoreDirections(puzzle *perspectivego.Puzzle, size uint32, sphere *perspectivego.Location, tested map[string]bool) int {
	limit := 0
	scores := make(chan int)
	posId := sphere.String()
	for _, dir := range directions {
		id := posId + dir.String()
		mutex.Lock()
		t, ok := tested[id]
		mutex.Unlock()
		if !ok || !t {
			mutex.Lock()
			tested[id] = true
			mutex.Unlock()
			limit++
			go func(p *perspectivego.Puzzle, s uint32, d *perspectivego.Location, l *perspectivego.Location, t map[string]bool) {
				scores <- ScoreDirection(p, s, d, l, t)
			}(puzzle, size, dir, &perspectivego.Location{
				X: sphere.X,
				Y: sphere.Y,
				Z: sphere.Z,
			}, tested)
		}
	}
	sum := 0
	count := 0
	for i := 0; i < limit; i++ {
		s := <-scores
		if s > 0 {
			sum += s
			count++
		}
	}
	if count > 0 {
		return sum / count
	}
	return sum
}

func ScoreDirection(puzzle *perspectivego.Puzzle, size uint32, direction *perspectivego.Location, sphere *perspectivego.Location, tested map[string]bool) int {
	// log.Println("Scoring Direction:", direction)
	score := 0
	portaled := false
	for {
		// log.Println("Sphere:", sphere)
		if Abs(sphere.X) > size || Abs(sphere.Y) > size || Abs(sphere.Z) > size {
			// log.Println("Out of Bounds")
			return score + BAD
		}
		for _, g := range puzzle.Goal {
			if g.Location.X == sphere.X && g.Location.Y == sphere.Y && g.Location.Z == sphere.Z {
				// log.Println("Goal")
				return score + GOOD
			}
		}
		for _, b := range puzzle.Block {
			if b.Location.X == (sphere.X+direction.X) && b.Location.Y == (sphere.Y+direction.Y) && b.Location.Z == (sphere.Z+direction.Z) {
				// log.Println("Block")
				s := ScoreDirections(puzzle, size, sphere, tested)
				if s <= 0 {
					return s
				}
				return s + score + GOOD
			}
		}
		if !portaled {
			for _, p := range puzzle.Portal {
				if p.Location.X == sphere.X && p.Location.Y == sphere.Y && p.Location.Z == sphere.Z {
					// log.Println("Portal")
					sphere.X = p.Link.X
					sphere.Y = p.Link.Y
					sphere.Z = p.Link.Z
					score += GOOD
					portaled = true
					break
				}
			}
			if portaled {
				continue
			}
		}
		sphere.X += direction.X
		sphere.Y += direction.Y
		sphere.Z += direction.Z
		portaled = false
	}
}

func Abs(a int32) uint32 {
	if a < 0 {
		return uint32(-a)
	}
	return uint32(a)
}
