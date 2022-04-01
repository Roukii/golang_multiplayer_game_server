package universe

import "time"

type World struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Level     int32     `json:"level"`
	Length    int64     `json:"length"`
	Width     int64     `json:"width"`
	Chunks    []Chunk   `json:"chunks"`
	Seed      string    `json:"seed"`
	Type      WorldType `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updatedAt"`
}

type WorldType int8

const (
	Base WorldType = iota
	Arid
	Green
)
