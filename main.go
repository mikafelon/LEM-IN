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
	// Open the file and parse the ant farm
	file, err := os.Open("txt/fourmie.txt") // Make sure the file path is correct
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Parse the ant farm from the file
	farm, err := chemain.ParseAntFarm(file)
	if err != nil {
		fmt.Println("Error parsing ant farm:", err)
		return
	}

	// Print the details of the entire ant farm
	chemain.PrintFarmDetails(farm)

	// Find all possible paths using BFS
	allPaths, err := antfarm.FindAllPathsBFS(farm)
	if err != nil {
		fmt.Println("Error finding paths:", err)
		return
	}
	independentPaths := antfarm.FindIndependentPaths(allPaths)

	// Print the independent path combinations
	for i, paths := range independentPaths {
		fmt.Printf("Combination %d:\n", i+1)
		for _, path := range paths {
			for _, room := range path {
				fmt.Printf("%s (ID %d) -> ", room.Name, room.ID)
			}
			fmt.Println("End of Path")
		}

	}

	// Print all the paths found
	chemain.PrintPaths(allPaths)
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
