package entity

import (
	"time"
)

type DynamicEntity struct {
	UUID       string            `json:"uuid"`
	Name       string            `json:"name"`
	Coordinate Position          `json:"position"`
	EntityType DynamicEntityType `json:"type"`
	Stats      Stats             `json:"stats"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updatedAt"`
}

type DynamicEntityType int8

const (
	Player DynamicEntityType = iota
	Creature
)
