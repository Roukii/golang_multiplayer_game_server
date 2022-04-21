package universe

import (
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
)

type World struct {
	UUID            string `json:"uuid"`
	Name            string `json:"name"`
	MaxPlayer       int    `json:"max_player"`
	Level           int    `json:"level"`
	Length          int    `json:"length"`
	Width           int    `json:"width"`
	ScaleXY         float32
	ScaleHeight         float32
	Chunks          map[int]map[int]Chunk   `json:"chunks"`
	SpawnPoints     []player.SpawnPoint     `json:"chunks"`
	DynamicEntities []entity.IDynamicEntity `json:"dynamicEntities"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdateAt        time.Time               `json:"updatedAt"`
	Seed            string                  `json:"seed"`
	Type            WorldType               `json:"type"`
	Octave          float64
	Persistance     float64
	Lacunarity      float64
	UseFallOff      bool
}

type WorldType int8

const (
	Base WorldType = iota
	Arid
	Green
)
