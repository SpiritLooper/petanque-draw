package main

import (
	"petanque-draw/tournament"
)

const (
	MAX_FIELD = 8
	NB_ROUNDS = 4
)

func main() {
	tournament := tournament.DrawRandomTournament(32, NB_ROUNDS, MAX_FIELD)
	tournament.Display()
}
