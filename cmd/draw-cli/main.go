package main

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"petanque-draw/tournament"
	"strconv"
	"time"
)

const (
	MAX_FIELD       = 9
	NB_ROUNDS       = 4
	NB_PLAYERS      = 16
	MAX_GENERATIONS = 100000000
)

func Display(f io.Writer, t tournament.Tournament) {
	for i, rounds := range t {
		fmt.Fprintf(f, "\tRonde %d\n", i+1)
		for j, game := range rounds {
			fmt.Fprintf(f, "Terrain %d : [ ", j+1)
			for _, player := range game.Team1 {
				fmt.Fprintf(f, "%d ", player+1)
			}
			fmt.Fprintf(f, "] | [ ")
			for _, player := range game.Team2 {
				fmt.Fprintf(f, "%d ", player+1)
			}
			fmt.Fprintf(f, "]\n")
		}
		fmt.Println("--------------------------")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	for i := 17; i < 41; i++ {
		f, err := os.Create("output_" + strconv.Itoa(i) + ".txt")
		if err != nil {
			panic(err)
		}
		var t tournament.Tournament
		var best tournament.Tournament
		bcollision := math.MaxInt
		for j := range MAX_GENERATIONS {
			t = tournament.DrawRandomTournament(i, NB_ROUNDS, MAX_FIELD)
			col := t.CountCollision()
			if bcollision > col {
				best = t
				bcollision = col
			}
			if bcollision == 0 {
				best = t
				break
			}
			if j%100000 == 0 {
				fmt.Printf("Nombre Joueur : %d -> Generation %d | Best Found : %d\n", i, j, bcollision/2)
			}
		}
		Display(f, best)
		fmt.Fprintf(f, "Nb Collision : %d\n", bcollision/2)
		f.Close()
	}
}
