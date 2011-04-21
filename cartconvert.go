// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

//target:github.com/the42/cartconvert

// This package provides a series of functions to deal with 
// conversion, transformation and projection of coordinate systems.
package cartconvert

import (
	"math"
	"strings"
	"strconv"
	"fmt"
	"os"
)

// A CartographyError is yielded when a literal can not be parsed as a bearing specifier.
// In this case the following values may be set and carry the meaning:
type CartographyError struct {
	Coord string   // a fragment of the coordinate literal containing the error.
	Val   float64  // The value parsed so far
	Index int      // Position at which the error occured
	Error os.Error // inherited error from an attempt of strconv to parse a number
}

func (ce CartographyError) String() string {
	return fmt.Sprintf("unable to parse fragment \"%s\". Partial value: %f. The additional error was: %s", ce.Coord, ce.Val, ce.Error.String())
}

// Set of common ellipsoidal models regularly found in cartography 
var (
	Bessel1841MGIEllipsoid = NewEllipsoid(6377397.155, 6356078.965)
	Bessel1841Ellipsoid    = NewEllipsoid(6377397.155, 6356078.962822)
	GRS80Ellipsoid         = NewEllipsoid(6378137, 6356752.31414)
	WGS84Ellipsoid         = NewEllipsoid(6378137, 6356752.31425)
	Airy1830Ellipsoid      = NewEllipsoid(6377563.396, 6356256.909)
	DefaultEllipsoid       = WGS84Ellipsoid
)

type Ellipsoid struct {
	a, b float64
}

// Holds latitude, longitude and ellipsoidal height, relative to El, the reference ellipsoid
type PolarCoord struct {
	Latitude, Longitude, Height float64
	El                          *Ellipsoid
}

// specifier for the string representation of a polar coordinate
type LatLongFormat int

const (
	LLFUnknown LatLongFormat = iota
	LLFdeg                   // format a lat/long coordinate in degrees using leading sign for negative bearings
	LLFdms                   // format a lat/long coordinate in degrees, minutes and seconds with prepended main directions N, S, E, W
)

func LatLongToString(pc *PolarCoord, format LatLongFormat) (pcs string) {

	var lat, long, latrem, longrem, latmin, longmin, latsec, longsec float64
	var latitude, longitude string

	switch format {
	case LLFdeg:
		pcs = "lat: " + f64toa(pc.Latitude, 6) + "°, long: " + f64toa(pc.Longitude, 6) + "°"
	case LLFdms:
		lat, latrem = math.Modf(pc.Latitude)

		if lat < 0 {
			lat *= -1
			latrem *= -1
		}

		long, longrem = math.Modf(pc.Longitude)
		if long < 0 {
			long *= -1
			longrem *= -1
		}

		latmin, latrem = math.Modf(latrem / 100 * 6000)
		longmin, longrem = math.Modf(longrem / 100 * 6000)

		latsec = latrem / 100 * 6000
		longsec = longrem / 100 * 6000

		if pc.Latitude < 0 {
			latitude = "S "

		} else {
			latitude = "N "
		}
		latitude += fmt.Sprintf("%d°", int(lat))
		if latmin != 0.0 || latsec != 0.0 {
			latitude += fmt.Sprintf("%d'", int(latmin))
		}
		if latsec != 0.0 {
			latitude += fmt.Sprintf("%s''", f64toa(latsec, 2))
		}

		if pc.Longitude < 0 {
			longitude = "W "
		} else {
			longitude = "E "
		}
		longitude += fmt.Sprintf("%d°", int(long))
		if longmin != 0.0 || longsec != 0.0 {
			longitude += fmt.Sprintf("%d'", int(longmin))
		}
		if longsec != 0.0 {
			longitude += fmt.Sprintf("%s''", f64toa(longsec, 2))
		}

		pcs = latitude + ", " + longitude
	}
	return
}

// Canonical representation of a lat/long bearing
func (pc *PolarCoord) String() string {
	return LatLongToString(pc, LLFdeg)
}


// A generic representation of easting (right, Y) and northing (Height,X) of a 2D projection
// relative to Ellipsoid El. The height H at Point X,Y is above defining ellipsoid 
type GeoPoint struct {
	X, Y, H float64
	El      *Ellipsoid
}

// A generic Cartesian, geocentric point. For ease of conversion between polar and Cartesian
// coordinates, the ellipsis might be included 
type CartPoint struct {
	X, Y, Z float64
	El      *Ellipsoid
}

func degtorad(deg float64) float64 {
	return math.Pi * deg / 180
}

func radtodeg(rad float64) float64 {
	return 180 * rad / math.Pi
}

