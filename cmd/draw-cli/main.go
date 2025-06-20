package main

import (
	"fmt"
	"log"
	"petanque-draw/draws"
	"petanque-draw/tournament"
	"petanque-draw/utils/pdfgen"
	"strconv"
)

func displayColision(col []tournament.PlayerSet) {
	fmt.Println("Collisions found :")
	for i, set := range col {
		if len(set) > 0 {
			fmt.Printf("\t%d -> ", i+1)
			for j := range set {
				fmt.Printf("%d ", j+1)
			}
			fmt.Printf("\n")
		}
	}
}

const MIN_PLAYER_TO_GEN = 16
const MAX_PLAYER_TO_GEN = 40

func main() {
	opts := draws.NewDefaultDrawOpts()
	opts.ALGO_TYPE = draws.ALGO_GREEDY
	opts.MAX_FIELDS = 9
	opts.NB_ROUNDS = 4
	opts.FOLLOWING_PLAYER = false
	opts.MAX_ITERATION = 200_000

	for i := MIN_PLAYER_TO_GEN; i <= MAX_PLAYER_TO_GEN; i++ {
		generator := pdfgen.NewTournamentPDFGenerator()

		opts.NB_PLAYER = i

		fmt.Printf("=================GENERATE=%d=PLAYERS=================\n", i)
		tour := draws.DrawTournamentBruteForce(opts)
		// tour := draws.DrawTournament(opts)

		tour.Display()
		fmt.Printf("Nb Collision : %d\n", tour.CountCollision())
		displayColision(tour.GetCollision())

		// Generer la premiere page avec le tournoi
		generator.GenerateTournamentPage(&tour, i)

		// Generer la seconde page avec les collisions
		generator.GenerateCollisionPage(&tour, i)

		err := generator.SavePDF("tournoi_petanque_" + strconv.Itoa(i) + ".pdf")
		if err != nil {
			log.Fatalf("Erreur lors de la generation du PDF: %v", err)
		}

		fmt.Printf("PDF genere avec succes: tournoi_petanque_%d.pdf\n", i)
	}
}
