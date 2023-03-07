package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	RandomPattern = "random"
	GliderPattern = "glider"
)

type Game struct {
	Field     [][]bool
	FieldSize int
}

func New(fieldSize int) *Game {
	field := createEmptyField(fieldSize)
	return &Game{FieldSize: fieldSize, Field: field}
}

func (g *Game) seed(pattern string) error {
	switch pattern {
	case RandomPattern:
		for x := range g.Field {
			for y := range g.Field[x] {
				g.Field[x][y] = rand.Intn(2) == 1
			}
		}

	case GliderPattern:
		center := g.FieldSize / 2

		g.Field[center][center] = true
		g.Field[center][center+1] = true
		g.Field[center][center-1] = true
		g.Field[center-1][center+1] = true
		g.Field[center-2][center] = true

	default:
		return fmt.Errorf("[%s] is unsupported game pattern", pattern)
	}

	return nil
}

func (g Game) draw() {
	for _, row := range g.Field {
		for y, cell := range row {
			if y == 0 {
				fmt.Print("| ")
			}

			fmt.Printf(" %s |", Marker(cell).Draw())
		}
		fmt.Print("\n")
	}
}

func (g *Game) makeGeneration() error {
	next := createEmptyField(g.FieldSize)

	for x := 0; x < g.FieldSize; x++ {
		for y := 0; y < g.FieldSize; y++ {
			neigbors := g.countOfLifeNeighbors(x, y)
			next[x][y] = g.Field[x][y]

			if g.Field[x][y] {
				if neigbors < 2 || neigbors > 3 {
					next[x][y] = false
				}
			} else {
				if neigbors == 3 {
					next[x][y] = true
				}
			}

		}
	}

	if isSameField(g.Field, next) {
		return errors.New("evolution stopped")
	}

	g.Field = next

	return nil
}

func (g *Game) countOfLifeNeighbors(x int, y int) int {
	c := 0

	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if x == i && j == y {
				// skip current cell
				continue
			}

			if i > 0 && i < g.FieldSize && j > 0 && j < g.FieldSize {
				if g.Field[i][j] {
					c++
				}
			}
		}
	}

	return c
}

func createEmptyField(size int) (f [][]bool) {
	field := make([][]bool, size)

	for x := range field {
		field[x] = make([]bool, size)
	}

	return field
}

func isSameField(f1 [][]bool, f2 [][]bool) bool {
	for i := range f1 {
		for j := range f1[i] {
			if f1[i][j] != f2[i][j] {
				return false
			}
		}
	}
	return true
}

type Marker bool

func (m Marker) Draw() string {
	if m {
		return "O"
	}

	return "."
}

/*
|--------------------------------------------------------------------------
| Basic Game Mechanics
|--------------------------------------------------------------------------
|
| Create game field X * X slots
| Apply zero generation according some pattern of live cells (Main Storage)
| Fill the game field
| Draw the field based on the state of the main storage
| Make a new generation (clone main storage for new generation)
| Apply generation changes(replace main storage state of temp storage state)
| Draw the field
|
| Rules of cell live or dead:
|  Any live cell with fewer than two live neighbors dies as if caused by underpopulation.
|  Any live cell with two or three live neighbors lives on to the next generation.
|  Any live cell with more than three live neighbors dies, as if by overcrowding.
|  Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction.
|
| Describe rules in simple "if statements" way:
|   Live (1) cell:
|       if have < 2 live neighbors
|               or
|               > 3 live neighbors
|               change state to died (0)
|   Died (0) cell:
|      if have == 3 live neighbors
|           change state to life (1)
|
*/
func main() {
	fieldSize, pattern, inputOk := getUserInput()

	if !inputOk {
		return
	}

	nGame := New(fieldSize)

	if err := nGame.seed(pattern); err != nil {
		fmt.Println(err)
		return
	}

	gameLoop(*nGame)
}

func gameLoop(g Game) {
	for {
		clear()
		g.draw()

		if err := g.makeGeneration(); err != nil {
			fmt.Println("The End.")
			fmt.Println("Survived those who survived :)")
			return
		}

		time.Sleep(260 * time.Millisecond)
	}
}

func clear() {
	clear := make(map[string]func())

	clear["linux"] = func() {
		cmd := exec.Command("clear") // Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") // Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	clear[runtime.GOOS]()
}

func getUserInput() (fs int, p string, ok bool) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please enter a field size from 10 to 25[default]:")
	scanner.Scan()

	fieldSize := 25
	sizeStr := strings.TrimSpace(scanner.Text())

	if len(sizeStr) > 0 {
		size, err := strconv.Atoi(sizeStr)

		if err != nil || size < 10 || size > 25 {
			fmt.Println("Invalid field size, please write number between 10 - 25.")

			return 0, "", false
		}

		fieldSize = size
	}

	fmt.Println("Please enter a game patern random[defaut], glider:")
	scanner.Scan()

	pattern := strings.TrimSpace(scanner.Text())

	if pattern == "" {
		pattern = RandomPattern
	} else if pattern != GliderPattern && pattern != RandomPattern {
		fmt.Println("Invalid pattern, please choose random or glider.")

		return 0, "", false
	}

	return fieldSize, pattern, true
}
