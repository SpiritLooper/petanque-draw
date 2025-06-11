package main

import (
	"fmt"
	"log"
	"petanque-draw/draws"
	"petanque-draw/tournament"
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

func main() {
	opts := draws.NewDefaultDrawOpts()
	opts.ALGO_TYPE = draws.ALGO_RANDOM
	opts.MAX_FIELDS = 9
	opts.NB_PLAYER = 16
	opts.NB_ROUNDS = 4
	opts.FOLLOWING_PLAYER = false

	for i := 16; i < opts.NB_PLAYER+1; i++ {
		generator := NewTournamentPDFGenerator()

		fmt.Printf("=================GENERATE=%d=PLAYERS=================\n", i)
		tour := draws.DrawTournament(opts)

		// Generer la premiere page avec le tournoi
		generator.GenerateTournamentPage(&tour, i)

		// Generer la seconde page avec les collisions
		generator.GenerateCollisionPage(&tour, i)

		tour.Display()
		fmt.Printf("Nb Collision : %d\n", tour.CountCollision())
		displayColision(tour.GetCollision())
		err := generator.SavePDF("tournoi_petanque_" + strconv.Itoa(i) + ".pdf")
		if err != nil {
			log.Fatalf("Erreur lors de la generation du PDF: %v", err)
		}

		fmt.Printf("PDF genere avec succes: tournoi_petanque_%d.pdf", i)
	}
}
