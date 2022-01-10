package main

import (
	"time"

	vector "github.com/xonmello/BotKoba/vector3"
	rotator "github.com/xonmello/BotKoba/rotator"
	math "github.com/chewxy/math32"
	RLBot "github.com/Trey2k/RLBotGo"
)

var lastjump int64

func initialSetup(koba *RLBot.PlayerInfo, opponent *RLBot.PlayerInfo, ball *RLBot.BallInfo) (*vector.Vector3, *rotator.Rotator, *vector.Vector3, *vector.Vector3, *rotator.Rotator, *vector.Vector3, *vector.Vector3, *vector.Vector3) {
	// Get self information into a useful format
	koba_pos := vector.New(koba.Physics.Location.X, koba.Physics.Location.Y, koba.Physics.Location.Z)
	koba_rot := rotator.New(koba.Physics.Rotation.Pitch, koba.Physics.Rotation.Yaw, koba.Physics.Rotation.Roll)
	koba_vel := vector.New(koba.Physics.Velocity.X, koba.Physics.Velocity.Y, koba.Physics.Velocity.Z)
	
	// Get opponent information into a useful format
	opponent_pos := vector.New(opponent.Physics.Location.X, opponent.Physics.Location.Y, opponent.Physics.Location.Z)
	opponent_rot := rotator.New(opponent.Physics.Rotation.Pitch, opponent.Physics.Rotation.Yaw, opponent.Physics.Rotation.Roll)
	opponent_vel := vector.New(opponent.Physics.Velocity.X, opponent.Physics.Velocity.Y, opponent.Physics.Velocity.Z)
	
	// Get ball information into a useful format
	ball_pos := vector.New(ball.Physics.Location.X, ball.Physics.Location.Y, ball.Physics.Location.Z)
	ball_vel := vector.New(ball.Physics.Velocity.X, ball.Physics.Velocity.Y, ball.Physics.Velocity.Z)

	// Flips coordinates when on orange team
	if koba.Team == 1 {
		koba_pos = koba_pos.MultiplyScalar(-1)
		koba_rot = koba_rot.RotateYaw(math.Pi)
		koba_vel = koba_vel.MultiplyScalar(-1)
		ball_pos = ball_pos.MultiplyScalar(-1)
		ball_vel = ball_vel.MultiplyScalar(-1)
		opponent_pos = opponent_pos.MultiplyScalar(-1)
		opponent_rot = opponent_rot.RotateYaw(math.Pi)
		opponent_vel = opponent_vel.MultiplyScalar(-1)
	}

	return koba_pos, koba_rot, koba_vel, opponent_pos, opponent_rot, opponent_vel, ball_pos, ball_vel
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
	steer = cap(steer, -1, 1)
	return steer
}

func flipToward(self_pos *vector.Vector3, jumped bool, self_rot *rotator.Rotator, target *vector.Vector3, PlayerInput *RLBot.ControllerState) *RLBot.ControllerState {
	local := target.Subtract(self_pos)
	localAngle := rotator.New(0,math.Atan2(local.Y,local.X),0).RotateYaw(-self_rot.Yaw).Yaw

	if !jumped {
		PlayerInput.Jump = true
		lastjump = currentTime()
	} else if jumped && currentTime() < lastjump + 70 {
		PlayerInput.Jump = true
	}

	if jumped && time.Now().UnixMilli() > lastjump + 110 {
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

func cap(x float32, low float32, high float32) (float32) {
	if x < low {
		return low
	} else if x > high {
		return high
	}
	return x
}

func currentTime() (int64) {
	return time.Now().UnixMilli()
}