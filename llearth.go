package main

/*
 * Space station in 1075 mi high, circular orbit at 66.5 deg inclination.
 *
 * ground track of 24 hours of orbits, non-rotating earth
 *
 */

import (
	"image"
	"image/color"
	"image/gif"
	"log"
	"math"
	"os"

	"github.com/jonas-p/go-shp"
)

func main() {

	img := makeMap()

	const G = 6.673e-11 // m^3/(kg s^2)

	Mearth := 5.97e24

	rtod := 360. / (2. * math.Pi)

	// Initial conditions - 1075 mi circular orbit
	X := 6.371e6 + 1686870.745000 // meters
	Y := 0.0
	Z := 0.0

	inclination := (66.5 / 360.) * 2.0 * math.Pi

	// Velocities in meters/second
	Vmag := math.Sqrt(G * Mearth / X)
	Vx := 0.0
	Vy := Vmag * math.Cos(inclination)
	Vz := Vmag * math.Sin(inclination)

	GM1 := G * Mearth

	// 8.101e6 m orbit radius
	// orbit circumference = 2*pi*8.101e6 = 5.09E7
	// t = 5.09e7/7012.5 = 7258 sec
	var t, r float64
	var intervalCount int
	dt := .250 // seconds

	for t = 0.0; t <= 86400; t += dt {

		r2 := X*X + Y*Y + Z*Z
		r = math.Sqrt(r2)

		longitude := rtod * math.Atan2(Y, X)
		latitude := rtod * math.Atan2(Z, math.Sqrt(X*X+Y*Y))

		if intervalCount%4 == 0 {
			img.SetColorIndex(int(10.*(longitude+180.)), int(10.*(180.-latitude)), 2)
		}
		intervalCount++

		// magnitude of attraction F = G*M1*m2/(r^2)
		Fmag := GM1 / r2
		Fx := (-X / r) * Fmag
		Fy := (-Y / r) * Fmag
		Fz := (-Z / r) * Fmag

		dVx := Fx * dt
		dVy := Fy * dt
		dVz := Fz * dt

		Vx += dVx
		Vy += dVy
		Vz += dVz

		X += Vx * dt
		Y += Vy * dt
		Z += Vz * dt

	}

	gif.Encode(os.Stdout, img, &gif.Options{NumColors: 3, Quantizer: nil, Drawer: nil})
}

func makeMap() *image.Paletted {
	shape, err := shp.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer shape.Close()

	var pal []color.Color
	red := color.RGBA{255, 0, 0, 255}

	pal = append(pal, image.White) // 0
	pal = append(pal, image.Black) // 1
	pal = append(pal, red)         // 2

	img := image.NewPaletted(image.Rectangle{image.Point{0, 0}, image.Point{3600, 3600}}, pal)
	scale := 10.0

	for shape.Next() {
		_, p := shape.Shape()
		switch p.(type) {
		case *shp.PolyLine:
			pl := p.(*shp.PolyLine)
			for _, pt := range pl.Points {
				img.SetColorIndex(int(scale*(pt.X+180.)), int(scale*(-pt.Y+180.)), 1)
			}
		}
	}

	for x := 0.; x <= 360.0; x += 0.01 {
		img.SetColorIndex(int(scale*x), int(180.*scale), 2)
	}

	return img
}
