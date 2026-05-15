package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMaze(t *testing.T) {
	maze := newMaze(Player{})

	assert.NotNil(t, maze.Id)
	assert.Equal(t, 1, len(maze.Grid))
	assert.Equal(t, 1, len(maze.Grid[0]))
	assert.Equal(t, "cell-0-0", maze.Grid[0][0].Id)
	assert.Equal(t, "Starting cell", maze.Grid[0][0].Description)
	assert.Equal(t, []string{}, maze.Grid[0][0].Items)
	assert.Equal(t, Position{Row: 0, Col: 0}, maze.Current)
	assert.Empty(t, maze.Player.Items)
	assert.Equal(t, []string{maze.Grid[0][0].Id}, maze.Visited)
}
