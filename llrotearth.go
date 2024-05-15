package main

/*
 * Space station in 1075 mi high, circular orbit at 66.5 deg inclination.
 *
 * Ground track of 24 hours of orbit on a rotating earth,
 * Map projection is longitude -> X, latitude -> Y
 */

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math"
	"os"

	"github.com/jonas-p/go-shp"
)

func main() {

	img, err := MakeMap(3600, 1800, os.Args[1], true)
	if err != nil {
		log.Fatal(err)
	}

	const G = 6.673e-11 // m^3/(kg s^2)

	Mearth := 5.97e24

	GM1 := G * Mearth

	// Initial conditions - 1075 mi circular orbit
	X := 6.371e6 + 1686870.745000 // meters
	Y := 0.0
	Z := 0.0

	inclination := (66.5 / 360.) * 2.0 * math.Pi

	// Velocities in meters/second
	Vmag := math.Sqrt(GM1 / X)
	Vx := 0.0
	Vy := Vmag * math.Cos(inclination)
	Vz := Vmag * math.Sin(inclination)

	// 8.101e6 m orbit radius
	// orbit circumference = 2*pi*8.101e6 = 5.09E7
	// t = 5.09e7/7012.5 = 7258 sec
	var t, r float64
	var intervalCount int
	dt := .250 // seconds

	max := 86400. // + 24.*36.

	for t = 0.0; t <= max; t += dt {

		r2 := X*X + Y*Y + Z*Z
		r = math.Sqrt(r2)

		if intervalCount%4 == 0 {
			img.Point(X, Y, Z, t)
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

	img.WriteImage(os.Stdout)
}

type globeImage struct {
	image   *image.Paletted
	palette []color.Color
	scale   float64
	offsetX float64 // add to longitude
	offsetY float64 // subtract latitude
}

// MakeMap creates a filled in *globalImage
// from the shapefile named by fileName argument
func MakeMap(width, height int, fileName string, verbose bool) (*globeImage, error) {
	if verbose {
		fmt.Fprintf(os.Stderr, "Shape file %q\n", fileName)
	}
	shape, err := shp.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer shape.Close()

	var pal []color.Color
	red := color.RGBA{255, 0, 0, 255}

	pal = append(pal, image.White) // 0
	pal = append(pal, image.Black) // 1
	pal = append(pal, red)         // 2

	if verbose {
		fmt.Fprintf(os.Stderr, "Creating image %d wide, %d high\n", width, height)
	}
	img := image.NewPaletted(image.Rectangle{image.Point{0, 0}, image.Point{width, height}}, pal)

	scale := float64(width) / 360.

	gi := &globeImage{
		image:   img,
		palette: pal,
		scale:   scale,
		offsetX: float64(width) / (2. * scale),
		offsetY: float64(height) / (2. * scale),
	}
	if verbose {
		gi.Describe(os.Stderr)
	}

	for shape.Next() {
		_, p := shape.Shape()
		switch p.(type) {
		case *shp.PolyLine:
			pl := p.(*shp.PolyLine)
			for _, pt := range pl.Points {
				gi.image.SetColorIndex(
					int(gi.scale*(pt.X+gi.offsetX)),
					int(gi.scale*(gi.offsetY-pt.Y)),
					1,
				)
			}
		}
	}

	for x := 0.; x <= 360.0; x += 0.01 {
		gi.image.SetColorIndex(int(gi.scale*x), int(gi.scale*gi.offsetY), 2)
	}

	return gi, nil
}

var rtod = 360. / (2. * math.Pi) // radians to degrees multiplier

func (gi *globeImage) Point(X, Y, Z, t float64) {
	longitude := rtod * math.Atan2(Y, X)
	latitude := rtod * math.Atan2(Z, math.Sqrt(X*X+Y*Y))
	dlong := t * 0.004166
	dlong = math.Remainder(dlong, 360.)
	longitude = math.Remainder(longitude-dlong, 360)
	gi.image.SetColorIndex(
		int(gi.scale*(longitude+gi.offsetX)),
		int(gi.scale*(gi.offsetY-latitude)),
		2,
	)
}

func (gi *globeImage) WriteImage(fout *os.File) {
	gif.Encode(
		fout,
		gi.image,
		&gif.Options{
			NumColors: len(gi.palette),
			Quantizer: nil,
			Drawer:    nil,
		},
	)
}
func (gi *globeImage) Describe(fout *os.File) {
	fmt.Fprintf(fout, "paletted image with %d colors\n", len(gi.palette))
	fmt.Fprintf(fout, "input long/lat scaled by %.03f\n", gi.scale)
	fmt.Fprintf(fout, "input long/lat offset by (%.03f, %.03f)\n", gi.offsetX, gi.offsetY)
}
