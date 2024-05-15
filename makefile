IMAGES =  velocity.png orbital_radius.png llearth.gif rottrack.gif \
		  	hammer_groundtrack.gif

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

llrotearth: llrotearth.go
	go build llrotearth.go

rottrack.gif: llrotearth ne_10m_coastline.shp
	./llrotearth ne_10m_coastline.shp > rottrack.gif

llrotearthhammer: llrotearthhammer.go
	go build llrotearthhammer.go

hammer_groundtrack.gif: llrotearthhammer ne_10m_coastline.shp
	./llrotearthhammer ne_10m_coastline.shp > hammer_groundtrack.gif

velocity.png: circ.dat velocity.load
	gnuplot < velocity.load
orbital_radius.png: circ.dat radius.load
	gnuplot <  radius.load

clean:
	- rm -rf circ.dat ll.dat $(IMAGES)
