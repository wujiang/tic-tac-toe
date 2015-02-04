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

func TestTTTCnameToRune(t *testing.T) {
	setup()
	assert.Equal(t, tttc.nameToRune(""), ttt.SpecialRune)
	teardown()
}

func TestTTTCcellIsPinnableFalse(t *testing.T) {
	setup()
	assert.False(t, tttc.cellIsPinnable(ttt.Position{ttt.RandInt(3),
		ttt.RandInt(3)}))
	teardown()
}

func TestTTTCcellIsPinnableTrue(t *testing.T) {
	setup()
	tttc.RoundID = "round-id"
	assert.True(t, tttc.cellIsPinnable(ttt.Position{ttt.RandInt(3),
		ttt.RandInt(3)}))
	assert.False(t, tttc.cellIsPinnable(ttt.Position{3, 0}))
	teardown()
}
