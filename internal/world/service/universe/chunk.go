package universe_service

import (
	"errors"
	"log"
	"math"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/internal/world/service/procedural_generation"
	"github.com/gocql/gocql"
)

func (us *UniverseService) loadWorldChunks(world *universe.World) error {
	chunks, err := us.chunkDao.LoadWorldChunk(world.UUID)
	if err != nil {
		log.Println("Couldn't load chunks : ", err)
		return err
	}
	log.Println("chunks length : ", len(chunks))
	if len(chunks) == 0 {
		return us.generateAndSaveWorldChunks(world)
	}
	world.Chunks = make(map[int]map[int]universe.Chunk)
	for _, chunk := range chunks {
		if world.Chunks[chunk.PositionX] == nil {
			world.Chunks[chunk.PositionX] = make(map[int]universe.Chunk)
		}
		world.Chunks[chunk.PositionX][chunk.PositionY] = *chunk
	}
	return err
}

func (us *UniverseService) generateAndSaveWorldChunks(world *universe.World) (err error) {
	generator := procedural_generation.NewWorldGenerator(world)
	us.WorldGenerators[world.UUID] = &generator

	log.Println("Start generate and save chunks")
	world.Chunks = make(map[int]map[int]universe.Chunk)
	for x := 0; x < world.Length; x++ {
		world.Chunks[x] = make(map[int]universe.Chunk)
		for y := 0; y < world.Width; y++ {
			chunk, err := us.generateChunk(world, entity.Vector2{x, y})
			world.Chunks[x][y] = *chunk
			if err != nil {
				log.Println("Error generating chunk : ", x, "/", y, " with error : ", err)
				return err
			}
		}
	}
	log.Println("Save chunks to database")
	err = us.saveWorldChunks(world)
	if err != nil {
		log.Println("Error saving chunks : ", err)
		return err
	}
	return nil
}

func (us *UniverseService) getChunksFromSpawnSpoint(spawnPoint player.SpawnPoint, viewDistance int) ([]*universe.Chunk, error) {
	var chunks []*universe.Chunk
	world, ok := us.Universe.Worlds[spawnPoint.WorldUUID]
	if !ok {
		return nil, errors.New("world not found : " + spawnPoint.WorldUUID)
	}
	if len(world.Chunks) == 0 {
		err := us.loadWorldChunks(&world)
		if err != nil {
			return nil, errors.New("world not found : " + spawnPoint.WorldUUID)
		}
	}
	spawnChunkPosX := int(math.Floor(float64(spawnPoint.Coordinate.Position.X / 25)))
	spawnChunkPosY := int(math.Floor(float64(spawnPoint.Coordinate.Position.Y / 25)))
	for x := -viewDistance; x <= viewDistance; x++ {
		for y := -viewDistance; y <= viewDistance; y++ {
			if chunk, ok := world.Chunks[spawnChunkPosX+x][spawnChunkPosY+y]; ok {
				chunks = append(chunks, &chunk)
			} else {
				log.Println("Couldn't load chunk from pos : ", spawnChunkPosX+x, "/", spawnChunkPosY+y)
			}
		}
	}
	return chunks, nil
}

func (us *UniverseService) saveWorldChunks(world *universe.World) (err error) {
	for _, chunks := range world.Chunks {
		for _, chunk := range chunks {
			err = us.chunkDao.Insert(world.UUID, &chunk)
			if err != nil {
				log.Println("can't save chunk : ", err)
			}
		}
	}
	return
}

func (us *UniverseService) generateChunk(world *universe.World, position entity.Vector2) (*universe.Chunk, error) {
	generator := us.WorldGenerators[world.UUID]
	chunk, err := generator.GenerateChunk(position.X, position.Y)
	if err != nil {
		log.Println("error generating chunk : ", err)
		return nil, err
	}
	chunk.UUID = gocql.TimeUUID().String()
	us.chunkDao.Insert(world.UUID, chunk)
	return chunk, nil
}
