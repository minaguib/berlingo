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

// Request represents the pending response to the game command
type Response struct {
	Moves []*ResponseMove
}

func NewResponse() (*Response, error) {
	return &Response{
		Moves: make([]*ResponseMove, 0, 10),
	}, nil
}

func (response *Response) ToJson() (response_json []byte, err error) {
	return json.Marshal(response.Moves)
}

func (response *Response) AddMove(from_node *Node, to_node *Node, num_soldiers int) {
	response.Moves = append(response.Moves, &ResponseMove{from_node.Id, to_node.Id, num_soldiers})
}
