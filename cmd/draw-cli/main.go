package main

import (
	"fmt"
	"petanque-draw/tournament"
)

const (
	MAX_FIELD = 8
	NB_ROUNDS = 3
)

func main() {
	tournament := tournament.DrawTournament(32, NB_ROUNDS, MAX_FIELD)

	for i, rounds := range tournament {
		fmt.Printf("\tRonde %d\n", i+1)
		for j, game := range rounds {
			fmt.Printf("Terrain %d : [ ", j+1)
			for _, player := range game.Team1 {
				fmt.Printf("%d ", player+1)
			}
			fmt.Printf("] | [ ")
			for _, player := range game.Team2 {
				fmt.Printf("%d ", player+1)
			}
			fmt.Printf("]\n")
		}
		fmt.Println("--------------------------")
	}
}
