package tournament

import "errors"

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

func (game Game) genEncounteredMatrix() PlayersTimeEncountered {
	res := make(PlayersTimeEncountered)
	for _, player := range game.Team1 {
		res[player] = make(map[Player]int)
		for _, opponent := range game.Team1 {
			if player != opponent {
				res[player][opponent] = 1
			}
		}
		for _, opponent := range game.Team2 {
			res[player][opponent] = 1
		}
	}
	for _, player := range game.Team2 {
		res[player] = make(map[Player]int)
		for _, opponent := range game.Team1 {
			res[player][opponent] = 1
		}
		for _, opponent := range game.Team2 {
			if player != opponent {
				res[player][opponent] = 1
			}
		}
	}
	return res
}

func (game Game) CollisionFoundIfIplaceThisPlayer(p Player, matchPlayed PlayersTimeEncountered) int {
	res := 0
	for _, opponent := range game.Team1 {
		if val, exist := matchPlayed[p][opponent]; exist {
			res += val
		}
	}
	for _, opponent := range game.Team2 {
		if val, exist := matchPlayed[p][opponent]; exist {
			res += val
		}
	}
	return res
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

func (game Game) CountCollision(matchPlayed PlayersTimeEncountered) int {
	res := 0
	for _, player := range game.Team1 {
		for _, opponent := range game.Team1 {
			if val, exist := matchPlayed[player][opponent]; exist && player != opponent {
				res += val
			}
		}
		for _, opponent := range game.Team2 {
			if val, exist := matchPlayed[player][opponent]; exist {
				res += val
			}
		}
	}
	for _, player := range game.Team2 {
		for _, opponent := range game.Team1 {
			if val, exist := matchPlayed[player][opponent]; exist {
				res += val
			}
		}
		for _, opponent := range game.Team2 {
			if val, exist := matchPlayed[player][opponent]; exist && player != opponent {
				res += val
			}
		}
	}
	return res
}

func (game Game) CountPlacedPlayer() int {
	return len(game.Team1) + len(game.Team2)
}
