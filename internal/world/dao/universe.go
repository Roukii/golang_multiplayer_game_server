package dao

import (
	"fmt"
	"strings"

	entity "github.com/Roukii/pock_multiplayer/internal/world/entity/universe"
	"github.com/gocql/gocql"
	"github.com/relops/cqlr"
)

// UniverseDao -.
type UniverseDao struct {
	session *gocql.Session
}

// New -.
func NewUniverseDao(session *gocql.Session) *UniverseDao {
	return &UniverseDao{session}
}

func (a UniverseDao) Insert(universe *entity.Universe) error {
	value := universe.GetJsonFields()
	colums := strings.Join(value, ", ")
	colums = strings.Trim(colums, "")
	query := "INSERT into universe(" + colums + ") VALUES ("
	for i := 0; i < len(value); i++ {
		if i != len(value)-1 {
			query = query + "?, "
		} else {
			query = query + "?)"
		}
	}
	fmt.Print(query)
	b := cqlr.Bind(query, universe)
	if err := b.Exec(a.session); err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}
