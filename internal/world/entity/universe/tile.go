package universe

import "time"

type Tile struct {
	TileType  TileType  `json:"type"`
	Elevation int8      `json:"elevation"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TileType int8

const (
	Dirt TileType = iota
	Grass
	Rock
	Forest
	Sand
	Snow
)
