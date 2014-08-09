package berlingo

// NodeType declares the type of a node (re-usable characteristics)
type NodeType struct {
	Name              string
	Points            int
	Soldiers_Per_Turn int
}

// Map represents the map and the nodes in it
type Map struct {
	Game      *Game
	Directed  bool
	NodeTypes map[string]*NodeType
	Nodes     map[int]*Node

	// Lazy caches
	freeNodes       []*Node
	ownedNodes      []*Node
	enemyNodes      []*Node
	controlledNodes []*Node
}

// NewMap initializes a new map
func NewMap(game *Game) (m *Map, err error) {

	request := game.Request

	m = &Map{
		Game:      game,
		Directed:  request.Infos.Directed,
		NodeTypes: make(map[string]*NodeType),
		Nodes:     make(map[int]*Node),
	}

	for _, rt := range request.Map.Types {
		m.NodeTypes[rt.Name] = &NodeType{
			Name:              rt.Name,
			Points:            rt.Points,
			Soldiers_Per_Turn: rt.Soldiers_Per_Turn,
		}
	}

	for _, rn := range request.Map.Nodes {
		node := NewNode(m)
		node.Id = rn.Id
		node.Type = m.NodeTypes[rn.Type]
		m.Nodes[rn.Id] = node
	}

	for _, rs := range request.State {
		node := m.Nodes[rs.Node_Id]
		node.Player_Id = rs.Player_Id
		node.Number_Of_Soldiers = rs.Number_Of_Soldiers
		node.reset()
	}

	for _, rp := range request.Map.Paths {
		m.Nodes[rp.From].linkTo(m.Nodes[rp.To])
		if m.Directed == false {
			m.Nodes[rp.To].linkTo(m.Nodes[rp.From])
		}
	}

	return m, nil
}

// FreeNodes returns an array of nodes on this map that are free
func (m *Map) FreeNodes() []*Node {
	if m.freeNodes != nil {
		return m.freeNodes
	}
	m.freeNodes = make([]*Node, 0, len(m.Nodes)/2)
	for _, node := range m.Nodes {
		if node.IsFree() {
			m.freeNodes = append(m.freeNodes, node)
		}
	}
	return m.freeNodes
}

// OwnedNodes returns an array of nodes on this map that are owned
func (m *Map) OwnedNodes() []*Node {
	if m.ownedNodes != nil {
		return m.ownedNodes
	}
	m.ownedNodes = make([]*Node, 0, len(m.Nodes)/2)
	for _, node := range m.Nodes {
		if node.IsOwned() {
			m.ownedNodes = append(m.ownedNodes, node)
		}
	}
	return m.ownedNodes
}

// EnemyNodes returns an array of nodes on this map that are enemy nodes
func (m *Map) EnemyNodes() []*Node {
	if m.enemyNodes != nil {
		return m.enemyNodes
	}
	m.enemyNodes = make([]*Node, 0, len(m.Nodes)/2)
	for _, node := range m.Nodes {
		if node.IsEnemy() {
			m.enemyNodes = append(m.enemyNodes, node)
		}
	}
	return m.enemyNodes
}

// ControlledNodes returns an array of nodes on this map that are controlled by the current player
func (m *Map) ControlledNodes() []*Node {
	if m.controlledNodes != nil {
		return m.controlledNodes
	}
	m.controlledNodes = make([]*Node, 0, len(m.Nodes)/2)
	for _, node := range m.Nodes {
		if node.IsControlled() {
			m.controlledNodes = append(m.controlledNodes, node)
		}
	}
	return m.controlledNodes
}
