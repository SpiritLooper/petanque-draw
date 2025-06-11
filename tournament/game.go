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

func (game Game) CollisionFoundIfIplaceThisPlayer(p Player, matchPlayed PlayerEncountered) int {
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

func (game Game) CountCollision(matchPlayed PlayerEncountered) int {
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
