package player

import (
	"log"

	entity "github.com/Roukii/pock_multiplayer/internal/world/entity"
)

type Player struct {
	entity.IDynamicEntity
	SpawnPoint SpawnPoint `json:"spawn_point"`
	CurrentWorldUUID string
}

func (p *Player) Update(elapstedTime int64) {
	log.Println("update player : ", p.UUID)
}