func removeblank(input string) (accu string) {

	for _, token := range input {
		switch token {
		case ' ':
			continue
		default:
			accu += string(token)
		}
	}

	return
}

// The function accepts a string representing a bearing in degree, minute and second.
// Minute and second are booth optional and the second may contain fractions.
// The bearing may be prepended by the literal 'N', 'E', 'S', 'W', representing the
// four main directions. 'S' and 'W' denotes negative bearing. Instead of the main directions,
// the signs '+' or '-' may be used.
//
// [N|E|S|W|+|-]ddd°[dd'[dd'']]
func ADegMMSSToNum(DegMMSS string) (degf float64, err os.Error) {

	var accu string
	var i, position int
	var token int
	var tf float64

	degree := strings.ToUpper(removeblank(DegMMSS))
	slen := len(degree)

	negate := false

	// parse the degree
L2:
	for i, token = range degree {
		switch token {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			accu += string(token)
		case 'S', 'W', '-':
			negate = true
		case 'N', 'E', '+':
			continue
		case '°':
			degf, err = strconv.Atof64(accu)

			if err != nil {
				return degf, err
			}

			i += len("°")
			degree = degree[i:]
			position += i
			accu = ""
			break L2
		default:
			return 0, CartographyError{Val: degf, Index: i, Coord: degree, Error: os.EINVAL}
		}
	}

L4:
	// parse the minute
	for i, token = range degree {
		switch token {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			accu += string(token)
		case '\'':
			tf, err = strconv.Atof64(accu)

			if err != nil {
				return tf, err
			}

			degf += tf / 60.0
			i += len("'")
			degree = degree[i:]
			position += i
			accu = ""
			break L4
		default:
			return 0, CartographyError{Val: degf, Index: i, Coord: degree, Error: os.EINVAL}
		}
	}

	// parse the second
L6:
	for i, token = range degree {
		switch token {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
			accu += string(token)
		case '\'':
			tf, err = strconv.Atof64(accu)

			if err != nil {
				return tf, err
			}

			degf += tf / (60.0 * 60.0)
			i += len("'")
			degree = degree[i:]
			position += i

			if !(position < slen && degree[0] == '\'') {
				return 0, CartographyError{Val: degf, Index: i, Coord: degree, Error: os.EINVAL}
			}
			break L6
		default:
			return 0, CartographyError{Val: degf, Index: i, Coord: degree, Error: os.EINVAL}
		}
	}

	if negate {
		degf = -degf
	}

	return
}

// The function accepts a string literal representing a bearing
// denoted as decimal degrees. The suffix is optional.
// The literal value must end with the degree mark '°'
// The bearing may be prepended by the literal 'N', 'E', 'S', 'W', representing the
// four main directions. 'S' and 'W' denotes negative bearing. Instead of the main directions,
// the signs '+' or '-' may be used.
//
// [N|E|S|W|+|-]ddd[.suffix]°
func ADegCommaToNum(DegComma string) (degf float64, err os.Error) {

	var accu string
	var i, token int
	var tf float64

	degree := strings.ToUpper(removeblank(DegComma))
	negate := false

	// parse the degree
L2:
	for i, token = range degree {
		switch token {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			accu += string(token)
		case 'S', 'W', '-':
			negate = true
		case 'N', 'E', '+':
			continue
		case '.', '°':
			degf, err = strconv.Atof64(accu)

			if err != nil {
				return degf, err
			}

			degree = degree[i+len(string(token)):]
			accu = ""
			break L2
		default:
			return 0, CartographyError{Val: degf, Index: i, Coord: degree, Error: os.EINVAL}
		}
	}
L4:
	// parse the suffix
	for i, token = range degree {
		switch token {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			accu += string(token)
		case '°':
			tf, err = strconv.Atof64("0." + accu)

			if err != nil {
				return tf, err
			}

			degf += tf
			accu = ""
			break L4
		default:
			return 0, CartographyError{Val: degf, Index: i, Coord: degree, Error: os.EINVAL}
		}
	}

	if len(accu) > 0 {
		return 0, CartographyError{Val: degf, Index: i, Coord: degree, Error: os.EINVAL}
	}

	if negate {
		degf = -degf
	}

	return
}

// ## Polar to Cartesian coordinate conversion and vice-versa

// Function accepts two bearing datum as Deg°MM'SS'' (typically northing and easting)
// and height at bearing relative to reference ellipsoid and returns a polar coordinate type.
// If the reference ellipsoid is nil, the DefaultEllipsoid will be set in the resulting polar coordinate.
func ADegMMSSToPolar(Northing, Easting string, Height float64, El *Ellipsoid) (*PolarCoord, os.Error) {

	northing, err := ADegMMSSToNum(Northing)

	if err == nil {
		easting, err := ADegMMSSToNum(Easting)

		if err == nil {

			el := El
			if el == nil {
				el = DefaultEllipsoid
			}

			return &PolarCoord{Latitude: northing, Longitude: easting, Height: Height, El: el}, nil
		}
	}
	return nil, err
}

