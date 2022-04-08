package dynamic_entity

import (
	"github.com/Roukii/pock_multiplayer/internal/world/entity"
)

type dynamicEntityNode struct {
	next *dynamicEntityNode
	prev *dynamicEntityNode
	key  DynamicEntityFrame
}

type dynamicEntityLinkedList struct {
	head  *dynamicEntityNode
	tail  *dynamicEntityNode
	count int
}

type DynamicEntityFrame struct {
	Position                  map[string]entity.Vector3f
	elapsedTimeSinceLastFrame int
}

func (L *dynamicEntityLinkedList) Push(key DynamicEntityFrame) {
	list := &dynamicEntityNode{
		prev: L.head,
		key:  key,
	}
	if L.head != nil {
		L.head.next = list
	}
	L.head = list
	if L.tail == nil {
		L.tail = list
	}
	L.count++
}

func (L *dynamicEntityLinkedList) PopLast() {
	if L.tail == nil {
		return
	} else if L.tail == L.head {
		L.head = nil
		L.tail = nil
		return
	}
	list := L.tail
	L.tail = list.prev
	list.prev = nil
	L.count--
}
