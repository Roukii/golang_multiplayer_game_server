package dao

import (
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/table"
)

// PlayerDao -.
type PlayerDao struct {
	session              *gocqlx.Session
	PlayerByUserMetadata *table.Table
}

type PlayerByUser struct {
	UserUuid   gocql.UUID
	PlayerUuid gocql.UUID
	Name       string
	Stats      entity.Stats
	SpawnPoint player.SpawnPoint
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// New -.
func NewPlayerDao(session *gocqlx.Session) *PlayerDao {
	PlayerByUserMetadata := table.New(table.Metadata{
		Name:    "game.players_by_user",
		Columns: []string{"user_uuid", "player_uuid", "name", "stats", "spawn_point", "created_at", "updated_at"},
		PartKey: []string{"user_uuid"},
		SortKey: []string{"player_uuid"},
	})

	return &PlayerDao{session, PlayerByUserMetadata}
}

func (a PlayerDao) Insert(userUuid string, p *player.Player) error {
	return a.PlayerByUserMetadata.InsertQuery(*a.session).BindStruct(PlayerByUser{
		UserUuid:   mustParseUUID(userUuid),
		PlayerUuid: mustParseUUID(p.UUID),
		Name:       p.Name,
		Stats:      p.Stats,
		SpawnPoint: p.SpawnPoint,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}).ExecRelease()
}

func (a PlayerDao) GetAllPlayersFromUserUUID(userUUID string) ([]*player.Player, error) {
	var players []*player.Player
	var playersByUser []*PlayerByUser

	query := a.PlayerByUserMetadata.SelectQuery(*a.session).BindStruct(&PlayerByUser{
		UserUuid: mustParseUUID(userUUID),
	})
	if err := query.Select(&playersByUser); err != nil {
		return players, err
	}
	for _, p := range playersByUser {
		players = append(players, &player.Player{
			UUID:       p.PlayerUuid.String(),
			Name:       p.Name,
			CreatedAt:  p.CreatedAt,
			UpdatedAt:  p.UpdatedAt,
			Stats:      p.Stats,
			SpawnPoint: p.SpawnPoint,
		})
	}
	return players, nil
}

func (a PlayerDao) GetPlayerFromUUID(playerUUID string) (*player.Player, error) {
	var p *PlayerByUser
	query := a.PlayerByUserMetadata.SelectQuery(*a.session).BindStruct(&PlayerByUser{
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
}

func mustParseUUID(s string) gocql.UUID {
	u, err := gocql.ParseUUID(s)
	if err != nil {
		panic(err)
	}
	return u
}
