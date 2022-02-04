package main

import (
	"fmt"

	RLBot "github.com/Trey2k/RLBotGo"
	math "github.com/chewxy/math32"
)

var state = StateInfo{current: ATBA, start: currentTime(), allotted: 0}

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
	if currentTime() > state.start+state.allotted {
		if !gameState.GameTick.GameInfo.IsRoundActive {
			statePrint("State Changed: Kickoff", Kickoff)
			state = StateInfo{current: Kickoff, start: currentTime(), allotted: 2000}
		} else if (math.Abs(koba.Physics.Rotation.Pitch) > math.Pi/2-0.1 || math.Abs(koba.Physics.Rotation.Roll) > math.Pi/2-0.1) && koba.HasWheelContact {
			statePrint("State Changed: OnWall", OnWall)
			state = StateInfo{current: OnWall, start: currentTime(), allotted: 200}
		} else if math.Abs(ball.Physics.Location.X) > 1300 && ball.Physics.Location.Y < -2250 {
			statePrint("State Changed: Defensive Corner", DefensiveCorner)
			state = StateInfo{current: DefensiveCorner, start: currentTime(), allotted: 0}
		} else if ball.Physics.Location.Z > 800 {
			statePrint("State Changed: Air", Air)
			state = StateInfo{current: Air, start: currentTime(), allotted: 0}
		} else if math.Abs(ball.Physics.Location.X) > 2000 && ball.Physics.Location.Y > 2250 {
			statePrint("State Changed: Offensive Corner", OffensiveCorner)
			state = StateInfo{current: OffensiveCorner, start: currentTime(), allotted: 0}
		} else {
			statePrint("State Changed: ATBA", ATBA)
			state = StateInfo{current: ATBA, start: currentTime(), allotted: 0}
		}
	}

	// Get Player Input given state
	var PlayerInput *RLBot.ControllerState
	switch state.current {
	case ATBA:
		PlayerInput = stateATBA(&koba, &opponent, &ball, gameState.BallPrediction)
	case Kickoff:
		// If round is still inactive, reset start time
		if !gameState.GameTick.GameInfo.IsRoundActive {
			state.start = currentTime()
		}
		PlayerInput = stateKickoff(&koba, &opponent, &ball)
	case DefensiveCorner:
		PlayerInput = stateDefensiveCorner(&koba, &opponent, &ball)
	case Air:
		PlayerInput = stateAir(&koba, &opponent, &ball, gameState.BallPrediction)
	case OffensiveCorner:
		PlayerInput = stateOffensiveCorner(&koba, &opponent, &ball)
	case OnWall:
		PlayerInput = stateOnWall(&koba, &opponent, &ball)
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
