package berlingo

// Node represents a single node on the Map
type Node struct {
	Map  *Map
	Id   int
	Type *NodeType

	Paths              map[int]*Node
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
		Paths:          make(map[int]*Node),
		Paths_Outbound: make(map[int]*Node),
		Paths_Inbound:  make(map[int]*Node),
	}
}

// Sets up a unidirectional link pointing from this node towards another
func (node *Node) link_to(other *Node) {

	node.Paths_Outbound[other.Id] = other
	node.Paths[other.Id] = other

	other.Paths_Inbound[node.Id] = node
	other.Paths[node.Id] = node
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

// Returns one of "free", "enemy", "owned"
func (node *Node) OwnershipStatus() string {
	if node.IsFree() {
		return "free"
	} else if node.IsEnemy() {
		return "enemy"
	} else if node.IsOwned() {
		return "owned"
	} else {
		panic("Illegal ownership status calculation")
	}
}

func (node *Node) IsControlled() bool {
	return node.IsOwned() && node.Number_Of_Soldiers > 0
}

func (node *Node) HasOutboundPathTo(other_node *Node) bool {
	_, ok := node.Paths_Outbound[other_node.Id]
	return ok
}

func (node *Node) HasInboundPathFrom(other_node *Node) bool {
	_, ok := node.Paths_Inbound[other_node.Id]
	return ok
}

func (node *Node) IsAdjacentTo(other_node *Node) bool {
	_, ok := node.Paths[other_node.Id]
	return ok
}

func (node *Node) AdjacentNodes() (nodes []*Node) {
	nodes = make([]*Node, 0, len(node.Paths))
	for _, node := range node.Paths {
		nodes = append(nodes, node)
	}
	return nodes
}
