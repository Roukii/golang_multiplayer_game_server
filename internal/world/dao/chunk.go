package dao

import (
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
		X:         chunk.PositionX,
		Y:         chunk.PositionY,
		CreatedAt: time.Now(),
	})
	return query.ExecRelease()
}

func (a ChunkDao) LoadChunckBetweenCoordinate(minX int, maxX int, minY int, maxY int) ([]*universe.Chunk, error) {
	var chunks []*universe.Chunk
	var chunksByWorld []*ChunksByWorld

	//query := qb.Select(a.ChunksByWorldMetadata.Name()).Where(qb.LtOrEqLit("x", string(maxX)), qb.LtOrEqLit("y", string(maxY)), qb.GtOrEqLit("x", string(minX)), qb.GtOrEqLit("y", string(minY)))
	//_ := query.Query(*a.session).Iter()

	for _, c := range chunksByWorld {
		chunks = append(chunks, &universe.Chunk{
			UUID:      c.ChunkUuid.String(),
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			PositionX: c.X,
			PositionY: c.Y,
			Tiles:     c.Tiles,
		})
	}
	return chunks, nil
}
