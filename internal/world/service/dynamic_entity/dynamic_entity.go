package dynamic_entity

import (
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity"
	"github.com/spf13/viper"
)

type DynamicEntityChange interface{}

type DynamicEntityService struct {
	DynamicEntities     map[string]entity.IDynamicEntity
	frames              dynamicEntityLinkedList
	maxSizeFrame        int
	lastUpdate          time.Time
	updateTicker        *time.Ticker
	tickerDone          chan bool
	DynamicEntityChange chan DynamicEntityChange
}

func NewDynamicEntityService() *DynamicEntityService {
	des := &DynamicEntityService{
		DynamicEntities: make(map[string]entity.IDynamicEntity),
	}
	des.maxSizeFrame = viper.GetInt("server.dynamic_entity.saved_frame_buffer_size")
	return des
}

func (des *DynamicEntityService) StopService() {
	des.updateTicker.Stop()
	des.tickerDone <- true
}

func (des *DynamicEntityService) UpdateService() {
	des.updateTicker = time.NewTicker(time.Duration(viper.GetInt("server.refresh_rate")))
	go func() {
		for {
			select {
			case <-des.tickerDone:
				return
			case <-des.updateTicker.C:
				elapsed := time.Since(des.lastUpdate)
				frame := DynamicEntityFrame{
					Position:                  make(map[string]entity.Vector3f),
					elapsedTimeSinceLastFrame: int(elapsed),
				}
				for _, entity := range des.DynamicEntities {
					entity.Update(elapsed.Milliseconds())
					frame.Position[entity.GetUUID()] = entity.Position.Position
				}
				des.frames.Push(frame)
				if des.frames.count > des.maxSizeFrame {
					des.frames.PopLast()
				}
				des.lastUpdate = time.Now()
			}
		}
	}()
}

func (des *DynamicEntityService) AddDynamicEntity(entity entity.IDynamicEntity) {
	des.DynamicEntities[entity.GetUUID()] = entity
}

func (des *DynamicEntityService) RemoveDynamicEntity(uuid string) {
	delete(des.DynamicEntities, uuid)
}
