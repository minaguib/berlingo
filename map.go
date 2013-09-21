package berlingo

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
	free_nodes       []*Node
	owned_nodes      []*Node
	enemy_nodes      []*Node
	controlled_nodes []*Node
}

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
		m.Nodes[rp.From].link_to(m.Nodes[rp.To])
		if m.Directed == false {
			m.Nodes[rp.To].link_to(m.Nodes[rp.From])
		}
	}

	return m, nil
}

func (m *Map) FreeNodes() []*Node {
	if m.free_nodes != nil {
		return m.free_nodes
	}
	m.free_nodes = make([]*Node, 0, len(m.Nodes)/2)
	for _, node := range m.Nodes {
		if node.IsFree() {
			m.free_nodes = append(m.free_nodes, node)
		}
	}
	return m.free_nodes
}

func (m *Map) OwnedNodes() []*Node {
	if m.owned_nodes != nil {
		return m.owned_nodes
	}
	m.owned_nodes = make([]*Node, 0, len(m.Nodes)/2)
	for _, node := range m.Nodes {
		if node.IsOwned() {
			m.owned_nodes = append(m.owned_nodes, node)
		}
	}
	return m.owned_nodes
}

func (m *Map) EnemyNodes() []*Node {
	if m.enemy_nodes != nil {
		return m.enemy_nodes
	}
	m.enemy_nodes = make([]*Node, 0, len(m.Nodes)/2)
	for _, node := range m.Nodes {
		if node.IsEnemy() {
			m.enemy_nodes = append(m.enemy_nodes, node)
		}
	}
	return m.enemy_nodes
}

func (m *Map) ControlledNodes() []*Node {
	if m.controlled_nodes != nil {
		return m.controlled_nodes
	}
	m.controlled_nodes = make([]*Node, 0, len(m.Nodes)/2)
	for _, node := range m.Nodes {
		if node.IsControlled() {
			m.controlled_nodes = append(m.controlled_nodes, node)
		}
	}
	return m.controlled_nodes
}
