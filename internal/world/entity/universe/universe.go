package universe

import (
	"log"
	"reflect"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
)

type Universe struct {
	UUID      string           `json:"uuid"`
	Name      string           `json:"name"`
	Worlds    map[string]World `json:"worlds"`
	Players   []player.Player  `json:"players"`
	CreatedAt time.Time        `json:"created_at"`
}

func (a Universe) GetJsonFields() []string {
	val := reflect.ValueOf(a)
	jsonkey := []string{}
	for i := 0; i < val.Type().NumField(); i++ {
		jsonkey = append(jsonkey, val.Type().Field(i).Tag.Get("json"))
	}
	return jsonkey
}

func (a Universe) GetValueByFieldName(name string) string {
	val := reflect.ValueOf(a)
	for i := 0; i < val.NumField(); i++ {
		log.Print(val.Field(i).String() + "\n")
		if val.Field(i).Type().Name() == name {
			return val.Field(i).String()
		}
	}

	return ""
}
