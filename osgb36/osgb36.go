// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// This package provides functions to deal with 
// conversion and transformations of coordinates in OSGB36,
// the UK National Grid.
//
// The conversion between WGS84 and OSGB36 uses a simple helmert transformation,
// which in the case of OSGB36 inconsistencies may result in an accuracy not exceeding +/- 5m.
// If higher accuracy is required, a set of helmert parameters must be used or the
// procedure described in http://www.ordnancesurvey.co.uk/gps/docs/Geomatics_world.pdf.
//
// For further info see http://gps.ordnancesurvey.co.uk/etrs89geo_natgrid.asp
package osgb36

import (
	"github.com/the42/cartconvert"
	"strings"
	"strconv"
	"fmt"
	"os"
	"math"
)

// A OSGB36 coordinate is specified by right-value (easting), height-value (northing)
// and zone. 
type OSGB36Coord struct {
	Easting, Northing uint
	RelHeight         float64
	Zone              string
	el                *cartconvert.Ellipsoid
	GridLen           OSGB36prec
}

// The precision of an OSGB36 coordinate can either be set explicitely from meter - resolutiuon (OSGB36_5)
// to the bare 100x100 km Zone OSGB36Min.
type OSGB36prec byte

const (
	OSGB36Min OSGB36prec = iota
	OSGB36_1
	OSGB36_2
	OSGB36_3
	OSGB36_4
	OSGB36_5
	OSGB36_Max  = OSGB36_5
	OSGB36Leave = 128
	OSGB36Auto  = OSGB36Leave + 1
)

// Canonical representation of a OSGB36 datum
func (coord *OSGB36Coord) String() string {
	if coord.GridLen > 0 {
		return fmt.Sprintf("%s%d%d", coord.Zone, coord.Easting, coord.Northing)
	}
	return coord.Zone
}

// Parses a string representation of an OSGB36 coordinate datum into a OSGB36 coordinate struct. The literal
// can be specified as follows:
//    ZO EA NO
//    ZO EANO
//    ZOEANO
// with ZO the two letter zone specifier, EA easting and NO northing to the accuracy of a meter.
//
// prec has the meaning:
//    OSGB36Leave - pass the northing and easting as is to the internal struct
//    OSGB36Auto - shorten trailing zeros. NN1665034570 becomes NN16653457 and NN16001700 becomes NN1617.
//    OSGB36_1 ... _5 - shorten northing and easting to the specified length, regardless of loss of precision
// The reference ellipsoid of an OSGB36 coordinate is always the Airy1830 ellipsoid.
func AOSGB36ToStruct(osgb36coord string, prec OSGB36prec) (*OSGB36Coord, os.Error) {

	compact := strings.ToUpper(strings.TrimSpace(osgb36coord))
	var zone, enn string
	var east, north int
	var err os.Error

L1:
	for _, item := range compact {
		switch {
		case item == ' ':
			continue L1
		case byte(item)-'0' >= 0 && byte(item)-'0' <= 9:
			enn += string(item)
		default:
			zone += string(item)
		}
	}

	zl := len(zone)
	if zl == 0 || zl > 2 {
		return nil, os.EINVAL
	}

	ennlen := OSGB36prec(len(enn))
	if ennlen > 0 {
		if ennlen%2 > 0 {
			return nil, os.EINVAL
		}

		ennlen /= 2

		east, err = strconv.Atoi(enn[:ennlen])
		if err != nil {
			return nil, err
		}
		north, err = strconv.Atoi(enn[ennlen:])
		if err != nil {
			return nil, err
		}
	}

	return NewOSGB36Coord(zone, uint(east), uint(north), prec, 0)
}

// Returns northing and easting based on OSGB36 zone specifier relative to false northing and easting
func OSGB36ZoneToRefCoords(coord *OSGB36Coord) (easting, northing uint) {

	// get numeric values of letter references, mapping A->0, B->1, C->2, etc:
	l1 := uint(coord.Zone[0] - 'A')
	l2 := uint(coord.Zone[1] - 'A')

	// shuffle down letters after 'I' since 'I' is not used in grid:
	if l1 > 7 {
		l1--
	}
	if l2 > 7 {
		l2--
	}

	// convert grid letters into 100km-square indexes from false origin (grid square SV):
	easting = (((l1-2)%5)*5 + l2%5) * 100000
	northing = ((19 - l1/5*5) - l2/5) * 100000

	// append numeric part of references to grid index:
	fact := uint(math.Pow(10, float64(OSGB36_Max-coord.GridLen)))
	easting += coord.Easting*fact + 5*(fact/10)
	northing += coord.Northing*fact + 5*(fact/10)

	return
}

