package universe

import (
	"fmt"
	"reflect"
	"time"

	entity "github.com/Roukii/pock_multiplayer/internal/world/entity/player"
)

type Universe struct {
	UUID      string          `json:"uuid"`
	Name      string          `json:"name"`
	Worlds    []World         `json:"worlds"`
	Players   []entity.Player `json:"players"`
	CreatedAt time.Time       `json:"createdAt"`
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
		fmt.Print(val.Field(i).String() + "\n")
		if val.Field(i).Type().Name() == name {
			return val.Field(i).String()
		}
	}

	return ""
}