// Convert polar coordinates to Cartesian. The polar coordinates must be in decimal degrees.
// The reference ellipsoid is copied verbatim to the result. 
// Inspired by http://www.movable-type.co.uk/scripts/latlong-convert-coords.html
func PolarToCartesian(gc *PolarCoord) *CartPoint {

	var p CartPoint

	el := gc.El

	lat := degtorad(gc.Latitude)
	long := degtorad(gc.Longitude)

	esq := (el.a*el.a - el.b*el.b) / (el.a * el.a)

	u := el.a / math.Sqrt(1-esq*math.Pow(math.Sin(lat), 2))

	p.X = (gc.Height + u) * math.Cos(lat) * math.Cos(long)
	p.Y = (gc.Height + u) * math.Cos(lat) * math.Sin(long)
	p.Z = ((1-esq)*u + gc.Height) * math.Sin(lat)

	p.El = el

	return &p
}

// Convert Cartesian coordinates to polar.
// The reference ellipsoid is copied verbatim to the result.
// The resulting polar coordinates are in decimal degrees.
// Inspired by http://www.movable-type.co.uk/scripts/latlong-convert-coords.html
func CartesianToPolar(pt *CartPoint) *PolarCoord {

	var gc PolarCoord

	el := pt.El

	esq := (el.a*el.a - el.b*el.b) / (el.a * el.a)
	p := math.Hypot(pt.X, pt.Y)

	lat := math.Atan2(pt.Z, p*(1-esq))
	lat0 := 2.0 * math.Pi

	precision := 4.0 / el.a
	var v float64
	for math.Fabs(lat-lat0) > precision {
		v = el.a / math.Sqrt(1-esq*math.Pow(math.Sin(lat), 2))

		lat0 = lat
		lat = math.Atan2(pt.Z+esq*v*math.Sin(lat), p)
	}

	gc.Height = p/math.Cos(lat) - v
	gc.Latitude = radtodeg(lat)
	gc.Longitude = radtodeg(math.Atan2(pt.Y, pt.X))

	gc.El = el

	return &gc
}

// ## Transverse Mercator Projection

// Direct transverse mercator projection: Projection of an ellipsoid onto the surface of
// of a cylinder. Also known as Gauss-Krüger projection. Input parameters:
//
//	gc *PolarCoord: Latitude and Longitude or point to be projected; in decimal degrees
//	latO, longO: Shifted origin of latitude and longitude in decimal degrees
//	fe, fn: False easting and northing respectively in meters
//	scale: Projection scaling; Dimensionless, typically 1 or little bellow
//
// This algorithm uses the algorithm described by Redfearn
// http://en.wikipedia.org/wiki/Transverse_Mercator:_Redfearn_series
//
// Taken from "OGP Publication 373-7-2 – Surveying and Positioning Guidance Note number 7, part 2 – November 2010",
// pp. 48 - 51
func DirectTransverseMercator(gc *PolarCoord, latO, longO, scale, fe, fn float64) *GeoPoint {

	var pt GeoPoint

	el := gc.El

	latOrad := degtorad(latO)
	longOrad := degtorad(longO)

	latrad := degtorad(gc.Latitude)
	longrad := degtorad(gc.Longitude)

	f := 1 - el.b/el.a
	esq := math.Sqrt(2.0*f - f*f)

	n := f / (2.0 - f)
	B := (el.a / (1 + n)) * (1 + n*n/4.0 + n*n*n*n/64.0)

	h1 := n/2.0 - (2.0/3.0)*(n*n) + (5.0/16.0)*(n*n*n) + (41.0/180.0)*(n*n*n*n)
	h2 := (13.0/48.0)*(n*n) - (3.0/5.0)*(n*n*n) + (557.0/1440.0)*(n*n*n*n)
	h3 := (61.0/240.0)*(n*n*n) - (103.0/140.0)*(n*n*n*n)
	h4 := (49561.0 / 161280.0) * (n * n * n * n)

	var SO float64

	if latOrad != 0.0 {
		QO := math.Asinh(math.Tan(latOrad)) - (esq * math.Atanh(esq*math.Sin(latOrad)))
		bO := math.Atan(math.Sinh(QO))
		xiO0 := bO // math.Asin(math.Sin(bO))

		xiO1 := h1 * math.Sin(2.0*xiO0)
		xiO2 := h2 * math.Sin(4.0*xiO0)
		xiO3 := h3 * math.Sin(6.0*xiO0)
		xiO4 := h4 * math.Sin(8.0*xiO0)

		xiO := xiO0 + xiO1 + xiO2 + xiO3 + xiO4

		SO = B * xiO
	}

	Q := math.Asinh(math.Tan(latrad)) - (esq * math.Atanh(esq*math.Sin(latrad)))
	b := math.Atan(math.Sinh(Q))

	eta0 := math.Atanh(math.Cos(b) * math.Sin(longrad-longOrad))
	xi0 := math.Asin(math.Sin(b) * math.Cosh(eta0))

	xi1 := h1 * math.Sin(2*xi0) * math.Cosh(2*eta0)
	xi2 := h2 * math.Sin(4*xi0) * math.Cosh(4*eta0)
	xi3 := h3 * math.Sin(6*xi0) * math.Cosh(6*eta0)
	xi4 := h4 * math.Sin(8*xi0) * math.Cosh(8*eta0)
	xi := xi0 + xi1 + xi2 + xi3 + xi4

	eta1 := h1 * math.Cos(2*xi0) * math.Sinh(2*eta0)
	eta2 := h2 * math.Cos(4*xi0) * math.Sinh(4*eta0)
	eta3 := h3 * math.Cos(6*xi0) * math.Sinh(6*eta0)
	eta4 := h4 * math.Cos(8*xi0) * math.Sinh(8*eta0)
	eta := eta0 + eta1 + eta2 + eta3 + eta4

	pt.X = fe + scale*B*eta
	pt.Y = fn + scale*(B*xi-SO)

	pt.El = el

	return &pt
}

