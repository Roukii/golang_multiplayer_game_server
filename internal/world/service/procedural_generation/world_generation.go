package procedural_generation

import (
	"crypto/md5"
	"encoding/binary"

	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/pkg/random"
	opensimplex "github.com/ojrac/opensimplex-go"
	// "github.com/pzsz/voronoi"
)

type (
	WorldGenerator struct {
		World             *universe.World
		ElevationNoiseMap opensimplex.Noise
		RainfallNoiseMap  opensimplex.Noise
	}
)

const (
	seedSize    int = 16
	chunkLength int = 25
)

func NewWorldGenerator(world *universe.World) WorldGenerator {
	tmp := WorldGenerator{
		World: world,
	}
	tmp.GenerateNoiseMap()
	return tmp
}

func GenerateSeed() string {
	return random.RandStringRunes(seedSize)
}

func (wg *WorldGenerator) GenerateChunk(positionX int, positionY int) (chunk *universe.Chunk, err error) {
	elevationHeightMap := generateHeightmap(positionX*chunkLength, positionY*chunkLength, wg.ElevationNoiseMap)
	rainfallHeightMap := generateHeightmap(positionX*chunkLength, positionY*chunkLength, wg.RainfallNoiseMap)

	chunk = &universe.Chunk{
		PositionX:       positionX,
		PositionY:       positionY,
	}
	tileNumber := len(elevationHeightMap)
	for i := 0; i < tileNumber; i++ {
		tile := universe.Tile{
			TileType:  getTileType(elevationHeightMap[i], rainfallHeightMap[i]),
			Elevation: elevationHeightMap[i],
		}
		chunk.Tiles = append(chunk.Tiles, tile)
	}
	return chunk, err
}

func (wg *WorldGenerator) GenerateNoiseMap() {
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

func generateHeightmap(startingPositionX int, startingPositionY int, noise opensimplex.Noise) []float64 {
	w, h := startingPositionX+chunkLength, startingPositionY+chunkLength
	heightmap := make([]float64, chunkLength*chunkLength)
	for y := 0; y+startingPositionY < h; y++ {
		for x := 0; x+startingPositionX < w; x++ {
			xFloat := float64(x+startingPositionX) * 0.15
			yFloat := float64(y+startingPositionY) * 0.15
			heightmap[(y*chunkLength)+x] = noise.Eval2(xFloat, yFloat)
		}
	}
	return heightmap
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
