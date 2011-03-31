// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// This package provides functions to deal with 
// conversion and transformations of coordinates in OSGB36,
// the UK National Grid
//
// The conversion between WGS84 and OSGB36 uses a simple helmert transformation,
// which in the case of OSGB36 inconsistencies may result in an accuracy not exceeding +/- 5m.
// If higher accuracy is required, a set of helmert parameters must be used or the
// procedure described in http://www.ordnancesurvey.co.uk/gps/docs/Geomatics_world.pdf.
//
// cf. http://gps.ordnancesurvey.co.uk/etrs89geo_natgrid.asp
package osgb36

import (
	"github.com/the42/cartconvert"
	"strings"
	"strconv"
	"fmt"
	"os"
)


// A OSGB36 coordinate is specified by right-value (easting), height-value (northing)
// and zone. 
type OSGB36Coord struct {
	Right, Height int
	RelHeight     float64
	Zone          string
	el            *cartconvert.Ellipsoid
}

// Canonical representation of a OSGB36 datum
func (bc *OSGB36Coord) String() string {

	return fmt.Sprintf("%s%d%d", bc.Zone, bc.Right, bc.Height)
}

// Parses a string representation of a BMN-Coordinate into a struct holding a BMN coordinate value.
// The reference ellipsoid of BMN coordinates is always the Bessel ellipsoid.
func AOSGB36ToStruct(osgb36coord string) (*OSGB36Coord, os.Error) {

	compact := strings.ToUpper(strings.TrimSpace(osgb36coord))
	var rights, heights, rnh string
	var zone string
	var right, height int
	var err os.Error

L1:
	for _, item := range compact {
		switch {
		case item == ' ':
			continue L1
		case byte(item)-'0' >= 0 && byte(item)-'0' <= 9:
			rnh += string(item)
		default:
			zone += string(item)
		}

	}

	zl := len(zone)
	if zl == 0 || zl > 2 {
		return nil, os.EINVAL
	}

	rnhlen := len(rnh)
	if rnhlen%2 > 0 {
		return nil, os.EINVAL
	}

	rights, heights = rnh[:rnhlen/2], rnh[rnhlen/2:]

	right, err = strconv.Atoi(rights)
	if err == nil {

		height, err = strconv.Atoi(heights)
		if err == nil {

			return &OSGB36Coord{Right: right, Height: height, Zone: zone, el: cartconvert.Airy1830Ellipsoid}, nil
		}
	}
	return nil, err
}
/*
// Transform a BMN coordinate value to a WGS84 based latitude and longitude coordinate. Function returns
// nil, if the meridian stripe of the bmn-coordinate is not set
func BMNToWGS84LatLong(bmncoord *BMNCoord) (*cartconvert.PolarCoord, os.Error) {

	var long0, fe float64

	switch bmncoord.Meridian {
	case BMNM28:
		long0 = 10.0 + 20.0/60.0
		fe = 150000
	case BMNM31:
		long0 = 13.0 + 20.0/60.0
		fe = 450000
	case BMNM34:
		long0 = 16.0 + 20.0/60.0
		fe = 750000
	default:
		return nil, os.EINVAL
	}

	gc := cartconvert.InverseTransverseMercator(
		&cartconvert.GeoPoint{Y: bmncoord.Height, X: bmncoord.Right, El: bmncoord.el},
		0,
		long0,
		1,
		fe,
		-5000000)

	cart := cartconvert.PolarToCartesian(gc)
	pt := cartconvert.HelmertWGS84ToMGI.InverseTransform(&cartconvert.Point3D{X: cart.X, Y: cart.Y, Z: cart.Z})

	return cartconvert.CartesianToPolar(&cartconvert.CartPoint{X: pt.X, Y: pt.Y, Z: pt.Z, El: cartconvert.WGS84Ellipsoid}), nil
}

// Transform a latitude / longitude coordinate datum into a BMN coordinate. Function returns
// nil, if the meridian stripe of the bmn-coordinate is not set.
//
// Important: The reference ellipsoid of the originating coordinate system will be assumed
// to be the WGS84Ellipsoid and will be set thereupon, regardless of the actually set reference ellipsoid.
func WGS84LatLongToBMN(gc *cartconvert.PolarCoord, meridian BMNMeridian) (*BMNCoord, os.Error) {

	var long0, fe float64

	// This sets the Ellipsoid to WGS84, regardless of the actual value set
	gc.El = cartconvert.WGS84Ellipsoid

	cart := cartconvert.PolarToCartesian(gc)
	pt := cartconvert.HelmertWGS84ToMGI.Transform(&cartconvert.Point3D{X: cart.X, Y: cart.Y, Z: cart.Z})
	polar := cartconvert.CartesianToPolar(&cartconvert.CartPoint{X: pt.X, Y: pt.Y, Z: pt.Z, El: cartconvert.BesselEllipsoid})

	// Determine meridian stripe based on longitude
	if meridian == BMNZoneDet {
		switch {
		case 11.0+0.5/6*10 >= polar.Longitude && polar.Longitude >= 8.0+0.5/6*10:
			meridian = BMNM28
		case 14.0+0.5/6*10 >= polar.Longitude && polar.Longitude >= 11.0+0.5/6*10:
			meridian = BMNM31
		case 17.0+0.5/6*10 >= polar.Longitude && polar.Longitude >= 14.0+0.5/6*10:
			meridian = BMNM34
		}
	}

	switch meridian {
	case BMNM28:
		long0 = 10.0 + 20.0/60.0
		fe = 150000
	case BMNM31:
		long0 = 13.0 + 20.0/60.0
		fe = 450000
	case BMNM34:
		long0 = 16.0 + 20.0/60.0
		fe = 750000
	default:
		return nil, os.EINVAL
	}

	gp := cartconvert.DirectTransverseMercator(
		polar,
		0,
		long0,
		1,
		fe,
		-5000000)

	return &BMNCoord{Meridian: meridian, Height: gp.Y, Right: gp.X, el: gp.El}, nil
}

func NewBMNCoord(Meridian BMNMeridian, Right, Height, RelHeight float64) *BMNCoord {
	return &BMNCoord{Right: Right, Height: Height, RelHeight: RelHeight, Meridian: Meridian, el: cartconvert.BesselEllipsoid}
}
*/
