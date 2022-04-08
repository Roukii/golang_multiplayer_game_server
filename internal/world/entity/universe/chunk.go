package universe

import (
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
)

type Chunk struct {
	UUID            string                 `json:"uuid"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updatedAt"`
	PositionX       int                    `json:"positionX"`
	PositionY       int                    `json:"positionY"`
	Tiles           []Tile                 `json:"tiles"`
	StaticEntities  []entity.StaticEntity  `json:"staticEntities"`
}
