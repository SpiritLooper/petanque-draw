package draws

import (
	"fmt"
	"math"
	"petanque-draw/tournament"
)

func DrawRandomRoundBruteForce(nbPlayer int, nbRound int, maxField int, MAX_GENERATION int) tournament.Tournament {
	// Tirer aléatoirement la première ronde
	bCollision := math.MaxInt
	bTournament := make(tournament.Tournament, 0, nbRound)
	for i := range MAX_GENERATION {
		tournament := make(tournament.Tournament, 0, nbRound)
		for range nbRound {
			tournament = append(tournament, drawRandomRound(nbPlayer, maxField))
		}
		collision := tournament.CountCollision()
		if collision < bCollision {
			bTournament = tournament
			bCollision = collision
		}
		if bCollision == 0 {
			return bTournament
		}
		if i%(MAX_GENERATION/100) == 0 || bCollision == 0 {
			fmt.Printf(".")
		}
	}
	return bTournament
}

func DrawRandomTournamentBruteForce(nbPlayer int, nbRound int, maxField int, MAX_GENERATION int) tournament.Tournament {
	// Tirer aléatoirement la première ronde
	bCollision := math.MaxInt
	bTournament := make(tournament.Tournament, 0, nbRound)
	for i := range MAX_GENERATION {
		tournament := make(tournament.Tournament, 0, nbRound)
		for range nbRound {
			tournament = append(tournament, drawRandomRound(nbPlayer, maxField))
		}
		collision := tournament.CountCollision()
		if collision < bCollision {
			bTournament = tournament
			bCollision = collision
		}
		if bCollision == 0 {
			return bTournament
		}
		if i%(MAX_GENERATION/100) == 0 || bCollision == 0 {
			fmt.Printf(".")
		}
	}
	return bTournament
}

func DrawRandomTournamentStepByStep(nbPlayer int, nbRound int, maxField int, MAX_GENERATION int) tournament.Tournament {
	// Tirer aléatoirement la première ronde
	tournament := make(tournament.Tournament, 0, nbRound)
	tournament = append(tournament, drawRandomRound(nbPlayer, maxField))
	encountredPlayer := tournament.GenEncounteredMatrix()
	bCollision := math.MaxInt
	bRound := drawRandomRound(nbPlayer, maxField)
	for i := range nbRound - 1 {
		fmt.Printf("------Generating-Round-%d------\n", i+1)
		bRound = drawRandomRound(nbPlayer, maxField)
		bCollision = math.MaxInt
		for j := 0; bCollision != 0 && j < MAX_GENERATION; j++ {
			round := drawRandomRound(nbPlayer, maxField)
			collision := round.CountCollision(encountredPlayer)
			if bCollision > collision {
				bCollision = collision
				bRound = round
			}
			if j%(MAX_GENERATION/100) == 0 || bCollision == 0 {
				fmt.Printf(".")
			}
		}
		fmt.Println("")
		tournament = append(tournament, bRound)
		encountredPlayer = tournament.GenEncounteredMatrix()
	}
	return tournament
}
