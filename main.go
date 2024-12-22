package main

import (
	"bufio"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	whiteImage = ebiten.NewImage(3, 3)
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

type triangle struct {
	p [3]vec3d
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

type mesh struct {
	tris []triangle
}

func (m *mesh) Load(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	var vertices []vec3d

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if line[0] == 'v' {
			line = line[2:]
			parts := strings.Split(line, " ")
			v := vec3d{}
			v.x, _ = strconv.ParseFloat(parts[0], 64)
			v.y, _ = strconv.ParseFloat(parts[1], 64)
			v.z, _ = strconv.ParseFloat(parts[2], 64)
			v.w = 1
			vertices = append(vertices, v)
		}

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
	trianglesToRaster []triangle
}

func (g *Game) Update() error {
	milliseconds := float64(time.Now().UnixMilli())
	delta := milliseconds - g.milliseconds
	g.milliseconds = milliseconds
	g.elapsedTime += delta
	// g.fTheta = 1.0 * (g.elapsedTime / 1000)

	g.rotZ.rotateZ(g.fTheta)
	g.rotX.rotateX(g.fTheta)
	g.trans.translate(0, 0, 5)

	g.matWorld = matrixMakeIdentity()
	g.matWorld = g.matWorld.multiplyMatrix(&g.rotZ)
	g.matWorld = g.matWorld.multiplyMatrix(&g.rotX)
	g.matWorld = g.matWorld.multiplyMatrix(&g.trans)

	return nil
}

func drawTriangle(screen *ebiten.Image, t *triangle) {
	path := &vector.Path{}
	path.MoveTo(t.X(0), t.Y(0))
	path.LineTo(t.X(1), t.Y(1))
	path.LineTo(t.X(2), t.Y(2))
	path.Close()

	vector.DrawFilledPath(screen, path, t, false, vector.FillRuleEvenOdd)
	// vector.StrokePath(screen, path, color.Black, false, &vector.StrokeOptions{Width: 1})
}

func getColor(lum float64) (uint32, uint32, uint32, uint32) {
	v := 64 * 1024 * lum
	return uint32(v), uint32(v), uint32(v), math.MaxUint32
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(clearColor)

	g.trianglesToRaster = nil

	// draw triangles
	for _, t := range g.mesh.tris {
		var triProjected triangle
		var triTransformed triangle

		triTransformed.p[0] = g.matWorld.matrixMultiplyVector(&t.p[0])
		triTransformed.p[1] = g.matWorld.matrixMultiplyVector(&t.p[1])
		triTransformed.p[2] = g.matWorld.matrixMultiplyVector(&t.p[2])

		// NORMAL
		line1 := triTransformed.p[1].Sub(&triTransformed.p[0])
		line2 := triTransformed.p[2].Sub(&triTransformed.p[0])
		normal := line1.CrossProduct(&line2)
		normal = normal.Normalize()

		vCameraRay := triTransformed.p[0].Sub(&g.vCamera)

		if normal.DotProduct(&vCameraRay) < 0.0 {

			light_direction := vec3d{0, 0, -1, 1}
			light_direction = light_direction.Normalize()

			dp := normal.x*light_direction.x + normal.y*light_direction.y + normal.z*light_direction.z

			rr, gg, bb, aa := getColor(dp)
			triTransformed.r = rr
			triTransformed.g = gg
			triTransformed.b = bb
			triTransformed.a = aa

			triProjected.p[0] = g.matProj.matrixMultiplyVector(&triTransformed.p[0])
			triProjected.p[1] = g.matProj.matrixMultiplyVector(&triTransformed.p[1])
			triProjected.p[2] = g.matProj.matrixMultiplyVector(&triTransformed.p[2])

			triProjected.p[0] = triProjected.p[0].Div(triProjected.p[0].w)
			triProjected.p[1] = triProjected.p[1].Div(triProjected.p[1].w)
			triProjected.p[2] = triProjected.p[2].Div(triProjected.p[2].w)

			triProjected.r = triTransformed.r
			triProjected.g = triTransformed.g
			triProjected.b = triTransformed.b
			triProjected.a = triTransformed.a

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

	// sort triangles from back to front
	sort.Slice(g.trianglesToRaster, func(i, j int) bool {
		z1 := (g.trianglesToRaster[i].p[0].z + g.trianglesToRaster[i].p[1].z + g.trianglesToRaster[i].p[2].z) / 3
		z2 := (g.trianglesToRaster[j].p[0].z + g.trianglesToRaster[j].p[1].z + g.trianglesToRaster[j].p[2].z) / 3
		return z1 > z2
	})

	for _, triProjected := range g.trianglesToRaster {
		drawTriangle(screen, &triProjected)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%.0f FPS", ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return w, h
}

func main() {
	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("3D Engine")

	cube := mesh{}
	cube.Load("./axis.obj")

	fNear := float64(0.1)
	fFar := float64(1000)
	fFov := float64(90)
	fAspectRatio := float64(h) / float64(w)
	fFovRad := 1.0 / math.Tan(fFov*0.5/180*math.Pi)

	projectionMatrix := matrixMakeProjection(fFovRad, fAspectRatio, fNear, fFar)

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
		vCamera: vec3d{
			x: 0,
			y: 0,
			z: 0,
			w: 1,
		},
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
