package tournament

import (
	"fmt"
	"maps"
)

type Tournament []Round

type PlayerSet map[Player]struct{}
type PlayerEncountered map[Player]PlayerSet

func (t *Tournament) Display() {
	for i, rounds := range *t {
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

func (tournament Tournament) GenEncounteredMatrix() PlayerEncountered {
	res := make(PlayerEncountered)
	for _, round := range tournament {
		roundMatrix := round.genEncounteredMatrix()
		for player, _ := range roundMatrix {
			if _, exist := res[player]; !exist {
				res[player] = PlayerSet{}
			}
			maps.Copy(res[player], roundMatrix[player])
		}
	}
	return res
}

func (tournament Tournament) countCollision(matchPlayed PlayerEncountered) int {
	res := 0
	for _, round := range tournament {
		res += round.CountCollision(matchPlayed)
	}
	return res
}

func (tournament Tournament) CountCollision() int {
	res := 0
	for i, round := range tournament[1:] {
		matchPlayed := tournament[:i+1].GenEncounteredMatrix()
		res += round.CountCollision(matchPlayed)
	}
	return res / 2
}

func (tournament Tournament) nbPlayer() int {
	maxId := 0
	for _, game := range tournament[0] {
		for _, player := range game.Team1 {
			maxId = max(maxId, int(player))
		}
		for _, player := range game.Team2 {
			maxId = max(maxId, int(player))
		}
	}
	return maxId + 1
}

func (tournament Tournament) GetCollision() []PlayerSet {
	n := tournament.nbPlayer()
	res := make([]PlayerSet, n)
	playedMatrix := make([][]int, n, n)
	for i := range playedMatrix {
		playedMatrix[i] = make([]int, n, n)
		res[i] = make(PlayerSet)
		for j := range playedMatrix[i] {
			playedMatrix[i][j] = 0
		}
	}

	for _, round := range tournament {
		for _, game := range round {
			for _, player := range game.Team1 {
				for _, opponent := range game.Team1 {
					if player != opponent {
						playedMatrix[player][opponent] += 1
					}
				}
				for _, opponent := range game.Team2 {
					playedMatrix[player][opponent] += 1
				}
			}
			for _, player := range game.Team2 {
				for _, opponent := range game.Team1 {
					playedMatrix[player][opponent] += 1
				}
				for _, opponent := range game.Team2 {
					if player != opponent {
						playedMatrix[player][opponent] += 1
					}
				}
			}
		}
	}

	for i := range playedMatrix {
		for j := range playedMatrix[i] {
			if playedMatrix[i][j] >= 2 {
				res[i][Player(j)] = struct{}{}
			}
		}
	}

	return res
}
