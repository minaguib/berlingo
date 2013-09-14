package berlingo

// Node represents a single node on the Map
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

// Sets up a unidirectional link pointing from this node towards another
func (node *Node) link_to(other *Node) {
	node.Paths_Outbound[other.Id] = other
	other.Paths_Inbound[node.Id] = node
}

func (node *Node) reset() {
	node.Incoming_Soldiers = 0
	node.Available_Soldiers = 0
	if node.IsOwned() {
		node.Available_Soldiers = node.Number_Of_Soldiers
	}
}

// IsFree returns whether the node is free, or owned by any player
func (node *Node) IsFree() bool {
	return node.Player_Id < 0
}

// IsOwned returns whether this node is owned by the current player
//
// Note - this deviates from the ruby client implementation, where ruby's owned? is essentially the opposite of free? - this is quite confusing as naturally asking a node.IsOwned() most likely indicates the caller wants to know if they own it themselves
//
// Callers who wish to mimick the owned? behavior of the ruby client may simply ask for !node.IsFree()
func (node *Node) IsOwned() bool {
	return node.IsOwnedBy(node.Map.Game.Player_Id)
}

func (node *Node) IsOwnedBy(player_id int) bool {
	return node.Player_Id == player_id
}

func (node *Node) IsEnemy() bool {
	return !node.IsFree() && !node.IsOwned()
}

func (node *Node) IsControlled() bool {
	return node.IsOwned() && node.Number_Of_Soldiers > 0
}
