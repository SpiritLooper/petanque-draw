package draws

import "petanque-draw/tournament"

func DrawTournamentGreedy(nbPlayer int, nbRound int, maxField int) tournament.Tournament {
	tournament := make(tournament.Tournament, 0, nbRound)
	tournament = append(tournament, drawRandomRound(nbPlayer, maxField))
	for range nbRound - 1 {

	}
	return tournament
}
