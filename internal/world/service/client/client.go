package client

import (
	"errors"
	"log"
	"time"

	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
)

type client struct {
	streamPlayerServer pb.PlayerService_StreamServer
	streamChunkServer  pb.ChunkService_StreamServer
	PlayerDone               chan error
	ChunkDone               chan error
	lastMessage        time.Time
	userUUID           string
	playerUUID         string
}

func (c *client) AddPlayerStream(stream pb.PlayerService_StreamServer) error {
	if c.streamPlayerServer != nil {
		return errors.New("stream already active")
	}
	c.streamPlayerServer = stream
	log.Println("Start new player stream for : ", c.playerUUID)
	return nil
}


func (c *client) RemovePlayerStream() {
	c.streamPlayerServer = nil
}

func (c *client) AddChunkStream(stream pb.ChunkService_StreamServer) error {
	if c.streamChunkServer != nil {
		return errors.New("stream already active")
	}
	c.streamChunkServer = stream
	log.Println("Start new chunk stream for : ", c.playerUUID)

	return nil
}

func (c *client) GetPlayerUUID() string {
	return c.playerUUID
}

func (c *client) Update() {
	c.lastMessage = time.Now()
}