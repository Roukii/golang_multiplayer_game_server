package client

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Roukii/pock_multiplayer/internal/world/entity/player"
	pb "github.com/Roukii/pock_multiplayer/internal/world/proto"
	player_action "github.com/Roukii/pock_multiplayer/internal/world/service/action/player"
	"github.com/Roukii/pock_multiplayer/internal/world/service/game"
)

const (
	clientTimeout = 15
	timeOutMessage = "timed out"
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

func (c *ClientService) AddClient(userUUID string, p *player.Player) {
	c.mu.Lock()
	c.clients[userUUID] = &client{
		lastMessage: time.Now(),
		PlayerDone:  make(chan error),
		ChunkDone:   make(chan error),
		playerUUID:  p.UUID,
		userUUID:    userUUID,
	}
	c.mu.Unlock()

	// TODO restrain connect player info
	c.game.PlayerChangeChannel <- player_action.ConnectPlayerChange{
		Player: p,
	}
}

func (c *ClientService) BroadcastPlayer(resp *pb.PlayerStreamResponse) {
	c.mu.Lock()
	for id, currentClient := range c.clients {
		if currentClient.streamPlayerServer == nil {
			continue
		}
		if err := currentClient.streamPlayerServer.Send(resp); err != nil {
			log.Printf("%s - broadcast error %v", id, err)
			currentClient.PlayerDone <- errors.New("failed to broadcast message")
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
			currentClient.ChunkDone <- errors.New("failed to broadcast message")
			continue
		}
		log.Printf("%s - broadcasted %+v", resp, id)
	}
	c.mu.Unlock()
}

// TODO change to disconnect only the stream to allow reconnect
func (c *ClientService) DisconnectClient(client *client, message string) {
	log.Println("removing client : ", client.userUUID)
	c.mu.Lock()
	playerUUID := c.clients[client.userUUID].playerUUID
	delete(c.clients, client.userUUID)
	c.mu.Unlock()
	c.removePlayer(playerUUID, message)
}

// TODO check if there is not a double sending of disconnect
func (c *ClientService) removePlayer(playerUUID string, message string) {
	c.game.PlayerService.DisconnectPlayer(playerUUID)

	c.game.SendPlayerChange(player_action.DisconnectPlayerChange{
		PlayerUUID: playerUUID,
		Message:    message,
	})
}

func (c *ClientService) watchTimeout() {
	timeoutTicker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			for _, client := range c.clients {
				if time.Now().Sub(client.lastMessage).Seconds() > clientTimeout {
					if client.streamPlayerServer != nil {
						client.streamPlayerServer.Send(&pb.PlayerStreamResponse{
							Uuid: client.playerUUID,
							Info: &pb.PlayerStreamResponse_DynamicEntity{
								DynamicEntity: pb.DynamicEntityType_PLAYER,
							},
							Action: &pb.PlayerStreamResponse_Disconnect{
								Disconnect: &pb.PlayerDisconnect{
									Message: timeOutMessage,
								},
							},
						})
						client.PlayerDone <- errors.New(timeOutMessage)
					}
					if client.streamChunkServer != nil {
						client.ChunkDone <- errors.New(timeOutMessage)
					}
					c.DisconnectClient(client, timeOutMessage)
					return
				}
			}
			<-timeoutTicker.C
		}
	}()
}
