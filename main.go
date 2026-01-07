package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	levels     [][][]rune
	levelIndex int
	level      [][]rune
	playerX    int
	playerY    int
	catTargetX int
	catTargetY int
	catEarned  bool
	catX       int
	catY       int
	width      int
	height     int
	win        bool
}

func initialModel() model {
	level1 := [][]rune{
		[]rune("####################"),
		[]rune("#..................#"),
		[]rune("#...####...........#"),
		[]rune("#...#..#...........#"),
		[]rune("#...#..#####.......#"),
		[]rune("#..................#"),
		[]rune("#.......######.....#"),
		[]rune("#.......#....#...X.#"),
		[]rune("#.......#....#.....#"),
		[]rune("#..................#"),
		[]rune("####################"),
	}

	level2 := [][]rune{
		[]rune("####################"),
		[]rune("#....X.............#"),
		[]rune("#...###............#"),
		[]rune("#..................#"),
		[]rune("#.####.########....#"),
		[]rune("#..................#"),
		[]rune("#.......#....#####.#"),
		[]rune("#.......#...#....-##"),
		[]rune("#.......#....#.###.#"),
		[]rune("#.......#..........#"),
		[]rune("####################"),
	}

	level3 := [][]rune{
		[]rune("####################"),
		[]rune("#..................#"),
		[]rune("#.###.....###......#"),
		[]rune("#.#.......#.#......#"),
		[]rune("#.#.......#.#......#"),
		[]rune("#.###.....###......#"),
		[]rune("#..................#"),
		[]rune("#.......#####......#"),
		[]rune("#.......#...#....X.#"),
		[]rune("#.......#####......#"),
		[]rune("####################"),
	}

	allLevels := [][][]rune{level1, level2, level3}

	return model{
		levels:     allLevels,
		levelIndex: 0,
		level:      allLevels[0],
		playerX:    1,
		playerY:    1,
		catX:       1,
		catY:       1,
		catEarned:  false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.catX, m.catY = m.catTargetX, m.catTargetY
		m.catTargetX, m.catTargetY = m.playerX, m.playerY
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "up", "k":
			if m.isWalkable(m.playerX, m.playerY-1) {
				m.playerY--
			}

		case "down", "j":
			if m.isWalkable(m.playerX, m.playerY+1) {
				m.playerY++
			}

		case "left", "h":
			if m.isWalkable(m.playerX-1, m.playerY) {
				m.playerX--
			}

		case "right", "l":
			if m.isWalkable(m.playerX+1, m.playerY) {
				m.playerX++
			}
		case "r":
			if m.levelIndex == len(m.levels) {
				m.levelIndex = 0
				m.level = m.levels[0]
				m.playerX = 1
				m.playerY = 1
				m.catEarned = true
				m.catX = 1
				m.catY = 1
				m.catTargetX = 1
				m.catTargetY = 1
			}
		}

		if m.levelChange(m.playerX, m.playerY) {
			m.levelIndex++
			if m.levelIndex != len(m.levels) {
				m.level = m.levels[m.levelIndex]
				m.catTargetX, m.catTargetY = m.playerX, m.playerY
				m.catX, m.catY = m.playerX, m.playerY 
				
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) levelChange(x, y int) bool {
	return m.level[y][x] == 'X'
}

func (m model) View() string {
	if m.levelIndex >= len(m.levels) {
		victoryStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Align(lipgloss.Center)

		endMessage := "\n🎉 VICTORY! 🎉\n\nYou completed all levels! "
		if !m.catEarned {
			endMessage += "However...\n\n\n^._.^\n\n  A mysterious cat watched your journey...\n"
			endMessage += "Play again to have them join you!\n\n"	
		} else {
			endMessage += "\n\n\nฅ^•ﻌ•^ฅ\n\n You and your faithful companion made it to safety!\n\n"
		}
		endMessage += "Press 'q' to quit | Press 'r' to restart"
		return victoryStyle.Render(endMessage)
	}

	// Define styles
	wallStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	floorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666"))

	playerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		Bold(true)

	enemyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		MarginTop(1)

	// Build the game view
	var sb strings.Builder

	for y, row := range m.level {
		for x, cell := range row {
			if x == m.playerX && y == m.playerY {
				sb.WriteString(playerStyle.Render("@"))
			} else if m.catEarned && x == m.catX && y == m.catY {
				sb.WriteString(playerStyle.Render("o"))
			} else if cell == '#' {
				sb.WriteString(wallStyle.Render(string(cell)))
			} else if cell == 'X' {
				sb.WriteString(enemyStyle.Render(string(cell)))
			} else {
				sb.WriteString(floorStyle.Render(string(cell)))
			}
		}
		sb.WriteString("\n")
	}

	// Add help text
	help := helpStyle.Render("Arrow keys or hjkl to move • q/esc to quit")

	return sb.String() + help
}

func (m model) isWalkable(x, y int) bool {
	if y < 0 || y >= len(m.level) || x < 0 || x >= len(m.level[y]) {
		return false
	}

	if m.level[y][x] == 'X' {
		return true
	}

	return m.level[y][x] == '.'
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
