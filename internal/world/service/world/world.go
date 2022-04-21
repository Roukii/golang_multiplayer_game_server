package world_service

import (
	"log"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/dao"
	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/internal/world/service/dynamic_entity"
	"github.com/Roukii/pock_multiplayer/internal/world/service/procedural_generation"
)

type WorldService struct {
	World                *universe.World
	Generator            *procedural_generation.WorldGenerator
	DynamicEntityService *dynamic_entity.DynamicEntityService
	chunkDao             *dao.ChunkDao
	PlayerCount          int
}

func NewWorldService(world *universe.World, chunkDao *dao.ChunkDao, needGeneration bool) *WorldService {
	generator := procedural_generation.NewWorldGenerator(world)
	dynamicEntityService := dynamic_entity.NewDynamicEntityService()
	worldService := WorldService{
		World:                world,
		Generator:            &generator,
		DynamicEntityService: dynamicEntityService,
		chunkDao:             chunkDao,
		PlayerCount:          0,
	}
	if needGeneration {
		worldService.generateSpawnPoints()
		err := worldService.generateAndSaveWorldChunks()
		if err != nil {
			log.Println("Coudln't generate chunks")
			return nil
		}
	}
	return &worldService
}

// TODO upgrade spawn point generation
func (ws *WorldService) generateSpawnPoints() {
	var spawnPoints []player.SpawnPoint
	spawnPoints = append(spawnPoints, player.SpawnPoint{
		WorldUUID: ws.World.UUID,
		Coordinate: entity.Position{
			Position: entity.Vector3f{
				X: 1,
				Y: 1,
				Z: 1,
			},
			Rotation: entity.Vector3f{
				X: 1,
				Y: 1,
				Z: 1,
			},
		},
		UpdatedAt: time.Now(),
	})
	ws.World.SpawnPoints = spawnPoints
}
