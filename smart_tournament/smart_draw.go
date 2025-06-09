package smart_tournament

import "petanque-draw/utils/intset"

func DrawSmartTournament(nbPlayer int, nbRound int, maxField int) Tournament {
	var tournament Tournament
	tournament.playerEncountered = make(playerEncountered)
	for range nbRound {
		round := Round{}
		set := intset.CreateIntSet(nbPlayer)

		for !set.IsEmpty() {
			integer, _ := set.RandomPop()
			round.placePlayer(Player(integer), set.Count()+1, &tournament.playerEncountered, maxField, nbPlayer)
		}

		tournament.Rounds = append(tournament.Rounds, round)
	}
	return tournament
}
