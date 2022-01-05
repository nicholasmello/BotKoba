package Vector3

import (
	math "github.com/chewxy/math32"
)

type Vector3 struct {
	X, Y, Z 	float32
}

func (v *Vector3) Dot(s *Vector3) float32 {
	return v.X * s.X +  v.Y * s.Y + v.Z * s.Z
}

func (v *Vector3) Cross(s *Vector3) *Vector3 {
	return New(v.Y*s.Z - v.Z*s.Y, v.Z*s.X - v.X*s.Z, v.X*s.Y - v.Y*s.X)
}

func (v *Vector3) Add(s *Vector3) *Vector3 {
	return New(v.X + s.X, v.Y + s.Y, v.Z + s.Z)
}

func (v *Vector3) AddScalar(num float32) *Vector3 {
	return New(v.X + num, v.Y + num, v.Z + num)
}

func (v *Vector3) Subtract(s *Vector3) *Vector3 {
	return New(v.X - s.X, v.Y - s.Y, v.Z - s.Z)
}

func (v *Vector3) SubtractScalar(num float32) *Vector3 {
	return New(v.X - num, v.Y - num, v.Z - num)
}

func (v *Vector3) Divide(s *Vector3) *Vector3 {
	return New(v.X / s.X, v.Y / s.Y, v.Z / s.Z)
}

func (v *Vector3) DivideScalar(num float32) *Vector3 {
	return New(v.X / num, v.Y / num, v.Z / num)
}

func (v *Vector3) MultiplyScalar(num float32) *Vector3 {
	return New(v.X * num, v.Y * num, v.Z * num)
}

func (v *Vector3) Magnitude() float32 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v *Vector3) Distance(s *Vector3) float32 {
	return (v.Subtract(s)).Magnitude()
}

func New(x float32, y float32, z float32) *Vector3 {
	return &Vector3{x,y,z}
}