package universe

import "github.com/scylladb/gocqlx/v2"

type Tile struct {
	gocqlx.UDT
	TileType  int `json:"tile_type"`
	Elevation float64  `json:"elevation"`
}

type TileType int

const (
	Dirt TileType = iota
	Water
	Grass
	Rock
	Forest
	Sand
	Snow
)
