package strucdure

import (
	"fmt"
	"strings"
)

type Room struct {
	Name      string
	X         int
	Y         int
	ID        int // Identifiant numérique unique pour chaque salle
	Adjacent  []*Room
	Type      string
	Occupants int // Ajouter un champ pour gérer le nombre de fourmis dans la salle
}

type AntFarm struct {
	Rooms        map[string]*Room
	Intermediate map[string]*Room
	Start        *Room
	End          *Room
	Ants         int
}

func (af *AntFarm) String() string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("AntFarm with %d ants\n", af.Ants))
	result.WriteString("Rooms:\n")
	for _, room := range af.Rooms {
		result.WriteString(fmt.Sprintf(" %s at (%d, %d) with connections to: ", room.Name, room.X, room.Y))
		for _, adj := range room.Adjacent {
			result.WriteString(fmt.Sprintf("%s (ID %d), ", adj.Name, adj.ID))
		}
		result.WriteString(fmt.Sprintf("Occupants: %d\n", room.Occupants))
	}
	if af.Start != nil {
		result.WriteString(fmt.Sprintf("Start Room: %s (ID %d)\n", af.Start.Name, af.Start.ID))
	}
	if af.End != nil {
		result.WriteString(fmt.Sprintf("End Room: %s (ID %d)\n", af.End.Name, af.End.ID))
	}
	return result.String()
}

// Les structures State, Item et PriorityQueue sont correctes et peuvent rester telles quelles.

// Méthode pour initialiser les fourmis dans la salle de départ.
func (af *AntFarm) InitializeAnts() {
	af.Start.Occupants = af.Ants
}

// Méthode pour simuler le déplacement des fourmis.
func (af *AntFarm) MoveAnts() {
	// La logique de déplacement des fourmis doit être améliorée pour gérer les déplacements de toutes les fourmis.
	// L'exemple suivant est une simplification.
	for _, room := range af.Rooms {
		for _, adj := range room.Adjacent {
			// Logique simplifiée de déplacement vers des salles adjacentes non occupées
			if room.Occupants > 0 && adj.Occupants == 0 && adj != af.Start && adj != af.End {
				moveCount := room.Occupants
				adj.Occupants += moveCount
				room.Occupants -= moveCount
				fmt.Printf("Moved %d ants from %s to %s\n", moveCount, room.Name, adj.Name)
			}
		}
	}
}

// Méthode pour afficher l'état actuel de la ferme de fourmis.
func (af *AntFarm) PrintCurrentState() {
	fmt.Println("État actuel de la fourmilière :")
	for _, room := range af.Rooms {
		fmt.Printf("Salle: %s, Occupants: %d\n", room.Name, room.Occupants)
	}
}
