package chemain

import (
	"bufio"
	"fmt"
	"lem/strucdure"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Helper function to convert string to int safely
func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func ParseAntFarm(file *os.File) (*strucdure.AntFarm, error) {
	scanner := bufio.NewScanner(file)
	farm := &strucdure.AntFarm{
		Rooms:        make(map[string]*strucdure.Room),
		Intermediate: make(map[string]*strucdure.Room),
		Start:        nil,
		End:          nil,
	}

	var numberAnts int
	var err error
	currentRoomType := ""
	firstLine := true
	var pendingLinks []string
	roomID := 1

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if firstLine {
			numberAnts, err = strconv.Atoi(line)
			if err != nil {
				return nil, fmt.Errorf("error converting number of ants: %v", err)
			}
			farm.Ants = numberAnts
			firstLine = false
			continue
		}

		if line == "" || line[0] == '#' {
			if line == "##start" || line == "##end" {
				currentRoomType = line[2:]
				continue
			}
		} else if strings.Contains(line, "-") {
			pendingLinks = append(pendingLinks, line)
		} else {
			parts := strings.Fields(line)
			if len(parts) < 3 {
				return nil, fmt.Errorf("invalid room format")
			}
			if _, exists := farm.Rooms[parts[0]]; exists {
				return nil, fmt.Errorf("duplicate room name found: %s", parts[0])
			}
			room := &strucdure.Room{
				Name:     parts[0],
				X:        atoi(parts[1]),
				Y:        atoi(parts[2]),
				Adjacent: []*strucdure.Room{},
				ID:       roomID,
			}
			farm.Rooms[room.Name] = room
			switch currentRoomType {
			case "start":
				if farm.Start != nil {
					return nil, fmt.Errorf("multiple start rooms defined")
				}
				farm.Start = room
			case "end":
				if farm.End != nil {
					return nil, fmt.Errorf("multiple end rooms defined")
				}
				farm.End = room
			default:
				farm.Intermediate[room.Name] = room
			}
			roomID++
			currentRoomType = ""
		}
	}

	if err = ApplyLinks(farm, pendingLinks); err != nil {
		return nil, err
	}

	if farm.Start == nil || farm.End == nil {
		return nil, fmt.Errorf("start or end room not defined")
	}

	return farm, nil
}

func ApplyLinks(farm *strucdure.AntFarm, pendingLinks []string) error {
	links := make(map[string]bool)
	for _, link := range pendingLinks {
		parts := strings.Split(link, "-")
		if len(parts) != 2 {
			return fmt.Errorf("invalid link format: %s", link)
		}
		linkKey := parts[0] + "-" + parts[1]
		if links[linkKey] {
			return fmt.Errorf("duplicate link found: %s", link)
		}
		room1, room2 := farm.Rooms[parts[0]], farm.Rooms[parts[1]]
		if room1 == nil || room2 == nil {
			return fmt.Errorf("undefined room referenced in link: %s", link)
		}
		room1.Adjacent = append(room1.Adjacent, room2)
		room2.Adjacent = append(room2.Adjacent, room1)
		links[linkKey] = true
		links[parts[1]+"-"+parts[0]] = true // Ensure link is considered both ways
	}
	return nil
}

var nextRoomID int = 1 // Commencez à 1 pour éviter d'initialiser à 0, sauf si 0 est une valeur valide dans votre contexte.

func NewRoom(name string, x, y int) *strucdure.Room {
	room := &strucdure.Room{
		Name: name,
		X:    x,
		Y:    y,
		ID:   nextRoomID, // Assurez-vous que cet ID est unique
	}
	nextRoomID++ // Incrémenter l'ID pour la prochaine salle créée
	return room
}

func PrintPaths(paths [][]*strucdure.Room) {
	for i, path := range paths {
		fmt.Printf("Chemin %d: ", i+1)
		for _, room := range path {
			fmt.Printf("%d ", room.ID) // Utiliser l'ID de la salle au lieu de son nom
		}
		fmt.Println()
	}
}

