package tournament

import "maps"

type Round []Game

const CREATE_NEW_GAME = -1

func (round Round) Clone() Round {
	cloned := make(Round, len(round))
	for i, game := range round {
		cloned[i] = game.Clone()
	}
	return cloned
}

func (round Round) genEncounteredMatrix() (PlayersTimeEncountered, PlayersTimeEncountered) {
	resPlayWith := make(PlayersTimeEncountered)
	resPlayAgainst := make(PlayersTimeEncountered)
	for _, game := range round {
		gameMatrixPlayWith, gameMatrixPlayAgainst := game.genEncounteredMatrix()
		for player := range gameMatrixPlayWith {
			if _, exist := resPlayWith[player]; !exist {
				resPlayWith[player] = make(map[Player]int)
			}
			maps.Copy(resPlayWith[player], gameMatrixPlayWith[player])
		}
		for player := range gameMatrixPlayAgainst {
			if _, exist := resPlayAgainst[player]; !exist {
				resPlayAgainst[player] = make(map[Player]int)
			}
			maps.Copy(resPlayAgainst[player], gameMatrixPlayAgainst[player])
		}
	}
	return resPlayWith, resPlayAgainst
}

func (round Round) CountCollision(matchPlayedWith PlayersTimeEncountered, matchPlayedAgainst PlayersTimeEncountered) int {
	res := 0
	for _, game := range round {
		res += game.CountCollision(matchPlayedWith, matchPlayedAgainst)
	}
	return res
}

func (round *Round) PlayerCanBePlacedInGamesIndex(totalPlayer int, maxField int) []int {
	if len(*round) == 0 {
		return []int{CREATE_NEW_GAME}
	}

	// On regarde quelles sont les parties libre en doublette et en triplette
	freePlaceInDoublette := []int{}
	freePlaceInTriplette := []int{}
	for idxGame, game := range *round {
		if !game.IsFull() {
			freePlaceInTriplette = append(freePlaceInTriplette, idxGame)
		}
		if !game.IsFullForDoublette() {
			freePlaceInDoublette = append(freePlaceInDoublette, idxGame)
		}
	}
	if len(*round) < totalPlayer/4 && len(*round) < maxField {
		freePlaceInDoublette = append(freePlaceInDoublette, CREATE_NEW_GAME)
		freePlaceInTriplette = append(freePlaceInTriplette, CREATE_NEW_GAME)
	}
	if len(freePlaceInDoublette) > 0 {
		return freePlaceInDoublette
	} else if len(freePlaceInTriplette) > 0 {
		restrictTriplette := []int{}
		for _, i := range freePlaceInTriplette {
			if (*round)[i].CountPlacedPlayer() == 5 {
				restrictTriplette = append(restrictTriplette, i)
			}
		}
		if (totalPlayer-round.CountPlacedPlayer())%2 == 1 && len(restrictTriplette) > 0 {
			return restrictTriplette
		}
		return freePlaceInTriplette
	} else {
		panic("player cant be placed !")
	}
}

func (round *Round) PlacePlayerLazy(p Player, remainPlayers int, maxField int) {
	// Aucunes partie ? On en créer une
	if len(*round) == 0 {
		newGame := Game{}
		newGame.PlacePlayerInGame(p)
		*round = append(*round, newGame)
		return
	}

	// On prend la dernière partie
	lastGame := &(*round)[len(*round)-1]
	// Si la partie est pleine pour une doublette (mais qu'on est pas sur le cas d'une triplette)
	// On ajoute le joueur dans une nouvelle
	if lastGame.IsFullForDoublette() && remainPlayers >= 4 && len(*round) < maxField {
		newGame := Game{}
		newGame.PlacePlayerInGame(p)
		*round = append(*round, newGame)
		return
	}

	// On regarde si on est pas dans un cas où faire une triplette est mieux
	if lastGame.IsFullForDoublette() && (remainPlayers < 4 || len(*round) >= maxField) {
		// Dans ce cas on parcours les parties déjà présente, jusqu'à trouver une ou le jour peut être placé
		i := 0
		game := &(*round)[i]

		for game.IsFull() {
			i += 1
			if i >= len(*round) {
				panic("no game found to place player. Add more fields")
			}
			game = &(*round)[i]
		}
		// On peut ajouter le joueur dans cette partie
		game.PlacePlayerInGame(p)
		return
	}

	// Dans les autres cas le joueur peut être placé sur la partie actuelle
	lastGame.PlacePlayerInGame(p)
}

func (round Round) CountPlacedPlayer() int {
	res := 0
	for _, game := range round {
		res += game.CountPlacedPlayer()
	}
	return res
}

func (round Round) ReArranged() Round {
	threeVsThree := make(Round, 0)
	threeVsTwo := make(Round, 0)
	twoVsTwo := make(Round, 0)
	errors := make(Round, 0)
	for _, game := range round {
		switch len(game.Team1) + len(game.Team2) {
		case 4:
			twoVsTwo = append(twoVsTwo, game)
		case 5:
			res := game.Clone()
			if len(game.Team1) < len(game.Team2) {
				tmp := game.Team1
				res.Team1 = res.Team2
				res.Team2 = tmp
			}
			threeVsTwo = append(threeVsTwo, res)
		case 6:
			threeVsThree = append(threeVsThree, game)
		default:
			errors = append(errors, game)
		}
	}
	return append(threeVsThree, append(threeVsTwo, append(twoVsTwo, errors...)...)...)
}
