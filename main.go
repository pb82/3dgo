package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "image/png"
)

var (
	whiteImage = ebiten.NewImage(3, 3)
	//go:embed ryu.png
	textureData []byte
)

func init() {
	whiteImage.Fill(color.White)
}

var (
	clearColor = color.RGBA{
		R: 32,
		G: 32,
		B: 32,
		A: 255,
	}

	w = int(256)
	h = int(256)
)

type TextureAtlas interface {
	W() int
	H() int
	ColorAt(x, y int) color.Color
}

type TextureAtlasImpl struct {
	w, h int
	img  image.Image
}

func (t *TextureAtlasImpl) W() int {
	return t.w
}

func (t *TextureAtlasImpl) H() int {
	return t.h
}

func (t *TextureAtlasImpl) ColorAt(x, y int) color.Color {
	return t.img.At(x, y)
}

func (t *TextureAtlasImpl) LoadTexture() {
	texture, _, _ := image.Decode(bytes.NewReader(textureData))
	t.img = texture
	t.w = texture.Bounds().Dx()
	t.h = texture.Bounds().Dy()
}

type triangle struct {
	p [3]vec3d
	t UVs
	r uint32
	g uint32
	b uint32
	a uint32
}

func (t *triangle) X(index int) float32 {
	return float32(t.p[index].x)
}

func (t *triangle) Y(index int) float32 {
	return float32(t.p[index].y)
}

func (t *triangle) RGBA() (r, g, b, a uint32) {
	return t.r, t.g, t.b, t.a
}

func (t *triangle) Scale() {
	t.p[0].ScaleW()
	t.p[1].ScaleW()
	t.p[2].ScaleW()
}

type mesh struct {
	tris []triangle
}

func (m *mesh) translateX(dx float64) {
	for i, _ := range m.tris {
		m.tris[i].p[0].x += dx
		m.tris[i].p[1].x += dx
		m.tris[i].p[2].x += dx
	}
}

func (m *mesh) translateZ(dx float64) {
	for i, _ := range m.tris {
		m.tris[i].p[0].z += dx
		m.tris[i].p[1].z += dx
		m.tris[i].p[2].z += dx
	}
}

func (m *mesh) LoadCube() {
	m.tris = []triangle{
		{p: [3]vec3d{{0.0, 0.0, 0.0, 1}, {0.0, 1.0, 0.0, 1}, {1.0, 1.0, 0.0, 1}}, t: [3]vec2d{{0, 1, 1}, {0, 0, 1}, {1, 0, 1}}},
		{p: [3]vec3d{{0.0, 0.0, 0.0, 1}, {1.0, 1.0, 0.0, 1}, {1.0, 0.0, 0.0, 1}}, t: [3]vec2d{{0, 1, 1}, {1, 0, 1}, {1, 1, 1}}},
	}
}

