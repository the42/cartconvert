// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// Automated tests for the cartconvert/osgb36 package
package osgb36

import (
	"fmt"
	"github.com/the42/cartconvert"
	"testing"
)

// ## BMNStringToStruct
type oSGB36StringToStructTest struct {
	in  string
	out *OSGB36Coord
}

var oSGB36StringToStructTestssuc = []oSGB36StringToStructTest{
	{
		"NN166712",
		&OSGB36Coord{Zone: "NN", Right: 166, Height: 712},
	},
}

func osgb36equal(osgb1, osgb2 *OSGB36Coord) bool {
	p1 := fmt.Sprintf("%s", osgb1)
	p2 := fmt.Sprintf("%s", osgb2)

	return p1 == p2
}

func TestOSGB36StringToStruct(t *testing.T) {
	for _, test := range oSGB36StringToStructTestssuc {
		out, err := AOSGB36ToStruct(test.in)

		if err != nil {
			t.Error(err)
		}

		if !osgb36equal(test.out, out) {
			t.Error("TestOSGB36StringToStruct")
		}
	}
}


// ## BMNToWGS84LatLong
type oSGB36ToWGS84LatLongTest struct {
	in  *OSGB36Coord
	out *cartconvert.PolarCoord
}


var oSGB36ToWGS84LatLongTests = []oSGB36ToWGS84LatLongTest{
	{
		NewOSGB36Coord("ST", 58982, 72915, 0),
		&cartconvert.PolarCoord{Latitude: 50.815243, Longitude: 0.137062},
	},
}

func latlongequal(pcp1, pcp2 *cartconvert.PolarCoord) bool {
	pp1s := fmt.Sprintf("%.5g %.5g", pcp1.Latitude, pcp1.Longitude)
	pp2s := fmt.Sprintf("%.5g %.5g", pcp2.Latitude, pcp2.Longitude)

	return pp1s == pp2s
}

func TestOSGB36ToWGS84LatLong(t *testing.T) {
	for cnt, test := range oSGB36ToWGS84LatLongTests {

		out, err := OSGB36ToWGS84LatLong(test.in)

		if err != nil {
			t.Errorf("OSGB36ToWGS84LatLong [%d]: Error: %s", cnt, err)
		} else {

			if !latlongequal(test.out, out) {
				t.Errorf("OSGB36ToWGS84LatLong [%d]: Expected: %s, got: %s", cnt, test.out, out)
			}
		}
	}
}

/*
// ## WGS84LatLongToBMN
type wGS84LatLongToBMNParam struct {
	gc       *cartconvert.PolarCoord
	meridian BMNMeridian
}

type wGS84LatLongToBMNTest struct {
	in  wGS84LatLongToBMNParam
	out *BMNCoord
}

var wGS84LatLongToBMNTests = []wGS84LatLongToBMNTest{
	{
		wGS84LatLongToBMNParam{
			gc:       &cartconvert.PolarCoord{Latitude: 47.570299, Longitude: 14.236188, El: cartconvert.WGS84Ellipsoid},
			meridian: BMNM34},
		NewBMNCoord(BMNM34, 592269, 272290.05, 0),
	},
	{
		wGS84LatLongToBMNParam{
			gc:       &cartconvert.PolarCoord{Latitude: 48.507001, Longitude: 15.698748, El: cartconvert.WGS84Ellipsoid},
			meridian: BMNZoneDet},
		NewBMNCoord(BMNM34, 703168, 374510, 0),
	},
}

func TestWGS84LatLongToBMN(t *testing.T) {
	for index, test := range wGS84LatLongToBMNTests {
		out, _ := WGS84LatLongToBMN(test.in.gc, test.in.meridian)
		if !bmnequal(test.out, out) {
			t.Errorf("WGS84LatLongToBMN [%d]: expected %s, got %s", index, test.out, out)
		}
	}
}
*/
