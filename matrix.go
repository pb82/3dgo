package main

import "math"

type mat4x4 struct {
	m [4][4]float64
}

func (m *mat4x4) multiplyMatrix(m2 *mat4x4) mat4x4 {
	matrix := mat4x4{}
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			matrix.m[r][c] = m.m[r][0]*m2.m[0][c] + m.m[r][1]*m2.m[1][c] + m.m[r][2]*m2.m[2][c] + m.m[r][3]*m2.m[3][c]
		}
	}
	return matrix
}

func makeMatrix() mat4x4 {
	return mat4x4{
		m: [4][4]float64{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
	}
}

func (m *mat4x4) matrixMultiplyVector(i *vec3d) vec3d {
	v := vec3d{}
	v.x = i.x*m.m[0][0] + i.y*m.m[1][0] + i.z*m.m[2][0] + i.w*m.m[3][0]
	v.y = i.x*m.m[0][1] + i.y*m.m[1][1] + i.z*m.m[2][1] + i.w*m.m[3][1]
	v.z = i.x*m.m[0][2] + i.y*m.m[1][2] + i.z*m.m[2][2] + i.w*m.m[3][2]
	v.w = i.x*m.m[0][3] + i.y*m.m[1][3] + i.z*m.m[2][3] + i.w*m.m[3][3]
	return v
}

func matrixMakeIdentity() mat4x4 {
	m := mat4x4{}
	m.m[0][0] = 1
	m.m[1][1] = 1
	m.m[2][2] = 1
	m.m[3][3] = 1
	return m
}

func (m *mat4x4) translate(x, y, z float64) {
	m.m[0][0] = 1
	m.m[1][1] = 1
	m.m[2][2] = 1
	m.m[3][3] = 1
	m.m[3][0] = x
	m.m[3][1] = y
	m.m[3][2] = z
}

func (m *mat4x4) rotateY(angleRad float64) {
	m.m[0][0] = math.Cos(angleRad)
	m.m[0][2] = math.Sin(angleRad)
	m.m[2][0] = -math.Sin(angleRad)
	m.m[1][1] = 1
	m.m[2][2] = math.Cos(angleRad)
	m.m[3][3] = 1
}

func (m *mat4x4) rotateZ(angleRad float64) {
	m.m[0][0] = math.Cos(angleRad)
	m.m[0][1] = math.Sin(angleRad)
	m.m[1][0] = -math.Sin(angleRad)
	m.m[1][1] = math.Cos(angleRad)
	m.m[2][2] = 1
	m.m[3][3] = 1
}

func (m *mat4x4) rotateX(angleRad float64) {
	m.m[0][0] = 1
	m.m[1][1] = math.Cos(angleRad)
	m.m[1][2] = math.Sin(angleRad)
	m.m[2][1] = -math.Sin(angleRad)
	m.m[2][2] = math.Cos(angleRad)
	m.m[3][3] = 1
}

func matrixMakeProjection(fFovRad, fAspectRatio, fNear, fFar float64) mat4x4 {
	matrix := mat4x4{}
	matrix.m[0][0] = fAspectRatio * fFovRad
	matrix.m[1][1] = fFovRad
	matrix.m[2][2] = fFar / (fFar - fNear)
	matrix.m[3][2] = (-fFar * fNear) / (fFar - fNear)
	matrix.m[2][3] = 1.0
	matrix.m[3][3] = 0.0
	return matrix
}

func makeRotateY(angleRad float64) mat4x4 {
	m := makeMatrix()
	m.m[0][0] = math.Cos(angleRad)
	m.m[0][2] = math.Sin(angleRad)
	m.m[2][0] = -math.Sin(angleRad)
	m.m[1][1] = 1
	m.m[2][2] = math.Cos(angleRad)
	m.m[3][3] = 1
	return m
}

func makeRotateZ(angleRad float64) mat4x4 {
	m := makeMatrix()
	m.m[0][0] = math.Cos(angleRad)
	m.m[0][1] = math.Sin(angleRad)
	m.m[1][0] = -math.Sin(angleRad)
	m.m[1][1] = math.Cos(angleRad)
	m.m[2][2] = 1
	m.m[3][3] = 1
	return m
}

func makeRotateX(angleRad float64) mat4x4 {
	m := makeMatrix()
	m.m[0][0] = 1
	m.m[1][1] = math.Cos(angleRad)
	m.m[1][2] = math.Sin(angleRad)
	m.m[2][1] = -math.Sin(angleRad)
	m.m[2][2] = math.Cos(angleRad)
	m.m[3][3] = 1
	return m
}

func makeTranslate(x, y, z float64) mat4x4 {
	m := makeMatrix()
	m.m[0][0] = 1
	m.m[1][1] = 1
	m.m[2][2] = 1
	m.m[3][3] = 1
	m.m[3][0] = x
	m.m[3][1] = y
	m.m[3][2] = z
	return m
}
