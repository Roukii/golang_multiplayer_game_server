package world_service

import (
	"errors"
	"log"
	"math"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/gocql/gocql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ws *WorldService) loadWorldChunks() error {
	chunks, err := ws.chunkDao.LoadWorldChunk(ws.World.UUID)
	if err != nil {
		log.Println("Couldn't load chunks : ", err)
		return err
	}
	log.Println("chunks length : ", len(chunks))
	if len(chunks) == 0 {
		return ws.generateAndSaveWorldChunks()
	}
	ws.World.Chunks = make(map[int]map[int]universe.Chunk)
	for _, chunk := range chunks {
		if ws.World.Chunks[chunk.PositionX] == nil {
			ws.World.Chunks[chunk.PositionX] = make(map[int]universe.Chunk)
		}
		ws.World.Chunks[chunk.PositionX][chunk.PositionY] = *chunk
	}
	return err
}

func (ws *WorldService) generateAndSaveWorldChunks() (err error) {
	world := ws.World
	log.Println("Start generate and save chunks")
	world.Chunks = make(map[int]map[int]universe.Chunk)
	for x := 0; x < world.Length; x++ {
		world.Chunks[x] = make(map[int]universe.Chunk)
		for y := 0; y < world.Width; y++ {
			chunk, err := ws.generateChunk(entity.Vector2{x, y})
			world.Chunks[x][y] = *chunk
			if err != nil {
				log.Println("Error generating chunk : ", x, "/", y, " with error : ", err)
				return err
			}
		}
	}
	log.Println("Save chunks to database")
	err = ws.saveWorldChunks(world)
	if err != nil {
		log.Println("Error saving chunks : ", err)
		return err
	}
	return nil
}


func (ws *WorldService) getAllChunks() ([]*universe.Chunk, error) {
	var chunks []*universe.Chunk
	world := ws.World
	if len(world.Chunks) == 0 {
		err := ws.loadWorldChunks()
		if err != nil {
			return nil, errors.New("world not found : " + ws.World.UUID)
		}
	}
	for x := 0; x < world.Length; x++ {
		for y := 0; y < world.Width; y++ {
			chunk := world.Chunks[x][y]
			chunks = append(chunks, &chunk)
		}
	}
	return chunks, nil
}
func (ws *WorldService) getChunksFromSpawnSpoint(spawnPoint player.SpawnPoint, viewDistance int) ([]*universe.Chunk, error) {
	var chunks []*universe.Chunk
	world := ws.World
	if len(world.Chunks) == 0 {
		err := ws.loadWorldChunks()
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

func (ws *WorldService) saveWorldChunks(world *universe.World) (err error) {
	for _, chunks := range world.Chunks {
		for _, chunk := range chunks {
			err = ws.chunkDao.Insert(world.UUID, &chunk)
			if err != nil {
				log.Println("can't save chunk : ", err)
			}
		}
	}
	return
}

func (ws *WorldService) generateChunk(position entity.Vector2) (*universe.Chunk, error) {
	chunk, err := ws.Generator.GenerateChunk(position.X, position.Y)
	if err != nil {
		log.Println("error generating chunk : ", err)
		return nil, err
	}
	chunk.UUID = gocql.TimeUUID().String()
	ws.chunkDao.Insert(ws.World.UUID, chunk)
	return chunk, nil
}

// TODO lock write with mutex
func (ws *WorldService) LoadChunksFromSpawnPoint(spawnPoint player.SpawnPoint) (chunks []*universe.Chunk, err error) {
	chunks, err = ws.getChunksFromSpawnSpoint(spawnPoint, 1)
	if err != nil {
		log.Println("failed to load chunks", err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to load chunks")
	}
	return chunks, nil
}

// TODO lock write with mutex
func (ws *WorldService) LoadAllChunks() (chunks []*universe.Chunk, err error) {
	chunks, err = ws.getAllChunks()
	if err != nil {
		log.Println("failed to load chunks", err)
		return nil, status.Errorf(codes.InvalidArgument, "failed to load chunks")
	}
	return chunks, nil
}

// TODO lock write with mutex
func (ws *WorldService) LoadSpecificChunks(coordinates []*pb.Vector2Int) ([]*universe.Chunk, error) {
	var chunks []*universe.Chunk
	world := ws.World
	if len(world.Chunks) == 0 {
		err := ws.loadWorldChunks()
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
