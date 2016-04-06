# Polygen - implementing an image efficiently with polygons (golang).

![logo](https://github.com/armhold/polygen/blob/master/logo.gif "polygen Logo")


Polygen uses a genetic algorithm to approximate an image with a small number of polygons. 

The algorithm is pretty simple:

1. Create an initial string of candidate "DNA" consisting of a set of polygons (a color and a set of points) 
via random number generation. 

1. Render the DNA to an image (the "phenotype"). Compute its fitness by comparing to the reference image.

1. Apply random mutations to the candidate (change color, move polygon points, juggle the z-order) to 
create a population of offspring.

1. Evaluate the offspring, and if their fitness is better, replace the parent.

1. Repeat for N generations.

You can generate a fairly good approximation with surprisingly few polygons. Here's a sample that
has only 50 polygons (~380K generations):


![starry-night](https://github.com/armhold/polygen/blob/master/images/starry.jpg "starry night (orig)")
![starry-night 50 polygons](https://github.com/armhold/polygen/blob/master/images/starry-polygons.jpg "starry night (50 polygons)")


Polygen includes a built-in web server, so you can watch the image evolve in more or less realtime.
Just point your browser to http://localhost:8080.

_









Credit: this code is my own, but credit goes to Roger Johansson for the original idea, 
which he documented [here](http://rogeralsing.com/2008/12/07/genetic-programming-evolution-of-mona-lisa). 

