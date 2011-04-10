// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// This package provides a series of functions to deal with 
// conversion and transformations of coordinates in the Swiss coordinate system.
//
// The Swiss coordinate system recently switched from lv03 to lv95, however the difference is
// generall below one meter and ignored by this package. For accuracy within 1cm the
//  FINELTRA-Transformation has to be applied.
//
// References:
//
// [DE]: http://www.swisstopo.admin.ch/internet/swisstopo/de/home/topics/survey/sys/refsys/switzerland.parsysrelated1.24280.downloadList.32633.DownloadFile.tmp/refsysd.pdf
package lv03p

import (
	"github.com/the42/cartconvert"
	"strings"
	"strconv"
	"fmt"
	"os"
)

// Coordinate type of Switzerland. Only affects string representation but not accuracy (The two systems diverge by about 1m)
type SwissCoordType byte

const (
	LV03 SwissCoordType = iota
	LV95
)

// A coordinate in Switzerland is specified by easting (right-value, x), and northing (height-value, y)
type SwissCoord struct {
	Easting, Northing, RelHeight float64
	CoordType                 SwissCoordType
	el                       *cartconvert.Ellipsoid
}

var coordliterals = [][]string{{"y:", "x:"}, {"E:", "N:"}}

// Canonical representation of a SwissCoord-value
func (bc *SwissCoord) String() (fs string) {

	var next float64

	for i := 0; i < 2; i++ {
		fs += coordliterals[bc.CoordType][i]
		switch i {
		case 0:
			next = bc.Easting
		case 1:
			next = bc.Northing
		}

		tmp := fmt.Sprintf("%.0f", next)
		n := len(tmp)
		for n > 0 && tmp[n-1] == '0' {
			n--
		}
		if n > 0 && tmp[n-1] == '.' {
			n--
		}
		fs = fs + tmp[:n]
	}
	return
}

// Parses a string representation of a LV++ coordinate into a struct holding a SwissCoord coordinate value.
// The reference ellipsoid of Swisscoord datum is always the GRS80 ellipsoid.
func ASwissCoordToStruct(coord string) (*SwissCoord, os.Error) {

	compact := strings.ToUpper(strings.TrimSpace(bmncoord))
	var rights, heights string
	var meridian BMNMeridian
	var right, height float64
	var err os.Error

L1:
	for i, index := 0, 0; i < 3; i++ {
		index = strings.Index(compact, " ")
		if index == -1 {
			index = len(compact)
		}
		switch i {
		case 0:
			switch compact[:index] {
			case "M28":
				meridian = BMNM28
			case "M31":
				meridian = BMNM31
			case "M34":
				meridian = BMNM34
			default:
				err = os.EINVAL
				break L1
			}
		case 1:
			rights = compact[:index]
		case 2:
			heights = compact[:index]
			break L1
		}
		compact = compact[index+len(" "):]
		compact = strings.TrimLeft(compact, " ")
	}

	if err == nil {

		right, err = strconv.Atof64(rights)
		if err == nil {

			height, err = strconv.Atof64(heights)
			if err == nil {

				return &BMNCoord{Easting: right, Northing: height, CoordType: ct, el: cartconvert.GRS80Ellipsoid}, nil
			}
		}
	}

	return nil, err
}

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
	polar := cartconvert.CartesianToPolar(&cartconvert.CartPoint{X: pt.X, Y: pt.Y, Z: pt.Z, El: cartconvert.Bessel1841MGIEllipsoid})

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
	return &BMNCoord{Right: Right, Height: Height, RelHeight: RelHeight, Meridian: Meridian, el: cartconvert.Bessel1841MGIEllipsoid}
}
