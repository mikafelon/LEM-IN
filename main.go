// Déclaration du package principal.
package main

// Importation des paquets nécessaires.
import (
	"bufio"             // Utilisé pour lire des données ligne par ligne.
	"fmt"               // Utilisé pour l'impression de texte à la console.
	"lem-min/algorithm" // Assurez-vous que ce chemin correspond à l'emplacement de votre paquet d'algorithme.
	"os"                // Utilisé pour l'interaction avec le système de fichiers.
	"strconv"           // Utilisé pour la conversion de chaînes en nombres.
	"strings"           // Utilisé pour les manipulations de chaînes de caractères.
)

// Définition de la structure `Ant` pour représenter une fourmi.
type Ant struct {
	ID       int               // Identifiant unique pour la fourmi.
	Position *algorithm.Room   // Position actuelle de la fourmi.
	Path     []*algorithm.Room // Chemin assigné à la fourmi à parcourir.
}

// parseFile lit et analyse le fichier spécifié pour configurer la simulation.
func parseFile(filename string) (int, map[string]*algorithm.Room, *algorithm.Room, *algorithm.Room, error) {
	file, err := os.Open(filename) // Ouvre le fichier spécifié.
	if err != nil {
		// Retourne une erreur si le fichier ne peut pas être ouvert.
		return 0, nil, nil, nil, err
	}
	defer file.Close() // S'assure que le fichier sera fermé à la fin de la fonction.

	scanner := bufio.NewScanner(file)         // Crée un scanner pour lire le fichier ligne par ligne.
	rooms := make(map[string]*algorithm.Room) // Map pour stocker les salles par leur nom.
	var startRoom, endRoom *algorithm.Room    // Variables pour stocker les salles de départ et d'arrivée.
	var antsCount int                         // Nombre de fourmis.
	firstLine := true                         // Indicateur pour traiter la première ligne différemment.

	// Lit le fichier ligne par ligne.
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text()) // Enlève les espaces en début et en fin de ligne.
		if firstLine {                            // Traite la première ligne pour obtenir le nombre de fourmis.
			antsCount, err = strconv.Atoi(line) // Convertit la première ligne en nombre.
			if err != nil {
				// Retourne une erreur si la première ligne ne peut pas être convertie en nombre.
				return 0, nil, nil, nil, fmt.Errorf("invalid number of ants: %v", err)
			}
			firstLine = false // Indique que la première ligne a été traitée.
			continue
		}

		// Traite les lignes définissant les salles et les connexions.
		if line == "" || strings.HasPrefix(line, "#") { // Ignore les lignes vides et les commentaires.
			if line == "##start" { // Marque la salle de départ.
				scanner.Scan() // Lit la ligne suivante pour obtenir le nom de la salle de départ.
				startLine := strings.TrimSpace(scanner.Text())
				parts := strings.Split(startLine, " ")
				startRoom = &algorithm.Room{Name: parts[0]}
				rooms[startRoom.Name] = startRoom
			} else if line == "##end" { // Marque la salle d'arrivée.
				scanner.Scan() // Lit la ligne suivante pour obtenir le nom de la salle d'arrivée.
				endLine := strings.TrimSpace(scanner.Text())
				parts := strings.Split(endLine, " ")
				endRoom = &algorithm.Room{Name: parts[0]}
				rooms[endRoom.Name] = endRoom
			}
			continue
		}

		// Traite les lignes définissant les salles et les tunnels.
		parts := strings.Split(line, " ")
		if len(parts) == 3 { // Traite une ligne définissant une salle.
			room := &algorithm.Room{Name: parts[0]}
			rooms[room.Name] = room
		} else if len(parts) == 1 { // Traite une ligne définissant un tunnel.
			tunnel := strings.Split(line, "-")
			room1, room2 := rooms[tunnel[0]], rooms[tunnel[1]]
			room1.Adjacent = append(room1.Adjacent, room2)
			room2.Adjacent = append(room2.Adjacent, room1)
		}
	}

	// Retourne le nombre de fourmis, la map des salles, la salle de départ, la salle d'arrivée et l'erreur si présente.
	return antsCount, rooms, startRoom, endRoom, nil
}

