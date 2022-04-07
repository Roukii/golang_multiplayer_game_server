package universe_service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/internal/world/service/procedural_generation"
	"github.com/gocql/gocql"
)

func (us *UniverseService) LoadWorlds() error {
	worlds, err := us.WorldDao.GetAllWorlds()
	if err != nil {
		return err
	}
	for _, world := range worlds {
		us.Universe.Worlds[world.UUID] = world
		generator := procedural_generation.NewWorldGenerator(&world)
		us.WorldGenerators[world.UUID] = &generator
	}
	if len(us.Universe.Worlds) == 0 {
		world, err := us.CreateWorld("tutorial")
		if err != nil {
			fmt.Println("error : ", err)
			return err
		}
		fmt.Println("create world with uuid :", world.UUID)
		us.Universe.Worlds[world.UUID] = world
	}

	return nil
}

func (us *UniverseService) CreateWorld(worldName string) (universe.World, error) {
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
	err := us.WorldDao.Insert(&world)
	if err != nil {
		fmt.Println("error : ", err)
		return world, err
	}
	err = us.generateAndSaveWorldChunks(&world)
	fmt.Println("Chunks : ", len(world.Chunks))
	return world, err
}

func (us *UniverseService) GetWorld(WorldUUID string) (*universe.World, error) {
	if world, ok := us.Universe.Worlds[WorldUUID]; ok {
		return &world, nil
	}
	return nil, errors.New("Can't find world")
}
