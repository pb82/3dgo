package main

import "math"

type vec3d struct {
	x, y, z, w float64
}

func (v *vec3d) Add(v2 *vec3d) vec3d {
	return vec3d{
		x: v.x + v2.x,
		y: v.y + v2.y,
		z: v.z + v2.z,
		w: v.w,
	}
}

func (v *vec3d) Sub(v2 *vec3d) vec3d {
	return vec3d{
		x: v.x - v2.x,
		y: v.y - v2.y,
		z: v.z - v2.z,
		w: v.w,
	}
}

func (v *vec3d) Mul(scalar float64) vec3d {
	return vec3d{
		x: v.x * scalar,
		y: v.y * scalar,
		z: v.z * scalar,
		w: v.w,
	}
}

func (v *vec3d) Div(scalar float64) vec3d {
	return vec3d{
		x: v.x / scalar,
		y: v.y / scalar,
		z: v.z / scalar,
		w: v.w,
	}
}

func (v *vec3d) DotProduct(v2 *vec3d) float64 {
	return v.x*v2.x + v.y*v2.y + v.z*v2.z
}

func (v *vec3d) Length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v *vec3d) Normalize() vec3d {
	l := v.Length()
	return vec3d{
		x: v.x / l,
		y: v.y / l,
		z: v.z / l,
		w: v.w,
	}
}

func (v *vec3d) CrossProduct(v2 *vec3d) vec3d {
	return vec3d{
		x: v.y*v2.z - v.z*v2.x,
		y: v.z*v2.x - v.x*v2.z,
		z: v.x*v2.y - v.y*v2.x,
		w: v.w,
	}
}
