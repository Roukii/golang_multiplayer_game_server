package universe

type Tile struct {
	TileType  TileType `json:"type"`
	Elevation int8     `json:"elevation"`
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
