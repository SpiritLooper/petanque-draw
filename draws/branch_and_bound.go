package draws

import (
	"fmt"
	"math"
	"petanque-draw/tournament"
	"slices"
)

const MAX_DEPTH_LOG = 5
const MAX_COLLISION_FOUND = 0

type RoundState struct {
	playersNotInGames []int
	round             tournament.Round
}

func (rs RoundState) createRoundUniqEncouterRecur(
	maxField int, playerEverEncountered tournament.PlayersTimeEncountered,
	bestCollision *int, bestResult *RoundState, depth int,
) {
	// Si on a plus de tirage, on retourne le résultat de collision
	if len(rs.playersNotInGames) <= 0 {
		*bestCollision = rs.round.CountCollision(playerEverEncountered)
		*bestResult = rs
		return
	}

	// Si il reste des joueurs
	// On regarde pour le placer dans une partie
	// --- Cas pas de partie créé
	if len(rs.round) <= 0 {
		// Pour chaque joueur non placés
		for i, iPlayer := range rs.playersNotInGames {
			// On créer sa partie avec lui uniquement
			newGame := tournament.Game{Team1: []tournament.Player{tournament.Player(iPlayer)}}
			// On le retire des joueurs déjà placés
			newPlayersNotPlaced := slices.Concat(rs.playersNotInGames[:i], rs.playersNotInGames[i+1:])
			// On prépare sa nouvelle feuille
			colision := math.MaxInt
			var leafResult RoundState
			leaf := RoundState{playersNotInGames: newPlayersNotPlaced, round: append(rs.round.Clone(), newGame)}
			// On parcourt l'arbre
			if depth < MAX_DEPTH_LOG {
				fmt.Printf("I'm in depth %d : Loop %d/%d\n", depth, i, len(rs.playersNotInGames))
			}
			leaf.createRoundUniqEncouterRecur(maxField, playerEverEncountered, &colision, &leafResult, depth+1)
			// On regarde les résultat
			if colision == 0 { // Meilleur qu'on puisse faire ça ne sert à rien de continuer
				*bestCollision = 0
				*bestResult = leafResult
				return
			} else if colision < *bestCollision { // Mise à jour du meilleur résultat
				*bestCollision = colision
				*bestResult = leafResult
			}
		}
		return
	}

	// On prend la dernière partie dispo
	lastGame := rs.round[len(rs.round)-1]

	// Si la partie est pleine pour une doublette (mais qu'on est pas sur le cas d'une triplette)
	// On ajoute le joueur dans une nouvelle
	if lastGame.IsFullForDoublette() && len(rs.playersNotInGames) >= 4 && len(rs.round) < maxField {
		for i, iPlayer := range rs.playersNotInGames {
			newGame := tournament.Game{Team1: make([]tournament.Player, 0, 3), Team2: make([]tournament.Player, 0, 3)}
			newGame.PlacePlayerInGame(tournament.Player(iPlayer))
			newPlayersNotPlaced := slices.Concat(rs.playersNotInGames[:i], rs.playersNotInGames[i+1:])

			// On prépare sa nouvelle feuille
			colision := math.MaxInt
			var leafResult RoundState
			leaf := RoundState{playersNotInGames: newPlayersNotPlaced, round: append(rs.round.Clone(), newGame.Clone())}
			// On parcourt l'arbre
			if depth < MAX_DEPTH_LOG {
				fmt.Printf("I'm in depth %d : Loop %d/%d\n", depth, i, len(rs.playersNotInGames))
			}
			leaf.createRoundUniqEncouterRecur(maxField, playerEverEncountered, &colision, &leafResult, depth+1)
			// On regarde les résultat
			if colision == 0 { // Meilleur qu'on puisse faire ça ne sert à rien de continuer
				*bestCollision = 0
				*bestResult = leafResult
				return
			} else if colision < *bestCollision { // Mise à jour du meilleur résultat
				*bestCollision = colision
				*bestResult = leafResult
			}
		}
		// Pas besoin de continuer
		return
	}

	// On regarde si on est pas dans un cas où faire une triplette est mieux
	if lastGame.IsFullForDoublette() && (len(rs.playersNotInGames) < 4 || len(rs.round) >= maxField) {
		// Dans ce cas on parcours les parties déjà présente, jusqu'à trouver une ou le jour peut être placé
		i := 0
		game := rs.round[i]

		for game.IsFull() {
			i += 1
			if i >= len(rs.round) {
				panic("no game found to place player. Add more fields")
			}
			game = rs.round[i]
		}

		// On peut ajouter le joueur dans cette partie
		for i, iPlayer := range rs.playersNotInGames {
			newPlayersNotPlaced := slices.Concat(rs.playersNotInGames[:i], rs.playersNotInGames[i+1:])
			newRound := rs.round.Clone()
			newGame := rs.round[i]
			newGame.PlacePlayerInGame(tournament.Player(iPlayer))
			newRound[i] = newGame.Clone()
			// On prépare sa nouvelle feuille
			colision := math.MaxInt
			var leafResult RoundState
			leaf := RoundState{playersNotInGames: newPlayersNotPlaced, round: newRound}
			// On parcourt l'arbre
			if depth < MAX_DEPTH_LOG {
				fmt.Printf("I'm in depth %d : Loop %d/%d\n", depth, i, len(rs.playersNotInGames))
			}
			leaf.createRoundUniqEncouterRecur(maxField, playerEverEncountered, &colision, &leafResult, depth+1)
			// On regarde les résultat
			if colision == 0 { // Meilleur qu'on puisse faire ça ne sert à rien de continuer
				*bestCollision = 0
				*bestResult = leafResult
				return
			} else if colision < *bestCollision { // Mise à jour du meilleur résultat
				*bestCollision = colision
				*bestResult = leafResult
			}
		}
		// Pas besoin de continuer
		return
	}

	// Cas habituel on peut placer le joueur sur la dernière partie
	for i, iPlayer := range rs.playersNotInGames {
		newRound := rs.round.Clone()
		newGame := rs.round[len(rs.round)-1]
		if newGame.CollisionFoundIfIplaceThisPlayer(tournament.Player(iPlayer), playerEverEncountered) <= MAX_COLLISION_FOUND {
			newGame.PlacePlayerInGame(tournament.Player(iPlayer))
			newRound[len(rs.round)-1] = newGame
			newPlayersNotPlaced := slices.Concat(rs.playersNotInGames[:i], rs.playersNotInGames[i+1:])
			// On prépare sa nouvelle feuille
			colision := math.MaxInt
			var leafResult RoundState
			leaf := RoundState{playersNotInGames: newPlayersNotPlaced, round: newRound}
			// On parcourt l'arbre
			if depth < MAX_DEPTH_LOG {
				fmt.Printf("I'm in depth %d : Loop %d/%d\n", depth, i, len(rs.playersNotInGames))
			}
			leaf.createRoundUniqEncouterRecur(maxField, playerEverEncountered, &colision, &leafResult, depth+1)
			// On regarde les résultat
			if colision == 0 { // Meilleur qu'on puisse faire ça ne sert à rien de continuer
				*bestCollision = 0
				*bestResult = leafResult
				return
			} else if colision < *bestCollision { // Mise à jour du meilleur résultat
				*bestCollision = colision
				*bestResult = leafResult
			}
		}
	}
}

func DrawRoundBranchAndBound(playerEverEncountered tournament.PlayersTimeEncountered, players []int, maxField int) tournament.Round {
	var res RoundState
	init := RoundState{playersNotInGames: players, round: make(tournament.Round, 0, maxField)}
	bestCol := math.MaxInt
	init.createRoundUniqEncouterRecur(maxField, playerEverEncountered, &bestCol, &res, 0)
	return res.round
}
