package dao

import (
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
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
	X         float64
	Y         float64
	Tiles     []universe.Tile
	CreatedAt time.Time
	UpdatedAt time.Time
}

// New -.
func NewChunkDao(session *gocqlx.Session) *ChunkDao {
	ChunksByWorldMetadata := table.New(table.Metadata{
		Name:    "game.chunks_by_world",
		Columns: []string{"chunk_uuid", "world_uuid", "x", "y", "z", "tiles", "created_at", "updated_at"},
		PartKey: []string{"chunk_uuid", "world_uuid"},
		SortKey: []string{"chunk_uuid"},
	})

	return &ChunkDao{session, ChunksByWorldMetadata}
}

func (a ChunkDao) Insert(worldUuid string, chunk *universe.Chunk) error {
	query := a.ChunksByWorldMetadata.InsertQuery(*a.session)
	query.BindStruct(ChunksByWorld{
		ChunkUuid: mustParseUUID(chunk.UUID),
		WorldUuid: mustParseUUID(worldUuid),
		X:         float64(chunk.PositionX),
		Y:         float64(chunk.PositionY),
		CreatedAt: time.Now(),
	})
	return query.ExecRelease()
}

func (a ChunkDao) LoadChunckBetweenCoordinate(minX int, maxX int, minY int, maxY int) ([]*universe.Chunk, error) {
	var chunk []*universe.Chunk
	query := a.ChunksByWorldMetadata.SelectQuery(*a.session).BindStruct(&PlayerByUser{
		PlayerUuid: mustParseUUID(playerUUID),
	})
	if err := query.Select(&p); err != nil {
		return nil, err
	}
	return &player.Player{
		UUID:       p.PlayerUuid.String(),
		Name:       p.Name,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		Stats:      p.Stats,
		SpawnPoint: p.SpawnPoint,
	}, nil
	return chunk, nil
}
