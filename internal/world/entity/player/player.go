package player

import (
	"time"

	entity "github.com/Roukii/pock_multiplayer/internal/world/entity"
)

type Player struct {
	UUID       string       `json:"uuid"`
	Name       string       `json:"name"`
	Level      int          `json:"level"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updatedAt"`
	Stats      entity.Stats `json:"stats"`
	SpawnPoint SpawnPoint   `json:"spawnPoint"`
}
