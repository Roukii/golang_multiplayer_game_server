package player

import (
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
)

type SpawnPoint struct {
	WorldUUID  string          `json:"world_uuid"`
	Coordinate entity.Position `json:"coordinate"`
	UpdatedAt  time.Time       `json:"updated_at"`
}
