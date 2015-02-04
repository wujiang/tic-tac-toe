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
