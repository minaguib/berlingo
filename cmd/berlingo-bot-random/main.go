package main

import (
	"github.com/minaguib/berlingo"
	"log"
	"math/rand"
	"time"
)

/*
This example mimicks the functionality in the berlin-ai demo ruby gem at:
https://github.com/thirdside/berlin-ai/
*/

type randomAI struct{}

func (ai *randomAI) GameStart(game *berlingo.Game) {
}

func (ai *randomAI) Turn(game *berlingo.Game) {
	for _, node := range game.Map.ControlledNodes() {
		for _, other_node := range node.AdjacentNodes() {
			if node.Available_Soldiers < 1 {
				break
			}
			soldiers := rand.Intn(node.Available_Soldiers)
			log.Println("Moving", soldiers, "soldiers from node", node.Id, "to node", other_node.Id)
			if err := game.AddMove(node, other_node, soldiers); err != nil {
				log.Println("Error moving:", err)
			}
		}
	}
}

func (ai *randomAI) GameOver(game *berlingo.Game) {
}

func (ai *randomAI) Ping(game *berlingo.Game) {
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	berlingo.Serve(&randomAI{})
}
