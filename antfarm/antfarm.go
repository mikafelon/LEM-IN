package antfarm

import (
	"fmt"
	structure "lem/strucdure" // Assurez-vous que le chemin du package est correct
	"sort"
)

// FindAllPathsBFS trouve tous les chemins de 'start' à 'end' à travers des salles intermédiaires en utilisant BFS.
func FindAllPathsBFS(farm *structure.AntFarm) ([][]*structure.Room, error) {
	if farm.Start == nil || farm.End == nil {
		return nil, fmt.Errorf("start or end room not defined")
	}

	var allPaths [][]*structure.Room
	paths, err := bfs(farm.Start, farm.End) // Appel à BFS avec un seul start et end
	if err != nil {
		return nil, err // Gérer l'erreur correctement
	}
	allPaths = append(allPaths, paths...)

	return allPaths, nil
}

func bfs(start, end *structure.Room) ([][]*structure.Room, error) {
	var allPaths [][]*structure.Room
	queue := [][]*structure.Room{{start}}

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]

		lastRoom := path[len(path)-1]
		if lastRoom == end {
			allPaths = append(allPaths, append([]*structure.Room(nil), path...))
			continue
		}

		for _, nextRoom := range lastRoom.Adjacent {
			if !isInPath(path, nextRoom) { // Check if nextRoom is already in the current path
				newPath := append([]*structure.Room(nil), path...)
				newPath = append(newPath, nextRoom)
				queue = append(queue, newPath)
			}
		}
	}
	return allPaths, nil
}

// filterAndSortPaths filtre et trie les chemins pour ne garder que ceux avec au maximum 7 étapes.
func filterAndSortPaths(paths [][]*structure.Room) [][]*structure.Room {
	var filteredPaths [][]*structure.Room

	// Filtrer les chemins
	for _, path := range paths {
		if len(path) <= 8 { // inclut la salle de départ et la salle d'arrivée dans le compte
			filteredPaths = append(filteredPaths, path)
		}
	}

	// Trier les chemins par longueur croissante
	sort.Slice(filteredPaths, func(i, j int) bool {
		return len(filteredPaths[i]) < len(filteredPaths[j])
	})

	return filteredPaths
}

// Helper function to check if room is already in path to prevent cycles
func isInPath(path []*structure.Room, room *structure.Room) bool {
	for _, p := range path {
		if p == room {
			return true
		}
	}
	return false
}

// printPaths affiche les chemins trouvés
func PrintPaths(paths [][]*structure.Room) {
	if len(paths) == 0 {
		fmt.Println("No paths meet the criteria.")
		return
	}
	for i, path := range paths {
		fmt.Printf("Path %d (length %d): ", i+1, len(path))
		for j, room := range path {
			if j > 0 {
				fmt.Print(" -> ")
			}
			fmt.Print(room.Name)
		}
		fmt.Println()
	}
}

// FindIndependentPaths returns all combinations of paths where no two paths share any room.
func FindIndependentPaths(allPaths [][]*structure.Room) [][][]*structure.Room {
	var result [][][]*structure.Room
	// Use a recursive function to generate combinations
	var findCombinations func(index int, current [][]*structure.Room)
	findCombinations = func(index int, current [][]*structure.Room) {
		// When a combination is formed, append it to result
		if len(current) > 1 {
			result = append(result, current)
		}
		for i := index; i < len(allPaths); i++ {
			// Check if the current path can be added without overlapping
			canAdd := true
			for _, existingPath := range current {
				if hasCommonRooms(existingPath, allPaths[i]) {
					canAdd = false
					break
				}
			}
			if canAdd {
				// Recursively add path and find further combinations
				findCombinations(i+1, append(current, allPaths[i]))
			}
		}
	}

	// Initialize recursive call with an empty set
	findCombinations(0, [][]*structure.Room{})
	return result
}

// hasCommonRooms checks if two paths have any room in common.
func hasCommonRooms(path1, path2 []*structure.Room) bool {
	roomSet := make(map[int]bool) // Use room ID for uniqueness
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
