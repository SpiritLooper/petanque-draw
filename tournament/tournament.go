package tournament

import (
	"fmt"
	"maps"
)

type Tournament []Round

type PlayerSet map[Player]struct{}
type PlayersTimeEncountered map[Player]map[Player]int

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

func (tournament Tournament) GenEncounteredMatrix() (PlayersTimeEncountered, PlayersTimeEncountered) {
	resPlayWith := make(PlayersTimeEncountered)
	resPlayAgainst := make(PlayersTimeEncountered)
	for _, round := range tournament {
		roundMatrixPlayWith, roundMatrixPlayAgainst := round.genEncounteredMatrix()
		for player, _ := range roundMatrixPlayWith {
			if _, exist := resPlayWith[player]; !exist {
				resPlayWith[player] = make(map[Player]int)
			}
			maps.Copy(resPlayWith[player], roundMatrixPlayWith[player])
		}
		for player, _ := range roundMatrixPlayAgainst {
			if _, exist := resPlayAgainst[player]; !exist {
				resPlayAgainst[player] = make(map[Player]int)
			}
			maps.Copy(resPlayAgainst[player], roundMatrixPlayAgainst[player])
		}
	}
	return resPlayWith, resPlayAgainst
}

func (tournament Tournament) countCollision(matchPlayedWith PlayersTimeEncountered, matchPlayedAgainst PlayersTimeEncountered) int {
	res := 0
	for _, round := range tournament {
		res += round.CountCollision(matchPlayedWith, matchPlayedAgainst)
	}
	return res
}

func (tournament Tournament) CountCollision() int {
	res := 0
	for i, round := range tournament[1:] {
		matchPlayedWith, matchPlayedAgainst := tournament[:i+1].GenEncounteredMatrix()
		res += round.CountCollision(matchPlayedWith, matchPlayedAgainst)
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
	playedMatrixWith := make([][]int, n, n)
	playedMatrixAgainst := make([][]int, n, n)
	for i := range playedMatrixWith {
		playedMatrixWith[i] = make([]int, n, n)
		playedMatrixAgainst[i] = make([]int, n, n)
		res[i] = make(PlayerSet)
		for j := range playedMatrixWith[i] {
			playedMatrixWith[i][j] = 0
		}
		for j := range playedMatrixAgainst[i] {
			playedMatrixAgainst[i][j] = 0
		}
	}

	for _, round := range tournament {
		for _, game := range round {
			for _, player := range game.Team1 {
				for _, opponent := range game.Team1 {
					if player != opponent {
						playedMatrixWith[player][opponent] += 1
					}
				}
				for _, opponent := range game.Team2 {
					playedMatrixAgainst[player][opponent] += 1
				}
			}
			for _, player := range game.Team2 {
				for _, opponent := range game.Team1 {
					playedMatrixAgainst[player][opponent] += 1
				}
				for _, opponent := range game.Team2 {
					if player != opponent {
						playedMatrixWith[player][opponent] += 1
					}
				}
			}
		}
	}

	for i := range playedMatrixWith {
		for j := range playedMatrixWith[i] {
			if playedMatrixWith[i][j] >= 2 {
				res[i][Player(j)] = struct{}{}
			}
			if playedMatrixAgainst[i][j] >= 2 {
				res[i][Player(j)] = struct{}{}
			}
		}
	}

	return res
}