// Inverse transverse mercator projection: Projection of an cylinder onto the surface of
// of an ellipsoid. Also known as reverse Gauss-Krüger projection. Input parameters:
//
//	pt *GeoPoint: Easting(Y) and Northing(X) of map point to be projected; in meters
//	latO, longO: Shifted origin of latitude and longitude in decimal degrees
//	fe, fn: False easting and northing respectively in meters
//	scale: Projection scaling; Dimensionless, typically 1 or little bellow
//
// This algorithm uses the algorithm described by Redfearn
// http://en.wikipedia.org/wiki/Transverse_Mercator:_Redfearn_series
//
// Taken from "OGP Publication 373-7-2 – Surveying and Positioning Guidance Note number 7, part 2 – November 2010",
// pp. 48 - 51
//
// More accurate, iterative but slower algorithmic implementation
func InverseTransverseMercator(pt *GeoPoint, latO, longO, scale, fe, fn float64) *PolarCoord {

	var gc PolarCoord

	el := pt.El

	latOrad := degtorad(latO)
	longOrad := degtorad(longO)

	f := 1 - el.b/el.a
	esq := math.Sqrt(2.0*f - f*f)

	n := f / (2.0 - f)
	B := (el.a / (1 + n)) * (1 + n*n/4.0 + n*n*n*n/64.0)

	var SO float64

	if latOrad != 0.0 {

		h1 := n/2.0 - (2.0/3.0)*n*n + (5.0/16.0)*n*n*n + (41.0/180.0)*n*n*n*n
		h2 := (13.0/48.0)*n*n - (3.0/5.0)*n*n*n + (557.0/1440.0)*n*n*n*n
		h3 := (61.0/240.0)*n*n*n - (103.0/140.0)*n*n*n*n
		h4 := (49561.0 / 161280.0) * n * n * n * n

		QO := math.Asinh(math.Tan(latOrad)) - (esq * math.Atanh(esq*math.Sin(latOrad)))
		bO := math.Atan(math.Sinh(QO))
		xiO0 := bO // math.Asin(math.Sin(bO))

		xiO1 := h1 * math.Sin(2.0*xiO0)
		xiO2 := h2 * math.Sin(4.0*xiO0)
		xiO3 := h3 * math.Sin(6.0*xiO0)
		xiO4 := h4 * math.Sin(8.0*xiO0)

		xiO := xiO0 + xiO1 + xiO2 + xiO3 + xiO4

		SO = B * xiO
	}

	h1i := n/2.0 - (2.0/3.0)*n*n + (37.0/96.0)*n*n*n - (1.0/360.0)*n*n*n*n
	h2i := (1.0/48.0)*n*n + (1.0/15.0)*n*n*n - (437.0/1440.0)*n*n*n*n
	h3i := (17.0/480.0)*n*n*n - (37.0/840.0)*n*n*n*n
	h4i := (4397.0 / 161280.0) * n * n * n * n

	etai := (pt.X - fe) / (B * scale)
	xii := ((pt.Y - fn) + scale*SO) / (B * scale)

	xi1i := h1i * math.Sin(2*xii) * math.Cosh(2*etai)
	xi2i := h2i * math.Sin(4*xii) * math.Cosh(4*etai)
	xi3i := h3i * math.Sin(6*xii) * math.Cosh(6*etai)
	xi4i := h4i * math.Sin(8*xii) * math.Cosh(8*etai)

	eta1i := h1i * math.Cos(2*xii) * math.Sinh(2*etai)
	eta2i := h2i * math.Cos(4*xii) * math.Sinh(4*etai)
	eta3i := h3i * math.Cos(6*xii) * math.Sinh(6*etai)
	eta4i := h4i * math.Cos(8*xii) * math.Sinh(8*etai)

	xi0i := xii - (xi1i + xi2i + xi3i + xi4i)
	eta0i := etai - (eta1i + eta2i + eta3i + eta4i)

	bi := math.Asin(math.Sin(xi0i) / math.Cosh(eta0i))

	Qi := math.Asinh(math.Tan(bi))
	Qiiold := Qi + (esq * math.Atanh(esq*math.Tanh(Qi)))
	Qii := Qi + (esq * math.Atanh(esq*math.Tanh(Qiiold)))

	for math.Fabs(Qiiold-Qii) > 1e-12 {
		Qiiold = Qii
		Qii = Qi + (esq * math.Atanh(esq*math.Tanh(Qiiold)))
	}

	gc.Latitude = radtodeg(math.Atan(math.Sinh(Qii)))
	gc.Longitude = radtodeg(longOrad + math.Asin(math.Tanh(eta0i)/math.Cos(bi)))

	gc.El = el

	return &gc
}

