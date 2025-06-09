package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"petanque-draw/tournament"
)

const (
	nPlayers     = 24
	nRounds      = 4
	nTerrains    = 9
	populationSz = 1000000
	nGenerations = 10
	mutationRate = 0.2
)

type EvaluatedTournament struct {
	Tournament tournament.Tournament
	Fitness    int
}

func generateInitialPopulation() []EvaluatedTournament {
	pop := make([]EvaluatedTournament, populationSz)
	for i := range pop {
		pop[i].Tournament = tournament.DrawTournament(nPlayers, nRounds, nTerrains)
		pop[i].Fitness = evaluateFitness(pop[i].Tournament)
	}
	return pop
}

func evaluateFitness(t tournament.Tournament) int {
	pairCount := make(map[[2]tournament.Player]int)
	count2v3 := 0
	countDoublettes := 0
	countTriplettes := 0

	for _, round := range t {
		for _, game := range round {
			teams := [][]tournament.Player{game.Team1, game.Team2}
			for _, team := range teams {
				for i := 0; i < len(team); i++ {
					for j := i + 1; j < len(team); j++ {
						a, b := team[i], team[j]
						if a > b {
							a, b = b, a
						}
						pairCount[[2]tournament.Player{a, b}]++
					}
				}
			}

			size1, size2 := len(game.Team1), len(game.Team2)
			if (size1 == 2 && size2 == 3) || (size1 == 3 && size2 == 2) {
				count2v3++
			}
			if size1 == 2 && size2 == 2 {
				countDoublettes++
			}
			if size1 == 3 && size2 == 3 {
				countTriplettes++
			}
		}
	}

	score := 0
	for _, count := range pairCount {
		if count == 1 {
			score += 10
		} else {
			score -= 5 * (count - 1)
		}
	}

	if count2v3 > 1 {
		score -= 50 * (count2v3 - 1)
	}

	score += 20 * countDoublettes
	score -= 5 * countTriplettes

	score += 200 - (10 * t.CountCollision())

	return score
}

func selectTop(pop []EvaluatedTournament) []EvaluatedTournament {
	sort.Slice(pop, func(i, j int) bool {
		return pop[i].Fitness > pop[j].Fitness
	})
	return pop[:populationSz/2]
}

func crossover(a, b tournament.Tournament) tournament.Tournament {
	cut := rand.Intn(len(a))
	off := make(tournament.Tournament, 0, len(a))
	off = append(off, a[:cut]...)
	off = append(off, b[cut:]...)
	return off
}

func mutate(t *tournament.Tournament) {
	if rand.Float64() > mutationRate {
		return
	}
	r := rand.Intn(len(*t))
	g := rand.Intn(len((*t)[r]))
	team1 := &((*t)[r][g].Team1)
	team2 := &((*t)[r][g].Team2)
	if len(*team1) == 0 || len(*team2) == 0 {
		return
	}
	i := rand.Intn(len(*team1))
	j := rand.Intn(len(*team2))
	(*team1)[i], (*team2)[j] = (*team2)[j], (*team1)[i]
}

func evolve(pop []EvaluatedTournament) tournament.Tournament {
	for g := 0; g < nGenerations; g++ {
		nextGen := selectTop(pop)
		for len(nextGen) < populationSz {
			a := nextGen[rand.Intn(len(nextGen))]
			b := nextGen[rand.Intn(len(nextGen))]
			child := crossover(a.Tournament, b.Tournament)
			mutate(&child)
			nextGen = append(nextGen, EvaluatedTournament{
				Tournament: child,
				Fitness:    evaluateFitness(child),
			})
		}
		pop = nextGen
		if g%100 == 0 {
			fmt.Printf("Generation %d: Best fitness = %d\n", g, pop[0].Fitness)
		}
	}
	return selectTop(pop)[0].Tournament
}

func binomial(n, k int) int {
	if k > n {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}
	if k > n-k {
		k = n - k
	}
	res := 1
	for i := 1; i <= k; i++ {
		res = res * (n - i + 1) / i
	}
	return res
}

// meilleureRepartitionMax4 essaie de maximiser terrains à 4 joueurs,
// puis complète avec terrains à 6 joueurs, enfin terrain à 5 si reste impair,
// sans forcément utiliser tous les terrains (certains peuvent rester vides)
func meilleureRepartitionMax4(joueurs, terrains int) []int {
	parties := make([]int, terrains)

	if joueurs > terrains*6 {
		return []int{}
	}

	nb4 := min(joueurs/4, terrains)
	for i := range nb4 {
		parties[i] = 4
	}

	restants := joueurs - 4*nb4
	nb6 := min(restants/2, terrains)
	var i int
	for i = range nb6 {
		parties[i] = 6
	}

	restants = restants - nb6*2

	if i+1 < terrains && restants == 1 {
		parties[i+1] = 5

	}

	return parties
}

func totalDoublettes(parties []int) int {
	total := 0
	for _, taille := range parties {
		if taille >= 2 {
			total += binomial(taille, 2)
		}
	}
	return total
}

func possibleSolutionStableMax4(joueurs, terrains, tours int) bool {
	parties := meilleureRepartitionMax4(joueurs, terrains)
	if len(parties) == 0 {
		fmt.Println("Impossible de répartir les joueurs avec la contrainte de maximiser terrains à 4 joueurs")
		return false
	}

	fmt.Printf("Répartition choisie : %v\n", parties)

	totalD := totalDoublettes(parties)

	// On utilise toujours le nombre total de joueurs ici !
	pairesMax := binomial(joueurs, 2)

	fmt.Printf("Joueurs totaux : %d, paires max possibles : %d, doublettes/tour : %d, total sur %d tours : %d\n",
		joueurs, pairesMax, totalD, tours, tours*totalD)

	return tours*totalD <= pairesMax
}

func main() {
	rand.Seed(time.Now().UnixNano())
	for i := 12; i < 41; i++ {
		if possibleSolutionStableMax4(i, nTerrains, nRounds) {
			fmt.Printf("%d : ✅\n", i)
		} else {
			fmt.Printf("%d : ❌\n", i)
		}
	}
	// pop := generateInitialPopulation()
	// best := evolve(pop)
	// fmt.Println("\nBest Tournament:")
	// best.Display()

	// fmt.Printf("Fitness: %d\t Collision: %d\n", evaluateFitness(best), best.CountCollision()/2)
}
