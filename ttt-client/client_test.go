package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wujiang/tic-tac-toe"
)

var tttc *TTTClient

func setup() {
	tttc = &TTTClient{
		Name: "Adam",
	}
}

func teardown() {
	tttc = &TTTClient{}
}

func TestisValidPosition(t *testing.T) {
	assert.True(t, isValidPosition(ttt.Position{1, 1}))
	assert.False(t, isValidPosition(ttt.Position{3, 1}))
}

func TestTTTCnameToRune(t *testing.T) {
	setup()
	assert.Equal(t, tttc.nameToRune(""), ttt.SpecialRune)
	teardown()
}
