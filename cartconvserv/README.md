cartconvserv - A RESTFul service for coordinate transformations
===============================================================

Purpose
-------

cartconvserv is a RESTFul service to the cartconv - package. It may be installed
as a stand-alone service or executed in the [Appengine](http://code.google.com/appengine/docs/go/) environment.


Functionality
-------------

* Conversion between coordinate systems and bearing representations. Valid representations are
  [Latitude and Longitude](http://en.wikipedia.org/wiki/Geographic_coordinate_system#Geographic_latitude_and_longitude),
  [UTM](http://en.wikipedia.org/wiki/Universal_Transverse_Mercator_coordinate_system),
  [geohash](http://en.wikipedia.org/wiki/Geohash),
  [Bundesmeldenetz](http://de.wikipedia.org/wiki/Bundesmeldenetz) used in Austria and
  [OSGB36, Ordnance Survey National Grid](http://en.wikipedia.org/wiki/OSGB) used in the UK.
* Serialization as XML or JSON by content negotiation.

Convention for this help:

* Binding: Base url to service (configurable)
* APIRoot: Root of the RESTFul API (configurable)


UTM - Conversions
-----------------

Base url for UTM operations:
   
    Binding/APIRoot/utm/<VALUE>.[xml|json]?outputformat=<latlongdeg|latlongcomma|utm|geohash|bmn|osgb>

Value is a coordinate in UTM representation. The reference ellipsoid is always
the WGS84Ellipsoid.

Examples of valid input values:

    17T 630084 4833438
    17T 630084.31 4833438.54

If the extension to value is empty or ".json", the result of the requested
output format is JSON-encoded.

If the extension to value is ".xml", the result of the requested output format
is XML-encoded.

The value to the parameter "outputformat" is one of

* latlongdeg: Latitude and longitude with fractions in degrees
* latlongcomma: Latitude and longitude with decimal fractions
* geohash: Geohash-encoded value of latitude and longitude
* bmn: Serialization of the value as BMN-coordinate
* osgb: Serialization of the value as OSGB36-coordinate


### Output requested as latitude and longitude in [arc degrees](http://en.wikipedia.org/wiki/Minute_of_arc)

Call
    http://localhost:1111/api/utm/17 630084 4833438.xml?outputformat=latlongdeg

Note: When running on GAE, the port may not be specified!

Output serialized as XML:

    <GEOConvertResponse>
      <Status/>
      <Code>0</Code>
      <Error>false</Error>
      <GEOConvertRequest>
        <Method>utm/</Method>
        <Values>
          <Key>outputformat</Key>
          <Values>latlongdeg</Values>
        </Values>
      </GEOConvertRequest>
      <Payload>
        <Lat>S 46°38'23.97''</Lat>
        <Long>W 175°18'1.13''</Long>
        <Fmt>LLFdms</Fmt>
        <LatLongString>lat: -46.639992°, long: -175.300313°</LatLongString>
      </Payload>
    </GEOConvertResponse>

Call
    http://localhost:1111/api/utm/17 630084 4833438.json?outputformat=latlongdeg

Output serialized as JSON:

    {"Status":"",
     "Code":0,
     "Error":false,
     "GEOConvertRequest": {"Method":"utm/","Values":[{"Key":"outputformat","Values":["latlongdeg"]}]},
     "Payload":{"Lat":"S 46°38'23.97''","Long":"W 175°18'1.13''","Fmt":"LLFdms","LatLongString":"lat: -46.639992°, long: -175.300313°"}}

* GEOConvertRequest Contains the request parameters and API method
* Payload contains the methods result:
* Lat, Long: Latitude, Longitude of converted coordinate, in degrees. Full
  degrees are separated by "°" from the minutes, which are separated by " ' " from
  the seconds, which are postfixed by " '' ". All zero minutes or seconds are
  omitted.
* Fmt: For serialization as arc seconds and arc minutes, the string "LLFdms"
  denotes, how "Lat" and "Long" shall be interpreted.
* LatLongString: A canonical representation of latitude and longitude as decimal
  degrees.


### Output requested as latitude and longitude in [degrees](http://en.wikipedia.org/wiki/Degree_(angle))

Call

    http://localhost:1111/api/utm/17T 630084 4833438.xml?outputformat=latlongcomma

Output serialized as XML:

    <GEOConvertResponse>
      <Status/>
      <Code>0</Code>
      <Error>false</Error>
      <GEOConvertRequest>
        <Method>utm/</Method>
        <Values>
          <Key>outputformat</Key>
          <Values>latlongcomma</Values>
        </Values>
      </GEOConvertRequest>
      <Payload>
        <Lat>43.642562</Lat>
        <Long>-79.387143</Long>
        <Fmt>LLFdeg</Fmt>
        <LatLongString>lat: 43.642562°, long: -79.387143°</LatLongString>
      </Payload>
    </GEOConvertResponse>

* The whole result is wrapped into a GEOConvertResponse-envelope (this is true for every XML-Response)
* GEOConvertRequest Contains the request parameters and API method
* Payload contains the methods result:
* Lat, Long: Latitude, Longitude of converted coordinate, in decimal degrees.
* Fmt: For serialization as a decimal, the string "LLFdeg"


### Output requested as [Geohash](http://en.wikipedia.org/wiki/Geohash)

Call

    http://localhost:1111/api/utm/17T 630084 4833438.json?outputformat=geohash

Output serialized as JSON:

    {"Status":"",
     "Code":0,
     "Error":false,
     "GEOConvertRequest":{"Method":"utm/","Values":[{"Key":"outputformat","Values":["geohash"]}]},
     "Payload":{"GeoHash":"dpz838bh37pv"}}


Call

    http://localhost:1111/api/utm/17T 630084 4833438.xml?outputformat=geohash

Output serialized as XML:

    <GEOConvertResponse>
      <Status/>
      <Code>0</Code>
      <Error>false</Error>
      <GEOConvertRequest>
        <Method>utm/</Method>
        <Values>
          <Key>outputformat</Key>
          <Values>geohash</Values>
        </Values>
      </GEOConvertRequest>
      <Payload>
        <GeoHash>dpz838bh37pv</GeoHash>
      </Payload>
    </GEOConvertResponse>


### Output requested as [BMN](http://homepage.ntlworld.com/anton.helm/bmn_mgi.html)
[Bundesmeldenetz](http://de.wikipedia.org/wiki/Bundesmeldenetz)

Call

    http://localhost:1111/api/utm/17T 630084 4833438.xml?outputformat=bmn

Output

This call returns with status code 500: Internal server error
and sets Error to true

    <GEOConvertResponse>
      <Status>value out of range</Status>
      <Code>0</Code>
      <Error>true</Error>
      <GEOConvertRequest>
        <Method>utm/</Method>
        <Values>
          <Key>outputformat</Key>
          <Values>bmn</Values>
        </Values>
      </GEOConvertRequest>
    </GEOConvertResponse>

In case of

    http://localhost:1111/api/utm/17T 630084 4833438.json?outputformat=bmn

the return would be

    {"Status":
     "value out of range",
     "Code":0,
     "Error":true,
     "GEOConvertRequest":{"Method":"utm/","Values":[{"Key":"outputformat","Values":["bmn"]}]},
     "Payload":null}

The reason is the requested serialization as a BMN bearing, which has only a
valid representation within the longitude of 8°50' and 17°50'. Errors get
serialized in the requested encoding, unless the serialization itself fails. In
that case the error is returned text/plain encoded. Valid calls would be:

Call

    http://localhost:1111/api/utm/33T 442552 5268825.xml?outputformat=bmn

Output serialized as XML:

    <GEOConvertResponse>
      <Status/>
      <Code>0</Code>
      <Error>false</Error>
      <GEOConvertRequest>
        <Method>utm/</Method>
        <Values>
          <Key>outputformat</Key>
          <Values>bmn</Values>
        </Values>
      </GEOConvertRequest>
      <Payload>
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
      </Payload>
    </GEOConvertResponse>

The reference ellipsoid of a BMN bearing is always the
[Bessel1841MGI](http://de.wikipedia.org/wiki/Geod%C3%A4tisches_Datum#Deutschland_und_.C3.96sterreich)
ellipsoid, which requires a helmert transformation from WGS84 to Bessel.

Call

    http://localhost:1111/api/utm/33T 442552 5268825.json?outputformat=bmn

Output serialized as JSON:

    {"Status":"",
     "Code":0,
     "Error":false,
     "GEOConvertRequest":{"Method":"utm/","Values":[{"Key":"outputformat","Values":["bmn"]}]},
     "Payload": {"BMNCoord":
                {"Right":517965.58808025334,"Height":270554.81500793993,"RelHeight":0,"Meridian":2,"El":{"CommonName":"Bessel1841MGI"}},
                 "BMNString":"M31 517966 270555"}}


### Output requested as [OSGB36](http://en.wikipedia.org/wiki/OSGB)

Call

    http://localhost:1111/api/utm/31U 365166 5684564.xml?outputformat=osgb

Output serialized as XML:

    <GEOConvertResponse>
      <Status/>
      <Code>0</Code>
      <Error>false</Error>
      <GEOConvertRequest>
        <Method>utm/</Method>
        <Values>
          <Key>outputformat</Key>
          <Values>osgb</Values>
        </Values>
      </GEOConvertRequest>
      <Payload>
        <OSGB36Coord>
          <Easting>13862</Easting>
          <Northing>59718</Northing>
          <RelHeight>0</RelHeight>
          <Zone>TR</Zone>
          <El>
            <CommonName>Airy1830</CommonName>
          </El>
        </OSGB36Coord>
        <OSGB36String>TR1386259718</OSGB36String>
      </Payload>
    </GEOConvertResponse>

The reference ellipsoid of a OSGB36 bearing is always the
[Airy1830](http://en.wikipedia.org/wiki/Ordnance_Survey_National_Grid#General)
ellipsoid, which requires a helmert transformation from WGS84 to Airy1830.


Geohash - Conversions
---------------------

All following conversions have the same output specifier as described by the UTM conversion.

* Output specifiers are utm, geohash, latlongdeg, latlongcomma, bmn or osgb
* Errors are encoded in the requested encoding (XML, JSON), unless the encoding itself fails,
  which means the error is encoded as text/plain.

Base url for Geohash operations:
   
    Binding/APIRoot/geohash/<VALUE>.[xml|json]?outputformat=<utm|geohash|latlongdeg|latlongcomma|bmn|osgb>

Value is a coordinate in Geohash representation. The reference ellipsoid is always
the WGS84Ellipsoid.

Examples of valid input values:

    u4pruydqqvj
    ezs42

Call

    http://localhost:1111/api/geohash/u4pruydqqvj.xml?outputformat=latlongdeg

Output serialized as XML:

    <GEOConvertResponse>
      <Status/>
      <Code>0</Code>
      <Error>false</Error>
      <GEOConvertRequest>
        <Method>geohash/</Method>
        <Values>
          <Key>outputformat</Key>
          <Values>latlongdeg</Values>
        </Values>
      </GEOConvertRequest>
      <Payload>
        <Lat>N 57°38'56.8''</Lat>
        <Long>E 10°24'26.78''</Long>
        <Fmt>LLFdms</Fmt>
        <LatLongString>lat: 57.649112°, long: 10.40744°</LatLongString>
      </Payload>
    </GEOConvertResponse>

Call

    http://localhost:1111/api/geohash/ezs42.json?outputformat=latlongcomma

Output serialized as JSON:

    {"Status":"",
     "Code":0,
     "Error":false,
     "GEOConvertRequest":
     {"Method":"geohash/","Values":[{"Key":"outputformat","Values":["latlongcomma"]}]},
      "Payload":{"Lat":"42.6","Long":"-5.6","Fmt":"LLFdeg","LatLongString":"lat: 42.6°, long: -5.6°"}}


Latitude / Longitude - Conversions
----------------------------------

Latitude and Longitude can't be represented as a single value, thus lat and long have to be specified as parameters.
Specifying a value will result in an error:

Call

    http://localhost:1111/api/latlong/23.json?lat=47.57°&long=14.23°&outputformat=utm

Output:

    {"Status":"Latlong doesn't accept an input value. Use the parameters 'lat' and 'long' instead",
     "Code":0,
     "Error":true,
     "GEOConvertRequest":{"Method":"latlong/","Values":[{"Key":"lat","Values":["47.57°"]},{"Key":"outputformat","Values":["utm"]},
     {"Key":"long","Values":["14.23°"]}]},
     "Payload":null}

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

* Output specifiers are utm, geohash, latlongdeg, latlongcomma, bmn or osgb
* Errors are encoded in the requested encoding (XML, JSON), unless the encoding itself fails,
  which means the error is encoded as text/plain.

Base url for latitude, longitude operations:
   
    Binding/APIRoot/latlong/.[xml|json]?outputformat=<utm|geohash|latlongdeg|latlongcomma|bmn|osgb>

Call

    http://localhost:1111/api/latlong/.json?lat=47.57°&long=14°0'27''&outputformat=latlongdeg
    http://localhost:1111/api/latlong/.json?lat=47.57°&long=14°0'27''&outputformat=utm
    http://localhost:1111/api/latlong/.xml?lat=47.57°&long=14°0'27''&outputformat=latlongcomma

Note the mix of fractions of degrees (long=14°0'27'') and decimal fractions
(lat=47.57°) in the specification of latitude and longitude.

Output:

    {"Status":"",
     "Code":0,
     "Error":false,
     "GEOConvertRequest":{"Method":"latlong/","Values":[{"Key":"lat","Values":["47.57°"]},{"Key":"outputformat","Values":["latlongdeg"]},{"Key":"long","Values":["14°0'27''"]}]},
     "Payload":{"Lat":"N 47°34'12''","Long":"E 14°0'27''","Fmt":"LLFdms","LatLongString":"lat: 47.57°, long: 14.0075°"}}

    {"Status":"",
     "Code":0,
     "Error":false,
     "GEOConvertRequest":{"Method":"latlong/","Values":[{"Key":"lat","Values":["47.57°"]},{"Key":"long","Values":["14°0'27''"]},{"Key":"outputformat","Values":["utm"]}]},
     "Payload":{"UTMCoord":{"Northing":5.268986550533157e+06,"Easting":425351.161981314,"Zone":"33T","El":{"CommonName":"WGS84"}},"UTMString":"33T 425351 5268987"}}

    <GEOConvertResponse>
      <Status/>
      <Code>0</Code>
      <Error>false</Error>
      <GEOConvertRequest>
        <Method>latlong/</Method>
        <Values>
          <Key>lat</Key>
          <Values>47.57°</Values>
        </Values>
        <Values>
          <Key>outputformat</Key>
          <Values>latlongcomma</Values>
        </Values>
        <Values>
          <Key>long</Key>
          <Values>14°0'27''</Values>
      </Values>
      </GEOConvertRequest>
      <Payload>
        <Lat>47.57</Lat>
        <Long>14.0075</Long>
        <Fmt>LLFdeg</Fmt>
        <LatLongString>lat: 47.57°, long: 14.0075°</LatLongString>
      </Payload>
    </GEOConvertResponse>


BMN - Conversions
-----------------

The Bundesmeldenetz is the former
[geodetic datum of Austria](http://georepository.com/datum_6805/Militar-Geographische-Institut-Ferro.html),
which was officially in use until about the end of 2009, however a lot of legacy
information including scientific reports, official maps of municipalities and
the Austrian provinces still carries this datum and is unlikely to get
converted. For more information see the [bmn package](https://github.com/the42/cartconvert/tree/master/cartconvert/bmn).
Coordinates are only valid within the longitude of 8°50' and 17°50'.

* Output specifiers are utm, geohash, latlongdeg, latlongcomma, bmn or osgb
* Errors are encoded in the requested encoding (XML, JSON), unless the encoding itself fails,
  which means the error is encoded as text/plain.

Base url for BMN operations:
   
    Binding/APIRoot/bmn/<VALUE>.[xml|json]?outputformat=<utm|geohash|latlongdeg|latlongcomma|bmn|osgb>

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

    {"Status":"",
     "Code":0,
     "Error":false,
     "GEOConvertRequest":{"Method":"bmn/","Values":[{"Key":"outputformat","Values":["latlongdeg"]}]},
     "Payload":{"Lat":"N 48°30'25.2''","Long":"E 15°41'55.49''","Fmt":"LLFdms","LatLongString":"lat: 48.507001°, long: 15.698748°"}}

Call

    http://localhost:1111/api/bmn/M34 703168 374510.xml?outputformat=utm

Output serialized as XML:

    <GEOConvertResponse>
      <Status/>
      <Code>0</Code>
      <Error>false</Error>
      <GEOConvertRequest>
        <Method>bmn/</Method>
        <Values>
          <Key>outputformat</Key>
          <Values>utm</Values>
        </Values>
      </GEOConvertRequest>
      <Payload>
        <UTMCoord>
          <Northing>5.372889492316671e+06</Northing>
          <Easting>551610.575844678</Easting>
          <Zone>33U</Zone>
          <El>
            <CommonName>WGS84</CommonName>
          </El>
        </UTMCoord>
        <UTMString>33U 551611 5372889</UTMString>
      </Payload>
    </GEOConvertResponse>


OSGB36 - Conversions
-----------------

Base url for osgb36 operations:
   
    Binding/APIRoot/osgb/<VALUE>.[xml|json]?outputformat=<utm|geohash|latlongdeg|latlongcomma|bmn|osgb>

Value is a coordinate in [OSGB36 representation](). The reference ellipsoid is always
the [Airy1830](http://en.wikipedia.org/wiki/Ordnance_Survey_National_Grid#cite_note-5) ellipsoid.

Examples of
[valid](http://en.wikipedia.org/wiki/Ordnance_Survey_National_Grid#Grid_letters)
[input](http://en.wikipedia.org/wiki/Ordnance_Survey_National_Grid#Grid_digits) values:


<table>
  <tr>
    <th>Input specifier</th>
    <th>Actual value</th>
    <th>Remarks</th>
  </tr>
  <tr>
    <td>NN or N N</td> <td>NN 0 0</td>
    <td>Point at NN 00000 00000 (30V 361994 6225065)</td>
  <tr>
    <td>NN11 or NN 1 1</td> <td>NN1000010000</td>
    <td>Point at NN1000010000 (30V 371845 6235206)</td>
  </tr>
  <tr>
    <td>NN1010 or NN 10 10</td> <td>NN1050010500</td>
    <td>Point at NN1050010500 (30V 372337 6235713)</td>
  </tr>
  <tr>
    <td>NN123123 or NN 123 123</td> <td>NN1235012350</td>
    <td>Point at NN1235012350 (30V 374160 6237590)</td>
  </tr>
  <tr>
    <td>NN12321232 or NN 1232 1232</td> <td>NN1232512325</td>
    <td>Point at NN1232512325 (30V 374135 6237564)</td>
  </tr>
  <tr>
    <td>NN1232112321 or NN 12321 12321</td> <td>NN1232112321</td>
    <td>Point at NN1232112321 (30V 374131 6237560)</td>
  </tr>
  <tr>
    <td>NN1230012300 or NN 12300 12300</td> <td>NN1230012300</td>
    <td>Point at NN1230012300 (30V 374110 6237539)</td>
  </tr>
</table>

In the case of `CCdddd`, `CCdddddd`, `CCdddddddd` the remaining digits for a fully qualified
OSGB36 datum are set to the middle of the rectangle, which is implicitly spanned by the inaccuracy
of the point specification. So `CC dd dd` gets `CC dd500 dd500`, `CC ddd ddd` gets `CC ddd50 ddd50`.
Trailing zeros prevent the automatic averaging towards the rectangles center:

NN 123 123 means NN 12350 12350 (30V 374160 6237590), whereas
NN 12300 12300 really specifies the point at NN 12300 12300 (30V 374110 6237539).

If only the
[letter designators](http://en.wikipedia.org/wiki/Ordnance_Survey_National_Grid#Grid_letters) `CC`
or one additional northing / easting `CC d d` is specified, automatic averaging towards the middle
of the rectangle is not performed.

Call

    http://localhost:1111/api/osgb/NN123123.json?outputformat=osgb
    http://localhost:1111/api/osgb/NN 123 123.json?outputformat=osgb

Output serialized as JSON:

    {"Status":"",
     "Code":0,
     "Error":false,
     "GEOConvertRequest":{"Method":"osgb/","Values":[{"Key":"outputformat","Values":["osgb"]}]},
     "Payload":{"OSGB36Coord":{"Easting":12350,"Northing":12350,"RelHeight":0,"Zone":"NN","El":{"CommonName":"Airy1830"}},"OSGB36String":"NN1235012350"}}

Call

    http://localhost:1111/api/osgb/NN1238812388.xml?outputformat=utm

Output serialized as XML:

    <GEOConvertResponse>
      <Status/>
      <Code>0</Code>
      <Error>false</Error>
      <GEOConvertRequest>
        <Method>osgb/</Method>
        <Values>
          <Key>outputformat</Key>
          <Values>utm</Values>
        </Values>
      </GEOConvertRequest>
      <Payload>
        <UTMCoord>
          <Northing>6.237628180629014e+06</Northing>
          <Easting>374197.17282885656</Easting>
          <Zone>30V</Zone>
          <El>
            <CommonName>WGS84</CommonName>
          </El>
        </UTMCoord>
        <UTMString>30V 374197 6237628</UTMString>
      </Payload>
    </GEOConvertResponse>


Configuration
-------------

### Stand alone application

Both `URL` (base url to service), `APIRoot` (root of the RESTFul API)
and DocRoot (root of documentation) can be configured. The default values are

* Binding: `:1111`, that is listening on TCP/IP port 1111
* APIRoot: `/api/`
* DocRoot: `/doc/`

which means the RESTFul handlers base Url listen at `:1111/api/` and documentation is served at `:1111/doc/`

By default, a JSON-encoded file named `config.json` is loaded within the directory in which the RESTFul service is started,
the contents gets parsed and the respective default values for APIRoot, Binding and DocRoot are replaced by
the corresponding values of keys named `APIRoot`, `Binding` and `DocRoot`. Example:

    {
        "APIRoot": "/myapi/",
        "Binding": "www.example.com:8080"
        "DocRoot": "/doc/"
    }

The command line parameter `--config=<filespec>` appended to the executable overides the location,
at which the configuration file gets loaded.

A call might look like

    http://www.example.com:8080/myapi/utm/17T 630084 4833438.xml?&outputformat=bmn

and documentation for UTM transformations is served as

    http://www.example.com:8080/doc/utm/

Both APIRoot and DocRoot must be valid, non-empty paths.

### Heroku
Not yet (20120214) tested on Heroku


Installation
------------

    go get github.com/the42/cartconv/cartconvserv

If you do not want documentation for the API functions, execute 

    go get -tags=nodoc github.com/the42/cartconv/cartconvserv


Test
----

go test github.com/the42/cartconv/cartconvserv
