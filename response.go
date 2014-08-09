package berlingo

import (
	"encoding/json"
)

// ResponseMove is a single move in the list of moves that make up a Response
type ResponseMove struct {
	From     int `json:"from"`
	To       int `json:"to"`
	Soldiers int `json:"number_of_soldiers"`
}

// Response represents the pending response to the game command
type Response struct {
	Moves []*ResponseMove
}

// NewResponse initializes a new empty response
func NewResponse() (*Response, error) {
	return &Response{
		Moves: make([]*ResponseMove, 0, 10),
	}, nil
}

// ToJSON returns a JSON representation of the response
func (response *Response) ToJSON() (responseJSON []byte, err error) {
	return json.Marshal(response.Moves)
}

// AddMove appends a move to the response
func (response *Response) AddMove(fromNode *Node, toNode *Node, numSoldiers int) {
	response.Moves = append(response.Moves, &ResponseMove{fromNode.Id, toNode.Id, numSoldiers})
}
