package main

import (
	"fmt"
	"lem/antfarm"
	"lem/chemain"
	structure "lem/strucdure"
	"os"
	"sort"
)

// PrintFarmDetails affiche des détails sur l'AntFarm fourni, y compris les liens entre les salles.
func PrintFarmDetails(farm *structure.AntFarm) {
	if farm == nil {
		fmt.Println("No farm data available.")
		return
	}

	fmt.Println("Detailed Information about the Ant Farm:")
	fmt.Printf("Total number of ants: %d\n", farm.Ants)
	fmt.Println("Rooms and their details:")

	for _, room := range farm.Rooms {
		fmt.Printf("Room: %s at (%d, %d)\n", room.Name, room.X, room.Y)
		printRoomConnections(room)
	}

	PrintLinks(farm)
}

// printRoomConnections affiche les connexions d'une salle spécifique.
func printRoomConnections(room *structure.Room) {
	if len(room.Adjacent) == 0 {
		fmt.Println("  No connections.")
		return
	}
	fmt.Print("  Connections to: ")
	for i, adj := range room.Adjacent {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%s", adj.Name)
	}
	fmt.Println()
}

// printLinks affiche tous les liens uniques entre les salles.
func PrintLinks(farm *structure.AntFarm) {
	fmt.Println("Unique Links between Rooms:")
	links := make(map[string]bool)
	for _, room := range farm.Rooms {
		for _, adj := range room.Adjacent {
			link := fmt.Sprintf("%s-%s", room.Name, adj.Name)
			reverseLink := fmt.Sprintf("%s-%s", adj.Name, room.Name)
			if _, exists := links[reverseLink]; !exists {
				links[link] = true
				fmt.Println(link)
			}
		}
	}
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run . <path_to_input_file>")
        return
    }

    inputFilePath := os.Args[1]

    file, err := os.Open(inputFilePath)
    if err != nil {
        fmt.Printf("Failed to open the file: %s\n", err)
        return
    }
    defer file.Close()

    farm, err := chemain.ParseAntFarm(file)
    if err != nil {
        fmt.Printf("Error parsing ant farm: %s\n", err)
        return
    }

    paths, err := antfarm.FindAllPathsBFS(farm)
    if err != nil {
        fmt.Printf("Error finding paths using BFS: %s\n", err)
        return
    }

    filteredPaths := filterAndSortPaths(paths)
    chemain.SimulateAnts(filteredPaths, farm.Ants)
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

// printPaths affiche les chemins trouvés
func printPaths(paths [][]*structure.Room) {
	if len(paths) == 0 {
		fmt.Println("No paths found.")
		return
	}
	for i, path := range paths {
		fmt.Printf("Path %d: ", i+1)
		for j, room := range path {
			if j > 0 {
				fmt.Print(" -> ")
			}
			fmt.Print(room.Name)
		}
		fmt.Println()
	}
}
