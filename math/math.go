package math

import (
	"fmt"
	"math"
)

type Vector2 struct {
	X, Y float64
}

func (v Vector2) ApproxEqual(ov Vector2) bool {
	const epsilon = 1e-16
	return math.Abs(v.X-ov.X) < epsilon && math.Abs(v.Y-ov.Y) < epsilon
}

func (v Vector2) String() string {
	return fmt.Sprintf("(%0.24f, %0.24f)", v.X, v.Y)
}

func (v Vector2) Norm() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vector2) Norm2() float64 {
	return v.Dot(v)
}

func (v Vector2) Normalize() Vector2 {
	n2 := v.Norm2()
	if n2 == 0 {
		return Vector2{0, 0}
	}

	return v.Mul(1 / math.Sqrt(n2))
}

func (v Vector2) IsUnit() bool {
	const epsilon = 5e-14
	return math.Abs(v.Norm2()-1) <= epsilon
}

func (v Vector2) Abs() Vector2 {
	return Vector2{math.Abs(v.X), math.Abs(v.Y)}
}

func (v Vector2) Add(ov Vector2) Vector2 {
	return Vector2{v.X + ov.X, v.Y + ov.Y}
}

func (v Vector2) Sub(ov Vector2) Vector2 {
	return Vector2{v.X - ov.X, v.Y - ov.Y}
}

func (v Vector2) Mul(m float64) Vector2 {
	return Vector2{m * v.X, m * v.Y}
}

func (v Vector2) Dot(ov Vector2) float64 {
	return v.X*ov.X + v.Y*ov.Y
}

func ClampFloat64(x, min, max float64) float64 {
	if x < min {
		return min
	} else if x > max {
		return max
	}

	return x
}

func ClampVec2(val, min, max Vector2) Vector2 {
	return Vector2{
		X: ClampFloat64(val.X, min.X, max.X),
		Y: ClampFloat64(val.Y, min.Y, max.Y),
	}
}