// ## UTM coordinate functions for parsing and conversion

// A UTM coordinate defined by Northin, Easting and relative origin by Zone
// The reference ellipsoid is typically the GRS80Ellipsoid or the WGS84Ellipsoid
type UTMCoord struct {
	Northing, Easting float64
	Zone              string
	El                *Ellipsoid
}

// Canonical representation of an UTM coordinate
func (utm *UTMCoord) String() string {
	return fmt.Sprintf("%s %.0f %.0f", utm.Zone, utm.Easting, utm.Northing)
}

// This function parses a string UTM coordinate literal of the format
//
//	"ZONE EASTING NORTHING"
//
// Zone is the UTM meridian zone specifier and must be specified in the unambiguous
// way of zone number and latitude band. Easting and northing are specified as decimal meters.
// If the reference ellipsoid is nil, the DefaultEllipsoid is assumed.  
func AUTMToStruct(utmcoord string, el *Ellipsoid) (*UTMCoord, os.Error) {

	var zone, northing, easting string
	var north, east float64
	var err os.Error

	compact := strings.TrimSpace(utmcoord)

L1:
	for i, index := 0, 0; i < 3; i++ {
		index = strings.Index(compact, " ")
		if index == -1 {
			index = len(compact)
		}
		switch i {
		case 0:
			zone = compact[:index]
		case 1:
			easting = compact[:index]
		case 2:
			northing = compact[:index]
			break L1
		}
		compact = compact[index+len(" "):]
		compact = strings.TrimLeft(compact, " ")
	}

	north, err = strconv.Atof64(northing)

	if err != nil {
		return nil, err
	}

	east, err = strconv.Atof64(easting)

	if err != nil {
		return nil, err
	}

	if el == nil {
		el = DefaultEllipsoid
	}
	return &UTMCoord{Northing: north, Easting: east, Zone: zone, El: el}, nil
}

// Convert from UTM 2D projection to 3D polar. If the UTM coordinates do not contain a
// reference ellipsoid, the WGS84Ellipsoid is assumed and copied to the resulting polar coordinates.
//
// Inspired by http://www.gpsy.com/gpsinfo/geotoutm/gantz/LatLong-UTMconversion.cpp.txt
func UTMToLatLong(coord *UTMCoord) (*PolarCoord, os.Error) {

	zonelength := len(coord.Zone)
	utmLetter := int8(coord.Zone[zonelength-1:][0])
	zonenumber, err := strconv.Atoui(coord.Zone[:zonelength-1])

	if err != nil {
		return nil, err
	}

	pt := &GeoPoint{Y: coord.Northing, X: coord.Easting, El: coord.El}

	if utmLetter-'N' < 0 {
		pt.Y -= 10000000.0
	}

	if pt.El == nil {
		pt.El = DefaultEllipsoid
	}

	gc := InverseTransverseMercator(pt, 0, (float64(zonenumber)-1)*6-180+3, 0.9996, 500000, 0)

	return gc, nil
}

