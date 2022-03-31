package entity

import (
	"time"
)

type StaticEntity struct {
	UUID             string           `json:"uuid"`
	Name             string           `json:"name"`
	PositionX        int8             `json:"x"`
	PositionY        int8             `json:"y"`
	EntityType       StaticEntityType `json:"type"`
	Stats            Stats            `json:"stats"`
	EntryToChunkUUID string           `json:"entryToChunkUuid"`
	CreatedAt        time.Time        `json:"createdAt"`
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
