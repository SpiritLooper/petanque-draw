package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"petanque-draw/tournament"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

// findSystemFont cherche une police sur le système NixOS
func findSystemFont(fontName string) (string, bool) {
	// Chemins typiques pour les polices sur NixOS
	fontPaths := []string{
		"/run/current-system/sw/share/fonts",
		"/etc/profiles/per-user/" + os.Getenv("USER") + "/share/fonts",
		os.Getenv("HOME") + "/.local/share/fonts",
		os.Getenv("HOME") + "/.fonts",
		"/usr/share/fonts", // Fallback
	}

	// Extensions possibles
	extensions := []string{".ttf", ".otf", ".TTF", ".OTF"}

	for _, basePath := range fontPaths {
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			continue
		}

		// Chercher récursivement dans tous les sous-dossiers
		err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Ignorer les erreurs et continuer
			}

			if !info.IsDir() {
				fileName := strings.ToLower(filepath.Base(path))
				fontNameLower := strings.ToLower(fontName)

				// Vérifier si le nom du fichier contient le nom de la police
				for _, ext := range extensions {
					if strings.Contains(fileName, fontNameLower) && strings.HasSuffix(fileName, strings.ToLower(ext)) {
						return fmt.Errorf("found: %s", path) // Utiliser erreur pour arrêter la recherche
					}
				}
			}
			return nil
		})

		if err != nil && strings.HasPrefix(err.Error(), "found: ") {
			return strings.TrimPrefix(err.Error(), "found: "), true
		}
	}

	return "", false
}

// setupFonts configure les polices du système
func setupFonts(pdf *gofpdf.Fpdf) {
	// Ordre de préférence des polices
	fontPreferences := []string{
		"Inter", "Roboto", "Ubuntu", "DejaVu Sans", "Liberation Sans",
		"Noto Sans", "Source Sans Pro", "Open Sans", "Lato", "Nunito",
	}

	var selectedFont string
	var fontPath string
	var found bool

	// Chercher la première police disponible
	for _, font := range fontPreferences {
		if fontPath, found = findSystemFont(font); found {
			selectedFont = font
			break
		}
	}

	if !found {
		log.Println("Aucune police système trouvée, utilisation de la police par défaut")
		return
	}

	log.Printf("Police sélectionnée: %s (%s)", selectedFont, fontPath)

	// Essayer d'ajouter la police
	pdf.AddUTF8Font(selectedFont, "", fontPath)

	// Chercher les variantes (gras, italique) dans le même répertoire
	fontDir := filepath.Dir(fontPath)
	fontBaseName := strings.ToLower(selectedFont)

	// Recherche des variantes
	boldVariants := []string{"bold", "b", "700", "black"}
	italicVariants := []string{"italic", "i", "oblique"}

	filepath.Walk(fontDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		fileName := strings.ToLower(filepath.Base(path))
		if !strings.Contains(fileName, fontBaseName) {
			return nil
		}

		// Vérifier pour les variantes en gras
		for _, variant := range boldVariants {
			if strings.Contains(fileName, variant) && !strings.Contains(fileName, "italic") {
				pdf.AddUTF8Font(selectedFont, "B", path)
				log.Printf("Police gras ajoutée: %s", path)
				return nil
			}
		}

		// Vérifier pour les variantes en italique
		for _, variant := range italicVariants {
			if strings.Contains(fileName, variant) && !strings.Contains(fileName, "bold") {
				pdf.AddUTF8Font(selectedFont, "I", path)
				log.Printf("Police italique ajoutée: %s", path)
				return nil
			}
		}

		return nil
	})
}

// TournamentPDFGenerator génère un PDF pour les tournois
type TournamentPDFGenerator struct {
	pdf *gofpdf.Fpdf
}

// NewTournamentPDFGenerator crée un nouveau générateur PDF avec polices système
func NewTournamentPDFGenerator() *TournamentPDFGenerator {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Configurer les polices système
	setupFonts(pdf)

	return &TournamentPDFGenerator{pdf: pdf}
}

// GenerateTournamentPage génère une page pour un tournoi donné
func (g *TournamentPDFGenerator) GenerateTournamentPage(tournament *tournament.Tournament, playerCount int) {
	g.pdf.AddPage()

	// En-tête principal
	g.pdf.SetFont("Arial", "B", 20)
	g.pdf.SetTextColor(44, 62, 80)
	g.pdf.CellFormat(0, 15, fmt.Sprintf("Tournoi de Petanque - %d Joueurs", playerCount), "", 1, "C", false, 0, "")
	g.pdf.Ln(5)

	// Sous-titre compact
	g.pdf.SetFont("Arial", "I", 10)
	g.pdf.SetTextColor(100, 100, 100)
	totalGames := len(*tournament) * len((*tournament)[0])
	g.pdf.CellFormat(0, 6, fmt.Sprintf("4 rondes - %d parties - %d terrains max", totalGames, len((*tournament)[0])), "", 1, "C", false, 0, "")
	g.pdf.Ln(8)

	// Génération des rondes en 2 colonnes
	g.addRoundsInTwoColumns(tournament)

	// Pied de page avec numéros des joueurs
	g.addPlayerNumbers(playerCount)
}

