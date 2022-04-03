package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/table"
)

// ChunkDao -.
type ChunkDao struct {
	session               *gocqlx.Session
	ChunksByWorldMetadata *table.Table
}

type ChunksByWorld struct {
	ChunkUuid gocql.UUID
	WorldUuid gocql.UUID
	X         int
	Y         int
	Tiles     []universe.Tile
	CreatedAt time.Time
}

// New -.
func NewChunkDao(session *gocqlx.Session) *ChunkDao {
	ChunksByWorldMetadata := table.New(table.Metadata{
		Name:    "game.chunks_by_world",
		Columns: []string{"world_uuid", "chunk_uuid", "x", "y", "tiles", "created_at"},
		PartKey: []string{"world_uuid"},
		SortKey: []string{},
	})

	return &ChunkDao{session, ChunksByWorldMetadata}
}

func (a ChunkDao) Insert(worldUuid string, chunk *universe.Chunk) error {
	query := a.ChunksByWorldMetadata.InsertQuery(*a.session)
	query.BindStruct(ChunksByWorld{
		ChunkUuid: mustParseUUID(chunk.UUID),
		WorldUuid: mustParseUUID(worldUuid),
		X:         chunk.PositionX,
		Y:         chunk.PositionY,
		Tiles:     chunk.Tiles,
		CreatedAt: time.Now(),
	})
	return query.ExecRelease()
}

func (a ChunkDao) LoadChunckBetweenCoordinate(worldUuid string, minX int, maxX int, minY int, maxY int) ([]*universe.Chunk, error) {
	var chunks []*universe.Chunk
	var chunksByWorld []*ChunksByWorld
	var ctx context.Context
	var arg []string
	arg = append(arg, worldUuid)
	fmt.Println(arg)
	if err := a.session.Query(`SELECT * FROM game.chunks_by_world WHERE world_uuid=33d34284-b2f0-11ec-911b-367dda4cfa8c AND x<=1 AND y<=1 AND x>=-1 AND y>=-1 ALLOW FILTERING`, arg).
		WithContext(ctx).Select(&chunksByWorld); err != nil {
		return nil, err
	}
	fmt.Print()
	for _, c := range chunksByWorld {
		chunks = append(chunks, &universe.Chunk{
			UUID:      c.ChunkUuid.String(),
			CreatedAt: c.CreatedAt,
			PositionX: c.X,
			PositionY: c.Y,
			Tiles:     c.Tiles,
		})
	}
	return chunks, nil
}
