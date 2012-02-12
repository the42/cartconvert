cartconvserv - A RESTFul service for coordinate transformations
===============================================================

Purpose
-------

cartconvserv is a RESTFul service to the cartconv - package. It may be installed
as a stand-alone service or executed in the [Appengine](http://code.google.com/appengine/docs/go/) environment.

Functionality
-------------

* Conversion between coordinate systems and bearing representations.
* Serialization as XML or JSON by content negotiation.

Convention for this help:

* URL: Base url to service
* APIRoot: Root of the RESTFul API (configurable)

UTM - Conversions
-----------------

Base Url for UTM operations:
   
    URL/APIRoot/utm/<VALUE>[.xml|.json]?outputformat=<utm|geohash|latlongdeg|latlongcomma|bmn>

Value is a coordinate in UTM representation. The reference ellipsoid is always
the WGS84Ellipsoid.

Examples of valid input values:

    17T 630084 4833438
    17T 630084.31 4833438.54

If the extension to value is empty or ".json", the result of the requested
output format is JSON-encoded.

If the extension to value is ".xml", the result of the requested output format
is XML-encoded.

The value to the parameter "outputformat" is either

* latlongdeg: Latitude and longitude with fractions in degrees
* latlongcomma: Latitude and longitude with decimal fractions
* geohash: Geohash-encoded value of latitude and longitude
* bmn: Serialization of the value as BMN-coordinate (see below)

### Examples
#### Output requested as latitude and longitude in [arc degrees](http://en.wikipedia.org/wiki/Minute_of_arc)

Call
    http://localhost:1111/api/utm/17 630084 4833438.xml?outputformat=latlongdeg

Output serialized as XML:

    <LatLong>
      <Lat>N 43°38'33.22''</Lat>
      <Long>W 79°23'13.71''</Long>
      <Fmt>LLFdms</Fmt>
      <LatLongString>lat: 43.642562°, long: -79.387143°</LatLongString>
    </LatLong>

Call
    http://localhost:1111/api/utm/17 630084 4833438.json?outputformat=latlongdeg

Output serialized as JSON:

    {"Lat":"N 43°38'33.22''",
     "Long":"W 79°23'13.71''",
     "Fmt":"LLFdms",
     "LatLongString":"lat: 43.642562°, long: -79.387143°"}

* Lat, Long: Latitude, Longitude of converted coordinate, in degrees. Full
  degrees are separated by "°" from the minutes, which are separated by " ' " from
  the seconds, which are postfixed by " '' ". All zero minutes or seconds are
  omitted.
* Fmt: For serialization as arc seconds and arc minutes, the string "LLFdms"
  denotes, how "Lat" and "Long" shall be interpreted.
* LatLongString: A canonical representation of latitude and longitude as decimal
  degrees.

#### Output requested as latitude and longitude in [degrees](http://en.wikipedia.org/wiki/Degree_(angle))

Call

    http://localhost:1111/api/utm/17T 630084 4833438.xml?outputformat=latlongcomma

Output serialized as XML:

    <LatLong>
      <Lat>43.642562</Lat>
      <Long>-79.387143</Long>
      <Fmt>LLFdeg</Fmt>
      <LatLongString>lat: 43.642562°, long: -79.387143°</LatLongString>
    </LatLong>

* Lat, Long: Latitude, Longitude of converted coordinate, in decimal degrees.
* Fmt: For serialization as a decimal, the string "LLFdeg"

#### Output requested as [Geohash](http://en.wikipedia.org/wiki/Geohash)

Call

    http://localhost:1111/api/utm/17T 630084 4833438.json?outputformat=geohash

Output serialized as JSON:

    {"GeoHash":"dpz838bh37pv"}

Call

    http://localhost:1111/api/utm/17T 630084 4833438.xml?outputformat=geohash

Output serialized as XML:

    <GeoHash>
      <GeoHash>dpz838bh37pv</GeoHash>
    </GeoHash>

#### Output requested as [BMN](http://homepage.ntlworld.com/anton.helm/bmn_mgi.html)
[Bundesmeldenetz](http://de.wikipedia.org/wiki/Bundesmeldenetz)

Call

    http://localhost:1111/api/utm/17T 630084 4833438.xml?outputformat=bmn

Output

This call returns with status code 500: Internal server error: Invalid argument

    <Error>
      <Error>invalid argument</Error>
    </Error>

In case of

    http://localhost:1111/api/utm/17T 630084 4833438.json?outputformat=bmn

the return would be

    {"Error":"invalid argument"}

The reason is the requested serialization as a BMN bearing, which has only a
valid representation within the longitude of 8°50' and 17°50'. Errors get
serialized in the requested encoding, unless the serialization itself fails. In
that case the error is returned text/plain encoded. Valid calls would be:

Call

    http://localhost:1111/api/utm/33T 442552 5268825.xml?outputformat=bmn

Output serialized as XML:

    <BMN>
      <BMNCoord>
        <Right>517965.58808025334</Right>
        <Height>270554.81500793993</Height>
        <RelHeight>0</RelHeight>
        <Meridian>2</Meridian>
        <El>
          <CommonName>Bessel1841MGI</CommonName>
        </El>
      </BMNCoord>
      <BMNString>M31 517966 270555</BMNString>
    </BMN>

The reference ellipsoid of a BMN bearing is always the
[Bessel1841MGI](http://de.wikipedia.org/wiki/Geod%C3%A4tisches_Datum#Deutschland_und_.C3.96sterreich)
ellipsoid, which requires a helmert transformation from WGS84 to Bessel.

Call

    http://localhost:1111/api/utm/33T 442552 5268825.json?outputformat=bmn

Output serialized as JSON:

    {"BMNCoord":{"Right":517965.58808025334,"Height":270554.81500793993,"RelHeight":0,"Meridian":2,"El":{"CommonName":"Bessel1841MGI"}},"BMNString":"M31 517966 270555"}

Geohash - Conversions
---------------------

All following conversions have the same output specifier as described by the UTM conversion.

* Output specifiers are utm, geohash, latlongdeg, latlongcomma bmn
* Errors are encoded in the requested encoding (XML, JSON), unless the encoding itself fails,
  which means the error is encoded as text/plain.

Base Url for Geohash operations:
   
    URL/APIRoot/geohash/<VALUE>[.xml|.json]?outputformat=<utm|geohash|latlongdeg|latlongcomma|bmn>

Value is a coordinate in Geohash representation. The reference ellipsoid is always
the WGS84Ellipsoid.

Examples of valid input values:

    u4pruydqqvj
    ezs42

Call

    http://localhost:1111/api/geohash/u4pruydqqvj.xml?outputformat=latlongdeg

Output serialized as XML:

    <LatLong>
      <Lat>N 57°38'56.8''</Lat>
      <Long>E 10°24'26.78''</Long>
      <Fmt>LLFdms</Fmt>
      <LatLongString>lat: 57.649112°, long: 10.40744°</LatLongString>
    </LatLong>

Call

    http://localhost:1111/api/geohash/ezs42.json?outputformat=latlongcomma

Output serialized as JSON:

    {"Lat":"42.6","Long":"-5.6","Fmt":"LLFdeg","LatLongString":"lat: 42.6°, long: -5.6°"}

Latitude / Longitude - Conversions
----------------------------------

Latitude and Longitude can't be represented as a single value, thus lat and long have to be specified as parameters.
Specifying a value will result in an error:

Call

    http://localhost:1111/api/latlong/23.json?lat=47.57°&long=14.23°&outputformat=utm

Output:

    {"Error":"Latlong doesn't accept an input value. Use parameters instead"}

Parameters:

    lat: Latitude in fraction of degrees or decimal fractions
    long: Longitude in fraction of degrees or decimal fractions

Both fraction of degrees or decimal fractions may be used. Examples:

Decimal fractions must end with the degree sign °. A negative sign " - "
prefixing latitude or longitude denotes eastern latitude or southern longitude.
Abbreviations of the main directions " N, S, E, W " may be used as well, in
which case " E " or " S " means a negative bearing.

    -179.50°
    -179°
    S50.50°
    N23.34567

As an alternative, fractions of degrees can be used to specify latitude and
longitude. In this case, whole parts of degrees must be postfixed with the
degree sign " ° ", minutes postfixed with the minute specifier " ' " and seconds
postfixed with the degrees seconds denominator sign " '' ". A negative sign " - "
prefixing latitude or longitude denotes eastern latitude or southern longitude.
Abbreviations of the main directions " N, S, E, W " may be used as well, in
which case " E " or " S " means a negative bearing.

<table>
  <tr>
    <th>Bearing in fraction of degrees</th>
    <th>Equivalent decimal fractions</th>
  </tr>
  <tr>
    <td>N 30 °56 ' 34.45 ''</td>
    <td>30.942903</td>
  </tr>
  <tr>
    <td>N30 ° 56 ' 34.45''</td>
    <td>30.942903</td>
  </tr>
  <tr>
    <td>N 170°56'34.45''</td>
    <td>170.942903</td>
  </tr>
  <tr>
    <td>170 °</td>
    <td>170.0</td>
  </tr>
  <tr>
    <td>- 359 ° 30'</td>
    <td>-359.5</td>
  </tr>
  <tr>
    <td>- 180°30'30''</td>
    <td>-180.508333</td>
  </tr>
  <tr>
    <td>- 180°0'30''</td>
    <td>-180.008333</td>
  </tr>
  <tr>
    <td>- 180°00'30''</td>
    <td>-180.008333</td>
  </tr>
  <tr>
    <td>- 180°0'0.5''</td>
    <td>-180.000139</td>
  </tr>
  <tr>
    <td>- 180°30'</td>
    <td>-180.5</td>
  </tr>
</table>

Call

    http://localhost:1111/api/latlong/.json?lat=47.57°&long=14°0'27''&outputformat=latlongdeg
    http://localhost:1111/api/latlong/.json?lat=47.57°&long=14°0'27''&outputformat=utm
    http://localhost:1111/api/latlong/.xml?lat=47.57°&long=14°0'27''&outputformat=latlongcomma

Note the mix of fractions of degrees (long=14°0'27'') and decimal fractions
(lat=47.57°) in the specification of latitude and longitude.

Output:

    {"Lat":"N 47°34'12''","Long":"E 14°0'23''","Fmt":"LLFdms","LatLongString":"lat: 47.57°, long: 14.006389°"}

    {"UTMCoord":{"Northing":5.268986550533157e+06,"Easting":425351.161981314,"Zone":"33T","El":{"CommonName":"WGS84"}},"UTMString":"33T 425351 5268987"}

    <LatLong>
      <Lat>47.57</Lat>
      <Long>14.0075</Long>
      <Fmt>LLFdeg</Fmt>
      <LatLongString>lat: 47.57°, long: 14.0075°</LatLongString>
    </LatLong>

BMN - Conversions
-----------------

The Bundesmeldenetz is the former
[geodetic datum of Austria](http://georepository.com/datum_6805/Militar-Geographische-Institut-Ferro.html),
which was officially in use until about the end of 2009, however a lot of legacy
information including scientific reports, official maps of municipalities and
the Austrian provinces still carries this datum and is unlikely to get
converted. For more information see the [bmn package](https://github.com/the42/cartconvert/tree/master/cartconvert/bmn).
Coordinates are only valid within the longitude of 8°50' and 17°50'.

* Output specifiers are utm, geohash, latlongdeg, latlongcomma or bmn
* Errors are encoded in the requested encoding (XML, JSON), unless the encoding itself fails,
  which means the error is encoded as text/plain.

Base Url for BMN operations:
   
    URL/APIRoot/bmn/<VALUE>[.xml|.json]?outputformat=<utm|geohash|latlongdeg|latlongcomma|bmn>

Value is a coordinate in BMN representation, which is always uses the . The reference ellipsoid is always
the [Bessel1841MGI](http://de.wikipedia.org/wiki/Geod%C3%A4tisches_Datum#Deutschland_und_.C3.96sterreich)
ellipsoid.

Examples of valid input values:

    M31 592269 272290
    M34 592269 272290
    M34 703168 374510
    M34 703168.99 374510 (fractions supported but unused in actual bearing)

Call

    http://localhost:1111/api/bmn/M34 703168 374510.json?outputformat=latlongdeg

Output serialized as JSON:

    {"Lat":"N 48°30'25.2''","Long":"E 15°41'55.49''","Fmt":"LLFdms","LatLongString":"lat: 48.507001°, long: 15.698748°"}

Call

    http://localhost:1111/api/bmn/M34 703168 374510.xml?outputformat=utm

Output serialized as XML:

    <UTMCoord>
      <UTMCoord>
        <Northing>5.372889492316671e+06</Northing>
        <Easting>551610.575844678</Easting>
        <Zone>33U</Zone>
        <El><CommonName>WGS84</CommonName></El>
      </UTMCoord>
      <UTMString>33U 551611 5372889</UTMString>
    </UTMCoord>

Installation
------------

go get github.com/the42/cartconv/cartconvserv

Test
----

go test github.com/the42/cartconv/cartconvserv
