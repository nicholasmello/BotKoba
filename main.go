package main

import (
	"fmt"

	math "github.com/chewxy/math32"
	RLBot "github.com/Trey2k/RLBotGo"
)

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

	// Determine State if needed
	if currentTime() > state.start + state.alloted {
		if !gameState.GameTick.GameInfo.IsRoundActive {
			statePrint("State Changed: Kickoff", Kickoff)
			state = StateInfo{current: Kickoff, start: currentTime(), alloted: 2000}
		} else if (math.Abs(koba.Physics.Rotation.Pitch) > math.Pi / 2 - 0.1 || math.Abs(koba.Physics.Rotation.Roll) > math.Pi / 2 - 0.1) && koba.HasWheelContact {
			statePrint("State Changed: OnWall", OnWall)
			state = StateInfo{current: OnWall, start: currentTime(), alloted: 200}
		} else if math.Abs(ball.Physics.Location.X) > 1300 && ball.Physics.Location.Y < -2250 {
			statePrint("State Changed: Defensive Corner", DefensiveCorner)
			state = StateInfo{current: DefensiveCorner, start: currentTime(), alloted: 0}
		} else {
			statePrint("State Changed: ATBA", ATBA)
			state = StateInfo{current: ATBA, start: currentTime(), alloted: 0}
		}
	}

	// Get Player Input given state
	var PlayerInput *RLBot.ControllerState
	switch state.current {
	case ATBA:
		PlayerInput = state_ATBA(&koba, &opponent, &ball)
	case Kickoff:
		// If round is still inactive, reset start time
		if !gameState.GameTick.GameInfo.IsRoundActive {
			state.start = currentTime()
		}
		PlayerInput = state_Kickoff(&koba, &opponent, &ball)
	case DefensiveCorner:
		PlayerInput = state_DefensiveCorner(&koba, &opponent, &ball)
	case OnWall:
		PlayerInput = state_OnWall(&koba, &opponent, &ball)
	}

	return PlayerInput
}

func statePrint(str string, sta State) {
	if sta != state.current {
		fmt.Println(str)
	}
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
