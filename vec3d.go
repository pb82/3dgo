package main

import (
	"math"
)

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

func (v1 *vec3d) DotProduct(v2 *vec3d) float64 {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}

func (v *vec3d) Length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v *vec3d) Normalize() *vec3d {
	l := v.Length()
	return &vec3d{
		x: v.x / l,
		y: v.y / l,
		z: v.z / l,
		w: v.w,
	}
}

func (v1 *vec3d) CrossProduct(v2 *vec3d) vec3d {
	v := vec3d{}
	v.x = v1.y*v2.z - v1.z*v2.y
	v.y = v1.z*v2.x - v1.x*v2.z
	v.z = v1.x*v2.y - v1.y*v2.x
	return v
}

func vectorIntersectPlane(plane_p, plane_n, line_start, line_end *vec3d) vec3d {
	plane_n = plane_n.Normalize()
	plane_d := -plane_n.DotProduct(plane_p)
	ad := line_start.DotProduct(plane_n)
	bd := line_end.DotProduct(plane_n)
	t := (-plane_d - ad) / (bd - ad)
	lineStartToEnd := line_end.Sub(line_start)
	lineToIntersect := lineStartToEnd.Mul(t)
	return line_start.Add(&lineToIntersect)
}

func triangleClipAgainstPlane(plane_p, plane_n vec3d, in_tri, out_tri1, out_tri2 *triangle) int {
	plane_n = *plane_n.Normalize()

	dist := func(p *vec3d) float64 {
		// n := p.Normalize()
		return plane_n.x*p.x + plane_n.y*p.y + plane_n.z*p.z - plane_n.DotProduct(&plane_p)
	}

	d0 := dist(&in_tri.p[0])
	d1 := dist(&in_tri.p[1])
	d2 := dist(&in_tri.p[2])

	nInsidePointCount := 0
	inside_points := [3]*vec3d{}

	nOutsidePointCount := 0
	outside_points := [3]*vec3d{}

	if d0 >= 0 {
		inside_points[nInsidePointCount] = &in_tri.p[0]
		nInsidePointCount += 1
	} else {
		outside_points[nOutsidePointCount] = &in_tri.p[0]
		nOutsidePointCount += 1
	}

	if d1 >= 0 {
		inside_points[nInsidePointCount] = &in_tri.p[1]
		nInsidePointCount += 1
	} else {
		outside_points[nOutsidePointCount] = &in_tri.p[1]
		nOutsidePointCount += 1
	}

	if d2 >= 0 {
		inside_points[nInsidePointCount] = &in_tri.p[2]
		nInsidePointCount += 1
	} else {
		outside_points[nOutsidePointCount] = &in_tri.p[2]
		nOutsidePointCount += 1
	}

	if nInsidePointCount == 0 {
		// All points lie on the outside of plane, so clip whole triangle
		// It ceases to exist

		return 0 // No returned triangles are valid
	}

	if nInsidePointCount == 3 {
		// All points lie on the inside of plane, so do nothing
		// and allow the triangle to simply pass through
		*out_tri1 = *in_tri

		return 1 // Just the one returned original triangle is valid
	}

	if nInsidePointCount == 1 && nOutsidePointCount == 2 {
		// Triangle should be clipped. As two points lie outside
		// the plane, the triangle simply becomes a smaller triangle

		// Copy appearance info to new triangle
		out_tri1.r = math.MaxUint32
		out_tri1.g = 0
		out_tri1.b = 0
		out_tri1.a = math.MaxUint32

		// The inside point is valid, so keep that...
		out_tri1.p[0] = *inside_points[0]

		// but the two new points are at the locations where the
		// original sides of the triangle (lines) intersect with the plane
		out_tri1.p[1] = vectorIntersectPlane(&plane_p, &plane_n, inside_points[0], outside_points[0])
		out_tri1.p[2] = vectorIntersectPlane(&plane_p, &plane_n, inside_points[0], outside_points[1])

		return 1 // Return the newly formed single triangle
	}

	if nInsidePointCount == 2 && nOutsidePointCount == 1 {
		// Triangle should be clipped. As two points lie inside the plane,
		// the clipped triangle becomes a "quad". Fortunately, we can
		// represent a quad with two new triangles

		// Copy appearance info to new triangles
		out_tri1.r = 0
		out_tri1.g = math.MaxUint32
		out_tri1.b = 0
		out_tri1.a = math.MaxUint32

		out_tri2.r = 0
		out_tri2.g = 0
		out_tri2.b = math.MaxUint32
		out_tri2.a = math.MaxUint32

		// The first triangle consists of the two inside points and a new
		// point determined by the location where one side of the triangle
		// intersects with the plane
		out_tri1.p[0] = *inside_points[0]
		out_tri1.p[1] = *inside_points[1]
		out_tri1.p[2] = vectorIntersectPlane(&plane_p, &plane_n, inside_points[0], outside_points[0])

		// The second triangle is composed of one of he inside points, a
		// new point determined by the intersection of the other side of the
		// triangle and the plane, and the newly created point above
		out_tri2.p[0] = *inside_points[1]
		out_tri2.p[1] = out_tri1.p[2]
		out_tri2.p[2] = vectorIntersectPlane(&plane_p, &plane_n, inside_points[1], outside_points[0])

		return 2 // Return two newly formed triangles which form a quad
	}

	return 0
}
