IMAGES =  velocity.png orbital_radius.png

all: $(IMAGES)

circularorbit: circularorbit.go
	go build circularorbit.go

circ.dat: circularorbit
	./circularorbit > circ.dat

velocity.png: circ.dat velocity.load
	gnuplot < velocity.load
orbital_radius.png: circ.dat radius.load
	gnuplot <  radius.load

clean:
	- rm -rf circ.dat $(IMAGES)
