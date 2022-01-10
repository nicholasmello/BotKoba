package main

import (
	vector "github.com/xonmello/BotKoba/vector3"
	// rotator "github.com/xonmello/BotKoba/rotator"
	math "github.com/chewxy/math32"
	RLBot "github.com/Trey2k/RLBotGo"
)

type State int

const (
	ATBA State = iota
	Kickoff
	DefensiveCorner
	OnWall
)

type StateInfo struct {
	current		State
	start		int64
	alloted		int64
}

// (Mostly) Always Toward Ball Agent
func state_ATBA(koba *RLBot.PlayerInfo, opponent *RLBot.PlayerInfo, ball *RLBot.BallInfo) *RLBot.ControllerState {
	PlayerInput := &RLBot.ControllerState{}

	koba_pos, koba_rot, _, _, _, _, ball_pos, _ := initialSetup(koba, opponent, ball)

	// If in air, point wheels down
	if !koba.HasWheelContact {
		PlayerInput.Roll = -1.0*koba_rot.Roll/math.Pi
	}

	// Go toward ball
	PlayerInput.Steer = steerToward(koba_pos, koba_rot, ball_pos)

	// Drift if close to the ball and not very aligned
	if koba_pos.Distance(ball_pos) < 500 && math.Abs(PlayerInput.Steer) > 0.3 {
		PlayerInput.Handbrake = true
	} 

	// If on the wrong side of the ball, go back
	if koba_pos.Y > ball_pos.Y + 400 {
		PlayerInput.Steer = -1 * (math.Pi / 2 + koba_rot.Yaw)
		PlayerInput.Boost = true
	}

	// Flip if close to the ball 
	if koba_pos.Distance(ball_pos) < 400 && koba_pos.Y < ball_pos.Y {
		PlayerInput = flipToward(koba_pos, koba.Jumped, koba_rot, ball_pos, PlayerInput)
	}

	// Go forward
	PlayerInput.Throttle = 1.0

	return PlayerInput
}

// Kickoff with 2 flips
func state_Kickoff(koba *RLBot.PlayerInfo, opponent *RLBot.PlayerInfo, ball *RLBot.BallInfo) *RLBot.ControllerState {
	PlayerInput := &RLBot.ControllerState{}

	koba_pos, koba_rot, koba_vel, _, _, _, ball_pos, _ := initialSetup(koba, opponent, ball)

	// Go toward ball
	PlayerInput.Steer = steerToward(koba_pos, koba_rot, ball_pos)

	// Starting flip to gain speed
	if koba_vel.Magnitude() > 1300 {
		PlayerInput = flipToward(koba_pos, koba.Jumped, koba_rot, ball_pos, PlayerInput)
	}

	// Flip when close to the ball 
	if koba_pos.Distance(ball_pos) < 500 && koba_pos.Y < ball_pos.Y {
		PlayerInput = flipToward(koba_pos, koba.Jumped, koba_rot, ball_pos, PlayerInput)
	}

	// Go forward and boost
	if koba.HasWheelContact {
		PlayerInput.Throttle = 1.0
		PlayerInput.Boost = true	
	}

	return PlayerInput
}

// Getting on the correct side of the ball when in the corner
func state_DefensiveCorner(koba *RLBot.PlayerInfo, opponent *RLBot.PlayerInfo, ball *RLBot.BallInfo) *RLBot.ControllerState {
	PlayerInput := &RLBot.ControllerState{}

	koba_pos, koba_rot, _, _, _, _, ball_pos, _ := initialSetup(koba, opponent, ball)

	// If in air, point wheels down
	if !koba.HasWheelContact {
		PlayerInput.Roll = -1.0*koba_rot.Roll/math.Pi
	}

	// Check if on the correct side of the ball
	if math.Abs(koba_pos.X) < math.Abs(ball_pos.X) && koba_pos.Y < ball_pos.Y {
		PlayerInput.Steer = steerToward(koba_pos, koba_rot, ball_pos)
		// Flip when close to the ball 
		if koba_pos.Distance(ball_pos) < 500 && koba_pos.Y < ball_pos.Y {
			PlayerInput = flipToward(koba_pos, koba.Jumped, koba_rot, ball_pos, PlayerInput)
		}
	} else {
		PlayerInput.Steer = steerToward(koba_pos, koba_rot, vector.New(0, -4500, 0))
	}

	// Drift if not very aligned
	if math.Abs(PlayerInput.Steer) > 0.5 {
		PlayerInput.Handbrake = true
	} 

	// Go forward
	PlayerInput.Throttle = 1.0

	return PlayerInput
}

// Jump off wall when on it
func state_OnWall(koba *RLBot.PlayerInfo, opponent *RLBot.PlayerInfo, ball *RLBot.BallInfo) *RLBot.ControllerState {
	PlayerInput := &RLBot.ControllerState{}

	koba_pos, koba_rot, _, _, _, _, ball_pos, _ := initialSetup(koba, opponent, ball)

	// Jump off wall
	PlayerInput.Jump = true

	// If in air, point wheels down and turn toward  ball
	if !koba.HasWheelContact {
		PlayerInput.Roll = -1.0*koba_rot.Roll/math.Pi
		PlayerInput.Yaw = steerToward(koba_pos, koba_rot, ball_pos)
	}

	return PlayerInput
}