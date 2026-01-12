package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"term64/internal/levelgen"
)

type model struct {
	levels      [][][]rune
	levelIndex  int
	level       [][]rune
	playerX     int
	playerY     int
	catTargetX  int
	catTargetY  int
	catEarned   bool
	catX        int
	catY        int
	width       int
	height      int
	win         bool
	hasKey      bool
	levelHasKey []bool
}

func initialModel() model {
	level0 := levelgen.GetMostCrowdedLevel(19, 11, 1, false)
	level1 := levelgen.GetMostCrowdedLevel(19, 11, 2, false)
	level2 := levelgen.GetMostCrowdedLevel(19, 11, 3, true)
	level3 := levelgen.GetMostCrowdedLevel(19, 11, 50, true)

	allLevels := [][][]rune{level0, level1, level2, level3}
	playerY, playerX := findPlayerStart(allLevels[0])

	return model{
		levels:      allLevels,
		levelIndex:  0,
		level:       allLevels[0],
		playerX:     playerX,
		playerY:     playerY,
		catX:        1,
		catY:        1,
		catEarned:   false,
		hasKey:      false,
		levelHasKey: []bool{false, false, true, true},
	}
}

func findPlayerStart(runeMap [][]rune) (row, col int) {
	/* Static start for now at 2,2
	for i := 0; i < len(runeMap); i++ {
			for j := 0; j < len(runeMap[i]); j++ {
					if runeMap[i][j] == 's' {
							return i, j
					}
			}
	}*/
	return 2, 2
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
				m.playerX = 2
				m.playerY = 2
				m.catEarned = true
				m.catX = 2
				m.catY = 2
				m.catTargetX = 2
				m.catTargetY = 2
				m.hasKey = false
			}
		}

		if m.level[m.playerY][m.playerX] == '⚷' {
			m.hasKey = true
		}

		if m.shouldChangeLevel(m.playerX, m.playerY) {
			if m.hasKey || (m.levelHasKey[m.levelIndex] == false) {
				m.levelIndex++
				if m.levelIndex < len(m.levels) {
					m.level = m.levels[m.levelIndex]
					m.playerY, m.playerX = findPlayerStart(m.level)
					m.catTargetX, m.catTargetY = m.playerX, m.playerY
					m.catX, m.catY = m.playerX, m.playerY
					m.hasKey = false
				}
				// else: player completed all levels, display victory message
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) shouldChangeLevel(x, y int) bool {
	return m.level[y][x] == '%'
}

func (m model) View() string {
	if m.levelIndex >= len(m.levels) {
		victoryStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Align(lipgloss.Center)

		catFaceStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Align(lipgloss.Center)
		endMessage := ""
		if !m.catEarned {
			endMessage += "\nYou notice a face peering from the darkness\n\n" + catFaceStyle.Render("^._.^")
			endMessage += "\n\nA mysterious cat was watching your journey.\n"
			endMessage += "Show them the way out?\n\n"
			endMessage += "Press 'r' to retrive cat"
		} else {
			endMessage += "\n🎉 VICTORY! 🎉\n\nYou saved the cat! "
			endMessage += "\n\n\n" + catFaceStyle.Render("ฅ^•ﻌ•^ฅ") + "\n\n You and your faithful companion made it!\n\n"
			endMessage += "Press 'q' to quit | Press 'r' to restart"
		}
		return victoryStyle.Render(endMessage)
	}

	// Define styles
	wallStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))
	floorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666"))
	stairsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Bold(true)
	lockedStairsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4A5568")).
		Bold(false)

	playerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		Bold(true)
	catStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFA500")).
		Bold(true)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C0C0C0")).
		Bold(false)

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
				sb.WriteString(catStyle.Render("o"))
			} else if cell == '#' {
				sb.WriteString(wallStyle.Render(string(cell)))
			} else if cell == '%' {
				if !m.hasKey && m.levelHasKey[m.levelIndex] {
					sb.WriteString(lockedStairsStyle.Render(string(cell)))
				} else {
					sb.WriteString(stairsStyle.Render(string(cell)))
				}
			} else if cell == '⚷' {
				if m.hasKey {
					sb.WriteString(floorStyle.Render("."))
				} else {
					sb.WriteString(keyStyle.Render("⚷"))
				}
			} else if cell == '@' {
				// Start of level marker
				sb.WriteString(floorStyle.Render("s"))
			} else {
				sb.WriteString(floorStyle.Render(string(cell)))
			}
		}
		sb.WriteString("\n")
	}

	// Add instruction text
	instruction := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00CC00")).
		Bold(true).
		Render("Get to the stairs (%) in all levels to escape!")

	// Add help text
	help := helpStyle.Render("Arrow keys or hjkl to move • q/esc to quit")

	if m.hasKey {
		help += helpStyle.Render("Inventory - ⚷")
	} else {
		help += helpStyle.Render("")
	}

	return sb.String() + help + "\n\n" + instruction
}

func (m model) isWalkable(x, y int) bool {
	if y < 0 || y >= len(m.level) || x < 0 || x >= len(m.level[y]) {
		return false
	}

	if m.level[y][x] == '⚷' || m.level[y][x] == '%' || m.level[y][x] == '@' || m.level[y][x] == 's' || m.level[y][x] == '-' {
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
