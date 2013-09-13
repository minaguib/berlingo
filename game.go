package berlingo

import (
	"errors"
	"io"
)

type AI interface {
	GameStart(*Game)
	Turn(*Game)
	GameOver(*Game)
	Ping(*Game)
}

type Game struct {
	Ai       AI
	Request  *Request
	Response *Response
	Map      *Map
}

func NewGame(ai AI, r io.Reader) (game *Game, err error) {

	request, err := NewRequest(r)
	if err != nil {
		return nil, err
	}

	m, err := NewMap(request)
	if err != nil {
		return nil, err
	}

	response, err := NewResponse()
	if err != nil {
		return nil, err
	}

	game = &Game{
		Ai:       ai,
		Request:  request,
		Map:      m,
		Response: response,
	}

	return game, nil

}

func (game *Game) Do() {
	switch game.Request.Action {
	case "game_start":
		game.Ai.GameStart(game)
	case "turn":
		game.Ai.Turn(game)
	case "game_over":
		game.Ai.GameOver(game)
	case "ping":
		game.Ai.Ping(game)
	}
}

func (game *Game) AddMove(from_node *Node, to_node *Node, num_soldiers int) (err error) {
	if from_node.Available_Soldiers < num_soldiers {
		return errors.New("Not enough available soldiers")
	}
	from_node.Available_Soldiers -= num_soldiers
	to_node.Incoming_Soldiers += num_soldiers
	game.Response.AddMove(from_node, to_node, num_soldiers)
	return nil
}
