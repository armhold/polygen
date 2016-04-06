# Polygen uses a genetic algorithm to approximate an image with a small number of polygons.

###

You can generate a fairly good approximation with surprisingly few polygons. Here's a sample that
has only 50 polygons (~380K generations):

![starry-night](https://github.com/armhold/polygen/blob/master/images/starry.jpg "starry night (orig)")
![starry-night 50 polygons](https://github.com/armhold/polygen/blob/master/images/starry-50-polygons.png "starry night (50 polygons)")

This one is 100 polygons:

![mona-lisa](https://github.com/armhold/polygen/blob/master/images/mona_lisa.jpg "mona lisa (orig)")
![mona-lisa 100 polygons](https://github.com/armhold/polygen/blob/master/images/mona_lisa-100-polygons.png "mona lisa (100 polygons)")

###

The algorithm is pretty simple:

1. Create an initial string of candidate "DNA" consisting of a set of polygons (a color and a set of points) 
via random number generation. 

1. Render the DNA to an image (the "phenotype"). Compute its fitness by comparing to the reference image.

1. Apply random mutations to the candidate (change color, move polygon points, juggle the z-order) to 
create a population of offspring.

1. Evaluate the offspring, and if their fitness is better, replace the parent.

1. Repeat for N generations.


### Usage

1. `$ go get github.com/llgcode/draw2d/draw2dimg`
1. `go run cmd/polygen/main.go -source images/mona_lisa.jpg -poly 50`
1. Let it run until you are happy with the output (in `output.png`), or until you notice that there is not much change
between generations.


Polygen includes a built-in web server, so you can watch the image evolve in more or less realtime.
Just point your browser to [http://localhost:8080](http://localhost:8080).


![logo](https://github.com/armhold/polygen/blob/master/images/logo.gif "polygen Logo")

###


### Credit

This code is my own, but credit goes to Roger Johansson for the original idea,
which he documented [here](http://rogeralsing.com/2008/12/07/genetic-programming-evolution-of-mona-lisa). 

