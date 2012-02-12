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

TODO

BMN - Conversions
-----------------

TODO

Installation
------------

go get github.com/the42/cartconv/cartconvserv

Test
----

go test github.com/the42/cartconv/cartconvserv
