package berlingo

import (
	"errors"
	"io"
)

// AI is the primary interface that must be implemented by an author to use the berlingo framework
type AI interface {
	GameStart(*Game)
	Turn(*Game)
	GameOver(*Game)
	Ping(*Game)
}

type Game struct {

	// The AI implementation that will play the game
	Ai AI

	// The json-parsed request received for the game.  Should normally not be needed by an AI author
	Request *Request
	// The pending response that will be returned for the current move
	Response *Response

	// General information on the game
	Id                      string
	Number_Of_Players       int
	Maximum_Number_Of_Turns int
	Player_Id               int
	Time_Limit_Per_Turn     int

	// Information on the current turn
	Current_Turn int
	Turns_Left   int

	// The game's parsed map
	Map *Map
}

func NewGame(ai AI, r io.Reader) (game *Game, err error) {

	request, err := NewRequest(r)
	if err != nil {
		return nil, err
	}

	response, err := NewResponse()
	if err != nil {
		return nil, err
	}

	game = &Game{
		Ai:                      ai,
		Request:                 request,
		Response:                response,
		Id:                      request.Infos.Game_Id,
		Number_Of_Players:       request.Infos.Number_Of_Players,
		Maximum_Number_Of_Turns: request.Infos.Maximum_Number_Of_Turns,
		Player_Id:               request.Infos.Player_Id,
		Time_Limit_Per_Turn:     request.Infos.Time_Limit_Per_Turn,
		Current_Turn:            request.Infos.Current_Turn,
		Turns_Left:              request.Infos.Maximum_Number_Of_Turns - request.Infos.Current_Turn,
	}

	m, err := NewMap(game)
	if err != nil {
		return nil, err
	}
	game.Map = m

	return game, nil

}

// Do invokes the appropriate AI method based on the Request Action
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

// AddMove adds to the response queue the requested move characteristics
func (game *Game) AddMove(from_node *Node, to_node *Node, num_soldiers int) (err error) {
	if from_node.Available_Soldiers < num_soldiers {
		return errors.New("Not enough available soldiers")
	}
	from_node.Available_Soldiers -= num_soldiers
	to_node.Incoming_Soldiers += num_soldiers
	game.Response.AddMove(from_node, to_node, num_soldiers)
	return nil
}
