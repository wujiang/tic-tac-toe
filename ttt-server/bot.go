package main

import (
	"errors"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/wujiang/tic-tac-toe"
)

var am = InitAIManager()

var BotNames = [50]string{
	"Bulbasaur",
	"Ivysaur",
	"Venusaur",
	"Charmander",
	"Charmeleon",
	"Charizard",
	"Squirtle",
	"Wartortle",
	"Blastoise",
	"Caterpie",
	"Metapod",
	"Butterfree",
	"Weedle",
	"Kakuna",
	"Beedrill",
	"Pidgey",
	"Pidgeotto",
	"Pidgeot",
	"Rattata",
	"Raticate",
	"Spearow",
	"Fearow",
	"Ekans",
	"Arbok",
	"Pikachu",
	"Raichu",
	"Sandshrew",
	"Sandslash",
	"Nidoran",
	"Nidorina",
	"Nidoqueen",
	"Nidoran",
	"Nidorino",
	"Nidoking",
	"Clefairy",
	"Clefable",
	"Vulpix",
	"Ninetales",
	"Jigglypuff",
	"Wigglytuff",
	"Zubat",
	"Golbat",
	"Oddish",
	"Gloom",
	"Vileplume",
	"Paras",
	"Parasect",
	"Venonat",
	"Venomoth",
	"Diglett",
}

type AIPlayer struct {
	Name       string
	ID         string
	Score      int
	Conn       *websocket.Conn
	VSID       string
	VSName     string
	VSScore    int
	RoundID    string
	Status     string
	Grid       ttt.Grid
	StatusChan chan *ttt.PlayerStatus
	QuitChan   chan bool
}

type AIManager struct {
	AIPlayers *map[string]*AIPlayer
}

func (am *AIManager) NewAIPlayer(id string) *AIPlayer {
	p := &AIPlayer{
		StatusChan: make(chan *ttt.PlayerStatus, BufferedChanLen),
		QuitChan:   make(chan bool, BufferedChanLen),
	}
	go p.Play()
	(*am.AIPlayers)[id] = p
	glog.Infoln("total AI players", len((*am.AIPlayers)))
	return p
}

func (am *AIManager) UpdatePlayer(s *ttt.PlayerStatus) error {
	p := (*am.AIPlayers)[s.PlayerID]
	if p == nil {
		glog.Warningln("Can not find such player")
		return errors.New("Can not find such player")
	}
	p.StatusChan <- s
	return nil
}

// Dispatch player status to players
func (am *AIManager) Dispatch() {
	for {
		select {
		case s := <-playerStatuses:
			am.UpdatePlayer(&s)
		}
	}
}

func (ai *AIPlayer) Update(s *ttt.PlayerStatus) error {
	if s.RoundID != "" && ai.RoundID != s.RoundID &&
		!ttt.IsOverStatus(ai.Status) {
		msg := "Round IDs do not match for AI player"
		glog.Warningln(msg)
		return errors.New(msg)
	} else {
		ai.RoundID = s.RoundID
	}
	ai.ID = s.PlayerID
	ai.Score = s.PlayerScore
	ai.VSID = s.VSID
	ai.VSName = s.VSName
	ai.VSScore = s.VSScore
	ai.Status = s.Status
	if s.GridSnap != nil {
		ai.Grid = *s.GridSnap
	} else {
		var grid ttt.Grid
		ai.Grid = grid
	}
	if ttt.IsAIOverStatus(ai.Status) {
		ai.QuitChan <- true
	}
	return nil
}

func (ai *AIPlayer) Move() {
	if ai.Status != ttt.StatusYourTurn {
		return
	}
	pos := ai.GetBestPosition()
	m := ttt.PlayerAction{
		RoundID:    ai.RoundID,
		PlayerID:   ai.ID,
		PlayerName: ai.Name,
		Pos:        pos,
		Cmd:        ttt.CmdMove,
	}
	// sleep a bit to make it look like human
	time.Sleep(time.Duration(ttt.RandInt(1000)) * time.Millisecond)
	playerActions <- m
}

func (ai *AIPlayer) GetBestPosition() ttt.Position {
	g := ttt.Game{
		CurrentPlayer: ai.ID,
		NextPlayer:    ai.VSID,
		Grd:           ai.Grid,
	}
	r := g.GetBestMove(ai.ID)
	return r.Pos
}

func (ai *AIPlayer) Play() {
	for {
		select {
		case s := <-ai.StatusChan:
			ai.Update(s)
			ai.Move()
		case <-ai.QuitChan:
			delete((*am.AIPlayers), ai.ID)
			break
		}
	}
}

func InitAIManager() *AIManager {
	players := make(map[string]*AIPlayer)
	am := &AIManager{
		AIPlayers: &players,
	}
	go am.Dispatch()
	return am
}
