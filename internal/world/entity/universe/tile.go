package universe

import "github.com/scylladb/gocqlx/v2"

type Tile struct {
	gocqlx.UDT
	TileType  TileType     `cql:"tile_type"`
	Elevation float64 `cql:"elevation"`
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
