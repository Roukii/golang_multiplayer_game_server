package procedural_generation

import (
	"crypto/md5"
	"encoding/binary"

	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/pkg/random"
	opensimplex "github.com/ojrac/opensimplex-go"
	"github.com/pzsz/voronoi"
)

type (
	worldGenerator struct {
		World    *universe.World
		NoiseMap opensimplex.Noise
	}

	WorldGenerator interface {
		GenerateChunk(x int, y int) (universe.Chunk, error)
		GenerateWorld() error
	}
)

const (
	seedSize  int = 16
	chunkSize int = 25
)

func NewWorldGenerator(world *universe.World) WorldGenerator {
	return &worldGenerator{
		World: world,
	}
}

func generateSeed() string {
	return random.RandStringRunes(seedSize)
}

func (wg *worldGenerator) GenerateChunk(positionX int, positionY int) (chunk universe.Chunk, err error) {
	return chunk, err
}

func (wg *worldGenerator) GenerateChunkHeightmap(startingPositionX int, startingPositionY int) []float64 {
	w, h := startingPositionX+chunkSize, startingPositionY+chunkSize
	heightmap := make([]float64, w*h)
	for y := startingPositionY; y < h; y++ {
		for x := startingPositionX; x < w; x++ {
			xFloat := float64(x) / float64(w)
			yFloat := float64(y) / float64(h)
			heightmap[(y*w)+x] = wg.NoiseMap.Eval2(xFloat, yFloat)
		}
	}
	return heightmap
}

func (wg *worldGenerator) GenerateWorld() error {
	if wg.World.Seed == "" {
		wg.World.Seed = generateSeed()
	}
	h := md5.New()
	var seed uint64 = binary.BigEndian.Uint64(h.Sum(nil))
	wg.NoiseMap = opensimplex.New(int64(seed))
	
	return nil
}

func useVoronoi() {
	// Sites of voronoi diagram
sites := []voronoi.Vertex{
voronoi.Vertex{4, 5},
voronoi.Vertex{6, 5},
}

// Create bounding box of [0, 20] in X axis
// and [0, 10] in Y axis
bbox := NewBBox(0, 20, 0, 10)

// Compute diagram and close cells (add half edges from bounding box)
diagram := NewVoronoi().Compute(sites, bbox, true)

// Iterate over cells
for _, cell := diagram.Cells {
for _, hedge := cell.Halfedges {
	 ...
}	
}

// Iterate over all edges
for _, edge := diagram.Edge {
 ...
}
}