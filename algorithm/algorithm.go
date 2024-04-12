// Déclaration du package 'algorithm'.
package algorithm

// Room définit la structure d'une salle dans la fourmilière.
type Room struct {
	Name     string  // Nom de la salle, utilisé comme identifiant unique.
	Adjacent []*Room // Liste des salles adjacentes, représentant les tunnels vers les autres salles.
	Occupied bool    // Indicateur pour savoir si la salle est occupée par une fourmi.
}

// FindPaths trouve tous les chemins possibles de la salle 'start' à la salle 'end'.
func FindPaths(start *Room, end *Room) [][]*Room {
	var result [][]*Room            // Résultat contenant tous les chemins possibles.
	visited := make(map[*Room]bool) // Map pour garder une trace des salles déjà visitées.
	queue := [][]*Room{{start}}     // File d'attente pour explorer les chemins, commence avec le chemin contenant seulement la salle de départ.

	// Continue tant qu'il y a des chemins à explorer dans la file d'attente.
	for len(queue) > 0 {
		path := queue[0]  // Prend le premier chemin de la file d'attente.
		queue = queue[1:] // Supprime ce chemin de la file d'attente.

		lastRoom := path[len(path)-1] // Obtient la dernière salle du chemin courant.
		if lastRoom == end {          // Vérifie si le chemin courant atteint la salle de fin.
			result = append(result, path) // Ajoute le chemin au résultat.
			continue                      // Continue avec le prochain chemin dans la file d'attente.
		}

		visited[lastRoom] = true // Marque la dernière salle comme visitée.

		// Parcourt les salles adjacentes à la dernière salle du chemin courant.
		for _, nextRoom := range lastRoom.Adjacent {
			if visited[nextRoom] { // Ignore les salles déjà visitées pour éviter les cycles.
				continue
			}
			newPath := append([]*Room{}, path...) // Crée un nouveau chemin en copiant le chemin courant.
			newPath = append(newPath, nextRoom)   // Ajoute la salle adjacente au nouveau chemin.
			queue = append(queue, newPath)        // Ajoute le nouveau chemin à la file d'attente pour exploration ultérieure.
		}
	}

	return result // Retourne tous les chemins trouvés de 'start' à 'end'.
}

// contains vérifie si une salle spécifique 'room' est déjà dans un chemin 'path'.
func contains(path []*Room, room *Room) bool {
	for _, r := range path { // Parcourt chaque salle dans le chemin.
		if r == room { // Vérifie si la salle courante est la salle spécifiée.
			return true // Retourne vrai si la salle est trouvée dans le chemin.
		}
	}
	return false // Retourne faux si la salle n'est pas dans le chemin.
}
