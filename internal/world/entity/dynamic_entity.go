package entity

import (
	"time"
)

type DynamicEntitiy interface {
	GetUUID() string
	GetName() string
	Update(elapstedTime int64)
	GetPosition() *Position
	SetPosition(pos *Position) *Position
	GetStats() *Stats
	GetType() DynamicEntitiy
	SetStats(stats *Stats)
}

type IDynamicEntity struct {
	DynamicEntitiy
	UUID       string            `json:"uuid"`
	Name       string            `json:"name"`
	Position   Position          `json:"position"`
	EntityType DynamicEntityType `json:"type"`
	Stats      Stats             `json:"stats"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updatedAt"`
}

type DynamicEntityType int8

const (
	Player DynamicEntityType = iota
	Creature
	Projectile
)

func (ade *IDynamicEntity) GetUUID() string {
	return ade.UUID
}

func (ade *IDynamicEntity) GetName() string {
	return ade.Name
}

func (ade *IDynamicEntity) GetPosition() Position {
	return ade.Position
}

func (ade *IDynamicEntity) SetPosition(position *Position) {
	ade.Position = *position
}

func (ade *IDynamicEntity) GetStats() *Stats {
	return &ade.Stats
}

func (ade *IDynamicEntity) SetStats(stats *Stats) {
	ade.Stats = *stats
}
