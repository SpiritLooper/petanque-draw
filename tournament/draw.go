package tournament

import (
	"petanque-draw/utils/intset"
)

func drawPlayer(set *intset.IntSet) Player {
	number, err := intset.RandomPop(set)
	if err != nil {
		panic(err)
	}
	return Player(number)
}

func DrawTournament(nbPlayer int, nbRound int, maxField int) Tournament {
	var tournament Tournament
	for range nbRound {
		set := intset.CreateIntSet(nbPlayer)

		var round Round
		for !intset.IsEmpty(set) {
			player := drawPlayer(&set)
			placePlayerInRound(player, &round, intset.Count(set), maxField)
		}

		tournament = append(tournament, round)
	}

	return tournament
}
