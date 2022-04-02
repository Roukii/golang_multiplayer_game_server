package entity

import (
	"time"
)

type StaticEntity struct {
	UUID             string           `json:"uuid"`
	Name             string           `json:"name"`
	Coordinate       Position         `json:"position"`
	EntityType       StaticEntityType `json:"type"`
	Stats            Stats            `json:"stats"`
	EntryToChunkUUID string           `json:"entryToChunkUuid"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updatedAt"`
}

type StaticEntityType int8

const (
	Empty StaticEntityType = iota
	Building
	ChunkEntry
	Ressource
	Item
)
