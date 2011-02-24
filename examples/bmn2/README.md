bmn2 - Konvertiert Bundesmeldenetz Daten zu Längen/Breiten
==========================================================

Funktion
--------

bmn2 liest geografische Daten im Bundesmeldenetz von stdin, konvertiert diese
Daten, und schreibt das Ergebnis nach stdout. Fehler werden nach stderr
geschrieben.

Verwendung
----------

Usage of ./bmn2:
  -of="deg": specify output format. Possible values are:  dms  geohash  utc  deg

Das Eingabeformat muss diese Konvention erfüllen:

M28|M31|M34 xxx[.yyy] xxxx[.yyy]

z.B.: "infile.txt":

M31 592269 272290
M28 592270 272290
M34 592269 272290
M34 703168 374510
M34 592269 272290.05

bmn2 -of="deg" < infile.txt
  Liest Koordinaten im BMN-Format aus der Datei "infile.txt" und schreibt das
  Ergebnis im Format Länge und Breite in Dezimalschreibweise nach stdout:

lat: 47.573851°, long: 15.223856°
lat: 47.439212°, long: 16.197434°
lat: 47.570299°, long: 14.236188°
lat: 48.507001°, long: 15.698748°
lat: 47.570299°, long: 14.236188°


bmn2 -of="dms" < infile.txt > outfile.txt
  Liest Koordinaten im BMN-Format aus der Datei "infile.txt" und schreibt das
  Ergebnis für Länge und Breite im Format Grad°Minuten'Sekunden.Komma'' in die
  Datei "outfile.txt"

N 47°34'25.86'', E 15°13'25.88''
N 47°26'21.16'', E 16°11'50.76''
N 47°34'13.07'', E 14°14'10.28''
N 48°30'25.20'', E 15°41'55.49''
N 47°34'13.08'', E 14°14'10.28''


bmn2 -of="utm" < infile.txt > outfile.txt
  Liest Koordinaten im BMN-Format aus der Datei "infile.txt" und schreibt das
  Ergebnis in UTM Koordinaten in die Datei "outfile.txt"

33T 516836 5268962
33T 590286 5254669
33T 442552 5268825
33U 551611 5372889
33T 442552 5268825

bmn2 -of="geohash" < infile.txt
  Liest Koordinaten im BMN-Format aus der Datei "infile.txt" und schreibt das
  Ergebnis als [Geohash](http://en.wikipedia.org/wiki/Geohash) nach stdout

u26ydkt9v8d5
u27tbs497w44
u26negymp4r5
u2e5vnrmz276
u26negymp4rn


Installation
------------

  make install


Lizenz
-------

Dieses Programm und der Quellcode wird unter der [Simplified BSD
License](http://www.freebsd.org/copyright/freebsd-license.html) Lizenz zur
Verfügung gestellt. Siehe "LICENSE"
