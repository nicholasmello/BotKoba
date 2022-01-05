package main

import (
	"fmt"
	"time"

	vector "github.com/xonmello/BotKoba/vector3"
	rotator "github.com/xonmello/BotKoba/rotator"
	math "github.com/chewxy/math32"
	RLBot "github.com/Trey2k/RLBotGo"
)

var lastjump int64
var lastkickoff int64

func getInput(gameState *RLBot.GameState, rlBot *RLBot.RLBot) *RLBot.ControllerState {
	PlayerInput := &RLBot.ControllerState{}

	// Get self information into a useful format
	koba := gameState.GameTick.Players[rlBot.PlayerIndex]
	koba_pos := vector.New(koba.Physics.Location.X, koba.Physics.Location.Y, koba.Physics.Location.Z)
	koba_rot := rotator.New(koba.Physics.Rotation.Pitch, koba.Physics.Rotation.Yaw, koba.Physics.Rotation.Roll)
	
	// Get opponent information into a useful format
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
	ball_pos := vector.New(ball.Physics.Location.X, ball.Physics.Location.Y, ball.Physics.Location.Z)

	// Flips coordinates when on orange team
	if koba.Team == 1 {
		koba_pos = koba_pos.MultiplyScalar(-1)
		koba_rot = koba_rot.RotateYaw(math.Pi)
		ball_pos = ball_pos.MultiplyScalar(-1)
		// opponent_pos = opponent_pos.MultiplyScalar(-1)
		// opponent_rot = opponent_rot.RotateYaw(math.Pi)
	}

	// Put wheels down if in the air
	if !koba.HasWheelContact {
		PlayerInput.Roll = -1.0*koba_rot.Roll/math.Pi
	}

	steer := steerToward(koba_pos, koba_rot, ball_pos)

	// Drift if close to the ball and not very aligned
	if koba_pos.Distance(ball_pos) < 500 && math.Abs(steer) > 0.3 {
		PlayerInput.Handbrake = true
	} 

	if koba_pos.Y > ball_pos.Y + 400 {
		steer = -1 * (math.Pi / 2 + koba_rot.Yaw)
		PlayerInput.Boost = true
	}

	if koba_pos.Distance(ball_pos) < 400 && koba_pos.Y < ball_pos.Y {
		PlayerInput = flipToward(koba_pos, koba.Jumped, koba_rot, ball_pos, PlayerInput)
	}

	// Put final calculation into player input
	PlayerInput.Steer = steer

	// Go forward
	PlayerInput.Throttle = 1.0

	// Boost on kickoff
	if !gameState.GameTick.GameInfo.IsRoundActive {
		lastkickoff = time.Now().UnixMilli()
	} else if time.Now().UnixMilli() < lastkickoff + 2000 {
		PlayerInput.Boost = true
	}

	fmt.Println(time.Now().UnixMilli() - lastkickoff)

	return PlayerInput
}

func steerToward(self_pos *vector.Vector3, self_rot *rotator.Rotator, target *vector.Vector3) float32 {
	// Center the car in the coordinate system
	local := target.Subtract(self_pos)
	toTargetAngle := math.Atan2(local.Y,local.X)
	
	// Steer toward the ball depending on our Yaw (direction we are facing)
	steer := toTargetAngle - self_rot.Yaw
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
	return steer
}

func flipToward(self_pos *vector.Vector3, jumped bool, self_rot *rotator.Rotator, target *vector.Vector3, PlayerInput *RLBot.ControllerState) *RLBot.ControllerState {
	local := target.Subtract(self_pos)
	localAngle := rotator.New(0,math.Atan2(local.Y,local.X),0).RotateYaw(-self_rot.Yaw).Yaw

	if !jumped {
		PlayerInput.Jump = true
		lastjump = time.Now().UnixMilli()
	}

	if jumped && time.Now().UnixMilli() > lastjump + 70 {
		PlayerInput.Jump = true
		if math.Abs(localAngle) <= 0.3 {
			PlayerInput.Pitch = -1
			PlayerInput.Yaw = 0
		} else if localAngle <= (math.Pi / 2) && 1.14 <= localAngle {
			PlayerInput.Pitch = 0
			PlayerInput.Yaw = 1
		} else if localAngle <= -1.14 && -(math.Pi / 2) <= localAngle {
			PlayerInput.Pitch = 0
			PlayerInput.Yaw = -1
		} else if localAngle <= 1.14 && 0.3 <= localAngle {
			PlayerInput.Pitch = -1
			PlayerInput.Yaw = 1
		} else if localAngle <= -0.3 && -1.14 <= localAngle {
			PlayerInput.Pitch = -1
			PlayerInput.Yaw = -1
		}
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
