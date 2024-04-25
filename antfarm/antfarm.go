package antfarm

import (
	"fmt"
	structure "lem/strucdure" // Assurez-vous que le chemin du package est correct
	"sort"
)

// FindAllPathsBFS trouve tous les chemins possibles de 'start' à 'end' en utilisant l'algorithme de recherche en largeur (BFS).
func FindAllPathsBFS(farm *structure.AntFarm) ([][]*structure.Room, error) {
	// Vérifie si la salle de départ ou d'arrivée n'est pas définie
	if farm.Start == nil || farm.End == nil {
		return nil, fmt.Errorf("salle de départ ou de fin non définie")
	}

	var allPaths [][]*structure.Room
	paths, err := bfs(farm.Start, farm.End) // Exécute BFS pour trouver les chemins
	if err != nil {
		return nil, err // Gère correctement les erreurs rencontrées
	}
	allPaths = append(allPaths, paths...) // Ajoute les chemins trouvés à la liste des tous les chemins

	return allPaths, nil
}

// bfs réalise une recherche en largeur à partir de la salle de départ jusqu'à la salle d'arrivée.
func bfs(start, end *structure.Room) ([][]*structure.Room, error) {
	var allPaths [][]*structure.Room
	queue := [][]*structure.Room{{start}} // Initialise la file avec la salle de départ

	for len(queue) > 0 {
		path := queue[0]  // Prend le premier chemin de la file
		queue = queue[1:] // Enlève ce chemin de la file

		lastRoom := path[len(path)-1] // Dernière salle du chemin
		if lastRoom == end {
			allPaths = append(allPaths, append([]*structure.Room(nil), path...)) // Ajoute le chemin complet si la fin est atteinte
			continue
		}

		// Parcourt les salles adjacentes à la dernière salle du chemin
		for _, nextRoom := range lastRoom.Adjacent {
			if !isInPath(path, nextRoom) { // Vérifie si la salle suivante n'est pas déjà dans le chemin pour éviter les cycles
				newPath := append([]*structure.Room(nil), path...) // Crée un nouveau chemin en ajoutant la salle suivante
				newPath = append(newPath, nextRoom)
				queue = append(queue, newPath) // Ajoute le nouveau chemin à la file
			}
		}
	}
	return allPaths, nil
}

// filterAndSortPaths filtre et trie les chemins pour ne garder que ceux qui ont au maximum 7 étapes.
func filterAndSortPaths(paths [][]*structure.Room) [][]*structure.Room {
	var filteredPaths [][]*structure.Room

	// Filtre les chemins selon leur longueur
	for _, path := range paths {
		if len(path) <= 8 { // Inclut la salle de départ et la salle d'arrivée
			filteredPaths = append(filteredPaths, path)
		}
	}

	// Trie les chemins filtrés par longueur croissante
	sort.Slice(filteredPaths, func(i, j int) bool {
		return len(filteredPaths[i]) < len(filteredPaths[j])
	})

	return filteredPaths
}

// isInPath vérifie si une salle est déjà présente dans le chemin pour éviter les cycles.
func isInPath(path []*structure.Room, room *structure.Room) bool {
	for _, p := range path {
		if p == room {
			return true
		}
	}
	return false
}

// printPaths affiche les chemins trouvés de manière lisible.
func PrintPaths(paths [][]*structure.Room) {
	if len(paths) == 0 {
		fmt.Println("Aucun chemin ne répond aux critères.")
		return
	}
	for i, path := range paths {
		fmt.Printf("Chemin %d (longueur %d) : ", i+1, len(path))
		for j, room := range path {
			if j > 0 {
				fmt.Print(" -> ")
			}
			fmt.Print(room.Name)
		}
		fmt.Println()
	}
}

// FindIndependentPaths retourne toutes les combinaisons de chemins où aucun chemin ne partage une salle avec un autre.
func FindIndependentPaths(allPaths [][]*structure.Room) [][][]*structure.Room {
	var result [][][]*structure.Room
	// Fonction récursive pour générer des combinaisons
	var findCombinations func(index int, current [][]*structure.Room)
	findCombinations = func(index int, current [][]*structure.Room) {
		// Quand une combinaison est formée, l'ajoute au résultat
		if len(current) > 1 {
			result = append(result, current)
		}
		for i := index; i < len(allPaths); i++ {
			// Vérifie si le chemin actuel peut être ajouté sans chevauchement
			canAdd := true
			for _, existingPath := range current {
				if hasCommonRooms(existingPath, allPaths[i]) {
					canAdd = false
					break
				}
			}
			if canAdd {
				// Ajoute le chemin et trouve plus de combinaisons
				findCombinations(i+1, append(current, allPaths[i]))
			}
		}
	}

	// Initialisation de l'appel récursif avec un ensemble vide
	findCombinations(0, [][]*structure.Room{})
	return result
}

// hasCommonRooms vérifie si deux chemins partagent des salles communes.
func hasCommonRooms(path1, path2 []*structure.Room) bool {
	roomSet := make(map[int]bool) // Utilise l'ID de la salle pour l'unicité
	for _, room := range path1 {
		roomSet[room.ID] = true
	}
	for _, room := range path2 {
		if roomSet[room.ID] {
			return true
		}
	}
	return false
}
