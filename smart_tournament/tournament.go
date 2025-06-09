package smart_tournament

import (
	"errors"
	"fmt"
)

type Tournament struct {
	Rounds            []Round
	playerEncountered playerEncountered
}
type Player int

type Round []Game
type Game []Player

type playerEncountered map[Player]playerSet
type playerSet map[Player]struct{}

func (t *Tournament) Display() {
	for i, rounds := range (*t).Rounds {
		fmt.Printf("\tRonde %d\n", i+1)
		for j, game := range rounds {
			fmt.Printf("Terrain %d : [ ", j+1)
			for _, player := range game {
				fmt.Printf("%d ", player+1)
			}
			fmt.Printf("]\n")
		}
		fmt.Println("--------------------------")
	}
}

func (round *Round) canCreateNewField(nbPlayer int, nbMaxField int) bool {
	return len(*round) <= min(nbPlayer/4, nbMaxField)
}

func (p Player) neverEncounteredNewPlayer(playerEncountered playerEncountered, newGame Game) bool {
	for _, opponent := range newGame {
		if _, exist := playerEncountered[p][opponent]; opponent != p && exist {
			return false
		}
	}
	return true
}

func (g *Game) gamePlayersIsFull() bool {
	return len(*g) >= 6
}

func (g *Game) gamePlayersIsFullForDoublette() bool {
	return len(*g) >= 4
}

func (round *Round) selectGame(playersRemain int, maxField int) *Game {
	if len(*round) == 0 {
		*round = append(*round, Game{})
		return &(*round)[len(*round)-1]
	}

	// Check Si toute les doublette son prise
	for _, game := range *round {
		if game.gamePlayersIsFullForDoublette() {
			return &game
		}
	}

	lastGame := &(*round)[len(*round)-1]
	if lastGame.gamePlayersIsFullForDoublette() && playersRemain >= 4 && len(*round) < maxField {
		*round = append(*round, Game{})
		return &(*round)[len(*round)-1]
	}

	if lastGame.gamePlayersIsFullForDoublette() && (playersRemain < 4 || len(*round) >= maxField) {
		i := 0
		game := &(*round)[i]

		for game.gamePlayersIsFull() {
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

func (round *Round) placePlayer(p Player, playerRemain int, encountered *playerEncountered, maxField int, nbPlayer int) {

	gameSelected := round.selectGame(playerRemain, maxField)
	if !p.neverEncounteredNewPlayer(*encountered, *gameSelected) {
		found := false
		for _, game := range *round {
			if p.neverEncounteredNewPlayer(*encountered, game) && !game.gamePlayersIsFull() {
				found = true
				gameSelected = &game
			}
		}

		if !found && round.canCreateNewField(nbPlayer, maxField) {
			*round = append(*round, Game{})
			gameSelected = &(*round)[len(*round)-1]
		} else if !found {
			panic("help !!!!")
		}
	}
	err := gameSelected.placePlayer(p, encountered)
	if err != nil {
		panic(err)
	}
}

func (g *Game) placePlayer(p Player, enc *playerEncountered) error {
	if g.gamePlayersIsFull() {
		return errors.New("can't place new player. game is full")
	}

	*g = append(*g, p)

	if _, exist := (*enc)[p]; !exist {
		(*enc)[p] = make(playerSet)
	}

	for _, opponent := range *g {
		(*enc)[opponent][p] = struct{}{}
		(*enc)[p][opponent] = struct{}{}
	}

	return nil
}
