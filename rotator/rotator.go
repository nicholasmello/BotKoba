package Rotator

import (
	math "github.com/chewxy/math32"
)

type Rotator struct {
	Pitch, Yaw, Roll 	float32
}

func (v *Rotator) Dot(s *Rotator) float32 {
	return v.Pitch * s.Pitch +  v.Yaw * s.Yaw + v.Roll * s.Roll
}

func (v *Rotator) Cross(s *Rotator) *Rotator {
	return New(v.Yaw*s.Roll - v.Roll*s.Yaw, v.Roll*s.Pitch - v.Pitch*s.Roll, v.Pitch*s.Yaw - v.Yaw*s.Pitch)
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