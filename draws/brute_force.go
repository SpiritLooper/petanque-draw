package draws

import (
	"fmt"
	"math"
	"petanque-draw/tournament"
	"sync"
)

type SafeTournament struct {
	mu          sync.Mutex
	bTournament tournament.Tournament
	bCollision  int
}

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
	encountredPlayerWith, encountredPlayerAgainst := tournament.GenEncounteredMatrix()
	bCollision := math.MaxInt
	bRound := drawRandomRound(nbPlayer, maxField)
	for i := range nbRound - 1 {
		fmt.Printf("------Generating-Round-%d------\n", i+1)
		bRound = drawRandomRound(nbPlayer, maxField)
		bCollision = math.MaxInt
		for j := 0; bCollision != 0 && j < MAX_GENERATION; j++ {
			round := drawRandomRound(nbPlayer, maxField)
			collision := round.CountCollision(encountredPlayerWith, encountredPlayerAgainst)
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
		encountredPlayerWith, encountredPlayerAgainst = tournament.GenEncounteredMatrix()
	}
	return tournament
}

func (bTour *SafeTournament) drawAndCompare(i int, opts TounamentDrawOpts) {
	defer func() { recover() }()
	bTour.mu.Lock()
	if (*bTour).bCollision == 0 {
		return
	}
	bTour.mu.Unlock()

	tour := DrawTournament(opts)
	col := tour.CountCollision()

	bTour.mu.Lock()
	if (*bTour).bCollision > col {
		(*bTour).bTournament = tour
		(*bTour).bCollision = col
		fmt.Printf("Found better tournament with %d collisions\n", bTour.bCollision)
		if col == 0 {
			return
		}
	}
	bTour.mu.Unlock()
	if i%(opts.MAX_ITERATION/100) == 0 {
		fmt.Printf(".")
	}
}

func DrawTournamentBruteForce(opts TounamentDrawOpts) tournament.Tournament {
	bTournament := SafeTournament{}
	bTournament.bCollision = math.MaxInt
	for i := range opts.MAX_ITERATION {
		go bTournament.drawAndCompare(i, opts)
	}
	return bTournament.bTournament
}
