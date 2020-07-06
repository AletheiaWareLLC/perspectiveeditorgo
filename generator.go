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

func Generate(puzzle *perspectivego.Puzzle, size uint32,
	goalCount int, goalMesh, goalColour, goalTexture, goalMaterial []string, goalShader string,
	sphereCount int, sphereMesh, sphereColour, sphereTexture, sphereMaterial []string, sphereShader string,
	blockCount int, blockMesh, blockColour, blockTexture, blockMaterial []string, blockShader string,
	portalCount int, portalMesh, portalColour, portalTexture, portalMaterial []string, portalShader string) *perspectivego.Puzzle {
	rand.Seed(time.Now().UnixNano())

	occupied := make(map[string]bool, goalCount+blockCount+sphereCount+portalCount)
	if goalCount > 0 {
		puzzle.Goal = make([]*perspectivego.Goal, 0, goalCount)
		for i := 0; i < goalCount; i++ {
			location := GenerateLocation(occupied, size)
			goal := &perspectivego.Goal{
				Name:     "g" + strconv.Itoa(i),
				Mesh:     goalMesh[i%len(goalMesh)],
				Colour:   goalColour[i%len(goalColour)],
				Location: location,
				Texture:  goalTexture[i%len(goalTexture)],
				Material: goalMaterial[i%len(goalMaterial)],
				Shader:   goalShader,
			}
			puzzle.Goal = append(puzzle.Goal, goal)
		}
	}
	if blockCount > 0 {
		puzzle.Block = make([]*perspectivego.Block, 0, blockCount)
		for i := 0; i < blockCount; i++ {
			location := GenerateLocation(occupied, size)
			block := &perspectivego.Block{
				Name:     "b" + strconv.Itoa(i),
				Mesh:     blockMesh[i%len(blockMesh)],
				Colour:   blockColour[i%len(blockColour)],
				Location: location,
				Texture:  blockTexture[i%len(blockTexture)],
				Material: blockMaterial[i%len(blockMaterial)],
				Shader:   blockShader,
			}
			puzzle.Block = append(puzzle.Block, block)
		}
	}
	if sphereCount > 0 {
		puzzle.Sphere = make([]*perspectivego.Sphere, 0, sphereCount)
		for i := 0; i < sphereCount; i++ {
			location := GenerateLocation(occupied, size)
			sphere := &perspectivego.Sphere{
				Name:     "s" + strconv.Itoa(i),
				Mesh:     sphereMesh[i%len(sphereMesh)],
				Colour:   sphereColour[i%len(sphereColour)],
				Location: location,
				Texture:  sphereTexture[i%len(sphereTexture)],
				Material: sphereMaterial[i%len(sphereMaterial)],
				Shader:   sphereShader,
			}
			puzzle.Sphere = append(puzzle.Sphere, sphere)
		}
	}
	if portalCount > 0 {
		puzzle.Portal = make([]*perspectivego.Portal, 0, portalCount)
		var previous *perspectivego.Portal
		for i := 0; i < portalCount; i++ {
			location := GenerateLocation(occupied, size)
			portal := &perspectivego.Portal{
				Name:     "p" + strconv.Itoa(i),
				Mesh:     portalMesh[i%len(portalMesh)],
				Colour:   portalColour[(i/2)%len(portalColour)],
				Location: location,
				Texture:  portalTexture[(i/2)%len(portalTexture)],
				Material: portalMaterial[(i/2)%len(portalMaterial)],
				Shader:   portalShader,
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

func GenerateLocation(occupied map[string]bool, size uint32) *perspectivego.Location {
	location := &perspectivego.Location{}
	var key string
	for {
		location.X = int32(RandomLocation(size))
		location.Y = int32(RandomLocation(size))
		location.Z = int32(RandomLocation(size))
		key = location.String()
		if !occupied[key] {
			occupied[key] = true
			return location
		}
	}
}

func RandomLocation(size uint32) int {
	return rand.Intn(int(size)) - int(size/2)
}
