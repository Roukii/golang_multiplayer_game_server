package universe

import "time"

type Tile struct {
	TileType  TileType  `json:"type"`
	Elevation float64   `json:"elevation"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TileType int8

const (
	Dirt TileType = iota
	Water
	Grass
	Rock
	Forest
	Sand
	Snow
)
