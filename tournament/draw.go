package tournament

import (
	"petanque-draw/utils/intset"
)

type playerSet map[Player]struct{}
type playerEncountered map[Player]playerSet

func (t Tournament) CountCollision() int {
	var playerEncountered = make(map[Player]playerSet)

	// Setup Encountered Matrix
	for _, game := range t[0] {
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
	nbCollision := 0
	// Checking double encontered
	for _, draw := range t[1:] {
		for _, game := range draw {
			for _, player := range game.Team1 {
				for _, otherPlayer := range game.Team1 {
					_, exist := playerEncountered[player][otherPlayer]
					if exist {
						nbCollision += 1
					}
				}
				for _, otherPlayer := range game.Team2 {
					if player != otherPlayer {
						_, exist := playerEncountered[player][otherPlayer]
						if exist {
							nbCollision += 1
						}
					}
				}
			}
			for _, player := range game.Team2 {
				for _, otherPlayer := range game.Team1 {
					_, exist := playerEncountered[player][otherPlayer]
					if exist {
						nbCollision += 1
					}
				}
				for _, otherPlayer := range game.Team2 {
					if player != otherPlayer {
						_, exist := playerEncountered[player][otherPlayer]
						if exist {
							nbCollision += 1
						}
					}
				}
			}
		}
	}

	return nbCollision
}

func drawPlayer(set *intset.IntSet) Player {
	number, err := set.RandomPop()
	if err != nil {
		panic(err)
	}
	return Player(number)
}

func DrawRandomTournament(nbPlayer int, nbRound int, maxField int) Tournament {
	var tournament Tournament
	for range nbRound {
		round := Round{}
		set := intset.CreateIntSet(nbPlayer)

		for !set.IsEmpty() {
			player := drawPlayer(&set)
			placePlayerInRound(player, &round, set.Count()+1, maxField)
		}

		tournament = append(tournament, round)
	}
	return tournament
}
