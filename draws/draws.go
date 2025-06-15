package draws

import (
	"fmt"
	"maps"
	"petanque-draw/tournament"
	"petanque-draw/utils"
	"petanque-draw/utils/intset"
	"slices"
)

const (
	ALGO_BRANCH_AND_BOUND  = 0
	ALGO_GREEDY            = 1
	ALGO_RANDOM            = 2
	ALGO_RANDOM_BRUTEFORCE = 3
	ALGO_RANDOM_BYSTEP     = 4
)

type TounamentDrawOpts struct {
	ALGO_TYPE        int
	FOLLOWING_PLAYER bool
	MAX_ITERATION    int
	MAX_FIELDS       int
	NB_PLAYER        int
	NB_ROUNDS        int
	VERBOSE          bool
}

func NewDefaultDrawOpts() TounamentDrawOpts {
	return TounamentDrawOpts{
		ALGO_TYPE:        ALGO_RANDOM,
		FOLLOWING_PLAYER: false,
		MAX_ITERATION:    100_000,
		MAX_FIELDS:       4,
		NB_PLAYER:        16,
		NB_ROUNDS:        4,
		VERBOSE:          false,
	}
}

func drawFollowingRound(nbPlayer int, maxField int) tournament.Round {
	var round tournament.Round
	for i := range nbPlayer {
		round.PlacePlayerLazy(tournament.Player(i), nbPlayer-i, maxField)
	}
	return round
}

func drawRandomRound(nbPlayer int, maxField int) tournament.Round {
	var round tournament.Round
	setInt := intset.CreateIntSet(nbPlayer)
	for !setInt.IsEmpty() {
		randInt, _ := setInt.RandomPop()
		round.PlacePlayerLazy(tournament.Player(randInt), setInt.Count()+1, maxField)
	}
	return round.ReArranged()
}

func DrawTournament(opts TounamentDrawOpts) tournament.Tournament {
	tournamentRes := make(tournament.Tournament, 0, opts.NB_ROUNDS)
	// Add first round
	if opts.VERBOSE {
		fmt.Printf("-----------Generate-Round-%d---------\n", 1)
	}
	if opts.FOLLOWING_PLAYER {
		tournamentRes = append(tournamentRes, drawFollowingRound(opts.NB_PLAYER, opts.MAX_FIELDS))
	} else {
		tournamentRes = append(tournamentRes, drawRandomRound(opts.NB_PLAYER, opts.MAX_FIELDS))
	}

	for roundIndex := range opts.NB_ROUNDS - 1 {
		if opts.VERBOSE {
			fmt.Printf("-----------Generate-Round-%d---------\n", roundIndex+2)
		}
		// Tirage des joueurs
		intPlayerVector := make([]int, opts.NB_PLAYER)
		if opts.FOLLOWING_PLAYER {
			for i := range opts.NB_PLAYER {
				intPlayerVector[i] = i
			}
			intPlayerVector = utils.LeftRotation(intPlayerVector, roundIndex)
		} else {
			intPlayerVector = slices.Collect(maps.Keys(intset.CreateIntSet(opts.NB_PLAYER)))
		}
		// Generation de la matrice de rencontre
		encounterMatrix := tournamentRes.GenEncounteredMatrix()

		round := make(tournament.Round, 0)

		// Execution de l'algo de tirage
		switch opts.ALGO_TYPE {
		case ALGO_BRANCH_AND_BOUND:
			round = DrawRoundBranchAndBound(encounterMatrix, intPlayerVector, opts.MAX_FIELDS)
		case ALGO_GREEDY:
			round = DrawRoundGreed(encounterMatrix, intPlayerVector, opts.MAX_FIELDS)
		default:
			if opts.FOLLOWING_PLAYER {
				round = drawFollowingRound(opts.NB_PLAYER, opts.MAX_FIELDS)
			} else {
				round = drawRandomRound(opts.NB_PLAYER, opts.MAX_FIELDS)
			}
		}

		tournamentRes = append(tournamentRes, round)
	}

	return tournamentRes
}
