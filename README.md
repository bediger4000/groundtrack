# Ground track of an earth orbiting satellite

I would like a ground track of Wernher von Braun's
1952 Space Station from [Across the Space Frontier]().
It's in a 1075 mile high, circular orbit inclined at 66.5&deg to the equator.

The interesting thing about a ground track is that the satellite's orbit
is stable in XYZ cartesian coordinates,
but the earth rotates underneath it.

## Process

1. Write a computer program to simulate a 1075 mile high, 66.5&deg; inclined, circular orbit.
See if the program is numerically stable for several days of simulated time.
    - [code](circularorbit.go)
2. Figure out how to map (X,Y,Z) orbital positions to longitude/latitude
3. Produce an image of a flattned-out ground track on a non-rotating earth.
May (longitude, latitude) to (X,Y) coordinates in the image.
4. Figure out how to produce a [Hammer projection]() of the ground track
on a non-rotating earth.
5. Include the effects of a rotating earth on the Hammer projection ground track
