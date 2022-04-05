package game

import (
	"errors"
	"fmt"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/internal/world/service/procedural_generation"
	"github.com/gocql/gocql"
)

func (g *GameService) CreateWorld(worldName string) (universe.World, error) {
	world := universe.World{
		UUID:      gocql.TimeUUID().String(),
		Name:      worldName,
		MaxPlayer: 20,
		Level:     1,
		Length:    10,
		Width:     10,
		Seed:      procedural_generation.GenerateSeed(),
		Type:      0,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}
	err := g.WorldDao.Insert(&world)
	if err != nil {
		fmt.Println("error : ", err)
		return world, err
	}
	err = g.generateAndSaveWorldChunks(&world)
	fmt.Println("Chunks : ", len(world.Chunks))
	return world, err
}

func (g *GameService) GetWorld(WorldUUID string) (*universe.World, error) {
	if world, ok := g.Universe.Worlds[WorldUUID]; ok {
		return &world, nil
	}
	return nil, errors.New("Can't find world")
}
