package game

import (
	"fmt"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/gocql/gocql"
)

func (g *GameService) GenerateChunk(world *universe.World, position entity.Vector2) (*universe.Chunk, error) {
	generator := g.WorldGenerators[world.UUID]
	chunk, err := generator.GenerateChunk(position.X, position.Y)
	if err != nil {
		fmt.Println("error generating chunk : ", err)
		return nil, err
	}
	chunk.UUID = gocql.TimeUUID().String()
	g.ChunkDao.Insert(world.UUID, &chunk) 
	return &chunk, nil
}