// Convert from 3D polar to UTM 2D projection. If the polar coordinates do not contain a
// reference ellipsoid, the WGS84Ellipsoid is assumed and copied to the resulting UTM coordinates.
//
// Inspired by http://www.gpsy.com/gpsinfo/geotoutm/gantz/LatLong-UTMconversion.cpp.txt
func LatLongToUTM(gcin *PolarCoord) *UTMCoord {

	var utm UTMCoord

	gc := *gcin // Make a copy as we might set the ellipsoid and we will not alter the input values
	zonenumber := uint((gc.Longitude+180)/6) + 1

	if gc.Latitude >= 56.0 && gc.Latitude < 64.0 && gc.Longitude >= 3.0 && gc.Longitude < 12.0 {
		zonenumber = 32
	}

	if gc.Latitude >= 72.0 && gc.Latitude < 84.0 {
		switch {
		case gc.Longitude >= 0.0 && gc.Longitude < 9.0:
			zonenumber = 31
		case gc.Longitude >= 9.0 && gc.Longitude < 21.0:
			zonenumber = 33
		case gc.Longitude >= 21.0 && gc.Longitude < 33.0:
			zonenumber = 35
		case gc.Longitude >= 33.0 && gc.Longitude < 42.0:
			zonenumber = 37
		}
	}

	if gc.El == nil {
		gc.El = DefaultEllipsoid
	}

	pt := DirectTransverseMercator(&gc, 0, (float64(zonenumber)-1)*6-180+3, 0.9996, 500000, 0)

	utm.Zone = strconv.Uitoa(zonenumber) + string(utmLetterDesignator(gc.Latitude))
	utm.Northing = pt.Y
	utm.Easting = pt.X

	if gc.Latitude < 0 {
		utm.Northing += 10000000
	}

	utm.El = pt.El

	return &utm
}

// This routine determines the correct UTM letter designator for the given latitude
// returns 'Z' if latitude is outside the UTM limits of 84N to 80S
func utmLetterDesignator(Lat float64) (LetterDesignator byte) {

	switch {
	case 84 >= Lat && Lat >= 72:
		LetterDesignator = 'X'
	case 72 > Lat && Lat >= 64:
		LetterDesignator = 'W'
	case 64 > Lat && Lat >= 56:
		LetterDesignator = 'V'
	case 56 > Lat && Lat >= 48:
		LetterDesignator = 'U'
	case 48 > Lat && Lat >= 40:
		LetterDesignator = 'T'
	case 40 > Lat && Lat >= 32:
		LetterDesignator = 'S'
	case 32 > Lat && Lat >= 24:
		LetterDesignator = 'R'
	case 24 > Lat && Lat >= 16:
		LetterDesignator = 'Q'
	case 16 > Lat && Lat >= 8:
		LetterDesignator = 'P'
	case 8 > Lat && Lat >= 0:
		LetterDesignator = 'N'
	case 0 > Lat && Lat >= -8:
		LetterDesignator = 'M'
	case -8 > Lat && Lat >= -16:
		LetterDesignator = 'L'
	case -16 > Lat && Lat >= -24:
		LetterDesignator = 'K'
	case -24 > Lat && Lat >= -32:
		LetterDesignator = 'J'
	case -32 > Lat && Lat >= -40:
		LetterDesignator = 'H'
	case -40 > Lat && Lat >= -48:
		LetterDesignator = 'G'
	case -48 > Lat && Lat >= -56:
		LetterDesignator = 'F'
	case -56 > Lat && Lat >= -64:
		LetterDesignator = 'E'
	case -64 > Lat && Lat >= -72:
		LetterDesignator = 'D'
	case -72 > Lat && Lat >= -80:
		LetterDesignator = 'C'
	default:
		LetterDesignator = 'Z' //  error flag to show that the Latitude is outside the UTM limits
	}

	return
}

// Base32 codeset for geohash as described in http://en.wikipedia.org/wiki/Geohash
const Base32GeohashCode = "0123456789bcdefghjkmnpqrstuvwxyz"

// The following functions deal with geohash encoding & decoding as described in http://en.wikipedia.org/wiki/Geohash
// Inspiration taken from
// - http://code.google.com/p/geospatialweb/source/browse/trunk/geohash/src/Geohash.java
// - http://blog.dixo.net/downloads/geohash-php-class/
// - https://github.com/kungfoo/geohash-java/blob/master/src/main/java/ch/hsr/geohash/GeoHash.java
func asciiindex(s string, ch byte) int16 {

	for index, val := range s {
		if byte(val) == ch {
			return int16(index)
		}
	}
	return -1
}

func abase32tobitset(sval, codeset string) (bitset string, err os.Error) {

	var bval int16
	sval = strings.TrimSpace(sval)

	for _, val := range sval {
		bval = asciiindex(codeset, byte(val))
		if bval == -1 {
			return "", os.EINVAL
		}
		bitset += fmt.Sprintf("%05b", bval)
	}
	return
}

