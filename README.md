cartconvert - a cartography package
===================================

cartconvert is a package providing a set of cartography functions for the
[Go](http://golang.org/) programming language to parse, convert, transform and
project coordinates. The set of functions is typically required when interacting
with online GIS systems as Google Maps, Bing Maps or OpenStreetMap.

Features
--------

The package provides the following functionality:

* Conversion between [polar
  coordinates](http://en.wikipedia.org/wiki/Polar_coordinate_system) to
  cartesian coordinates
* Supports a set of standard [reference
  ellipsoids](http://en.wikipedia.org/wiki/Reference_ellipsoid) (WGS84, Airy,
  Bessel) as well as user defined ones
* [Direct Transverse Mercator
  Projection](http://en.wikipedia.org/wiki/Transverse_Mercator_projection) and
  inverse thereof for the projection of a Geoid (model of the earth) onto the
  surface of a cylinder (map projection); Also know as Gauss-Krüger projection.
* [UTM coordinates](http://en.wikipedia.org/wiki/UTM_coordinate_system) to
  Latitude / Longitude
* [Geohashing:](http://en.wikipedia.org/wiki/Geohash) Latitude, Longitude to
  geohash and vice-versa
* [Helmert transformation](http://en.wikipedia.org/wiki/Helmert_transformation)
  to convert coordinates of one reference ellipsoidal model to another
* Various functions to parse different geodetic coordinate datums from string to
  internal data representations

Installation
------------

  goinstall github.com/the42/cartconvert

or alternatively download the package as tar file, extract the files into an
empty directory and run

  make install

Usage
-----

All features provided by the package are covered by test cases.

License
-------

The package is released under the [Simplified BSD
License](http://www.freebsd.org/copyright/freebsd-license.html) See file
"LICENSE"

Room for improvement
---------------------
The Geohashing implementation uses a well working yet inefficient implementation
of bitsets using strings. It would be good to replace this with a set of bitwise
operations on integers or use a (yet to emerge?) bit package

Implementation details
----------------------
The implementation of the direct transversal mercator projection uses the
iterative redfearn series algorithm. A very well explanation can be found in
"OGP Publication 373-7-2 – Surveying and Positioning Guidance Note number 7,
part 2 – November 2010"


Testing
-------

To run the tests:

  make test