package main

/*
 * Space station in 1075 mi high, circular orbit at 66.5 deg inclination.
 */

import (
	"fmt"
	"math"
)

func main() {

	const G = 6.673e-11 // m^3/(kg s^2)

	Mearth := 5.97e24

	// Initial conditions - 1075 mi circular orbit
	X := 6.371e6 + 1730044.745 // meters
	Y := 0.0
	Z := 0.0

	inclination := (66.5 / 360.) * 2.0 * math.Pi

	// Velocities in meters/second
	Vx := 0.0
	Vy := 7012.6 * math.Cos(inclination)
	Vz := 7012.6 * math.Sin(inclination)

	GM1 := G * Mearth

	// 8.101e6 m orbit radius
	// orbit circumference = 2*pi*8.101e6 = 5.09E7
	// t = 5.09e7/7012.5 = 7258 sec
	var t, r float64
	var intervalCount int
	dt := .250 // seconds

	fmt.Printf("# t\tVx\tVy\tVz\tx\ty\tz\n")
	for t = 0.0; t <= 86400; t += dt {

		r2 := X*X + Y*Y + Z*Z
		r = math.Sqrt(r2)

		if intervalCount%4 == 0 {
			fmt.Printf("%f\t%f\t%f\t%f\t%f\t%f\t%f\t%f\n", t, Vx, Vy, Vz, X, Y, Z, r)
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
}
