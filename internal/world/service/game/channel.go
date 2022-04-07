package game

type PlayerAction interface {
	Perform(game *GameService)
}

type PlayerChange interface{}

func (game *GameService) SendPlayerChange(change PlayerChange) {
	select {
	case game.PlayerChangeChannel <- change:
	default:
	}
}
