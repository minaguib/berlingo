package berlingo

type Node struct {
	Map                *Map
	Id                 int
	Type               *NodeType
	Paths_Outbound     map[int]*Node
	Paths_Inbound      map[int]*Node
	Player_Id          int
	Number_Of_Soldiers int
	Incoming_Soldiers  int
	Available_Soldiers int
}

func NewNode(m *Map) *Node {
	return &Node{
		Map:            m,
		Paths_Outbound: make(map[int]*Node),
		Paths_Inbound:  make(map[int]*Node),
	}
}

// Sets up a unidirectional link
func (node *Node) link_to(other *Node) {
	node.Paths_Outbound[other.Id] = other
	other.Paths_Inbound[node.Id] = node
}

func (node *Node) IsOwned() bool {
	return node.Player_Id == node.Map.My_Player_Id
}
func (node *Node) reset() {
	if node.IsOwned() {
		node.Available_Soldiers = node.Number_Of_Soldiers
		node.Incoming_Soldiers = 0
	} else {
		node.Available_Soldiers = 0
		node.Incoming_Soldiers = 0
	}
}
