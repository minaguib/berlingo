package berlingo

type NodeType struct {
	Name               string
	Points             int
	Number_Of_Soldiers int
}

// Map represents the map and the nodes in it
type Map struct {
	Game             *Game
	NodeTypes        map[string]*NodeType
	Nodes            map[int]*Node
	controlled_nodes []*Node
}

func NewMap(game *Game) (m *Map, err error) {

	request := game.Request

	m = &Map{
		Game:      game,
		NodeTypes: make(map[string]*NodeType),
		Nodes:     make(map[int]*Node),
	}

	for _, rt := range request.Map.Types {
		m.NodeTypes[rt.Name] = &NodeType{
			Name:               rt.Name,
			Points:             rt.Points,
			Number_Of_Soldiers: rt.Number_Of_Soldiers,
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
		if request.Infos.Directed == false {
			m.Nodes[rp.To].link_to(m.Nodes[rp.From])
		}
	}

	return m, nil
}

func (m *Map) ControlledNodes() []*Node {
	if m.controlled_nodes == nil {
		m.controlled_nodes = make([]*Node, 0, len(m.Nodes)/2)
		for _, node := range m.Nodes {
			if node.Player_Id == m.Game.Player_Id {
				m.controlled_nodes = append(m.controlled_nodes, node)
			}
		}
	}
	return m.controlled_nodes
}
