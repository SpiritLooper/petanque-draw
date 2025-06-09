package tournament

import (
	"fmt"
	"math"
	"petanque-draw/utils/intset"
	"slices"
)

const MAX_DEPTH_LOG = 2
const MAX_COLLISION_FOUND = 2

type RoundState struct {
	playersNotInGames []int
	round             Round
}

func (rs RoundState) createRoundUniqEncouterRecur(
	maxField int, playerEverEncountered PlayerEncountered,
	bestCollision *int, bestResult *RoundState, depth int,
) {
	// Si on a plus de tirage, on retourne le résultat de collision
	if len(rs.playersNotInGames) <= 0 {
		*bestCollision = rs.round.countCollision(playerEverEncountered)
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
			newGame := Game{Team1: []Player{Player(iPlayer)}}
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
			newGame := Game{Team1: make([]Player, 0, 3), Team2: make([]Player, 0, 3)}
			newGame.placePlayerInGame(Player(iPlayer))
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
			newGame.placePlayerInGame(Player(iPlayer))
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
		if newGame.collisionFoundIfIplaceThisPlayer(Player(iPlayer), playerEverEncountered) <= MAX_COLLISION_FOUND {
			newGame.placePlayerInGame(Player(iPlayer))
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

func DrawRoundUniqEncounter(playerEverEncountered PlayerEncountered, nbPlayer int, maxField int) Round {

	// Creations ordre de joueur aléatoire

	// players := slices.Collect(maps.Keys(intset.CreateIntSet(nbPlayer)))

	players := make([]int, nbPlayer, nbPlayer)
	for i := range nbPlayer {
		players[i] = i
	}

	var res RoundState
	init := RoundState{playersNotInGames: players, round: make(Round, 0, maxField)}
	bestCol := math.MaxInt
	init.createRoundUniqEncouterRecur(maxField, playerEverEncountered, &bestCol, &res, 0)
	return res.round
}

func (round *Round) placePlayerLazy(p Player, remainPlayers int, maxField int) {
	// Aucunes partie ? On en créer une
	if len(*round) == 0 {
		newGame := Game{}
		newGame.placePlayerInGame(p)
		*round = append(*round, newGame)
		return
	}

	// On prend la dernière partie
	lastGame := &(*round)[len(*round)-1]
	// Si la partie est pleine pour une doublette (mais qu'on est pas sur le cas d'une triplette)
	// On ajoute le joueur dans une nouvelle
	if lastGame.IsFullForDoublette() && remainPlayers >= 4 && len(*round) < maxField {
		newGame := Game{}
		newGame.placePlayerInGame(p)
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
		game.placePlayerInGame(p)
		return
	}

	// Dans les autres cas le joueur peut être placé sur la partie actuelle
	lastGame.placePlayerInGame(p)
}

func DrawNotRandomRound(nbPlayer int, maxField int) Round {
	var round Round
	for i := range nbPlayer {
		round.placePlayerLazy(Player(i), nbPlayer, maxField)
	}
	return round
}

func DrawRandomRound(nbPlayer int, maxField int) Round {
	var round Round
	setInt := intset.CreateIntSet(nbPlayer)
	for !setInt.IsEmpty() {
		randInt, _ := setInt.RandomPop()
		round.placePlayerLazy(Player(randInt), setInt.Count()+1, maxField)
	}
	return round
}

func DrawRandomTournament(nbPlayer int, nbRound int, maxField int) Tournament {
	var tournament Tournament
	// Tirer aléatoirement la première ronde
	for range nbRound {
		tournament = append(tournament, DrawRandomRound(nbPlayer, maxField))
	}
	return tournament
}

func DrawTournament(nbPlayer int, nbRound int, maxField int) Tournament {
	var tournament Tournament
	// Tirer aléatoirement la première ronde
	fmt.Printf("-----------Generate-Round-%d---------\n", 1)
	tournament = append(tournament, DrawNotRandomRound(nbPlayer, maxField))

	for i := range nbRound - 1 {
		fmt.Printf("-----------Generate-Round-%d---------\n", i+2)
		playerEncounteredMatrix := tournament.genEncounteredMatrix()
		tournament = append(tournament, DrawRoundUniqEncounter(playerEncounteredMatrix, nbPlayer, maxField))
	}
	return tournament
}
