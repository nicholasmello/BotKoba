package Rotator

import (
	math "github.com/chewxy/math32"
)

type Rotator struct {
	Pitch, Yaw, Roll 	float32
}

func (v *Rotator) RotatePitch(num float32) *Rotator {
	newPitch := v.Pitch + num
	if newPitch > math.Pi {
		newPitch -= 2 * math.Pi
	} else if newPitch < -math.Pi {
		newPitch += 2 * math.Pi
	}
	return New(newPitch, v.Yaw, v.Roll)
}

func (v *Rotator) RotateYaw(num float32) *Rotator {
	newYaw := v.Yaw + num
	if newYaw > math.Pi {
		newYaw -= 2 * math.Pi
	} else if newYaw < -math.Pi {
		newYaw += 2 * math.Pi
	}
	return New(v.Pitch, newYaw, v.Roll)
}

func (v *Rotator) RotateRoll(num float32) *Rotator {
	newRoll := v.Roll + num
	if newRoll > math.Pi {
		newRoll -= 2 * math.Pi
	} else if newRoll < -math.Pi {
		newRoll += 2 * math.Pi
	}
	return New(v.Pitch, v.Yaw, newRoll)
}

func (v *Rotator) Add(s *Rotator) *Rotator {
	return New(v.Pitch + s.Pitch, v.Yaw + s.Yaw, v.Roll + s.Roll)
}

func (v *Rotator) AddScalar(num float32) *Rotator {
	return New(v.Pitch + num, v.Yaw + num, v.Roll + num)
}

func (v *Rotator) Subtract(s *Rotator) *Rotator {
	return New(v.Pitch - s.Pitch, v.Yaw - s.Yaw, v.Roll - s.Roll)
}

func (v *Rotator) SubtractScalar(num float32) *Rotator {
	return New(v.Pitch - num, v.Yaw - num, v.Roll - num)
}

func (v *Rotator) Divide(s *Rotator) *Rotator {
	return New(v.Pitch / s.Pitch, v.Yaw / s.Yaw, v.Roll / s.Roll)
}

func (v *Rotator) DivideScalar(num float32) *Rotator {
	return New(v.Pitch / num, v.Yaw / num, v.Roll / num)
}

func (v *Rotator) Magnitude() float32 {
	return math.Sqrt(v.Pitch*v.Pitch + v.Yaw*v.Yaw + v.Roll*v.Roll)
}

func (v *Rotator) Distance(s *Rotator) float32 {
	return (v.Subtract(s)).Magnitude()
}

func New(pitch float32, yaw float32, roll float32) *Rotator {
	return &Rotator{pitch,yaw,roll}
}