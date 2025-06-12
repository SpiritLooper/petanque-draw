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

func (round Round) genEncounteredMatrix() PlayersTimeEncountered {
	res := make(PlayersTimeEncountered)
	for _, game := range round {
		gameMatrix := game.genEncounteredMatrix()
		for player, _ := range gameMatrix {
			if _, exist := res[player]; !exist {
				res[player] = make(map[Player]int)
			}
			maps.Copy(res[player], gameMatrix[player])
		}
	}
	return res
}

func (round Round) CountCollision(matchPlayed PlayersTimeEncountered) int {
	res := 0
	for _, game := range round {
		res += game.CountCollision(matchPlayed)
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
