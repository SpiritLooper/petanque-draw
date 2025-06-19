package pdfgen

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

const cellHeight = 5.0 * 2
const margin = 10.0

// findSystemFont cherche une police sur le systeme NixOS
func findSystemFont(fontName string) (string, bool) {
	// Chemins typiques pour les polices sur NixOS
	fontPaths := []string{
		"fonts/Roboto",
		"/run/current-system/sw/share/X11/fonts",
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

		// Chercher recursivement dans tous les sous-dossiers
		err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Ignorer les erreurs et continuer
			}

			if !info.IsDir() {
				fileName := strings.ToLower(filepath.Base(path))
				fontNameLower := strings.ToLower(fontName)

				// Verifier si le nom du fichier contient le nom de la police
				for _, ext := range extensions {
					if strings.Contains(fileName, fontNameLower) && strings.HasSuffix(fileName, strings.ToLower(ext)) {
						return fmt.Errorf("found: %s", path) // Utiliser erreur pour arreter la recherche
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

// setupFonts configure les polices du systeme
func setupFonts(pdf *gofpdf.Fpdf) {
	// Ordre de preference des polices
	fontPreferences := []string{
		"Roboto", "Ubuntu", "DejaVu Sans", "Roboto Sans",
		"Noto Sans", "Source Sans Pro", "Open Sans", "Lato", "Nunito",
	}

	var selectedFont string
	var fontPath string
	var found bool

	// Chercher la premiere police disponible
	for _, font := range fontPreferences {
		if fontPath, found = findSystemFont(font); found {
			selectedFont = font
			break
		}
	}

	if !found {
		// log.Println("Aucune police systeme trouvee, utilisation de la police par defaut")
		return
	}

	log.Printf("Police selectionnee: %s (%s)", selectedFont, fontPath)

	// Essayer d'ajouter la police
	pdf.AddUTF8Font(selectedFont, "", fontPath)

	// Chercher les variantes (gras, italique) dans le meme repertoire
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

		// Verifier pour les variantes en gras
		for _, variant := range boldVariants {
			if strings.Contains(fileName, variant) && !strings.Contains(fileName, "italic") {
				pdf.AddUTF8Font(selectedFont, "", path)
				log.Printf("Police gras ajoutee: %s", path)
				return nil
			}
		}

		// Verifier pour les variantes en italique
		for _, variant := range italicVariants {
			if strings.Contains(fileName, variant) && !strings.Contains(fileName, "bold") {
				pdf.AddUTF8Font(selectedFont, "I", path)
				log.Printf("Police italique ajoutee: %s", path)
				return nil
			}
		}

		return nil
	})
}

// TournamentPDFGenerator genere un PDF pour les tournois
type TournamentPDFGenerator struct {
	pdf *gofpdf.Fpdf
}

// NewTournamentPDFGenerator cree un nouveau generateur PDF avec polices systeme
func NewTournamentPDFGenerator() *TournamentPDFGenerator {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Configurer les polices systeme
	setupFonts(pdf)

	return &TournamentPDFGenerator{pdf: pdf}
}

// GenerateTournamentPage genere une page pour un tournoi donne
func (g *TournamentPDFGenerator) GenerateTournamentPage(tournament *tournament.Tournament, playerCount int) {
	g.pdf.AddPage()

	// En-tete principal
	g.pdf.SetFont("Roboto", "", 20)
	g.pdf.SetTextColor(44, 62, 80)
	g.pdf.CellFormat(0, 15, fmt.Sprintf("Tournoi de Pétanque - %d Joueurs", playerCount), "", 1, "C", false, 0, "")
	g.pdf.Ln(5)

	// Sous-titre compact
	g.pdf.SetFont("Roboto", "I", 10)
	g.pdf.SetTextColor(100, 100, 100)
	totalGames := len(*tournament) * len((*tournament)[0])
	g.pdf.CellFormat(0, 6, fmt.Sprintf("4 tirages - %d parties - %d terrains max", totalGames, len((*tournament)[0])), "", 1, "C", false, 0, "")
	g.pdf.Ln(8)

	// Generation des tours en 2 colonnes
	g.addRoundsInTwoColumns(tournament)

	// Collisions footer
	g.generateCollisionFooter(tournament)
}

func (g *TournamentPDFGenerator) generateCollisionFooter(tour *tournament.Tournament) {
	cellWidth := 10.0

	// En-tete principal
	g.pdf.SetFont("Roboto", "", cellWidth*1.4)
	g.pdf.SetTextColor(60, 60, 60)
	g.pdf.Ln(margin)

	g.pdf.CellFormat(2*cellWidth, cellWidth, "", "", 0, "C", false, 0, "")
	for i := range len(*tour) {
		g.pdf.CellFormat(cellWidth, cellWidth, fmt.Sprintf("%d", i+1), "", 0, "C", false, 0, "")
	}
	g.pdf.Ln(cellWidth)
	everDo := make([][]bool, (*tour)[0].CountPlacedPlayer())
	for i := range (*tour)[0].CountPlacedPlayer() {
		everDo[i] = make([]bool, (*tour)[0].CountPlacedPlayer())
		for j := range (*tour)[0].CountPlacedPlayer() {
			everDo[i][j] = false
		}
	}
	collisions := tour.GetCollision()
	for i, playerSet := range collisions {
		if len(playerSet) > 0 {
			for j := range playerSet {
				if _, exist := collisions[i][j]; exist && !everDo[tournament.Player(i)][j] && !everDo[j][tournament.Player(i)] {
					everDo[tournament.Player(i)][j] = true
					everDo[j][tournament.Player(i)] = true
					g.pdf.SetTextColor(10, 10, 10)
					g.pdf.CellFormat(cellWidth*2, cellWidth, fmt.Sprintf("%d-%d", i+1, j+1), "", 0, "C", false, 0, "")
					// On parcours les parties trouvés avec la collisions
					for _, round := range *tour {
						found := false
						for _, game := range round {
							if game.IsPlayerPlayWith(tournament.Player(i), j) {
								found = true
								g.pdf.SetTextColor(20, 200, 80)
								if game.IsContainsTriplette() {
									g.pdf.CellFormat(cellWidth, cellWidth, "T", "1", 0, "C", false, 0, "")
								} else {
									g.pdf.CellFormat(cellWidth, cellWidth, "D", "1", 0, "C", false, 0, "")
								}
								break
							} else if game.IsPlayerPlayAgainst(tournament.Player(i), j) {
								found = true
								g.pdf.SetTextColor(255, 40, 40)
								if game.IsContainsTriplette() {
									g.pdf.CellFormat(cellWidth, cellWidth, "T", "1", 0, "C", false, 0, "")
								} else {
									g.pdf.CellFormat(cellWidth, cellWidth, "D", "1", 0, "C", false, 0, "")
								}
								break
							}
						}
						if !found {
							g.pdf.CellFormat(cellWidth, cellWidth, "", "1", 0, "C", false, 0, "")
						}
					}
					g.pdf.Ln(cellWidth)
				}
			}
		}
	}
}

// GenerateCollisionPage genere une page avec les collisions
func (g *TournamentPDFGenerator) GenerateCollisionPage(tournament *tournament.Tournament, playerCount int) {
	g.pdf.AddPage()

	// En-tete principal
	g.pdf.SetFont("Roboto", "", 20)
	g.pdf.SetTextColor(44, 62, 80)
	g.pdf.CellFormat(0, 15, "Analyse des Collisions", "", 1, "C", false, 0, "")
	g.pdf.Ln(5)

	// Statistiques generales
	collisionCount := tournament.CountCollision()
	g.pdf.SetFont("Roboto", "", 12)
	g.pdf.SetTextColor(52, 73, 94)
	g.pdf.CellFormat(0, 8, fmt.Sprintf("Nombre total de collisions: %d", collisionCount), "", 1, "L", false, 0, "")
	g.pdf.Ln(5)

	// Recuperer les collisions
	collisions := tournament.GetCollision()

	if collisionCount == 0 {
		// Aucune collision
		g.pdf.SetFont("Roboto", "I", 12)
		g.pdf.SetTextColor(39, 174, 96)
		g.pdf.CellFormat(0, 10, "Aucune collision détectée ! Tournoi parfaitement équilibré.", "", 1, "C", false, 0, "")
	} else {
		// Affichage des collisions
		g.pdf.SetFont("Roboto", "", 10)
		g.pdf.SetTextColor(0, 0, 0)
		g.pdf.CellFormat(0, 8, "Détail des collisions par joueur:", "", 1, "L", false, 0, "")
		g.pdf.Ln(3)

		// En-tetes du tableau
		g.pdf.SetFont("Roboto", "", 10)
		g.pdf.SetFillColor(236, 240, 241)
		g.pdf.SetTextColor(52, 73, 94)
		g.pdf.CellFormat(30, 8, "Joueur", "1", 0, "C", true, 0, "")
		g.pdf.CellFormat(160, 8, "Collisions avec", "1", 1, "C", true, 0, "")

		// Contenu du tableau
		g.pdf.SetFont("Roboto", "", 9)
		g.pdf.SetTextColor(0, 0, 0)

		for i, playerSet := range collisions {
			if len(playerSet) > 0 {
				// Alternance de couleurs
				if i%2 == 1 {
					g.pdf.SetFillColor(248, 249, 250)
				} else {
					g.pdf.SetFillColor(255, 255, 255)
				}

				// Numero du joueur
				g.pdf.CellFormat(30, 6, fmt.Sprintf("Joueur %d", i+1), "1", 0, "C", true, 0, "")

				// Liste des joueurs en collision
				var collisionList []string
				for j := range playerSet {
					collisionList = append(collisionList, fmt.Sprintf("J%d", j+1))
				}
				collisionStr := strings.Join(collisionList, ", ")

				g.pdf.CellFormat(160, 6, collisionStr, "1", 1, "L", true, 0, "")
			}
		}

		// Explication
		g.pdf.Ln(5)
		g.pdf.SetFont("Roboto", "I", 8)
		g.pdf.SetTextColor(100, 100, 100)
		g.pdf.MultiCell(0, 4, "Une collision indique que deux joueurs se retrouvent dans la même équipe plus d'une fois durant le tournoi. L'objectif est de minimiser ces répétitions pour un tournoi équilibré.", "", "", false)
	}

	// Recommandations
	g.pdf.Ln(8)
	g.pdf.SetFont("Roboto", "", 10)
	g.pdf.SetTextColor(52, 73, 94)
	g.pdf.CellFormat(0, 6, "Recommandations:", "", 1, "L", false, 0, "")

	g.pdf.SetFont("Roboto", "", 9)
	g.pdf.SetTextColor(0, 0, 0)

	if collisionCount == 0 {
		g.pdf.MultiCell(0, 5, "- Ce tournoi est optimal, aucune modification nécessaire", "", "", false)
	} else if collisionCount <= 5 {
		g.pdf.MultiCell(0, 5, "- Niveau de collision acceptable pour un tournoi de cette taille", "", "", false)
		g.pdf.MultiCell(0, 5, "- Vous pouvez procéder avec cette configuration", "", "", false)
	} else {
		g.pdf.MultiCell(0, 5, "- Niveau de collision élevé, considérez regénérer le tournoi", "", "", false)
		g.pdf.MultiCell(0, 5, "- Vérifiez si le nombre de terrains est adapte au nombre de joueurs", "", "", false)
	}
}

// addRoundsInTwoColumns ajoute les tours en 2 colonnes
func (g *TournamentPDFGenerator) addRoundsInTwoColumns(tournament *tournament.Tournament) {
	startX := margin
	startY := g.pdf.GetY()
	Xmarge := startX
	width, _ := g.pdf.GetPageSize()
	columnWidth := (width - 3.0*Xmarge) / 2.0

	for i, round := range *tournament {
		// Calculer la position
		col := i % 2
		row := i / 2

		x := Xmarge*float64(col+1) + float64(col)*columnWidth
		y := startY + float64(row)*(cellHeight*float64(len((*tournament)[0]))+8.0+margin) // Espacement vertical entre les rangees

		g.pdf.SetXY(x, y)
		g.addCompactRound(round, i+1, columnWidth)
	}
}

// addCompactRound ajoute une ronde compacte
func (g *TournamentPDFGenerator) addCompactRound(round tournament.Round, roundNumber int, width float64) {
	currentX := g.pdf.GetX()
	// currentY := g.pdf.GetY()

	// En-tete de la ronde
	g.pdf.SetFont("Roboto", "", 12)
	g.pdf.SetTextColor(52, 73, 94)
	g.pdf.SetFillColor(236, 240, 241)
	g.pdf.CellFormat(width+0.4, 8, fmt.Sprintf("Parties n°%d", roundNumber), "1", 1, "C", true, 0, "")
	g.pdf.CellFormat(width+0.4, 0.2, "", "1", 1, "C", true, 0, "")

	// Retour à la position X de depart pour les lignes suivantes
	g.pdf.SetX(currentX)

	// Lignes des matchs compactes
	g.pdf.SetFont("Roboto", "", cellHeight*1.8)
	g.pdf.SetTextColor(0, 0, 0)

	for gameIndex, game := range round {
		g.pdf.SetX(currentX)

		// Terrain (plus petit)
		g.pdf.SetTextColor(60, 60, 60)
		g.pdf.SetFillColor(236, 240, 241)
		g.pdf.CellFormat(12, cellHeight, fmt.Sprintf("T%d", gameIndex+1), "1", 0, "C", true, 0, "")
		g.pdf.CellFormat(0.2, cellHeight, "", "1", 0, "C", true, 0, "")

		// Alternance de couleurs
		if gameIndex%2 == 1 {
			g.pdf.SetFillColor(243, 244, 245)
		} else {
			g.pdf.SetFillColor(255, 255, 255)
		}

		g.pdf.SetTextColor(200, 0, 0)
		// Match complet sur le reste de la largeur
		team1Width := (width - 12) / 2 / float64(len(game.Team1))
		for _, player := range game.Team1 {
			g.pdf.CellFormat(team1Width, cellHeight, strconv.Itoa(int(player)+1), "1", 0, "C", true, 0, "")
		}
		g.pdf.CellFormat(0.2, cellHeight, "", "1", 0, "C", true, 0, "")

		g.pdf.SetTextColor(0, 0, 0)
		team2Width := (width - 12) / 2 / float64(len(game.Team2))
		for _, player := range game.Team2 {
			g.pdf.CellFormat(team2Width, cellHeight, strconv.Itoa(int(player)+1), "1", 0, "C", true, 0, "")
		}
		g.pdf.SetY(g.pdf.GetY() + cellHeight)
	}
}

// SavePDF sauvegarde le PDF
func (g *TournamentPDFGenerator) SavePDF(filename string) error {
	return g.pdf.OutputFileAndClose(filename)
}
