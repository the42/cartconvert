// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

//target:github.com/the42/cartconvert/bmn

// This package provides a series of functions to deal with 
// conversion and transformations of coordinates in the Datum Austria
//
// http://de.wikipedia.org/wiki/Datum_Austria
//
// and here specifically of the Bundesmeldenetz, the former federal cartographic datum
// of Austria. The Bundesmeldenetz is already widely replaced by UTM coordinates but much legacy
// data is still encoded in BMN coordinates. Unlike UTM, the BMN uses the Bessel reference ellipsoid
// and ues lat0 at Hierro (canary islands), which makes transformations tedious.
// For more information see
//
// [DE]: http://www.topsoft.at/pstrainer/entwicklung/algorithm/karto/oek/austria_oek.htm#bmn
// [EN]: http://www.asprs.org/resources/grids/03-2004-austria.pdf
package bmn

import (
	"github.com/the42/cartconvert"
	"strings"
	"strconv"
	"fmt"
	"os"
)

// Meridian Coordinates of the Bundesmeldenetz, three values describing false easting and false northing
// The meridian specification of BMN plays the same role as the zone specifier of UTM
type BMNMeridian int

const (
	BMNMunknown BMNMeridian = iota
	BMNM28
	BMNM31
	BMNM34
)

var bmnStrings = map[BMNMeridian]string{BMNM28: "M28", BMNM31: "M31", BMNM34: "M34"}

// A BMN coordinate is specified by right-value (Easting), height-value (Northing)
// and the meridian stripe, 28°, 31° or 34° West of Hierro 
type BMNCoord struct {
	Right, Height, RelHeight float64
	Meridian                 BMNMeridian
	el                       *cartconvert.Ellipsoid
}

// Canoncial representation of a BMN-value
func (bc *BMNCoord) String() (fs string) {

	fs = bmnStrings[bc.Meridian]
	var next float64

	for i := 0; i < 2; i++ {
		fs += " "
		switch i {
		case 0:
			next = bc.Right
		case 1:
			next = bc.Height
		}

		fs += fmt.Sprintf("%.2f", next)
		n := len(fs)
		for n > 0 && fs[n-1] == '0' {
			n--
		}
		if n > 0 && fs[n-1] == '.' {
			n--
		}
		fs = fs[:n]
	}
	return
}

// Parses a string representation of a BMN-Coordinate into a struct holding a BMN coordinate value
// The reference ellipsoid of BMN coordinates is always the Bessel ellipsoid
func BMNStringToStruct(bmncoord string) (*BMNCoord, os.Error) {

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

				return &BMNCoord{Right: right, Height: height, Meridian: meridian, el: cartconvert.BesselEllipsoid}, nil
			}
		}
	}

	return nil, err
}

// Transform a BMN coordinate value to a WGS84 based latitude and longitude coordinate. Function returns
// nil, if the meridian stripe of the bmn-coordinate is not set
func BMNToWGS84LatLong(bmncoord *BMNCoord) *cartconvert.PolarCoord {

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
		return nil
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

	return cartconvert.CartesianToPolar(&cartconvert.CartPoint{X: pt.X, Y: pt.Y, Z: pt.Z, El: cartconvert.WGS84Ellipsoid})
}

// Transform a latitute / longitude coordinate datum into a BMN coordinate. Function returns
// nil, if the meridian stripe of the bmn-coordinate is not set
//
// Important: The reference ellipsoid of the originating coordinate system will be assumed
// to be the WGS84Ellipsoid and will be set thereupon, regardless of the actually set reference ellipsoid
func WGS84LatLongToBMN(gc *cartconvert.PolarCoord, meridian BMNMeridian) *BMNCoord {

	var long0, fe float64

	// This sets the Ellipsoid to WGS84, regardless of the actual value set
	gc.El = cartconvert.WGS84Ellipsoid

	cart := cartconvert.PolarToCartesian(gc)
	pt := cartconvert.HelmertWGS84ToMGI.Transform(&cartconvert.Point3D{X: cart.X, Y: cart.Y, Z: cart.Z})
	polar := cartconvert.CartesianToPolar(&cartconvert.CartPoint{X: pt.X, Y: pt.Y, Z: pt.Z, El: cartconvert.BesselEllipsoid})

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
		return nil
	}

	gp := cartconvert.DirectTransverseMercator(
		polar,
		0,
		long0,
		1,
		fe,
		-5000000)

	return &BMNCoord{Meridian: meridian, Height: gp.Y, Right: gp.X, el: gp.El}
}

func NewBMNCoord(Meridian BMNMeridian, Right, Height, RelHeight float64) *BMNCoord {
	return &BMNCoord{Right: Right, Height: Height, RelHeight: RelHeight, Meridian: Meridian, el: cartconvert.BesselEllipsoid}
}
