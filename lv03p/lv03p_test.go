// Copyright 2011 Johann Höchtl. All rights reserved.
// Use of this source code is governed by a Modified BSD License
// that can be found in the LICENSE file.

// Automated tests for the cartconvert/lv03p package
package lv03p

import (
	"os"
	// "fmt"
	//"github.com/the42/cartconvert"
	"testing"
)

type swissCoordRepresentation struct {
	in  SwissCoord
	out string
}

var swissCoordRepresentationTest = []swissCoordRepresentation{
	{
		SwissCoord{Easting: 235.5, Northing: 20.0, CoordType: LV03}, "y:235.5 x:20",
	},
}

func TestSwissCoordRepresentation(t *testing.T) {
	for cnt, test := range swissCoordRepresentationTest {
		out := test.in.String()

		if test.out != out {
			t.Errorf("TestSwissCoordRepresentation [%d]: Expected: %s, got: %s", cnt, test.out, out)
		}
	}
}

// ## ASwissCoordToStruct
type aSwissCoordToStructretparam struct {
	coord *SwissCoord
	err   os.Error
}

func (val *aSwissCoordToStructretparam) String() (fs string) {

	if val.coord != nil {
		fs = val.coord.String()
	}

	if val.err != nil {
		fs += " " + val.err.String()
	}
	return
}

type aSwissCoordToStruct struct {
	in  string
	out aSwissCoordToStructretparam
}

var aSwissCoordToStructTests = []aSwissCoordToStruct{
	{
		in: "x:25 y:34.3", out: aSwissCoordToStructretparam{coord: &SwissCoord{Easting: 34.3, Northing: 25, CoordType: LV03}, err: nil},
	},
	{
		in: "x:25.0 N:34.3", out: aSwissCoordToStructretparam{coord: nil, err: os.EINVAL},
	},
}

func aswisscoordtostructequal(coord1, coord2 aSwissCoordToStructretparam) bool {
	if coord1.coord != nil && coord2.coord != nil {
		return coord1.coord.String() == coord2.coord.String()
	}
	return coord1.err == coord2.err
}

func TestASwissCoordToStruct(t *testing.T) {
	for cnt, test := range aSwissCoordToStructTests {

		out, erro := ASwissCoordToStruct(test.in)
		retval := aSwissCoordToStructretparam{coord: out, err: erro}

		if !aswisscoordtostructequal(test.out, retval) {
			t.Errorf("ASwissCoordToStruct [%d]: expected %v, got %v", cnt, test.out, retval)
		}
	}
}

/*
// ## BMNToWGS84LatLong
type bMNToWGS84LatLongTest struct {
	in  *BMNCoord
	out *cartconvert.PolarCoord
}

func bMNStringToStructHelper(coord string) (bmncoord *BMNCoord) {
	bmncoord, _ = ABMNToStruct(coord)
	return
}

var bMNToWGS84LatLongTests = []bMNToWGS84LatLongTest{
	{
		NewBMNCoord(BMNM28, 592270.0, 272290, 0),
		&cartconvert.PolarCoord{Latitude: 47.439212, Longitude: 16.197434},
	},
	{ // TODO: Ist das möglich??
		bMNStringToStructHelper("M34 592269 272290"),
		&cartconvert.PolarCoord{Latitude: 47.570299, Longitude: 14.236188},
	},
	{
		bMNStringToStructHelper("M34 703168 374510"),
		&cartconvert.PolarCoord{Latitude: 48.507001, Longitude: 15.698748},
	},
}

func latlongequal(pcp1, pcp2 *cartconvert.PolarCoord) bool {
	pp1s := fmt.Sprintf("%.5g %.5g", pcp1.Latitude, pcp1.Longitude)
	pp2s := fmt.Sprintf("%.5g %.5g", pcp2.Latitude, pcp2.Longitude)

	return pp1s == pp2s
}

func TestBMNToWGS84LatLong(t *testing.T) {
	for _, test := range bMNToWGS84LatLongTests {

		out, _ := BMNToWGS84LatLong(test.in)

		if !latlongequal(test.out, out) {
			t.Error("BMNToWGS84LatLong")
		}
	}
}

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
