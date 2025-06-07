package tournament

import "errors"

type Tournament []Round
type Round []Game
type Game struct {
	Team1 []Player
	Team2 []Player
}
type Player int

func gamePlayersIsFull(game Game) bool {
	return len(game.Team1) >= 3 && len(game.Team2) >= 3
}

func gamePlayersIsFullForDoublette(game Game) bool {
	return len(game.Team1) >= 2 && len(game.Team2) >= 2
}

func selectGame(round *Round, playersRemain int, maxField int) *Game {
	if len(*round) == 0 {
		*round = append(*round, Game{})
		return &(*round)[len(*round)-1]
	}

	lastGame := &(*round)[len(*round)-1]
	if gamePlayersIsFullForDoublette(*lastGame) && playersRemain >= 4 && len(*round) < maxField {
		*round = append(*round, Game{})
		return &(*round)[len(*round)-1]
	}

	if gamePlayersIsFullForDoublette(*lastGame) && (playersRemain < 4 || len(*round) >= maxField) {
		i := 0
		game := &(*round)[i]

		for gamePlayersIsFull(*game) {
			i += 1
			if i >= len(*round) {
				panic("no game found to place player")
			}
			game = &(*round)[i]
		}
		return game
	}

	return lastGame
}

func placePlayerInRound(p Player, round *Round, playerRemain int, maxField int) {
	err := placePlayerInGame(p, selectGame(round, playerRemain, maxField))
	if err != nil {
		panic(err)
	}
}

func placePlayerInGame(p Player, game *Game) error {
	if gamePlayersIsFull(*game) {
		return errors.New("can't place new player. game is full")
	}

	if len(game.Team1) > len(game.Team2) {
		(*game).Team2 = append((*game).Team2, p)
	} else {
		(*game).Team1 = append((*game).Team1, p)
	}
	return nil
}
