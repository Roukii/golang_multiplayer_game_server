package entity

import "github.com/scylladb/gocqlx/v2"

type Stats struct {
	gocqlx.UDT
	Level int `json:"level"`
	Maxhp int `json:"maxhp"`
	Hp    int `json:"hp"`
	Maxmp int `json:"maxmp"`
	Mp    int `json:"mp"`
}
