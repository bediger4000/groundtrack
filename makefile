IMAGES =  velocity.png orbital_radius.png

all: $(IMAGES)

circularorbit: circularorbit.go
	go build circularorbit.go

circ.dat: circularorbit
	./circularorbit > circ.dat

longlat: longlat.go
	go build longlat.go

ll.dat: longlat
	./longlat > ll.dat

llearth: llearth.go
	go build llearth.go

llearth.gif: llearth ne_10m_coastline.shp
	./llearth ne_10m_coastline.shp > llearth.gif

velocity.png: circ.dat velocity.load
	gnuplot < velocity.load
orbital_radius.png: circ.dat radius.load
	gnuplot <  radius.load

clean:
	- rm -rf circ.dat ll.dat $(IMAGES)
