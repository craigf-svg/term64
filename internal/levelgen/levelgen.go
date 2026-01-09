package levelgen

import (
	"math/rand"
)

func GetMostCrowdedLevel(width, height, amount int) [][]rune {
	levelsSlice := make([][][]rune, amount)
	for i :=0; i< amount; i++ {
		levelsSlice[i] = GenerateLevel(width, height)
	}

	maxWalls := 0
	maxIndex := 0
	for i := 0; i < amount; i++ {
		wallCount := CountWalls(levelsSlice[i])
		if wallCount > maxWalls {
			maxWalls = wallCount
			maxIndex = i
		}
	}
	
	return levelsSlice[maxIndex]
}

func CountWalls(level [][]rune) int {
	count := 0
	for i := 0; i < len(level); i++ {
		for j := 0; j < len(level[i]); j++ {
			if level[i][j] == '#' {
				count++
			}
		}
	}
	return count
}

func GenerateLevel(width, height int) [][]rune {
	// Initialize with all walls
	level := make([][]rune, height)
	for y := 0; y < height; y++ {
		level[y] = make([]rune, width)
		for x := 0; x < width; x++ {
			level[y][x] = '#'
		}
	}

	// Place start and end positions
	startX, startY := 2, 2
	endX, endY := width-3, height-3
	level[startY][startX] = '.'
	// Start position is static for now 
	// For dynamic generate the start coords and place @ symbol
	// level[startY][startX] = '@'
	level[endY][endX] = '%'

	// Randomly choose to carve right or down from start
	if rand.Intn(2) == 0 {
		level[2][3] = '-' // Guarantee right exit
	} else {
		level[3][2] = '-' // Guarantee down exit
	}

	// Drunkards Walk
	// Rules: 
	// 1. Place a floor tile
	// 2. Pick a random direction you can move in
	// 3. Move onto that tile and place a floor tile 
	// 4. Continue until final location is chosen  

	x, y := startX, startY
	for !(x == endX && y == endY) {  // Loop until we reach the goal
		// Purely leave walkable - for debugging
		if level[y][x] != '-' {
			level[y][x] = '.'  // Place floor tile
		}
		
		// Ensure within bounds
		if x <= 1 {
			x++  // Must move right
		} else if x >= width-2 {
			x--  // Must move left
		} else if y <= 1 {
			y++  // Must move down
		} else if y >= height-2 {
			y--  // Must move up
		} else {
			// Safe to move in random cardinal direction
			// Place floor tile next loop after landing
			direction := rand.Intn(4)
			if direction == 0 {
				x++
			} else if direction == 1 {
				x--
			} else if direction == 2 {
				y++
			} else if direction == 3 {
				y--
			}
		}
	}

	return level
}

