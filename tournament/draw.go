package tournament

import (
	"petanque-draw/utils/intset"
)

type playerSet map[Player]struct{}
type playerEncountered map[Player]playerSet

func lastRoundDrawConsistent(actualDraw Round, previousDraws Tournament, numberPlayer int) bool {
	if len(previousDraws) == 0 {
		return true
	}

	var playerEncountered = make(map[Player]playerSet)

	// Setup Encountered Matrix
	for _, game := range actualDraw {
		for _, player := range game.Team1 {
			playerEncountered[player] = playerSet{}
			for _, otherPlayer := range game.Team1 {
				if player != otherPlayer {
					playerEncountered[player][otherPlayer] = struct{}{}
				}
			}
			for _, otherPlayer := range game.Team2 {
				if player != otherPlayer {
					playerEncountered[player][otherPlayer] = struct{}{}
				}
			}
		}
		for _, player := range game.Team2 {
			playerEncountered[player] = playerSet{}
			for _, otherPlayer := range game.Team1 {
				if player != otherPlayer {
					playerEncountered[player][otherPlayer] = struct{}{}
				}
			}
			for _, otherPlayer := range game.Team2 {
				if player != otherPlayer {
					playerEncountered[player][otherPlayer] = struct{}{}
				}
			}
		}
	}

	// Checking double encontered
	for _, draw := range previousDraws {
		for _, game := range draw {
			for _, player := range game.Team1 {
				for _, otherPlayer := range game.Team1 {
					_, exist := playerEncountered[player][otherPlayer]
					if exist {
						return false
					}
				}
				for _, otherPlayer := range game.Team2 {
					if player != otherPlayer {
						_, exist := playerEncountered[player][otherPlayer]
						if exist {
							return false
						}
					}
				}
			}
			for _, player := range game.Team2 {
				for _, otherPlayer := range game.Team1 {
					_, exist := playerEncountered[player][otherPlayer]
					if exist {
						return false
					}
				}
				for _, otherPlayer := range game.Team2 {
					if player != otherPlayer {
						_, exist := playerEncountered[player][otherPlayer]
						if exist {
							return false
						}
					}
				}
			}
		}
	}

	return true
}

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
		var round Round
		for {
			round = Round{}
			set := intset.CreateIntSet(nbPlayer)

			for !intset.IsEmpty(set) {
				player := drawPlayer(&set)
				placePlayerInRound(player, &round, intset.Count(set)+1, maxField)
			}

			if lastRoundDrawConsistent(round, tournament, nbPlayer) {
				break
			}
		}
		tournament = append(tournament, round)
	}
	return tournament
}