func PrintFarmDetails(farm *strucdure.AntFarm) {
	if farm == nil {
		fmt.Println("No farm data available.")
		return
	}

	// Affichage des informations de base
	fmt.Printf("Total Ants: %d\n", farm.Ants)
	fmt.Println("Rooms:")
	for name, room := range farm.Rooms {
		fmt.Printf("Room: %s, Coordinates: (%d, %d)\n", name, room.X, room.Y)
	}

	// Affichage des salles intermédiaires, si nécessaire
	fmt.Println("Intermediate Rooms:")
	for name, room := range farm.Intermediate {
		fmt.Printf("Intermediate Room: %s, Coordinates: (%d, %d)\n", name, room.X, room.Y)
	}

	// Affichage des salles de départ et d'arrivée
	if farm.Start != nil {
		fmt.Printf("Start Room: %s, Coordinates: (%d, %d)\n", farm.Start.Name, farm.Start.X, farm.Start.Y)
	}
	if farm.End != nil {
		fmt.Printf("End Room: %s, Coordinates: (%d, %d)\n", farm.End.Name, farm.End.X, farm.End.Y)
	}
}

// PathSlice est un type pour faciliter le tri des chemins
type PathSlice [][]*strucdure.Room

// Implement sort.Interface for PathSlice
func (p PathSlice) Len() int {
	return len(p)
}

func (p PathSlice) Less(i, j int) bool {
	return len(p[i]) < len(p[j]) // Trie les chemins par longueur croissante
}

func (p PathSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// GetAndSortPaths récupère et trie tous les chemins de la ferme
func GetAndSortPaths(farm *strucdure.AntFarm) (PathSlice, error) {
	paths, err := StoreAllPaths(farm)
	if err != nil {
		return nil, err
	}

	sortedPaths := PathSlice(paths)
	sort.Sort(sortedPaths)
	return sortedPaths, nil
}

// storeAllPaths trouve tous les chemins de la salle 'start' à la salle 'end' en utilisant BFS.
func StoreAllPaths(farm *strucdure.AntFarm) ([][]*strucdure.Room, error) {
	if farm.Start == nil || farm.End == nil {
		return nil, fmt.Errorf("start or end room not defined")
	}

	allPaths := [][]*strucdure.Room{}
	queue := [][]*strucdure.Room{{farm.Start}}

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]

		lastRoom := path[len(path)-1]
		if lastRoom == farm.End {
			allPaths = append(allPaths, path)
			continue
		}

		for _, nextRoom := range lastRoom.Adjacent {
			if !Contains(path, nextRoom) {
				newPath := append([]*strucdure.Room(nil), path...)
				newPath = append(newPath, nextRoom)
				queue = append(queue, newPath)
			}
		}
	}

	return allPaths, nil
}

// contains vérifie si la salle est déjà dans le chemin pour éviter les cycles.
func Contains(path []*strucdure.Room, room *strucdure.Room) bool {
	for _, p := range path {
		if p == room {
			return true
		}
	}
	return false
}
func PrintRoomDetails(room *strucdure.Room) {
	fmt.Printf("%s at (%d, %d) with connections to: ", room.Name, room.X, room.Y)
	for _, adjacent := range room.Adjacent {
		fmt.Printf("%s (ID %d), ", adjacent.Name, adjacent.ID)
	}
	fmt.Println()
}

// simulateAnts simule les déplacements des fourmis et imprime chaque mouvement.
func SimulateAnts(paths [][]*strucdure.Room, numAnts int) {
	if len(paths) == 0 || numAnts == 0 {
		fmt.Println("No paths or ants available for simulation.")
		return
	}

	// Initialiser la position des fourmis; -1 signifie que la fourmi n'a pas encore démarré.
	antPositions := make([]int, numAnts)
	for i := range antPositions {
		antPositions[i] = -1
	}

	// Continuer jusqu'à ce que toutes les fourmis aient atteint la fin de leurs chemins respectifs.
	completed := false
	turn := 0

	for !completed {
		completed = true
		moves := []string{}

		for i := 0; i < numAnts; i++ {
			if antPositions[i] < len(paths[i%len(paths)])-1 {
				antPositions[i]++
				currentRoom := paths[i%len(paths)][antPositions[i]]
				moves = append(moves, fmt.Sprintf("L%d-%s", i+1, currentRoom.Name))
				if antPositions[i] < len(paths[i%len(paths)])-1 {
					completed = false
				}
			}
		}

		// Afficher les mouvements pour ce tour, si des mouvements ont eu lieu.
		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		}

		turn++
	}
}
