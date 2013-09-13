package main

import (
	"fmt"
	"github.com/minaguib/berlingo"
	"math/rand"
)

/*
This example mimicks the functionality in the berlin-ai demo ruby gem at:
https://github.com/thirdside/berlin-ai/
*/

type AI1 struct{}

func (ai *AI1) GameStart(game *berlingo.Game) {
}

func (ai *AI1) Turn(game *berlingo.Game) {
	for _, node := range game.Map.ControlledNodes() {
		for _, other_node := range node.Paths_Outbound {
			if node.Available_Soldiers < 1 {
				continue
			}
			soldiers := rand.Intn(node.Available_Soldiers)
			fmt.Println("Moving", soldiers, "soldiers from node", node.Id, "to node", other_node.Id)
			if err := game.AddMove(node, other_node, soldiers); err != nil {
				fmt.Println("Error moving:", err)
			}
		}
	}
}

func (ai *AI1) GameOver(game *berlingo.Game) {
}

func (ai *AI1) Ping(game *berlingo.Game) {
}