func (m *mesh) Load(filename string, hasTexture bool) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	var vertices []vec3d
	var texs []vec2d

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if line[0] == 'v' {
			if line[1] == 't' {
				line = line[3:]
				v := vec2d{}
				parts := strings.Split(line, " ")
				v.u, _ = strconv.ParseFloat(parts[0], 64)
				v.v, _ = strconv.ParseFloat(parts[1], 64)
				v.w = 1
				texs = append(texs, v)
			} else {
				line = line[2:]
				parts := strings.Split(line, " ")
				v := vec3d{}
				v.x, _ = strconv.ParseFloat(parts[0], 64)
				v.y, _ = strconv.ParseFloat(parts[1], 64)
				v.z, _ = strconv.ParseFloat(parts[2], 64)
				v.w = 1
				vertices = append(vertices, v)
			}
		}

		if !hasTexture {
			if line[0] == 'f' {
				line = line[2:]
				var _x int64
				var _y int64
				var _z int64
				parts := strings.Split(line, " ")
				_x, _ = strconv.ParseInt(parts[0], 10, 32)
				_y, _ = strconv.ParseInt(parts[1], 10, 32)
				_z, _ = strconv.ParseInt(parts[2], 10, 32)

				m.tris = append(m.tris, triangle{
					p: [3]vec3d{vertices[_x-1], vertices[_y-1], vertices[_z-1]},
				})
			}
		} else {
			if line[0] == 'f' {
				line = line[2:]
				parts := strings.Split(line, " ")
				if len(parts) == 4 {
					pparts1 := strings.Split(parts[0], "/")
					vertex1, _ := strconv.ParseInt(pparts1[0], 10, 32)
					texture1, _ := strconv.ParseInt(pparts1[1], 10, 32)

					pparts2 := strings.Split(parts[1], "/")
					vertex2, _ := strconv.ParseInt(pparts2[0], 10, 32)
					texture2, _ := strconv.ParseInt(pparts2[1], 10, 32)

					pparts3 := strings.Split(parts[2], "/")
					vertex3, _ := strconv.ParseInt(pparts3[0], 10, 32)
					texture3, _ := strconv.ParseInt(pparts3[1], 10, 32)

					pparts4 := strings.Split(parts[3], "/")
					vertex4, _ := strconv.ParseInt(pparts4[0], 10, 32)
					texture4, _ := strconv.ParseInt(pparts4[1], 10, 32)

					m.tris = append(m.tris, triangle{
						p: [3]vec3d{vertices[vertex1-1], vertices[vertex2-1], vertices[vertex3-1]},
						t: [3]vec2d{texs[texture1-1], texs[texture2-1], texs[texture3-1]},
					})

					m.tris = append(m.tris, triangle{
						p: [3]vec3d{vertices[vertex1-1], vertices[vertex3-1], vertices[vertex4-1]},
						t: [3]vec2d{texs[texture1-1], texs[texture3-1], texs[texture4-1]},
					})
				} else if len(parts) == 3 {
					pparts1 := strings.Split(parts[0], "/")
					vertex1, _ := strconv.ParseInt(pparts1[0], 10, 32)
					texture1, _ := strconv.ParseInt(pparts1[1], 10, 32)

					pparts2 := strings.Split(parts[1], "/")
					vertex2, _ := strconv.ParseInt(pparts2[0], 10, 32)
					texture2, _ := strconv.ParseInt(pparts2[1], 10, 32)

					pparts3 := strings.Split(parts[2], "/")
					vertex3, _ := strconv.ParseInt(pparts3[0], 10, 32)
					texture3, _ := strconv.ParseInt(pparts3[1], 10, 32)

					m.tris = append(m.tris, triangle{
						p: [3]vec3d{vertices[vertex1-1], vertices[vertex2-1], vertices[vertex3-1]},
						t: [3]vec2d{texs[texture1-1], texs[texture2-1], texs[texture3-1]},
					})
				}
			}
		}
	}

	return true
}

type Game struct {
	mesh              mesh
	matProj           mat4x4
	matWorld          mat4x4
	milliseconds      float64
	elapsedTime       float64
	fTheta            float64
	vCamera           vec3d
	rotX              mat4x4
	rotY              mat4x4
	rotZ              mat4x4
	trans             mat4x4
	matView           mat4x4
	fYaw              float64
	trianglesToRaster []triangle
	tex               TextureAtlas
	depthBuffer       []float64
}

func (g *Game) Update() error {
	milliseconds := float64(time.Now().UnixMilli())
	delta := milliseconds - g.milliseconds
	g.milliseconds = milliseconds
	g.elapsedTime += delta

	msPassed := delta / 1000

	// g.mesh.translateZ(0.01)
	g.fTheta = 1.0 * (g.elapsedTime / 1000)

	// g.rotZ.rotateZ(g.fTheta)
	// g.rotX.rotateX(g.fTheta)
	// g.rotY.rotateY(g.fTheta)
	g.trans.translate(0, 0, 5)

	g.matWorld = matrixMakeIdentity()
	g.matWorld = g.matWorld.multiplyMatrix(&g.rotX)
	g.matWorld = g.matWorld.multiplyMatrix(&g.rotY)
	g.matWorld = g.matWorld.multiplyMatrix(&g.rotZ)
	g.matWorld = g.matWorld.multiplyMatrix(&g.trans)

	up := vec3d{0, 1, 0, 0}
	target := vec3d{0, 0, 1, 0}

	matCameraRot := matrixMakeIdentity()
	matCameraRot.rotateY(g.fYaw)

	vLookDirection := matCameraRot.matrixMultiplyVector(&target)

	target = g.vCamera.Add(&vLookDirection)

	camera := matrixMakeIdentity()
	camera.pointAt(&g.vCamera, &target, &up)
	g.matView = matrixQuickInverse(&camera)

	vForward := vLookDirection.Mul(8 * msPassed)

	keys := inpututil.AppendPressedKeys([]ebiten.Key{ebiten.KeyUp, ebiten.KeyDown, ebiten.KeyLeft, ebiten.KeyRight})
	for _, key := range keys {
		if key == ebiten.KeyW {
			g.vCamera = g.vCamera.Add(&vForward)
		}
		if key == ebiten.KeyS {
			g.vCamera = g.vCamera.Sub(&vForward)
		}
		if key == ebiten.KeyA {
			g.fYaw -= 1 * msPassed
		}
		if key == ebiten.KeyD {
			g.fYaw += 1 * msPassed
		}
		if key == ebiten.KeyUp {
			g.vCamera.y += 4 * msPassed
		}
		if key == ebiten.KeyDown {
			g.vCamera.y -= 4 * msPassed
		}
		if key == ebiten.KeyLeft {
			g.vCamera.x -= 4 * msPassed
		}
		if key == ebiten.KeyRight {
			g.vCamera.x += 4 * msPassed
		}
	}

	return nil
}