// actually a general purpose round
func round(val float64, prec int) float64 {

	var rounder float64
	intermed := val * math.Pow(10, float64(prec))

	if val >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / math.Pow(10, float64(prec))
}

// 
func decodegeohashbitset(bitset string, floor, ceiling float64) float64 {
	mid := 0.0
	// calculate the mean error for improved rounding 
	err := (ceiling - floor) / 2.0

	for _, val := range bitset {
		mid = (floor + ceiling) / 2.0
		err /= 2.0
		switch val {
		case '1':
			floor = mid
		case '0':
			ceiling = mid
		}
	}
	// rounding according to 
	return round(mid, int(math.Fmax(1.0, -round(math.Log10(err), 0))-1))
}

// Return latitude & longitude from a geohash-encoded string.
// If the reference ellipsoid is nil, the default Ellipsoid will be returned.
// If the string is not a geohash, err will be set to EINVAL.
func GeoHashToLatLong(geohash string, el *Ellipsoid) (pc *PolarCoord, err os.Error) {

	var bitfield, lats, longs string

	// TODO: This uses a very inefficient implemenation of strings for bitfields

	// base32 to bitfield
	bitfield, err = abase32tobitset(geohash, Base32GeohashCode)
	if err != nil {
		return nil, err
	}

	// split the bitfield into latitude and longitude component 
	for index, bit := range bitfield {
		if index%2 == 0 {
			longs += string(bit)
		} else {
			lats += string(bit)
		}
	}

	if el == nil {
		el = DefaultEllipsoid
	}

	return &PolarCoord{
		Latitude:  decodegeohashbitset(lats, -90.0, 90.0),
		Longitude: decodegeohashbitset(longs, -180.0, 180.0),
		El:        el},
		nil
}

func geohashbitset(val, floor, ceiling float64, prec byte) string {
	var accu string
	var mid float64

	for i := byte(0); i < prec; i++ {
		mid = (floor + ceiling) / 2
		if val >= mid {
			accu += string('1')
			floor = mid
		} else {
			accu += string('0')
			ceiling = mid
		}
	}
	return accu

}

func f64toa(val float64, prec int) (sval string) {

	sval = fmt.Sprintf("%.*f", prec, val)
	n := len(sval)
	for n > 0 && sval[n-1] == '0' {
		n--
	}
	if n > 0 && sval[n-1] == '.' {
		n--
	}
	sval = sval[:n]
	return
}

func precision(val string) float64 {

	prec := 0

	if idx := strings.Index(val, "."); idx > 0 {
		prec = -(len(val) - idx - 1)
	}
	return math.Pow(10, float64(prec)) / 2.0
}

// Return a geohased representation of a latitude & longitude bearing point.
func LatLongToGeoHash(pc *PolarCoord) string {

	var err float64
	var bits byte

	// the suffix of a float64 may not be accurately representable as a string value due to
	// IEEE bit representation. We limit it to 6 places right of comma. Geohash-Encoding
	// doesn't support more precision either 
	plat := precision(f64toa(pc.Latitude, 6))
	latbits := byte(1)
	for err = 45.0; err > plat; err /= 2.0 {
		latbits++
	}

	plong := precision(f64toa(pc.Longitude, 6))
	longbits := byte(1)
	for err = 90.0; err > plong; err /= 2.0 {
		longbits++
	}

	if latbits > longbits {
		bits = latbits
	} else {
		bits = longbits
	}

	return LatLongToGeoHashBits(pc, bits)
}

// Return a geohased representation of a latitude & longitude bearing point using bits precision.
// If bits is 0 or larger than 30, it is set to a maximum of 30 bits
func LatLongToGeoHashBits(pc *PolarCoord, bits byte) string {

	if bits == 0 || bits > 30 {
		bits = 30
	}

	var accu, bitstring string
	latbits, longbits := bits, bits

	for addlong := byte(1); (longbits+latbits)%5 != 0; addlong = (addlong + 1) % 2 {
		longbits += addlong % 2
		latbits += (addlong + 1) % 2
	}

	latbitstring := geohashbitset(pc.Latitude, -90.0, 90.0, latbits)
	longbitstring := geohashbitset(pc.Longitude, -180.0, 180.0, longbits)

	for startlong := true; len(latbitstring)+len(longbitstring) > 0; startlong = !startlong {
		switch startlong {
		case true:
			bitstring += longbitstring[:1]
			longbitstring = longbitstring[1:]
		case false:
			bitstring += latbitstring[:1]
			latbitstring = latbitstring[1:]
		}
	}
	var base32fragment string
	var base32index uint64

	for ; len(bitstring) > 0; bitstring = bitstring[5:] {
		base32fragment = bitstring[:5]

		// An error is entirely the algorithms fault, so we ignore it here
		base32index, _ = strconv.Btoui64(base32fragment, 2)
		accu += string(Base32GeohashCode[uint8(base32index)])
	}
	return accu
}

