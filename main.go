package main

import (
	"fmt"

	vector "github.com/xonmello/BotKoba/vector3"
	rotator "github.com/xonmello/BotKoba/rotator"
	math "github.com/chewxy/math32"
	RLBot "github.com/Trey2k/RLBotGo"
)

func getInput(gameState *RLBot.GameState, rlBot *RLBot.RLBot) *RLBot.ControllerState {
	PlayerInput := &RLBot.ControllerState{}

	// Get self information into a useful format
	koba := gameState.GameTick.Players[rlBot.PlayerIndex]
	koba_pos := vector.New(koba.Physics.Location.X, koba.Physics.Location.Y, koba.Physics.Location.Z)
	koba_rot := rotator.New(koba.Physics.Rotation.Pitch, koba.Physics.Rotation.Yaw, koba.Physics.Rotation.Roll)
	
	// Get opponent information into a useful format
	// Opponent information unused right now
	// var opponent RLBot.PlayerInfo
	// if rlBot.PlayerIndex == 0 {
	// 	opponent := gameState.GameTick.Players[1]
	// } else {
	// 	opponent := gameState.GameTick.Players[0]
	// }
	// opponent_pos := vector.New(opponent.Physics.Location.X, opponent.Physics.Location.Y, opponent.Physics.Location.Z)
	// opponent_rot := rotator.New(opponent.Physics.Rotation.Pitch, opponent.Physics.Rotation.Yaw, opponent.Physics.Rotation.Roll)
	
	// Get ball information into a useful format
	ball := gameState.GameTick.Ball
	ballLocation := vector.New(ball.Physics.Location.X, ball.Physics.Location.Y, ball.Physics.Location.Z)

	// Put wheels down if in the air
	if !koba.HasWheelContact {
		PlayerInput.Roll = -1.0*koba_rot.Roll/math.Pi
	}

	// Center the car in the coordinate system
	local := ballLocation.Subtract(koba_pos)
	toBallAngle := math.Atan2(local.Y,local.X)
	
	// Steer toward the ball depending on our Yaw (direction we are facing)
	steer := toBallAngle - koba_rot.Yaw
	if steer < -math.Pi {
		steer += math.Pi * 2.0;
	} else if steer >= math.Pi {
		steer -= math.Pi * 2.0;
	}

	// If angle is greater than 1 radian, limit to full turn
	if (steer > 1) {
		steer = 1
	} else if (steer < -1) {
		steer = -1
	}

	// Put final calculation into player input
	PlayerInput.Steer = steer

	// Drift if close to the ball and not very aligned
	if koba_pos.Distance(ballLocation) < 300 && math.Abs(steer) > 0.3 {
		PlayerInput.Handbrake = true
	} 

	// Go forward
	PlayerInput.Throttle = 1.0

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
