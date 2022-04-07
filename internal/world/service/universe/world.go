package universe_service

import (
	"errors"
	"log"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/internal/world/service/procedural_generation"
	"github.com/gocql/gocql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
)

func (us *UniverseService) LoadWorlds() error {
	worlds, err := us.worldDao.GetAllWorlds()
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
			log.Println("error : ", err)
			return err
		}
		log.Println("create world with uuid :", world.UUID)
		us.Universe.Worlds[world.UUID] = world
	}

	return nil
}

func (us *UniverseService) CreateWorld(worldName string) (universe.World, error) {
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
	err := us.worldDao.Insert(&world)
	if err != nil {
		log.Println("error : ", err)
		return world, err
	}
	err = us.generateAndSaveWorldChunks(&world)
	us.generateSpawnPoints(&world)
	log.Println("Chunks : ", len(world.Chunks))
	return world, err
}

// TODO upgrade spawn point generation
func (us *UniverseService) generateSpawnPoints(world *universe.World) {
	var spawnPoints []player.SpawnPoint
	spawnPoints = append(spawnPoints, player.SpawnPoint{
		WorldUUID:  world.UUID,
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
		UpdatedAt:  time.Now(),
	})
	world.SpawnPoints = spawnPoints
}

func (us *UniverseService) GetWorld(WorldUUID string) (*universe.World, error) {
	if world, ok := us.Universe.Worlds[WorldUUID]; ok {
		return &world, nil
	}
	return nil, errors.New("Can't find world")
}

func (us *UniverseService) GetWorlds() []*universe.World {
	worlds := make([]*universe.World, 0, len(us.Universe.Worlds))

	for  _, value := range us.Universe.Worlds {
		 worlds = append(worlds, &value)
	}
	return worlds
}

// TODO lock write with mutex
func (us *UniverseService) LoadChunksFromSpawnPoint(spawnPoint player.SpawnPoint) (chunks []*universe.Chunk, err error) {
	chunks, err = us.getChunksFromSpawnSpoint(spawnPoint, 1)
	if err != nil {
		log.Println("failed to load chunks", err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to load chunks")
	}
	return chunks, nil
}

// TODO lock write with mutex
func (us *UniverseService) LoadSpecificChunks(world *universe.World, coordinates []*pb.Vector2Int) ([]*universe.Chunk, error) {
	var chunks []*universe.Chunk
	if len(world.Chunks) == 0 {
		err := us.loadWorldChunks(world)
		if err != nil {
			return nil, errors.New("chunks not found for world : " + world.UUID)
		}
	}
	for _, coordinate := range coordinates {
		if chunk, ok := world.Chunks[int(coordinate.X)][int(coordinate.Y)]; ok {
			chunks = append(chunks, &chunk)
		} else {
			log.Println("Couldn't load chunk from pos : ", coordinate.X, "/", coordinate.Y)
		}
	}
	return chunks, nil
}
