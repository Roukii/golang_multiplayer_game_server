package procedural_generation

import (
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
)

type (
	entityGeneration struct {
		World *universe.World
	}

	EntityGeneration interface {
		UpdateChunk(x int, y int) (universe.Chunk, error)
	}
)

func NewEntityGeneration(world *universe.World) EntityGeneration {
	tmp := &entityGeneration{
		World: world,
	}
	return tmp
}

func (eg *entityGeneration) UpdateChunk(x int, y int) (chunk universe.Chunk, err error) {
	return chunk, err
}
