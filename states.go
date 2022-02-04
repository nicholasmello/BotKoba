package main

import (
	RLBot "github.com/Trey2k/RLBotGo"
	vector "github.com/xonmello/BotKoba/vector3"
	// rotator "github.com/xonmello/BotKoba/rotator"
	math "github.com/chewxy/math32"
)

type State int

const (
	ATBA State = iota
	Kickoff
	DefensiveCorner
	Air
	OffensiveCorner
	OnWall
)

type StateInfo struct {
	current  State
	start    int64
	allotted int64
}

// (Mostly) Always Toward Ball Agent
func stateATBA(koba *RLBot.PlayerInfo, opponent *RLBot.PlayerInfo, ball *RLBot.BallInfo, ballPrediction *RLBot.BallPrediction) *RLBot.ControllerState {
	PlayerInput := &RLBot.ControllerState{}

	kobaPos, kobaRot, _, opponentPos, _, _, ballPos, _ := initialSetup(koba, opponent, ball)

	// If in air, point wheels down
	if !koba.HasWheelContact {
		PlayerInput.Roll = -1.0 * kobaRot.Roll / math.Pi
	}

	// Go toward ball
	PlayerInput.Steer = steerToward(kobaPos, kobaRot, ballPos)

	// Drift if close to the ball and not very aligned
	if kobaPos.Distance(ballPos) < 500 && math.Abs(PlayerInput.Steer) > 0.3 {
		PlayerInput.Handbrake = true
	}

	// If on the wrong side of the ball, go back
	if kobaPos.Y > ballPos.Y+400 {
		PlayerInput.Steer = -1 * (math.Pi/2 + kobaRot.Yaw)
		PlayerInput.Boost = true
	} else if opponentPos.Distance(ballPos) > kobaPos.Distance(ballPos) {
		PlayerInput.Boost = true
	}

	// Flip if close to the ball
	if kobaPos.Distance(ballPos) < 400 && kobaPos.Y < ballPos.Y {
		predictionPosition := ballPrediction.Slices[5].Physics.Location
		predictionPositionVector := vector.New(predictionPosition.X, predictionPosition.Y, predictionPosition.Z)
		PlayerInput = flipToward(kobaPos, koba.Jumped, kobaRot, predictionPositionVector, PlayerInput)
	}

	// Go forward
	PlayerInput.Throttle = 1.0

	return PlayerInput
}

// Kickoff with 2 flips
func stateKickoff(koba *RLBot.PlayerInfo, opponent *RLBot.PlayerInfo, ball *RLBot.BallInfo) *RLBot.ControllerState {
	PlayerInput := &RLBot.ControllerState{}

	kobaPos, kobaRot, kobaVel, _, _, _, ballPos, _ := initialSetup(koba, opponent, ball)

	// Go toward ball
	PlayerInput.Steer = steerToward(kobaPos, kobaRot, ballPos)

	// Starting flip to gain speed
	if kobaVel.Magnitude() > 1300 {
		PlayerInput = flipToward(kobaPos, koba.Jumped, kobaRot, ballPos, PlayerInput)
	}

	// Flip when close to the ball
	if kobaPos.Distance(ballPos) < 500 && kobaPos.Y < ballPos.Y {
		PlayerInput = flipToward(kobaPos, koba.Jumped, kobaRot, ballPos, PlayerInput)
	}

	// Go forward and boost
	if koba.HasWheelContact {
		PlayerInput.Throttle = 1.0
		PlayerInput.Boost = true
	}

	return PlayerInput
}

// Getting on the correct side of the ball when in the corner
func stateDefensiveCorner(koba *RLBot.PlayerInfo, opponent *RLBot.PlayerInfo, ball *RLBot.BallInfo) *RLBot.ControllerState {
	PlayerInput := &RLBot.ControllerState{}

	kobaPos, kobaRot, _, _, _, _, ballPos, _ := initialSetup(koba, opponent, ball)

	// If in air, point wheels down
	if !koba.HasWheelContact {
		PlayerInput.Roll = -1.0 * kobaRot.Roll / math.Pi
	}

	// Check if on the correct side of the ball
	if math.Abs(kobaPos.X) < math.Abs(ballPos.X) && kobaPos.Y < ballPos.Y+500 {
		PlayerInput.Steer = steerToward(kobaPos, kobaRot, ballPos)
		// Flip when close to the ball
		if kobaPos.Distance(ballPos) < 500 && kobaPos.Y < ballPos.Y {
			PlayerInput = flipToward(kobaPos, koba.Jumped, kobaRot, ballPos, PlayerInput)
		}
	} else {
		PlayerInput.Steer = steerToward(kobaPos, kobaRot, vector.New(0, -4500, 0))
	}

	// Drift if not very aligned
	if math.Abs(PlayerInput.Steer) > 0.5 {
		PlayerInput.Handbrake = true
	}

	// Go forward
	PlayerInput.Throttle = 1.0

	return PlayerInput
}

// General state for when the ball is in the air
func stateAir(koba *RLBot.PlayerInfo, opponent *RLBot.PlayerInfo, ball *RLBot.BallInfo, ballPrediction *RLBot.BallPrediction) *RLBot.ControllerState {
	PlayerInput := &RLBot.ControllerState{}

	kobaPos, kobaRot, _, _, _, _, ballPos, _ := initialSetup(koba, opponent, ball)

	// Target is in front of where the ball will land
	target := lowestPrediction(ballPrediction)
	target.Y = ballPos.Y - 500

	// Go toward ball
	PlayerInput.Steer = steerToward(kobaPos, kobaRot, target)

	// Drift if not very aligned
	if math.Abs(PlayerInput.Steer) > 0.5 {
		PlayerInput.Handbrake = true
	}

	// Go forward
	PlayerInput.Throttle = 1.0

	return PlayerInput
}

// Getting in position to hit it if it goes center
// TODO Needs work
func stateOffensiveCorner(koba *RLBot.PlayerInfo, opponent *RLBot.PlayerInfo, ball *RLBot.BallInfo) *RLBot.ControllerState {
	PlayerInput := &RLBot.ControllerState{}

	kobaPos, kobaRot, _, _, _, _, ballPos, _ := initialSetup(koba, opponent, ball)

	// Go to the side of the field the ball is in
	// Wait for the ball to go in front of goal or out of corner
	target := vector.New(0, 200, 0)
	if ballPos.X > 0 {
		target.X = 2000
	} else {
		target.X = -2000
	}
	PlayerInput.Steer = steerToward(kobaPos, kobaRot, target)

	// Slow down if close to target location
	PlayerInput.Throttle = 1
	if kobaPos.Distance(target) < 1250 {
		PlayerInput.Throttle = 0.2
	}

	// Drift if not very aligned
	if math.Abs(PlayerInput.Steer) > 0.5 {
		PlayerInput.Handbrake = true
	}

	return PlayerInput
}

// Jump off wall when on it
func stateOnWall(koba *RLBot.PlayerInfo, opponent *RLBot.PlayerInfo, ball *RLBot.BallInfo) *RLBot.ControllerState {
	PlayerInput := &RLBot.ControllerState{}

	kobaPos, kobaRot, _, _, _, _, ballPos, _ := initialSetup(koba, opponent, ball)

	// Jump off wall
	PlayerInput.Jump = true

	// If in air, point wheels down and turn toward  ball
	if !koba.HasWheelContact {
		PlayerInput.Roll = -1.0 * kobaRot.Roll / math.Pi
		PlayerInput.Yaw = steerToward(kobaPos, kobaRot, ballPos)
	}

	return PlayerInput
}