// addRoundsInTwoColumns ajoute les rondes en 2 colonnes
func (g *TournamentPDFGenerator) addRoundsInTwoColumns(tournament *tournament.Tournament) {
	startX := 10.0
	startY := g.pdf.GetY()
	columnWidth := 95.0

	for i, round := range *tournament {
		// Calculer la position
		col := i % 2
		row := i / 2

		x := startX + float64(col)*columnWidth
		y := startY + float64(row)*70 // Espacement vertical entre les rangées

		g.pdf.SetXY(x, y)
		g.addCompactRound(round, i+1, columnWidth)
	}
}

// addCompactRound ajoute une ronde compacte
func (g *TournamentPDFGenerator) addCompactRound(round tournament.Round, roundNumber int, width float64) {
	currentX := g.pdf.GetX()
	// currentY := g.pdf.GetY()

	// En-tête de la ronde
	g.pdf.SetFont("Arial", "B", 12)
	g.pdf.SetTextColor(52, 73, 94)
	g.pdf.SetFillColor(236, 240, 241)
	g.pdf.CellFormat(width, 8, fmt.Sprintf("Ronde %d", roundNumber), "", 1, "C", true, 0, "")

	// Retour à la position X de départ pour les lignes suivantes
	g.pdf.SetX(currentX)

	// Lignes des matchs compactes
	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetTextColor(0, 0, 0)

	for gameIndex, game := range round {
		g.pdf.SetX(currentX)

		// Alternance de couleurs
		if gameIndex%2 == 1 {
			g.pdf.SetFillColor(248, 249, 250)
		} else {
			g.pdf.SetFillColor(255, 255, 255)
		}

		// Terrain (plus petit)
		g.pdf.CellFormat(12, 6, fmt.Sprintf("T%d", gameIndex+1), "1", 0, "C", true, 0, "")

		// Match complet sur le reste de la largeur
		team1Str := g.formatCompactTeam(game.Team1)
		team2Str := g.formatCompactTeam(game.Team2)
		matchStr := fmt.Sprintf("%s vs %s", team1Str, team2Str)

		g.pdf.CellFormat(width-12, 6, matchStr, "1", 1, "C", true, 0, "")
	}
}

// formatCompactTeam formate l'affichage compact d'une équipe
func (g *TournamentPDFGenerator) formatCompactTeam(team []tournament.Player) string {
	var players []string
	for _, player := range team {
		players = append(players, strconv.Itoa(int(player+1)))
	}
	return strings.Join(players, "-")
}

// addPlayerNumbers ajoute la liste des numéros de joueurs en pied de page
func (g *TournamentPDFGenerator) addPlayerNumbers(playerCount int) {
	g.pdf.SetY(280)
	g.pdf.SetFont("Arial", "I", 8)
	g.pdf.SetTextColor(100, 100, 100)

	// Générer la liste des joueurs de manière plus compacte
	var playerNumbers []string
	for i := 1; i <= playerCount; i++ {
		playerNumbers = append(playerNumbers, strconv.Itoa(i))
	}

	playerList := strings.Join(playerNumbers, " - ")
	g.pdf.CellFormat(0, 5, fmt.Sprintf("Joueurs: %s", playerList), "", 1, "C", false, 0, "")
}

// SavePDF sauvegarde le PDF
func (g *TournamentPDFGenerator) SavePDF(filename string) error {
	return g.pdf.OutputFileAndClose(filename)
}

const (
	MAX_FIELD  = 9
	NB_ROUNDS  = 4
	NB_PLAYERS = 18
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
	generator := NewTournamentPDFGenerator()
	tour := tournament.DrawTournament(NB_PLAYERS, NB_ROUNDS, MAX_FIELD)
	for len(tour[len(tour)-1]) == 0 {
		println("Last round not generated, do it again")
		tour = tournament.DrawTournament(NB_PLAYERS, NB_ROUNDS, MAX_FIELD)
	}
	generator.GenerateTournamentPage(&tour, NB_PLAYERS)
	tour.Display()
	fmt.Printf("Nb Collision : %d\n", tour.CountCollision())
	displayColision(tour.GetCollision())
	err := generator.SavePDF("tournoi_petanque_" + strconv.Itoa(NB_PLAYERS) + ".pdf")
	if err != nil {
		log.Fatalf("Erreur lors de la génération du PDF: %v", err)
	}

	fmt.Printf("PDF généré avec succès: tournoi_petanque_%d.pdf", NB_PLAYERS)
}
