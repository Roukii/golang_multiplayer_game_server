package entity

import "time"

type SpawnPoint struct {
	WorldUUID string    `json:"world_uuid"`
	ChunkUUID string    `json:"chunk_uuid"`
	PositionX int8      `json:"position_x"`
	PositionY int8      `json:"position_y"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
