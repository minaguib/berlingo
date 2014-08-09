package berlingo

import (
	"encoding/json"
	"io"
)

// Request represents the request, semi-verbatim, as received from the berlin server
// AI writers are not expected to have to deal with Request, but rather higher-level parsed
// representations such as Game, Map and Node
type Request struct {
	Action string
	Infos  struct {
		Game_Id                 string
		Current_Turn            int
		Maximum_Number_Of_Turns int
		Number_Of_Players       int
		Time_Limit_Per_Turn     int
		Directed                bool
		Player_Id               int
	}
	Map struct {
		Types []struct {
			Name              string
			Points            int
			Soldiers_Per_Turn int
		}
		Nodes []struct {
			Id   int
			Type string
		}
		Paths []struct {
			From int
			To   int
		}
	}
	State []struct {
		Node_Id            int
		Player_Id          int  `json:"-"`
		Player_Id_Ptr      *int `json:"Player_Id"`
		Number_Of_Soldiers int
	}
}

// NewRequest initializes a new request from the given io.Reader
func NewRequest(r io.Reader) (request *Request, err error) {

	request = new(Request)
	dec := json.NewDecoder(r)
	err = dec.Decode(request)
	if err != nil {
		return nil, err
	}

	/* Incoming player 0 and null is indistinguishable, so we use a *Ptr
	in the struct to capture it to differentiate, and set Player_Id properly or -1
	*/
	for i := range request.State {
		if request.State[i].Player_Id_Ptr == nil {
			request.State[i].Player_Id = -1
		} else {
			request.State[i].Player_Id = *request.State[i].Player_Id_Ptr
		}
	}

	return request, nil
}
