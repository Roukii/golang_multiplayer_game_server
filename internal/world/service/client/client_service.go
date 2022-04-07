package client

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
	"github.com/Roukii/pock_multiplayer/pkg/helper"
)

const (
	clientTimeout = 15
)

type ClientService struct {
	clients map[string]*client
	game    *game.GameService
	mu      sync.Mutex
}

func NewClientService(game *game.GameService) *ClientService {
	clientService := &ClientService{
		clients: make(map[string]*client),
		game:    game,
		mu:      sync.Mutex{},
	}
	clientService.watchTimeout()
	return clientService
}

func (c *ClientService) GetClient(userUUID string) (*client, bool) {
	client, ok := c.clients[userUUID]
	return client, ok
}

func (c *ClientService) DisconnectClient(client *client, message string) {
	log.Printf("%s - removing client", client.userUUID)
	c.mu.Lock()
	playerUUID := c.clients[client.userUUID].playerUUID
	delete(c.clients, client.userUUID)
	c.mu.Unlock()
	c.removePlayer(playerUUID, message)
}

func (c *ClientService) removePlayer(playerUUID string, message string) {
	c.game.PlayerService.Mu.Lock()
	c.game.PlayerService.DisconnectPlayer(playerUUID)
	c.game.PlayerService.Mu.Unlock()

	resp := pb.PlayerStreamResponse{
		Uuid: playerUUID,
		Action: &pb.PlayerStreamResponse_Disconnect{
			Disconnect: &pb.PlayerDisconnect{
				Message: message,
			},
		},
	}
	c.BroadcastPlayer(&resp)
}

func (c *ClientService) AddClient(userUUID string, p *player.Player) {
	c.mu.Lock()
	c.clients[userUUID] = &client{
		lastMessage: time.Now(),
		Done:        make(chan error),
		playerUUID:  p.UUID,
		userUUID:    userUUID,
	}
	c.mu.Unlock()

	// TODO restrain connect player info
	resp := pb.PlayerStreamResponse{
		Uuid: p.UUID,
		Action: &pb.PlayerStreamResponse_Connect{
			Connect: &pb.PlayerConnect{
				Player: helper.PlayerTypeToProto(p),
			},
		},
	}
	c.BroadcastPlayer(&resp)
}

func (c *ClientService) BroadcastPlayer(resp *pb.PlayerStreamResponse) {
	c.mu.Lock()
	for id, currentClient := range c.clients {
		if currentClient.streamPlayerServer == nil {
			continue
		}
		if err := currentClient.streamPlayerServer.Send(resp); err != nil {
			log.Printf("%s - broadcast error %v", id, err)
			currentClient.Done <- errors.New("failed to broadcast message")
			continue
		}
		log.Printf("%s - broadcasted %+v", resp, id)
	}
	c.mu.Unlock()
}

func (c *ClientService) BroadcastChunk(resp *pb.ChunkStreamResponse) {
	c.mu.Lock()
	for id, currentClient := range c.clients {
		if currentClient.streamChunkServer == nil {
			continue
		}
		if err := currentClient.streamChunkServer.Send(resp); err != nil {
			log.Printf("%s - broadcast error %v", id, err)
			currentClient.Done <- errors.New("failed to broadcast message")
			continue
		}
		log.Printf("%s - broadcasted %+v", resp, id)
	}
	c.mu.Unlock()
}

func (c *ClientService) watchTimeout() {
	timeoutTicker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			for _, client := range c.clients {
				if time.Now().Sub(client.lastMessage).Minutes() > clientTimeout {
					client.Done <- errors.New("you have been timed out")
					return
				}
			}
			<-timeoutTicker.C
		}
	}()
}
