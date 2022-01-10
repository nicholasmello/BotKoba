package main

import (
	"fmt"
	"time"

	RLBot "github.com/Trey2k/RLBotGo"
)

var lastjump int64
var lastkickoff int64
var state StateInfo = StateInfo{current: ATBA, start: currentTime(), alloted: 0}

func getInput(gameState *RLBot.GameState, rlBot *RLBot.RLBot) *RLBot.ControllerState {
	// Get Players and Ball data
	koba := gameState.GameTick.Players[rlBot.PlayerIndex]
	var opponent RLBot.PlayerInfo
	if rlBot.PlayerIndex == 0 {
		opponent = gameState.GameTick.Players[1]
	} else {
		opponent = gameState.GameTick.Players[0]
	}
	ball := gameState.GameTick.Ball

	// TO:DO Determine State if needed

	// Get Player Input given state
	var PlayerInput *RLBot.ControllerState
	switch state.current {
	case ATBA:
		PlayerInput = state_ATBA(&koba, &opponent, &ball)
	case Kickoff:
		PlayerInput = state_Kickoff(&koba, &opponent, &ball)
	}

	// Boost on kickoff
	if !gameState.GameTick.GameInfo.IsRoundActive {
		lastkickoff = time.Now().UnixMilli()
	} else if time.Now().UnixMilli() < lastkickoff + 2000 {
		PlayerInput.Boost = true
	}

	return PlayerInput
}

func main() {

	// connect to RLBot
	rlBot, err := RLBot.Connect(23234)
	if err != nil {
		panic(err)
	}

	// Send ready message
	err = rlBot.SendReadyMessage(true, true, true)
	if err != nil {
		panic(err)
	}

	// Set our tick handler
	err = rlBot.SetGetInput(getInput)
	fmt.Println(err.Error())

}
