package entity

import (
	"time"
)

type DynamicEntity struct {
	UUID       string            `json:"uuid"`
	Name       string            `json:"name"`
	PositionX  int8              `json:"positionX"`
	PositionY  int8              `json:"positionY"`
	EntityType DynamicEntityType `json:"type"`
	Stats      Stats             `json:"stats"`
	CreatedAt  time.Time         `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
}

type DynamicEntityType int8

const (
	Player DynamicEntityType = iota
	Creature
)