func swap[T comparable](a, b *T) {
	*a, *b = *b, *a
}

// 				screen.Set(int(j), int(i), sampleColour(tex_u/tex_w, tex_v/tex_w))

func drawTriangle(screen *ebiten.Image, t *triangle) {
	path := &vector.Path{}
	path.MoveTo(t.X(0), t.Y(0))
	path.LineTo(t.X(1), t.Y(1))
	path.LineTo(t.X(2), t.Y(2))
	path.Close()

	vector.DrawFilledPath(screen, path, t, false, vector.FillRuleEvenOdd)
	vector.StrokePath(screen, path, color.White, false, &vector.StrokeOptions{Width: 1})
}

func getColor(lum float64) (uint32, uint32, uint32, uint32) {
	if lum < 0.1 {
		lum = 0.1
	}
	v := 64 * 1024 * lum
	return uint32(v), uint32(v), uint32(v), math.MaxUint32
}

func TNormal(t *triangle) vec3d {
	line1 := t.p[1].Sub(&t.p[0])
	line2 := t.p[2].Sub(&t.p[0])
	normal := line1.CrossProduct(&line2)
	normal.Normalize()
	return normal
}

func (g *Game) Draw(screen *ebiten.Image) {
	t_start := time.Now()

	screen.Fill(clearColor)
	for i := range g.depthBuffer {
		g.depthBuffer[i] = 0
	}

	g.trianglesToRaster = nil

	// draw triangles
	for _, t := range g.mesh.tris {
		var triProjected triangle
		var triTransformed triangle
		var triViewed triangle

		triTransformed.p[0] = g.matWorld.matrixMultiplyVector(&t.p[0])
		triTransformed.p[1] = g.matWorld.matrixMultiplyVector(&t.p[1])
		triTransformed.p[2] = g.matWorld.matrixMultiplyVector(&t.p[2])
		triTransformed.t = t.t.Copy()

		// NORMAL
		normal := TNormal(&triTransformed)

		vCameraRay := triTransformed.p[0].Sub(&g.vCamera)
		dp := normal.DotProduct(&vCameraRay)

		if dp < 0 {
			light_direction := vec3d{0, 1, -1, 1}
			light_direction.Normalize()

			// dp := normal.x*light_direction.x + normal.y*light_direction.y + normal.z*light_direction.z

			triViewed.r = 20000
			triViewed.g = 2000
			triViewed.b = 2000
			triViewed.a = 2000

			// convert world space to view space
			triViewed.p[0] = g.matView.matrixMultiplyVector(&triTransformed.p[0])
			triViewed.p[1] = g.matView.matrixMultiplyVector(&triTransformed.p[1])
			triViewed.p[2] = g.matView.matrixMultiplyVector(&triTransformed.p[2])
			triViewed.t = triTransformed.t.Copy()

			// clip viewed triangle
			clipped := [2]triangle{}
			nClippedTriangles := triangleClipAgainstPlane(vec3d{0, 0, 0.1, 1}, vec3d{0, 0, 2.1, 1}, &triViewed, &clipped[0], &clipped[1])

			for n := 0; n < nClippedTriangles; n++ {
				// project from 3d to 2d
				triProjected.p[0] = g.matProj.matrixMultiplyVector(&clipped[n].p[0])
				triProjected.p[1] = g.matProj.matrixMultiplyVector(&clipped[n].p[1])
				triProjected.p[2] = g.matProj.matrixMultiplyVector(&clipped[n].p[2])
				triProjected.t[0] = clipped[n].t[0]
				triProjected.t[1] = clipped[n].t[1]
				triProjected.t[2] = clipped[n].t[2]

				triProjected.t.Scale(&triProjected)

				triProjected.r = clipped[n].r
				triProjected.g = clipped[n].g
				triProjected.b = clipped[n].b
				triProjected.a = clipped[n].a

				triProjected.Scale()

				// X/Y are inverted so put them back
				triProjected.p[0].x *= -1.0
				triProjected.p[1].x *= -1.0
				triProjected.p[2].x *= -1.0
				triProjected.p[0].y *= -1.0
				triProjected.p[1].y *= -1.0
				triProjected.p[2].y *= -1.0

				offsetView := vec3d{
					x: 1,
					y: 1,
					z: 0,
					w: 1,
				}

				triProjected.p[0] = triProjected.p[0].Add(&offsetView)
				triProjected.p[1] = triProjected.p[1].Add(&offsetView)
				triProjected.p[2] = triProjected.p[2].Add(&offsetView)

				triProjected.p[0].x *= 0.5 * float64(w)
				triProjected.p[0].y *= 0.5 * float64(h)
				triProjected.p[1].x *= 0.5 * float64(w)
				triProjected.p[1].y *= 0.5 * float64(h)
				triProjected.p[2].x *= 0.5 * float64(w)
				triProjected.p[2].y *= 0.5 * float64(h)

				g.trianglesToRaster = append(g.trianglesToRaster, triProjected)
			}
		}
	}

	// sort triangles from back to front
	/*
		sort.Slice(g.trianglesToRaster, func(i, j int) bool {
			z1 := (g.trianglesToRaster[i].p[0].z + g.trianglesToRaster[i].p[1].z + g.trianglesToRaster[i].p[2].z) / 3.0
			z2 := (g.trianglesToRaster[j].p[0].z + g.trianglesToRaster[j].p[1].z + g.trianglesToRaster[j].p[2].z) / 3.0
			return z1 > z2
		})
	*/

	trianglesDrawn := 0

	for _, triToRaster := range g.trianglesToRaster {
		clipped := [2]triangle{}
		var listTriangles []triangle
		listTriangles = append(listTriangles, triToRaster)
		newTriangles := 1

		for p := 0; p < 4; p++ {
			trisToAdd := 0
			for newTriangles > 0 {
				test := listTriangles[0]
				listTriangles = listTriangles[1:]
				newTriangles--

				// nClippedTriangles := triangleClipAgainstPlane(vec3d{0, 0, 0.1, 1}, vec3d{0, 0, 2.1, 1}, &triViewed, &clipped[0], &clipped[1])

				switch p {
				case 0:
					trisToAdd = triangleClipAgainstPlane(vec3d{0, 0, 0, 1}, vec3d{0, 1, 0, 1}, &test, &clipped[0], &clipped[1])
					break
				case 1:
					trisToAdd = triangleClipAgainstPlane(vec3d{0, float64(h - 1), 0, 1}, vec3d{0, -1, 0, 1}, &test, &clipped[0], &clipped[1])
					break
				case 2:
					trisToAdd = triangleClipAgainstPlane(vec3d{0, 0, 0, 1}, vec3d{1, 0, 0, 1}, &test, &clipped[0], &clipped[1])
					break
				case 3:
					trisToAdd = triangleClipAgainstPlane(vec3d{float64(w - 1), 0, 0, 1}, vec3d{-1, 0, 0, 1}, &test, &clipped[0], &clipped[1])
					break
				}

				for w := 0; w < trisToAdd; w++ {
					listTriangles = append(listTriangles, clipped[w])
				}
			}
			newTriangles = len(listTriangles)
		}

		for _, t := range listTriangles {
			// drawTriangle(screen, &t)
			g.texturedTriangle(
				int(t.p[0].x), int(t.p[0].y), t.t[0].u, t.t[0].v,
				int(t.p[1].x), int(t.p[1].y), t.t[1].u, t.t[1].v,
				int(t.p[2].x), int(t.p[2].y), t.t[2].u, t.t[2].v,
				t.t[0].w, t.t[1].w, t.t[2].w, screen)
			trianglesDrawn++
		}
	}

	t_duration := time.Since(t_start).Milliseconds()

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%.0f FPS, %d tris, %d rt", ebiten.ActualFPS(), trianglesDrawn, t_duration))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return w, h
}

