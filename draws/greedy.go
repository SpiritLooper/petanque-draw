package draws

import (
	"math"
	"math/rand"
	"petanque-draw/tournament"
)

func placePlayerInRoundGreedy(
	player tournament.Player, playerEverEncounterd tournament.PlayersTimeEncountered,
	totalPlayer int, maxField int, actualRound tournament.Round) tournament.Round {
	if len(actualRound) == 0 {
		game := tournament.Game{}
		game.PlacePlayerInGame(player)
		return append(actualRound, game)
	}

	resRound := actualRound.Clone()

	canCreateGame := false
	bIdx := -1
	bCollision := math.MaxInt
	var bGameToPlace tournament.Game
	// Cas ou les parties sont déjà existante
	for _, idxGame := range actualRound.PlayerCanBePlacedInGamesIndex(totalPlayer, maxField) {
		if idxGame == tournament.CREATE_NEW_GAME {
			canCreateGame = true
			continue
		}
		game := resRound[idxGame]

		if !game.IsFull() {
			col := game.CollisionFoundIfIplaceThisPlayer(player, playerEverEncounterd)
			// Meilleur cas
			if !game.IsFullForDoublette() && col == 0 {
				game.PlacePlayerInGame(player)
				resRound[idxGame] = game
				bIdx = idxGame
				return resRound
			}

			if bCollision > col || (bCollision == col) && rand.Int()%2 == 0 {
				bCollision = col
				bGameToPlace = game
				bIdx = idxGame
			}
		}
	}

	if canCreateGame && bCollision > 0 {
		bGameToPlace = tournament.Game{}
		bGameToPlace.PlacePlayerInGame(player)
		resRound = append(resRound, bGameToPlace)
		return resRound
	}

	bGameToPlace.PlacePlayerInGame(player)
	resRound[bIdx] = bGameToPlace
	return resRound
}

func DrawRoundGreed(playerEverEncountered tournament.PlayersTimeEncountered, players []int, maxField int) tournament.Round {
	var round tournament.Round
	for _, iPlayer := range players {
		round = placePlayerInRoundGreedy(tournament.Player(iPlayer), playerEverEncountered, len(players), maxField, round)
	}
	return round
}
