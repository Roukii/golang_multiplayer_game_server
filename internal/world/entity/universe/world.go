package universe

import (
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
)

type World struct {
	UUID        string                `json:"uuid"`
	Name        string                `json:"name"`
	MaxPlayer   int                   `json:"max_player"`
	Level       int                   `json:"level"`
	Length      int                   `json:"length"`
	Width       int                   `json:"width"`
	Chunks      map[int]map[int]Chunk `json:"chunks"`
	SpawnPoints []player.SpawnPoint   `json:"chunks"`
	Seed        string                `json:"seed"`
	Type        WorldType             `json:"type"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdateAt    time.Time             `json:"updatedAt"`
}

type WorldType int8

const (
	Base WorldType = iota
	Arid
	Green
)