// Convert an OSGB36 coordinate value to a WGS84 based latitude and longitude coordinate. Important: A OSGB36 datum like
// NN1745 will be internally expanded to NN1750045500 to point to the middle of the zone. For the point at
// NN1700045000 it is necessary to fully qualify northing and easting.
func OSGB36ToWGS84LatLong(coord *OSGB36Coord) (*cartconvert.PolarCoord, os.Error) {

	easting, northing := OSGB36ZoneToRefCoords(coord)

	gc := cartconvert.InverseTransverseMercator(
		&cartconvert.GeoPoint{Y: float64(northing), X: float64(easting), El: coord.el},
		49,
		-2,
		0.9996012717,
		400000,
		-100000)

	cart := cartconvert.PolarToCartesian(gc)
	pt := cartconvert.HelmertWGS84ToOSGB36.InverseTransform(&cartconvert.Point3D{X: cart.X, Y: cart.Y, Z: cart.Z})

	return cartconvert.CartesianToPolar(&cartconvert.CartPoint{X: pt.X, Y: pt.Y, Z: pt.Z, El: cartconvert.WGS84Ellipsoid}), nil
}

func SanitizeOSGB36CoordToPrec(easting, northing *uint, prec *OSGB36prec) os.Error {

	if *easting+*northing == 0 {
		*prec = 0
		return nil
	}

	inputpreclen := byte(uintlen(*easting))
	if inputpreclen != byte(uintlen(*northing)) {
		return os.EINVAL
	}

	desiredprec := *prec

	switch *prec {
	case OSGB36Auto:
		northprec, eastprec := inputpreclen, inputpreclen
		easttmp, northtmp := *easting, *northing

		for easttmp%10 == 0 && eastprec > 0 {
			easttmp /= 10
			eastprec--
		}

		if eastprec < inputpreclen {
			for northtmp%10 == 0 && northprec > 0 {
				northtmp /= 10
				northprec--
			}
			desiredprec = OSGB36prec(max(int(northprec), int(eastprec)))
		} else {
			break
		}
		fallthrough
	case OSGB36_1, OSGB36_2, OSGB36_3, OSGB36_4, OSGB36_5:
		if byte(desiredprec) < inputpreclen {
			fact := uint(math.Pow(10, float64(inputpreclen-byte(desiredprec))))
			*easting /= fact
			*northing /= fact

		} else {
			fact := uint(math.Pow(10, float64(byte(desiredprec)-inputpreclen)))
			*easting *= fact
			*northing *= fact

		}
		*prec = desiredprec
	case OSGB36Leave:
		*prec = OSGB36prec(inputpreclen)
	}
	return nil
}

// Build OSGB36 coordinate from easting and northing relative to Grid. Reduce precision to prec positions.
func GridRefNumToLet(easting, northing uint, height float64, prec OSGB36prec) (*OSGB36Coord, os.Error) {
	// get the 100km-grid indices
	easting100k := easting / 100000
	northing100k := northing / 100000

	if easting100k < 0 || easting100k > 6 || northing100k < 0 || northing100k > 12 {
		return nil, os.EINVAL
	}

	// translate those into numeric equivalents of the grid letters
	l1 := byte((19 - northing100k) - (19-northing100k)%5 + (easting100k+10)/5)
	l2 := byte((19-northing100k)*5%25 + easting100k%5)

	// compensate for skipped 'I' and build grid letter-pairs
	if l1 > 7 {
		l1++
	}
	if l2 > 7 {
		l2++
	}

	zone := string(l1+'A') + string(l2+'A')
	easting %= 100000
	northing %= 100000

	return NewOSGB36Coord(zone, easting, northing, prec, 0)
}

// Transform a latitude / longitude coordinate datum into a OSGB36 coordinate.
//
// Important: The reference ellipsoid of the originating coordinate system will be assumed
// to be the WGS84Ellipsoid and will be set thereupon, regardless of the actually set reference ellipsoid.
func WGS84LatLongToOSGB36(gc *cartconvert.PolarCoord) (*OSGB36Coord, os.Error) {
	// This sets the Ellipsoid to WGS84, regardless of the actual value set
	gc.El = cartconvert.WGS84Ellipsoid

	cart := cartconvert.PolarToCartesian(gc)
	pt := cartconvert.HelmertWGS84ToOSGB36.Transform(&cartconvert.Point3D{X: cart.X, Y: cart.Y, Z: cart.Z})
	polar := cartconvert.CartesianToPolar(&cartconvert.CartPoint{X: pt.X, Y: pt.Y, Z: pt.Z, El: cartconvert.Airy1830Ellipsoid})

	gp := cartconvert.DirectTransverseMercator(
		polar,
		49,
		-2,
		0.9996012717,
		400000,
		-100000)

	return GridRefNumToLet(uint(gp.X+0.5), uint(gp.Y+0.5), 0, OSGB36Auto)

}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func uintlen(x uint) (len uint) {
	for ; x > 0; x /= 10 {
		len++
	}
	return
}

func NewOSGB36Coord(Zone string, easting, northing uint, prec OSGB36prec, relheight float64) (*OSGB36Coord, os.Error) {
	if err := SanitizeOSGB36CoordToPrec(&easting, &northing, &prec); err != nil {
		return nil, err
	}
	return &OSGB36Coord{Easting: easting, Northing: northing, RelHeight: relheight, Zone: Zone, GridLen: prec, el: cartconvert.Airy1830Ellipsoid}, nil
}
