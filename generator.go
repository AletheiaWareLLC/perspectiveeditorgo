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
	"math/rand"
	"strconv"
	"time"
)

func Generate(puzzle *perspectivego.Puzzle, size uint32, goalCount int, goalMesh, goalColour string, sphereCount int, sphereMesh, sphereColour string, blockCount int, blockMesh, blockColour string, portalCount int, portalMesh, portalColour string) *perspectivego.Puzzle {
	rand.Seed(time.Now().UnixNano())

	occupied := make([]*perspectivego.Location, 0, goalCount+blockCount+sphereCount+portalCount)
	if goalCount > 0 {
		puzzle.Goal = make([]*perspectivego.Goal, 0, goalCount)
		for i := 0; i < goalCount; i++ {
			location := GenerateLocation(occupied, size)
			occupied = append(occupied, location)
			goal := &perspectivego.Goal{
				Name:     "g" + strconv.Itoa(i),
				Mesh:     goalMesh,
				Colour:   goalColour,
				Location: location,
			}
			puzzle.Goal = append(puzzle.Goal, goal)
		}
	}
	if blockCount > 0 {
		puzzle.Block = make([]*perspectivego.Block, 0, blockCount)
		for i := 0; i < blockCount; i++ {
			location := GenerateLocation(occupied, size)
			occupied = append(occupied, location)
			block := &perspectivego.Block{
				Name:     "b" + strconv.Itoa(i),
				Mesh:     blockMesh,
				Colour:   blockColour,
				Location: location,
			}
			puzzle.Block = append(puzzle.Block, block)
		}
	}
	if sphereCount > 0 {
		puzzle.Sphere = make([]*perspectivego.Sphere, 0, sphereCount)
		for i := 0; i < sphereCount; i++ {
			location := GenerateLocation(occupied, size)
			occupied = append(occupied, location)
			sphere := &perspectivego.Sphere{
				Name:     "s" + strconv.Itoa(i),
				Mesh:     sphereMesh,
				Colour:   sphereColour,
				Location: location,
			}
			puzzle.Sphere = append(puzzle.Sphere, sphere)
		}
	}
	if portalCount > 0 {
		puzzle.Portal = make([]*perspectivego.Portal, 0, portalCount)
		var previous *perspectivego.Portal
		for i := 0; i < portalCount; i++ {
			location := GenerateLocation(occupied, size)
			occupied = append(occupied, location)
			portal := &perspectivego.Portal{
				Name:     "p" + strconv.Itoa(i),
				Mesh:     portalMesh,
				Colour:   portalColour,
				Location: location,
			}
			if previous == nil {
				previous = portal
			} else {
				portal.Link = previous.Location
				previous.Link = portal.Location
				previous = nil
			}
			puzzle.Portal = append(puzzle.Portal, portal)
		}
	}
	return puzzle
}

func GenerateLocation(occupied []*perspectivego.Location, size uint32) *perspectivego.Location {
	location := &perspectivego.Location{}
	for {
		location.X = int32(RandomLocation(size))
		location.Y = int32(RandomLocation(size))
		location.Z = int32(RandomLocation(size))
		if !IsOccupied(occupied, location) {
			return location
		}
	}
}

func IsOccupied(occupied []*perspectivego.Location, location *perspectivego.Location) bool {
	for _, l := range occupied {
		if location.X == l.X && location.Y == l.Y && location.Z == l.Z {
			return true
		}
	}
	return false
}

func RandomLocation(size uint32) int {
	return rand.Intn(int(size)) - int(size/2)
}
