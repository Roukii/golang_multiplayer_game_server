package dao

import (
	"fmt"

	entity "github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

// UniverseDao -.
type UniverseDao struct {
	session          *gocqlx.Session
	universeMetadata *table.Table
}

// New -.
func NewUniverseDao(session *gocqlx.Session) *UniverseDao {
	universeMetadata := table.New(table.Metadata{
		Name:    "game.universe",
		Columns: []string{"uuid", "name", "worlds", "players", "created_at"},
		PartKey: []string{"uuid"},
		SortKey: []string{},
	})

	return &UniverseDao{session, universeMetadata}
}

func (a UniverseDao) Insert(universe *entity.Universe) error {
	query := a.universeMetadata.InsertQuery(*a.session)
	query.BindStruct(universe)
	return query.ExecRelease()
}

func (a UniverseDao) Query() (entity.Universe, error) {
	var universe entity.Universe
	q := a.universeMetadata.SelectBuilder().Where(qb.EqLit("name", "uuid")).Query(*a.session)
	fmt.Print(q.String())
	var jsonString string
	err := q.Get(&jsonString)
	return universe, err
}
