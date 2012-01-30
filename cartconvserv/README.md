cartconvserv - A RESTFul service for coordinate transformation
==============================================================

UTMto (comma, deg, geohash):	/utm/<value>.{json|xml}?outputformat={=latlongcomma|latlongdeg|geohash|bmn}
latlongto (utm, geohash): 	/latlong/<value>.{json|xml}/?lat=&long=&outputformat={=utm|geohash|bmn}
geohashtolatlong: 		/geohash/<hash>.{json.xml}[?outputformat={=latlongdeg|latlongcomma}]
bmnto (latlong, utm, geohash):	/bmn/<value>.{json|xml}[?outputformat={=latlongdeg|latlongcomma|utm|geohash}]
