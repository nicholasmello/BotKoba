package main

import (
	// vector "github.com/xonmello/BotKoba/vector3"
	// rotator "github.com/xonmello/BotKoba/rotator"
	math "github.com/chewxy/math32"
	RLBot "github.com/Trey2k/RLBotGo"
)

type State int

const (
	ATBA State = iota
	Kickoff
)

type StateInfo struct {
	current		State
	start		int64
	alloted		int64
}

// (Mostly) Always Toward Ball Agent
func state_ATBA(koba *RLBot.PlayerInfo, opponent *RLBot.PlayerInfo, ball *RLBot.BallInfo) *RLBot.ControllerState {
	PlayerInput := &RLBot.ControllerState{}

	koba_pos, koba_rot, _, _, ball_pos := initialSetup(koba, opponent, ball)

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

// TO:DO Kickoff
func state_Kickoff(koba *RLBot.PlayerInfo, opponent *RLBot.PlayerInfo, ball *RLBot.BallInfo) *RLBot.ControllerState {
	return &RLBot.ControllerState{}
}