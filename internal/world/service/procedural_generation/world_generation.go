package procedural_generation

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/pkg/advmath"
	"github.com/Roukii/pock_multiplayer/pkg/random"
	opensimplex "github.com/ojrac/opensimplex-go"
	// "github.com/pzsz/voronoi"
)

type (
	WorldGenerator struct {
		World             *universe.World
		ElevationNoiseMap opensimplex.Noise
		RainfallNoiseMap  opensimplex.Noise
		FallOffMap        []float64
	}
)

// TODO move scale somewhere else
const (
	seedSize     int     = 16
	chunkLength  int     = 25
	noiseScale   float64 = 0.15
	fallOffStart float64 = 0.94
	fallOffEnd   float64 = 1
)

func NewWorldGenerator(world *universe.World) WorldGenerator {
	tmp := WorldGenerator{
		World: world,
	}
	world.Persistance = 0.5
	world.Lacunarity = 2
	world.Octave = 4
	world.UseFallOff = true
	tmp.generateFallOffMap()
	tmp.generateNoiseMap()
	return tmp
}

func GenerateSeed() string {
	return random.RandStringRunes(seedSize)
}

func (wg *WorldGenerator) GenerateChunk(positionX int, positionY int) (chunk *universe.Chunk, err error) {
	elevationHeightMap := wg.generateHeightmap(positionX*(chunkLength-1), positionY*(chunkLength-1), wg.ElevationNoiseMap)
	//rainfallHeightMap := wg.generateHeightmap(positionX*chunkLength, positionY*chunkLength, wg.RainfallNoiseMap)

	chunk = &universe.Chunk{
		PositionX: positionX,
		PositionY: positionY,
	}
	i := 0
	fmt.Println("chunk pos : ", positionX, "/", positionY)
	mapWidth := chunkLength * wg.World.Length
	offsetPositionX := positionX * chunkLength
	offsetPositionY := positionY * (mapWidth * chunkLength)
	for y := 0; y < chunkLength; y++ {
		for x := 0; x < chunkLength; x++ {
			elevation := elevationHeightMap[i]
			if wg.World.UseFallOff {
				elevation = advmath.CircIn(advmath.ClampFloat64(elevation-wg.FallOffMap[(offsetPositionY+(y*mapWidth))+offsetPositionX+x], 0, 1)) * float64(wg.World.ScaleHeight)
			}
			tile := universe.Tile{
				// TileType:  getTileType(elevationHeightMap[i], rainfallHeightMap[i]),
				TileType:  universe.Dirt,
				Elevation: elevation,
			}
			chunk.Tiles = append(chunk.Tiles, tile)
			i++
			fmt.Print(fmt.Sprintf("%.2f", elevation), " ")
		}
		fmt.Print("\n")
	}
	fmt.Print("\n\n\n")
	return chunk, err
}

func (wg *WorldGenerator) generateNoiseMap() {
	h := md5.New()
	var elevationSeed uint64 = binary.BigEndian.Uint64(h.Sum([]byte(wg.World.Seed)))
	var temperatureSeed uint64 = binary.BigEndian.Uint64(h.Sum([]byte(reverseString(wg.World.Seed))))
	wg.ElevationNoiseMap = opensimplex.New(int64(elevationSeed))
	wg.RainfallNoiseMap = opensimplex.New(int64(temperatureSeed))
}

func reverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// TODO make a cleaner implementation of this shit
func getTileType(elevation float64, rainfall float64) universe.TileType {
	if elevation <= -0.2 {
		return universe.Water
	} else if elevation <= -0.24 {
		return universe.Sand
	} else if elevation >= 0.5 {
		return universe.Snow
	}
	if rainfall <= -0.4 {
		if elevation <= 0 {
			return universe.Sand
		} else {
			return universe.Rock
		}
	} else if rainfall <= 0 {
		if elevation <= 0.2 {
			return universe.Grass
		} else {
			return universe.Rock
		}
	} else {
		return universe.Forest
	}
}

func (wg *WorldGenerator) generateHeightmap(startingPositionX int, startingPositionY int, noise opensimplex.Noise) []float64 {
	fmt.Println("starting pos : ", startingPositionX, "/", startingPositionY)
	w, h := startingPositionX+chunkLength, startingPositionY+chunkLength
	heightmap := make([]float64, chunkLength*chunkLength)
	minNoiseHeight, maxNoiseHeight := math.MaxFloat64, math.SmallestNonzeroFloat64
	for y := 0; y+startingPositionY < h; y++ {
		for x := 0; x+startingPositionX < w; x++ {
			amplitude := 1.0
			frequency := 1.0
			noiseHeight := 0.0
			var noiseValue float64
			for i := 0; i < int(wg.World.Octave); i++ {
				xFloat := float64(x+startingPositionX) * noiseScale * frequency
				yFloat := float64(y+startingPositionY) * noiseScale * frequency
				noiseValue = noise.Eval2(xFloat, yFloat)
				noiseHeight += noiseValue * amplitude

				amplitude *= wg.World.Persistance
				frequency *= wg.World.Lacunarity
			}
			if noiseHeight > maxNoiseHeight {
				maxNoiseHeight = noiseHeight
			} else if noiseHeight < minNoiseHeight {
				minNoiseHeight = noiseHeight
			}
			heightmap[(y*chunkLength)+x] = noiseHeight
		}
		// for x := 0; x < chunkLength; x++ {
		// 	// InverseLerp
		// 	heightmap[(y*chunkLength)+x] = advmath.InverseLerpFloat64(minNoiseHeight, maxNoiseHeight, heightmap[(y*chunkLength)+x])
		// }
	}
	return heightmap
}

// Taken from https://youtu.be/XpG3YqUkCTY?t=92
func (wg *WorldGenerator) generateFallOffMap() {
	mapWidth := chunkLength * wg.World.Length
	wg.FallOffMap = make([]float64, mapWidth*mapWidth)
	for i := 0; i < mapWidth; i++ {
		for j := 0; j < mapWidth; j++ {
			x := float64(i)/float64(mapWidth)*2 - 1
			y := float64(j)/float64(mapWidth)*2 - 1

			t := math.Max(math.Abs(float64(x)), math.Abs(float64(y)))
			if t < fallOffStart {
				wg.FallOffMap[i*mapWidth+j] = 0
			} else if t > fallOffEnd {
				wg.FallOffMap[i*mapWidth+j] = 1
			} else {
				wg.FallOffMap[i*mapWidth+j] = advmath.Smoothstep(0, 1, advmath.InverseLerpFloat64(fallOffStart, fallOffEnd, t))
			}
		}
	}
}

// TODO upgrade world creation with voronoi diagram for smoother biome generation
// func GenerateVoronoiDiagram(w float64, h float64, r int) {
// 	bbox := voronoi.NewBBox(0, w, 0, h)
// 	sites := []voronoi.Vertex{}

// 	// Compute voronoi diagram.
// 	d := voronoi.ComputeDiagram(sites, bbox, true)

// 	// Max number of iterations is 16
// 	if r > 16 {
// 		r = 16
// 	}

// 	// Relax using Lloyd's algorithm
// 	for i := 0; i < r; i++ {
// 		sites = utils.LloydRelaxation(d.Cells)
// 		d = voronoi.ComputeDiagram(sites, bbox, true)
// 	}

// 	center := voronoi.Vertex{float64(w / 2), float64(h / 2)}

// 	return &voronoi.Diagram{d, center}
// }
