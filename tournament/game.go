package tournament

import (
	"errors"
	"math"
)

type Game struct {
	Team1 []Player
	Team2 []Player
}
type Player int

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

func (game Game) IsFull() bool {
	return len(game.Team1) >= 3 && len(game.Team2) >= 3
}

func (game Game) IsFullForDoublette() bool {
	return len(game.Team1) >= 2 && len(game.Team2) >= 2
}

func (game Game) IsContainsTriplette() bool {
	return len(game.Team1) >= 3 || len(game.Team2) >= 3
}

func (game *Game) PlacePlayerInGame(p Player) error {
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

func (game *Game) PlacePlayerInGameGreedyTeam(p Player, playWith PlayersTimeEncountered, playAgainst PlayersTimeEncountered) error {
	if game.IsFull() {
		return errors.New("can't place new player. game is full")
	}

	if len(game.Team1) >= 3 {
		(*game).Team2 = append((*game).Team2, p)
		return nil
	}
	if len(game.Team2) >= 3 {
		(*game).Team1 = append((*game).Team1, p)
		return nil
	}

	colTeam1 := game.countCollisionIfPlaceTeam1(p, playWith, playAgainst)
	colTeam2 := game.countCollisionIfPlaceTeam2(p, playWith, playAgainst)

	if colTeam1 > colTeam2 {
		(*game).Team2 = append((*game).Team2, p)
		return nil
	} else if colTeam1 < colTeam2 {
		(*game).Team1 = append((*game).Team1, p)
		return nil
	}

	if len(game.Team1) > len(game.Team2) {
		(*game).Team2 = append((*game).Team2, p)
	} else {
		(*game).Team1 = append((*game).Team1, p)
	}
	return nil
}

func (game Game) countCollisionIfPlaceTeam1(p Player, playWith PlayersTimeEncountered, playAgainst PlayersTimeEncountered) int {
	if len(game.Team1) >= 3 {
		return math.MaxInt
	}
	res := 0
	for _, playerWith := range game.Team1 {
		if playWith[p][playerWith] >= 1 {
			res += 1
		}
	}
	for _, playerAgainst := range game.Team2 {
		if playAgainst[p][playerAgainst] >= 1 {
			res += 1
		}
	}
	return res
}

func (game Game) countCollisionIfPlaceTeam2(p Player, playWith PlayersTimeEncountered, playAgainst PlayersTimeEncountered) int {
	if len(game.Team2) >= 3 {
		return math.MaxInt
	}
	res := 0
	for _, playerWith := range game.Team2 {
		if playWith[p][playerWith] >= 1 {
			res += 1
		}
	}
	for _, playerAgainst := range game.Team1 {
		if playAgainst[p][playerAgainst] >= 1 {
			res += 1
		}
	}
	return res
}

func (game Game) genEncounteredMatrix() (PlayersTimeEncountered, PlayersTimeEncountered) {
	playerPlayWith := make(PlayersTimeEncountered)
	playerPlayAgainst := make(PlayersTimeEncountered)
	for _, player := range game.Team1 {
		playerPlayWith[player] = make(map[Player]int)
		playerPlayAgainst[player] = make(map[Player]int)
		for _, opponent := range game.Team1 {
			if player != opponent {
				playerPlayWith[player][opponent] = 1
			}
		}
		for _, opponent := range game.Team2 {
			playerPlayAgainst[player][opponent] = 1
		}
	}
	for _, player := range game.Team2 {
		playerPlayWith[player] = make(map[Player]int)
		playerPlayAgainst[player] = make(map[Player]int)
		for _, opponent := range game.Team1 {
			playerPlayAgainst[player][opponent] = 1
		}
		for _, opponent := range game.Team2 {
			if player != opponent {
				playerPlayWith[player][opponent] = 1
			}
		}
	}
	return playerPlayWith, playerPlayAgainst
}

func (game Game) CollisionFoundIfIplaceThisPlayer(p Player, matchPlayedWith PlayersTimeEncountered, matchPlayedAgainst PlayersTimeEncountered) int {
	resTeam1 := 0
	resTeam2 := 0
	for _, opponent := range game.Team1 {
		if val, exist := matchPlayedWith[p][opponent]; exist {
			resTeam1 += val
		}
	}
	for _, opponent := range game.Team2 {
		if val, exist := matchPlayedAgainst[p][opponent]; exist {
			resTeam1 += val
		}
	}
	for _, opponent := range game.Team2 {
		if val, exist := matchPlayedWith[p][opponent]; exist {
			resTeam2 += val
		}
	}
	for _, opponent := range game.Team1 {
		if val, exist := matchPlayedAgainst[p][opponent]; exist {
			resTeam2 += val
		}
	}
	return min(resTeam1, resTeam2)
}

func (game Game) IsPlayerPlayWith(p1 Player, p2 Player) bool {
	p1InTeam1 := false
	p1InTeam2 := false
	p2InTeam1 := false
	p2InTeam2 := false
	for _, p := range game.Team1 {
		if p == p1 {
			p1InTeam1 = true
		}
		if p == p2 {
			p2InTeam1 = true
		}
	}
	for _, p := range game.Team2 {
		if p == p1 {
			p1InTeam2 = true
		}
		if p == p2 {
			p2InTeam2 = true
		}
	}

	return (p1InTeam1 && p2InTeam1) || (p1InTeam2 && p2InTeam2)
}

func (game Game) IsPlayerPlayAgainst(p1 Player, p2 Player) bool {
	p1InTeam1 := false
	p1InTeam2 := false
	p2InTeam1 := false
	p2InTeam2 := false
	for _, p := range game.Team1 {
		if p == p1 {
			p1InTeam1 = true
		}
		if p == p2 {
			p2InTeam1 = true
		}
	}
	for _, p := range game.Team2 {
		if p == p1 {
			p1InTeam2 = true
		}
		if p == p2 {
			p2InTeam2 = true
		}
	}

	return (p1InTeam1 && p2InTeam2) || (p1InTeam2 && p2InTeam1)
}

func (game Game) CountCollision(matchPlayedWith PlayersTimeEncountered, matchPlayedAgainst PlayersTimeEncountered) int {
	res := 0
	for _, player := range game.Team1 {
		for _, opponent := range game.Team1 {
			if val, exist := matchPlayedWith[player][opponent]; exist && player != opponent {
				if val >= 1 {
					res += 1
				}
			}
		}
		for _, opponent := range game.Team2 {
			if val, exist := matchPlayedAgainst[player][opponent]; exist {
				if val >= 1 {
					res += 1
				}
			}
		}
	}
	for _, player := range game.Team2 {
		for _, opponent := range game.Team1 {
			if val, exist := matchPlayedAgainst[player][opponent]; exist {
				if val >= 1 {
					res += 1
				}
			}
		}
		for _, opponent := range game.Team2 {
			if val, exist := matchPlayedWith[player][opponent]; exist && player != opponent {
				if val >= 1 {
					res += 1
				}
			}
		}
	}
	return res
}

func (game Game) CountPlacedPlayer() int {
	return len(game.Team1) + len(game.Team2)
}
