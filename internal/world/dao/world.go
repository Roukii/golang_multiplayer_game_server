package dao

import (
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

// WorldDao -.
type WorldDao struct {
	session       *gocqlx.Session
	WorldMetadata *table.Table
}

type World struct {
	WorldUuid   gocql.UUID
	Name        string
	Seed        string
	Length      int
	Width       int
	MaxPlayer   int
	CreatedAt   time.Time
	SpawnPoints []SpawnPointType
}

// New -.
func NewWorldDao(session *gocqlx.Session) *WorldDao {
	world := table.New(table.Metadata{
		Name:    "game.world",
		Columns: []string{"world_uuid", "name", "seed", "length", "width", "max_player", "created_at", "spawn_points"},
		PartKey: []string{"world_uuid"},
		SortKey: []string{},
	})

	return &WorldDao{session, world}
}

func (a WorldDao) Insert(world *universe.World) error {
	var spawns []SpawnPointType
	for _, spawn := range world.SpawnPoints {
		spawns = append(spawns, SpawnPointType{
			WorldUUID: spawn.WorldUUID,
			X:         spawn.Coordinate.Position.X,
			Y:         spawn.Coordinate.Position.Y,
			Z:         spawn.Coordinate.Position.Z,
			UpdatedAt: time.Now(),
		})
	}
	query := a.WorldMetadata.InsertQuery(*a.session)
	query.BindStruct(World{
		WorldUuid:   mustParseUUID(world.UUID),
		Name:        world.Name,
		Seed:        world.Seed,
		Length:      world.Length,
		Width:       world.Width,
		MaxPlayer:   world.MaxPlayer,
		SpawnPoints: spawns,
		CreatedAt:   time.Now(),
	})
	return query.ExecRelease()
}

func (a WorldDao) GetAllWorlds() ([]universe.World, error) {
	var worldsEntity []universe.World
	var worlds []*World
	if err := qb.Select(a.WorldMetadata.Name()).Query(*a.session).Select(&worlds); err != nil {
		return nil, err
	}

	for _, w := range worlds {
		var spawns []player.SpawnPoint
		for _, spawn := range w.SpawnPoints {
			spawns = append(spawns, player.SpawnPoint{
				WorldUUID: spawn.WorldUUID,
				Coordinate: entity.Position{
					Position: entity.Vector3f{
						X: spawn.X,
						Y: spawn.Y,
						Z: spawn.Z,
					},
				},
				UpdatedAt: spawn.UpdatedAt,
			})
		}
		worldsEntity = append(worldsEntity, universe.World{
			UUID:        w.WorldUuid.String(),
			Name:        w.Name,
			Seed:        w.Seed,
			Length:      w.Length,
			Width:       w.Width,
			MaxPlayer:   w.MaxPlayer,
			CreatedAt:   w.CreatedAt,
			SpawnPoints: spawns,
		})
	}
	return worldsEntity, nil
}