// HELMERT Transformation - http://en.wikipedia.org/wiki/Helmert_transformation

// Returns a new ellipsoid by the given major axis a and major axis b in meters.
func NewEllipsoid(a, b float64) *Ellipsoid {
	return &Ellipsoid{a: a, b: b}
}

// ## Helmert transformation 

type transformer struct {
	dx, dy, dz, dM, drx, dry, drz float64
	datum                         string
}

// A generic Cartesian point to represent a 3D datum; used by the helmert-transformation
type Point3D struct {
	X, Y, Z float64
}

// A set of 3D datum transformations for the helmert transformation
var (
	// http://de.wikipedia.org/wiki/Datum_Austria
	HelmertWGS84ToMGI    = NewHelmertTransformer(-577.326, -90.129, -463.919, -2.4232, 5.1366, 1.4742, 5.2970, "WGS84toMGI")
	HelmertWGS84ToOSGB36 = NewHelmertTransformer(-446.448, 125.157, -542.060, 20.4894, -0.1502, -0.2470, -0.8421, "WGS84toOSGB36")
	// "Granit87" parameters
	HelmertLV03ToWGS84Granit87 = NewHelmertTransformer(660.077, 13.551, 369.3444, 5.66, 2.2356, 1.6047, 2.6451, "LV03toWGS84")
)

// Method to perform the Helmert transformation on a generic 3D datum and return a new datum.
// Instances of helmert transformations might be created by calls to NewHelmertTransformer
//
// Inspired by http://www.movable-type.co.uk/scripts/latlong-convert-coords.html
func (hp *transformer) Transform(ip *Point3D) *Point3D {

	var tp Point3D

	s := 1 + hp.dM/1e6

	rx := degtorad(hp.drx / 3600)
	ry := degtorad(hp.dry / 3600)
	rz := degtorad(hp.drz / 3600)

	tp.X = hp.dx + s*ip.X - rz*ip.Y + ry*ip.Z
	tp.Y = hp.dy + rx*ip.X + s*ip.Y - rx*ip.Z
	tp.Z = hp.dz - ry*ip.X + rx*ip.Y + s*ip.Z

	return &tp
}

// Method to perform the inverse helmert transformation on a generic 3D datum and return a new datum.
// The inverse might be easily computed by inverting all helmert parameters. The inversion is moderately
// accurate if the points of the orginal datum (x0, y0, z0) do not substantially diverge.   
// Instances of helmert transformations might be created by calls to NewHelmertTransformer
func (hp *transformer) InverseTransform(pt *Point3D) *Point3D {

	var ihp transformer

	ihp.dM = -1 * hp.dM
	ihp.drx = -1 * hp.drx
	ihp.dry = -1 * hp.dry
	ihp.drz = -1 * hp.drz
	ihp.dx = -1 * hp.dx
	ihp.dy = -1 * hp.dy
	ihp.dz = -1 * hp.dz

	return ihp.Transform(pt)
}

// Returns a canoncial representation of the helmert parameters
func (tp *transformer) String() string {
	return fmt.Sprintf("Helmert[%s](dx,dy,dz,dM,drx, dry,drz): (%f, %f, %f, %f, %f, %f, %f)", tp.datum, tp.dx, tp.dy, tp.dz, tp.dM, tp.drx, tp.dry, tp.drz)
}

// Get the well known text (WKT) for the helmert transformation as defined in
// http://www.geoapi.org/2.0/javadoc/org/opengis/referencing/doc-files/WKT.html
func (tp *transformer) WellKnownString() string {
	return fmt.Sprintf("TOWGS84[\"%f\", \"%f\", \"%f\", \"%f\", \"%f\", \"%f\", \"%f\"]", tp.dx, tp.dy, tp.dz, tp.dM, tp.drx, tp.dry, tp.drz)
}

// Create a new instance of helmert parameters.
//
//	dx, dy, dz: delta of coordinate origin in meters
//	drx, dry, drz: delta of coordinate bearing in rads
//	dM:  scale correction to be made to the position vector in the source coordinate reference
//		system, in parts of per million
//  
// Attention: Unlike the other functions dealing with bearings and coordinates in this package,
// the angular helmert parameters have to be specified in rad [-pi;pi]  
func NewHelmertTransformer(dx, dy, dz, dM, drx, dry, drz float64, datum string) *transformer {
	return &transformer{dx: dx, dy: dy, dz: dz, dM: dM, drx: drx, dry: dry, drz: drz, datum: datum}
}
