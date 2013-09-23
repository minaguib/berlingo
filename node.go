package berlingo

// Node represents a single node on the Map
type Node struct {
	Map  *Map
	Id   int
	Type *NodeType

	Player_Id          int
	Number_Of_Soldiers int
	Incoming_Soldiers  int
	Available_Soldiers int
	// Analysis is not populated by berlingo.  It's an area where your own AI may assign
	// custom analysis values to nodes
	Analysis interface{}

	paths          map[int]*Node
	paths_outbound map[int]*Node
	paths_inbound  map[int]*Node
}

func NewNode(m *Map) *Node {
	return &Node{
		Map:            m,
		paths:          make(map[int]*Node),
		paths_outbound: make(map[int]*Node),
		paths_inbound:  make(map[int]*Node),
	}
}

// Sets up a unidirectional link pointing from this node towards another
func (node *Node) link_to(other *Node) {

	node.paths_outbound[other.Id] = other
	node.paths[other.Id] = other

	other.paths_inbound[node.Id] = node
	other.paths[node.Id] = node
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
	_, ok := node.paths_outbound[other_node.Id]
	return ok
}

func (node *Node) HasInboundPathFrom(other_node *Node) bool {
	_, ok := node.paths_inbound[other_node.Id]
	return ok
}

func (node *Node) IsAdjacentTo(other_node *Node) bool {
	_, ok := node.paths[other_node.Id]
	return ok
}

func (node *Node) AdjacentNodes() (nodes []*Node) {
	nodes = make([]*Node, 0, len(node.paths))
	for _, node := range node.paths {
		nodes = append(nodes, node)
	}
	return nodes
}