func (g *Game) texturedTriangle(x1, y1 int, u1, v1 float64,
	x2, y2 int, u2, v2 float64,
	x3, y3 int, u3, v3 float64,
	w1, w2, w3 float64, screen *ebiten.Image) {

	if y2 < y1 {
		y1, y2 = y2, y1
		x1, x2 = x2, x1
		u1, u2 = u2, u1
		v1, v2 = v2, v1
		w1, w2 = w2, w1
	}

	if y3 < y1 {
		y1, y3 = y3, y1
		x1, x3 = x3, x1
		u1, u3 = u3, u1
		v1, v3 = v3, v1
		w1, w3 = w3, w1
	}

	if y3 < y2 {
		y2, y3 = y3, y2
		x2, x3 = x3, x2
		u2, u3 = u3, u2
		v2, v3 = v3, v2
		w2, w3 = w3, w2
	}

	dy1 := y2 - y1
	dx1 := x2 - x1
	dv1 := v2 - v1
	du1 := u2 - u1
	dw1 := w2 - w1

	dy2 := y3 - y1
	dx2 := x3 - x1
	dv2 := v3 - v1
	du2 := u3 - u1
	dw2 := w3 - w1

	var dax_step = float64(0)
	var dbx_step = float64(0)

	var du1_step = float64(0)
	var dv1_step = float64(0)
	var dw1_step = float64(0)
	var du2_step = float64(0)
	var dv2_step = float64(0)
	var dw2_step = float64(0)

	if dy1 >= 0 {
		dax_step = float64(dx1) / math.Abs(float64(dy1))
	}

	if dy2 >= 0 {
		dbx_step = float64(dx2) / math.Abs(float64(dy2))
	}

	if dy1 >= 0 {
		du1_step = float64(du1) / math.Abs(float64(dy1))
	}
	if dy1 >= 0 {
		dv1_step = float64(dv1) / math.Abs(float64(dy1))
	}
	if dy1 >= 0 {
		dw1_step = float64(dw1) / math.Abs(float64(dy1))
	}

	if dy2 >= 0 {
		du2_step = float64(du2) / math.Abs(float64(dy2))
	}
	if dy2 >= 0 {
		dv2_step = float64(dv2) / math.Abs(float64(dy2))
	}
	if dy2 >= 0 {
		dw2_step = float64(dw2) / math.Abs(float64(dy2))
	}

	if dy1 >= 0 {
		for i := y1; i <= y2; i++ {
			ax := float64(x1) + float64(i-y1)*dax_step
			bx := float64(x1) + float64(i-y1)*dbx_step

			tex_su := u1 + float64(i-y1)*du1_step
			tex_sv := v1 + float64(i-y1)*dv1_step
			tex_sw := w1 + float64(i-y1)*dw1_step

			tex_eu := u1 + float64(i-y1)*du2_step
			tex_ev := v1 + float64(i-y1)*dv2_step
			tex_ew := w1 + float64(i-y1)*dw2_step

			if ax > bx {
				ax, bx = bx, ax
				tex_su, tex_eu = tex_eu, tex_su
				tex_sv, tex_ev = tex_ev, tex_sv
				tex_sw, tex_ew = tex_ew, tex_sw
			}

			tex_u := tex_su
			tex_v := tex_sv
			tex_w := tex_sw

			tstep := 1.0 / (bx - ax)
			t := 0.0

			for j := ax; j < bx; j++ {
				tex_u = (1.0-t)*tex_su + t*tex_eu
				tex_v = (1.0-t)*tex_sv + t*tex_ev
				tex_w = (1.0-t)*tex_sw + t*tex_ew

				www := float64(g.tex.W() - 1)
				hhh := float64(g.tex.H() - 1)

				if tex_w > g.depthBuffer[i*w+int(j)] {
					screen.Set(int(j), i, g.tex.ColorAt(int((tex_u/tex_w)*www), int((1-tex_v/tex_w)*hhh)))
					g.depthBuffer[i*w+int(j)] = tex_w
				}

				t += tstep
			}
		}
	}

	dy1 = y3 - y2
	dx1 = x3 - x2
	dv1 = v3 - v2
	du1 = u3 - u2
	dw1 = w3 - w2

	if dy1 >= 0 {
		dax_step = float64(dx1) / math.Abs(float64(dy1))
	}

	if dy2 >= 0 {
		dbx_step = float64(dx2) / math.Abs(float64(dy2))
	}

	du1_step = 0
	dv1_step = 0

	if dy1 >= 0 {
		du1_step = du1 / math.Abs(float64(dy1))
	}
	if dy1 >= 0 {
		dv1_step = dv1 / math.Abs(float64(dy1))
	}
	if dy1 >= 0 {
		dw1_step = dw1 / math.Abs(float64(dy1))
	}

	if dy1 >= 0 {
		for i := y2; i <= y3; i++ {
			ax := float64(x2) + float64(i-y2)*dax_step
			bx := float64(x1) + float64(i-y1)*dbx_step

			tex_su := u2 + float64(i-y2)*du1_step
			tex_sv := v2 + float64(i-y2)*dv1_step
			tex_sw := w2 + float64(i-y2)*dw1_step

			tex_eu := u1 + float64(i-y1)*du2_step
			tex_ev := v1 + float64(i-y1)*dv2_step
			tex_ew := w1 + float64(i-y1)*dw2_step

			if ax > bx {
				ax, bx = bx, ax
				tex_su, tex_eu = tex_eu, tex_su
				tex_sv, tex_ev = tex_ev, tex_sv
				tex_sw, tex_ew = tex_ew, tex_sw
			}

			tex_u := tex_su
			tex_v := tex_sv
			tex_w := tex_sw

			tstep := 1.0 / (bx - ax)
			t := 0.0

			for j := ax; j < bx; j++ {
				tex_u = (1.0-t)*tex_su + t*tex_eu
				tex_v = (1.0-t)*tex_sv + t*tex_ev
				tex_w = (1.0-t)*tex_sw + t*tex_ew

				www := float64(g.tex.W() - 1)
				hhh := float64(g.tex.H() - 1)

				// Draw(j, i, tex->SampleGlyph(tex_u / tex_w, tex_v / tex_w), tex->SampleColour(tex_u / tex_w, tex_v / tex_w));
				if tex_w > g.depthBuffer[i*w+int(j)] {
					screen.Set(int(j), i, g.tex.ColorAt(int((tex_u/tex_w)*www), int((1-tex_v/tex_w)*hhh)))
					g.depthBuffer[i*w+int(j)] = tex_w
				}

				t += tstep
			}
		}
	}
}

func main() {
	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("3D Engine")

	cube := mesh{}
	cube.LoadCube()
	// cube.Load("./cube4.obj", true)

	fNear := float64(0.1)
	fFar := float64(1000)
	fFov := float64(90)
	fAspectRatio := float64(h) / float64(w)
	fFovRad := 1.0 / math.Tan(fFov*0.5/180*math.Pi)

	projectionMatrix := matrixMakeProjection(fFovRad, fAspectRatio, fNear, fFar)

	textureAtlas := &TextureAtlasImpl{}
	textureAtlas.LoadTexture()

	g := &Game{
		mesh:         cube,
		matProj:      projectionMatrix,
		milliseconds: float64(time.Now().UnixMilli()),
		elapsedTime:  0,
		fTheta:       0,
		matWorld:     matrixMakeIdentity(),
		rotX:         matrixMakeIdentity(),
		rotY:         matrixMakeIdentity(),
		rotZ:         matrixMakeIdentity(),
		trans:        matrixMakeIdentity(),
		matView:      matrixMakeIdentity(),
		tex:          textureAtlas,
		vCamera: vec3d{
			x: 0.5,
			y: 0.5,
			z: 4.5,
			w: 1,
		},
		depthBuffer: make([]float64, w*h),
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
