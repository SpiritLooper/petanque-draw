package tournament

import (
	"errors"
	"fmt"
	"maps"
)

type Tournament []Round
type Round []Game
type Game struct {
	Team1 []Player
	Team2 []Player
}
type Player int

type PlayerSet map[Player]struct{}
type PlayerEncountered map[Player]PlayerSet

func (g Game) Clone() Game {
	newTeam1 := make([]Player, len(g.Team1))
	copy(newTeam1, g.Team1)
	newTeam2 := make([]Player, len(g.Team2))
	copy(newTeam2, g.Team2)
	return Game{
		Team1: newTeam1,
		Team2: newTeam2,
	}
}

func (round Round) Clone() Round {
	cloned := make(Round, len(round))
	for i, game := range round {
		cloned[i] = game.Clone()
	}
	return cloned
}

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

func (game Game) IsFull() bool {
	return len(game.Team1) >= 3 && len(game.Team2) >= 3
}

func (game Game) IsFullForDoublette() bool {
	return len(game.Team1) >= 2 && len(game.Team2) >= 2
}

func (game *Game) placePlayerInGame(p Player) error {
	if game.IsFull() {
		return errors.New("can't place new player. game is full")
	}

	if len(game.Team1) > len(game.Team2) {
		(*game).Team2 = append((*game).Team2, p)
	} else {
		(*game).Team1 = append((*game).Team1, p)
	}
	return nil
}

func (game Game) genEncounteredMatrix() PlayerEncountered {
	res := make(PlayerEncountered)
	for _, player := range game.Team1 {
		res[player] = make(PlayerSet)
		for _, opponent := range game.Team1 {
			if player != opponent {
				res[player][opponent] = struct{}{}
			}
		}
		for _, opponent := range game.Team2 {
			res[player][opponent] = struct{}{}
		}
	}
	for _, player := range game.Team2 {
		res[player] = make(PlayerSet)
		for _, opponent := range game.Team1 {
			res[player][opponent] = struct{}{}
		}
		for _, opponent := range game.Team2 {
			if player != opponent {
				res[player][opponent] = struct{}{}
			}
		}
	}
	return res
}

func (round Round) genEncounteredMatrix() PlayerEncountered {
	res := make(PlayerEncountered)
	for _, game := range round {
		gameMatrix := game.genEncounteredMatrix()
		for player, _ := range gameMatrix {
			if _, exist := res[player]; !exist {
				res[player] = PlayerSet{}
			}
			maps.Copy(res[player], gameMatrix[player])
		}
	}
	return res
}

func (tournament Tournament) genEncounteredMatrix() PlayerEncountered {
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

func (game Game) collisionFoundIfIplaceThisPlayer(p Player, matchPlayed PlayerEncountered) int {
	res := 0
	for _, opponent := range game.Team1 {
		if _, exist := matchPlayed[p][opponent]; exist {
			res++
		}
	}
	for _, opponent := range game.Team2 {
		if _, exist := matchPlayed[p][opponent]; exist {
			res++
		}
	}
	return res
}

func (game Game) countCollision(matchPlayed PlayerEncountered) int {
	res := 0
	for _, player := range game.Team1 {
		for _, opponent := range game.Team1 {
			if _, exist := matchPlayed[player][opponent]; exist && player != opponent {
				res++
			}
		}
		for _, opponent := range game.Team2 {
			if _, exist := matchPlayed[player][opponent]; exist {
				res++
			}
		}
	}
	for _, player := range game.Team2 {
		for _, opponent := range game.Team1 {
			if _, exist := matchPlayed[player][opponent]; exist {
				res++
			}
		}
		for _, opponent := range game.Team2 {
			if _, exist := matchPlayed[player][opponent]; exist && player != opponent {
				res++
			}
		}
	}
	return res
}

func (round Round) countCollision(matchPlayed PlayerEncountered) int {
	res := 0
	for _, game := range round {
		res += game.countCollision(matchPlayed)
	}
	return res
}

func (tournament Tournament) countCollision(matchPlayed PlayerEncountered) int {
	res := 0
	for _, round := range tournament {
		res += round.countCollision(matchPlayed)
	}
	return res
}

func (tournament Tournament) CountCollision() int {
	res := 0
	for i, round := range tournament[1:] {
		matchPlayed := tournament[:i+1].genEncounteredMatrix()
		res += round.countCollision(matchPlayed)
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
