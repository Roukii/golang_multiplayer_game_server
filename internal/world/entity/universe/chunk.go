package universe

import (
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
)

type Chunk struct {
	UUID            string                 `json:"uuid"`
	Name            string                 `json:"name"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
	PositionX       int                    `json:"positionX"`
	PositionY       int                    `json:"positionY"`
	Tiles           []Tile                 `json:"tiles"`
	StaticEntities  []entity.StaticEntity  `json:"staticEntities"`
	DynamicEntities []entity.DynamicEntity `json:"dynamicEntities"`
	state           ChunkState             `json:"state"`
}

type ChunkState int8

const (
	Normal ChunkState = iota
	Combat
)
