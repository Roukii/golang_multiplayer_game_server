package player

import (
	"time"

	entity "github.com/Roukii/pock_multiplayer/internal/world/entity"
)

type Player struct {
	UUID       string       `json:"player_uuid"`
	Name       string       `json:"name"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	Stats      entity.Stats `json:"stats"`
	SpawnPoint SpawnPoint   `json:"spawn_point"`
}