// simulateAnts simule le mouvement des fourmis de la salle de départ à la salle d'arrivée.
func simulateAnts(antsCount int, startRoom, endRoom *algorithm.Room, rooms map[string]*algorithm.Room) {
	paths := algorithm.FindPaths(startRoom, endRoom) // Trouve tous les chemins possibles de départ à arrivée.
	if len(paths) == 0 {                             // Vérifie s'il existe des chemins.
		fmt.Println("Aucun chemin trouvé")
		return
	}

	// Initialise chaque salle comme non occupée.
	for _, room := range rooms {
		room.Occupied = false
	}
	startRoom.Occupied = true // La salle de départ est considérée comme occupée.

	// Initialise les fourmis et les assigne à des chemins.
	ants := make([]*Ant, antsCount)
	for i := range ants {
		pathIndex := i % len(paths) // Assigne un chemin à chaque fourmi.
		ants[i] = &Ant{ID: i + 1, Position: startRoom, Path: paths[pathIndex]}
		if ants[i].Position != startRoom && ants[i].Position != endRoom {
			ants[i].Position.Occupied = true // Marque la position actuelle de la fourmi comme occupée.
		}
	}

	// Débute la simulation du mouvement des fourmis.
	for !allAntsFinished(ants, endRoom) {
		var moves []string // Pour stocker les mouvements des fourmis à chaque tour.
		for _, ant := range ants {
			if ant.Position != endRoom { // Continue si la fourmi n'est pas encore arrivée à la fin.
				nextPosIndex := getNextPositionIndex(ant) // Obtient l'index de la prochaine position sur le chemin.
				if nextPosIndex < len(ant.Path) {         // Vérifie si la fourmi peut avancer.
					nextRoom := ant.Path[nextPosIndex]             // Obtient la prochaine salle sur le chemin.
					if !nextRoom.Occupied || nextRoom == endRoom { // Vérifie si la prochaine salle est libre ou est la salle de fin.
						if ant.Position != startRoom { // Libère la salle actuelle si ce n'est pas la salle de départ.
							ant.Position.Occupied = false
						}
						ant.Position = nextRoom                                                 // Déplace la fourmi vers la prochaine salle.
						ant.Position.Occupied = true                                            // Marque la nouvelle salle comme occupée.
						moves = append(moves, fmt.Sprintf("L%d-%s", ant.ID, ant.Position.Name)) // Ajoute le mouvement à la liste.
					}
				}
			}
		}

		// Affiche les mouvements des fourmis pour ce tour.
		fmt.Println(strings.Join(moves, " "))
		if len(moves) == 0 { // Si aucune fourmi ne bouge, la simulation est terminée.
			break
		}
	}
}

// getNextPositionIndex trouve l'indice de la prochaine position pour la fourmi sur son chemin.
func getNextPositionIndex(ant *Ant) int {
	for idx, room := range ant.Path { // Parcourt le chemin de la fourmi.
		if room == ant.Position { // Trouve la position actuelle de la fourmi sur le chemin.
			return idx + 1 // Retourne l'indice de la prochaine position.
		}
	}
	return -1 // Retourne -1 si la fourmi est à la fin du chemin.
}

// allAntsFinished vérifie si toutes les fourmis sont arrivées à la salle de fin.
func allAntsFinished(ants []*Ant, endRoom *algorithm.Room) bool {
	for _, ant := range ants { // Parcourt toutes les fourmis.
		if ant.Position != endRoom { // Vérifie si une fourmi n'est pas encore arrivée à la fin.
			return false // Retourne false si une fourmi n'est pas arrivée.
		}
	}
	return true // Retourne true si toutes les fourmis sont arrivées.
}

// La fonction main est le point d'entrée du programme.
func main() {
	filename := "txt/fourmie.txt"                                    // Définit le chemin du fichier à lire.
	antsCount, rooms, startRoom, endRoom, err := parseFile(filename) // Appelle parseFile pour lire et analyser le fichier.
	if err != nil {
		fmt.Println("Error reading file:", err) // Affiche une erreur si le fichier ne peut pas être lu.
		return
	}

	// Démarre la simulation du mouvement des fourmis.
	simulateAnts(antsCount, startRoom, endRoom, rooms)
}
