package universe_service

import (
	"errors"
	"log"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/dao"
	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/internal/world/service/procedural_generation"
	world_service "github.com/Roukii/pock_multiplayer/internal/world/service/world"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

// TODO maybe split universe service into world service and chunk service
type UniverseService struct {
	Universe      universe.Universe
	chunkDao      *dao.ChunkDao
	worldDao      *dao.WorldDao
	WorldServices map[string]*world_service.WorldService
}

func NewUniverseService(session *gocqlx.Session) *UniverseService {
	u := UniverseService{
		Universe: universe.Universe{
			Worlds: make(map[string]universe.World),
		},
		chunkDao:      dao.NewChunkDao(session),
		worldDao:      dao.NewWorldDao(session),
		WorldServices: make(map[string]*world_service.WorldService),
	}
	return &u
}

// TODO maybe remove univer.World
func (us *UniverseService) LoadWorlds() error {
	worlds, err := us.worldDao.GetAllWorlds()
	if err != nil {
		return err
	}
	for _, world := range worlds {
		us.Universe.Worlds[world.UUID] = world
		us.WorldServices[world.UUID] = world_service.NewWorldService(&world, us.chunkDao, false)
	}
	if len(us.Universe.Worlds) == 0 {
		worldService, err := us.createWorld("tutorial")
		if err != nil {
			log.Println("error : ", err)
			return err
		}
		log.Println("create world with uuid :", worldService.World.UUID)
		us.Universe.Worlds[worldService.World.UUID] = *worldService.World
		us.WorldServices[worldService.World.UUID] = worldService
	}

	return nil
}

func (us *UniverseService) createWorld(worldName string) (*world_service.WorldService, error) {
	world := universe.World{
		UUID:        gocql.TimeUUID().String(),
		Name:        worldName,
		MaxPlayer:   20,
		Level:       1,
		Length:      10,
		Width:       10,
		SpawnPoints: []player.SpawnPoint{},
		Seed:        procedural_generation.GenerateSeed(),
		Type:        0,
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
	}
	worldService := world_service.NewWorldService(&world, us.chunkDao, true)
	err := us.worldDao.Insert(&world)
	if err != nil {
		log.Println("error : ", err)
		return worldService, err
	}
	log.Println("Chunks : ", len(world.Chunks))
	return worldService, err
}

func (us *UniverseService) GetWorldService(WorldUUID string) (*world_service.WorldService, error) {
	if world, ok := us.WorldServices[WorldUUID]; ok {
		return world, nil
	}
	return nil, errors.New("Can't find world")
}

func (us *UniverseService) GetWorlds() []*universe.World {
	worlds := make([]*universe.World, 0, len(us.Universe.Worlds))

	for _, value := range us.Universe.Worlds {
		worlds = append(worlds, &value)
	}
	return worlds
}

func (us *UniverseService) AddPlayerToWorld(player *player.Player, world *world_service.WorldService) {
	world.DynamicEntityService.AddDynamicEntity(player)
	test := entity.ICreature{}
	world.DynamicEntityService.AddDynamicEntity(test)
}

func (us *UniverseService) MovePlayerToAnotherWorld() {

}