package main

import (
	"time"

	RLBot "github.com/Trey2k/RLBotGo"
	math "github.com/chewxy/math32"
	rotator "github.com/xonmello/BotKoba/rotator"
	vector "github.com/xonmello/BotKoba/vector3"
)

var lastJump int64

func initialSetup(koba *RLBot.PlayerInfo, opponent *RLBot.PlayerInfo, ball *RLBot.BallInfo) (*vector.Vector3, *rotator.Rotator, *vector.Vector3, *vector.Vector3, *rotator.Rotator, *vector.Vector3, *vector.Vector3, *vector.Vector3) {
	// Get self information into a useful format
	kobaPos := vector.New(koba.Physics.Location.X, koba.Physics.Location.Y, koba.Physics.Location.Z)
	kobaRot := rotator.New(koba.Physics.Rotation.Pitch, koba.Physics.Rotation.Yaw, koba.Physics.Rotation.Roll)
	kobaVel := vector.New(koba.Physics.Velocity.X, koba.Physics.Velocity.Y, koba.Physics.Velocity.Z)

	// Get opponent information into a useful format
	opponentPos := vector.New(opponent.Physics.Location.X, opponent.Physics.Location.Y, opponent.Physics.Location.Z)
	opponentRot := rotator.New(opponent.Physics.Rotation.Pitch, opponent.Physics.Rotation.Yaw, opponent.Physics.Rotation.Roll)
	opponentVel := vector.New(opponent.Physics.Velocity.X, opponent.Physics.Velocity.Y, opponent.Physics.Velocity.Z)

	// Get ball information into a useful format
	ballPos := vector.New(ball.Physics.Location.X, ball.Physics.Location.Y, ball.Physics.Location.Z)
	ballVel := vector.New(ball.Physics.Velocity.X, ball.Physics.Velocity.Y, ball.Physics.Velocity.Z)

	// Flip coordinates when on orange team
	if koba.Team == 1 {
		kobaPos = kobaPos.MultiplyScalar(-1)
		kobaRot = kobaRot.RotateYaw(math.Pi)
		kobaVel = kobaVel.MultiplyScalar(-1)
		ballPos = ballPos.MultiplyScalar(-1)
		ballVel = ballVel.MultiplyScalar(-1)
		opponentPos = opponentPos.MultiplyScalar(-1)
		opponentRot = opponentRot.RotateYaw(math.Pi)
		opponentVel = opponentVel.MultiplyScalar(-1)
	}

	return kobaPos, kobaRot, kobaVel, opponentPos, opponentRot, opponentVel, ballPos, ballVel
}

func steerToward(selfPos *vector.Vector3, selfRot *rotator.Rotator, target *vector.Vector3) float32 {
	// Center the car in the coordinate system
	local := target.Subtract(selfPos)
	toTargetAngle := math.Atan2(local.Y, local.X)

	// Steer toward the ball depending on our Yaw (direction we are facing)
	steer := toTargetAngle - selfRot.Yaw
	if steer < -math.Pi {
		steer += math.Pi * 2.0
	} else if steer >= math.Pi {
		steer -= math.Pi * 2.0
	}

	// If angle is greater than 1 radian, limit to full turn
	steer = bound(steer, -1, 1)
	return steer
}

func flipToward(selfPos *vector.Vector3, jumped bool, selfRot *rotator.Rotator, target *vector.Vector3, PlayerInput *RLBot.ControllerState) *RLBot.ControllerState {
	local := target.Subtract(selfPos)
	localAngle := rotator.New(0, math.Atan2(local.Y, local.X), 0).RotateYaw(-selfRot.Yaw).Yaw

	if !jumped {
		PlayerInput.Jump = true
		lastJump = currentTime()
	} else if jumped && currentTime() < lastJump+70 {
		PlayerInput.Jump = true
	}

	if jumped && time.Now().UnixMilli() > lastJump+110 {
		PlayerInput.Jump = true
		if math.Abs(localAngle) <= 0.3 {
			PlayerInput.Pitch = -1
			PlayerInput.Yaw = 0
		} else if localAngle <= (math.Pi/2) && 1.14 <= localAngle {
			PlayerInput.Pitch = 0
			PlayerInput.Yaw = 1
		} else if localAngle <= -1.14 && -(math.Pi/2) <= localAngle {
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

func bound(x float32, low float32, high float32) float32 {
	if x < low {
		return low
	} else if x > high {
		return high
	}
	return x
}

func currentTime() int64 {
	return time.Now().UnixMilli()
}

func lowestPrediction(ballPrediction *RLBot.BallPrediction) *vector.Vector3 {
	targetZ := ballPrediction.Slices[0].Physics.Location.Z
	targetIndex := 0
	for i := 0; i < len(ballPrediction.Slices); i++ {
		if targetZ > ballPrediction.Slices[i].Physics.Location.Z {
			targetZ = ballPrediction.Slices[i].Physics.Location.Z
			targetIndex = i
		}
	}
	physicsSlice := ballPrediction.Slices[targetIndex].Physics.Location
	return vector.New(physicsSlice.X, physicsSlice.Y, physicsSlice.Z)
}
