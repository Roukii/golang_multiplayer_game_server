package game

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/Roukii/pock_multiplayer/internal/world/service/procedural_generation"
	"github.com/gocql/gocql"
)

func (g *GameService) CreateWorld(worldName string) (*universe.World, error) {
	world := &universe.World{
		UUID:      gocql.TimeUUID().String(),
		Name:      worldName,
		MaxPlayer: 20,
		Level:     1,
		Length:    10,
		Width:     10,
		Seed:      procedural_generation.GenerateSeed(),
		Type:      0,
		Chunks:    make(map[int]map[int]universe.Chunk),
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}
	err := g.WorldDao.Insert(world)
	if err != nil {
		fmt.Println("error : ", err)
		return nil, err
	}
	return world, nil
}

func (g *GameService) GetWorld(WorldUUID string) (*universe.World, error) {
	for _, world := range g.Universe.Worlds {
		if world.UUID == WorldUUID {
			return &world, nil
		}
	}
	return nil, errors.New("Can't find world")
}

// TODO Found another way to do this because right now it sucks 
func (g *GameService) LoadChunksFromSpawnSpoint(spawnPoint player.SpawnPoint, viewDistance int) ([]*universe.Chunk, error) {
	var chunks []*universe.Chunk
	world, ok := g.Universe.Worlds[spawnPoint.WorldUUID]
	if !ok {
		return nil, errors.New("world not found : " + spawnPoint.WorldUUID)
	}
	chunkPosX := int(math.Floor(float64(spawnPoint.Coordinate.Position.X / 25)))
	chunkPosY := int(math.Floor(float64(spawnPoint.Coordinate.Position.Y / 25)))
	var missingChunks []entity.Vector2
	for x := -viewDistance; x <= viewDistance; x++ {
		for y := -viewDistance; y <= viewDistance; y++ {
			if chunk, ok := world.Chunks[x][y]; ok {
				fmt.Println("add chunk from ram : ", chunk.PositionX,"/", chunk.PositionY)
				chunks = append(chunks, &chunk)
			} else {
				missingChunks = append(missingChunks, entity.Vector2{
					X: chunkPosX + x,
					Y: chunkPosY + y,
				})
			}
		}
	}
	fmt.Println("missing chunks : ", missingChunks)
	if len(missingChunks) == 0 {
		return chunks, nil
	}
	dbChunks, err := g.ChunkDao.LoadChunckBetweenCoordinate(spawnPoint.WorldUUID, chunkPosX-viewDistance, chunkPosX+viewDistance, chunkPosY-viewDistance, chunkPosY+viewDistance)
	
	for _, dbChunk := range dbChunks {
		needToAddChunkToList := false
		for _, chunk := range chunks {
			if dbChunk.UUID == chunk.UUID {
				fmt.Println("already in list : ", dbChunk.PositionX,"/", dbChunk.PositionY)
				needToAddChunkToList = true
			}
		}
		if needToAddChunkToList {
			fmt.Println("need to add to list : ", dbChunk.PositionX,"/", dbChunk.PositionY)
			chunks = append(chunks, dbChunk)
			var missingChunksTmp []entity.Vector2
			for index, missing := range missingChunks {
				if dbChunk.PositionX == missing.X && dbChunk.PositionY == missing.Y {
					missingChunksTmp = append(missingChunksTmp, missingChunks[index])
				}
			}
			missingChunks = missingChunksTmp
		}
	}
	fmt.Println("missing chunks len : ", len(missingChunks))
	fmt.Println("missing chunks : ", missingChunks)
	for _, missingPos := range missingChunks {
		chunk , err := g.WorldGenerators[spawnPoint.WorldUUID].GenerateChunk(missingPos.X, missingPos.Y)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, &chunk)
	}
	if err != nil {
		fmt.Println("error loading chunk in db : ", err)
		return chunks, err
	}
	fmt.Println("len : ", len(chunks))
	go g.SaveChunks(spawnPoint.WorldUUID, chunks)
	return chunks, nil
}

func (g *GameService) SaveChunks(worldUUID string ,chunks []*universe.Chunk) {
	for _, chunk := range chunks {
		err := g.ChunkDao.Insert(worldUUID,chunk)
		if err != nil {
			fmt.Println("can't save chunk : ", err)
		}
		fmt.Println("insert to db chunk : ", chunk.PositionX,"/", chunk.PositionY)

	}
}