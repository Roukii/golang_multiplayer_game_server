package dao

import (
	"log"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
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
	SpawnPoint SpawnPointType
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type SpawnPointType struct {
	gocqlx.UDT
	WorldUUID string    `cql:"world_uuid"`
	X         float32   `cql:"x"`
	Y         float32   `cql:"y"`
	Z         float32   `cql:"z"`
	UpdatedAt time.Time `cql:"updated_at"`
}

// New -.
func NewPlayerDao(session *gocqlx.Session) *PlayerDao {
	PlayerByUserMetadata := table.New(table.Metadata{
		Name:    "game.players_by_user",
		Columns: []string{"user_uuid", "player_uuid", "name", "stats", "spawn_point", "created_at", "updated_at"},
		PartKey: []string{"user_uuid", "player_uuid"},
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
		SpawnPoint: SpawnPointType{
			WorldUUID: p.SpawnPoint.WorldUUID,
			X:         p.SpawnPoint.Coordinate.Position.X,
			Y:         p.SpawnPoint.Coordinate.Position.Y,
			Z:         p.SpawnPoint.Coordinate.Position.Z,
			UpdatedAt: time.Now(),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}).ExecRelease()
}

func (a PlayerDao) Update(p *player.Player) error {
	q := qb.Update(a.PlayerByUserMetadata.Name()).
		Set("spawn_point", "stats", "updated_at").
		Where(qb.Eq("player_uuid")).
		Query(*a.session).
		SerialConsistency(gocql.Serial).
		BindStruct(PlayerByUser{
			PlayerUuid: mustParseUUID(p.UUID),
			Name:       p.Name,
			Stats:      p.Stats,
			SpawnPoint: SpawnPointType{
				WorldUUID: p.SpawnPoint.WorldUUID,
				X:         p.SpawnPoint.Coordinate.Position.X,
				Y:         p.SpawnPoint.Coordinate.Position.Y,
				Z:         p.SpawnPoint.Coordinate.Position.Z,
				UpdatedAt: time.Now(),
			},
			CreatedAt: p.CreatedAt,
			UpdatedAt: time.Now(),
		})

	return q.ExecRelease()
}

func (a PlayerDao) GetAllPlayersFromUserUUID(userUUID string) ([]*player.Player, error) {
	var players []*player.Player
	var playersByUser []*PlayerByUser

	if err := qb.Select(a.PlayerByUserMetadata.Name()).Where(qb.EqLit("user_uuid", userUUID)).Query(*a.session).Select(&playersByUser); err != nil {
		return nil, err
	}
	for _, p := range playersByUser {
		players = append(players, &player.Player{
			IDynamicEntity: entity.IDynamicEntity{
				UUID:       p.PlayerUuid.String(),
				Name:       p.Name,
				CreatedAt:  p.CreatedAt,
				UpdatedAt:  p.UpdatedAt,
				Stats:      p.Stats,
				EntityType: entity.Player,
			},
			CurrentWorldUUID: p.SpawnPoint.WorldUUID,
			SpawnPoint: player.SpawnPoint{WorldUUID: p.SpawnPoint.WorldUUID, Coordinate: entity.Position{Position: entity.Vector3f{X: p.SpawnPoint.X, Y: p.SpawnPoint.Y, Z: p.SpawnPoint.Z}}, UpdatedAt: p.SpawnPoint.UpdatedAt},
		})
	}
	return players, nil
}

func (a PlayerDao) GetPlayerFromUUID(userUUID string, playerUUID string) (*player.Player, error) {
	var p PlayerByUser
	q := qb.Select(a.PlayerByUserMetadata.Name()).Where(qb.EqLit("user_uuid", userUUID), qb.EqLit("player_uuid", playerUUID)).Query(*a.session)
	log.Println(q.Statement())
	if err := qb.Select(a.PlayerByUserMetadata.Name()).Where(qb.EqLit("user_uuid", userUUID), qb.EqLit("player_uuid", playerUUID)).Query(*a.session).Get(&p); err != nil {
		return nil, err
	}

	return &player.Player{
		IDynamicEntity: entity.IDynamicEntity{
			UUID:       p.PlayerUuid.String(),
			Name:       p.Name,
			CreatedAt:  p.CreatedAt,
			UpdatedAt:  p.UpdatedAt,
			Stats:      p.Stats,
			EntityType: entity.Player,
		},
		CurrentWorldUUID: p.SpawnPoint.WorldUUID,
		SpawnPoint: player.SpawnPoint{
			WorldUUID: p.SpawnPoint.WorldUUID,
			Coordinate: entity.Position{
				Position: entity.Vector3f{
					X: p.SpawnPoint.X,
					Y: p.SpawnPoint.Y,
					Z: p.SpawnPoint.Z,
				},
			},
			UpdatedAt: p.SpawnPoint.UpdatedAt,
		},
	}, nil
}

// TODO clean
func mustParseUUID(s string) gocql.UUID {
	u, err := gocql.ParseUUID(s)
	if err != nil {
		return gocql.TimeUUID()
	}
	return u
}
