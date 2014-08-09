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

// Game represents the gamestate at a given turn
type Game struct {

	// The AI implementation that will play the game
	Ai AI

	// The json-parsed request received for the game.  Should normally not be needed by an AI author
	Request *Request
	// The pending response that will be returned for the current action
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

// NewGame initializes a new game
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

// DoAction invokes the appropriate AI method based on the Request Action
func (game *Game) DoAction() {
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
func (game *Game) AddMove(fromNode *Node, toNode *Node, numSoldiers int) (err error) {
	if !fromNode.IsOwned() {
		return errors.New("Cannot move soldiers from a node you don't own")
	} else if fromNode.Available_Soldiers < numSoldiers {
		return errors.New("Not enough available soldiers")
	}
	fromNode.Available_Soldiers -= numSoldiers
	toNode.Incoming_Soldiers += numSoldiers
	game.Response.AddMove(fromNode, toNode, numSoldiers)
	return nil
}
