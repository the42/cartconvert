{{define "DocSetup"}}<!DOCTYPE HTML>
<html>
  <head>
    <title>Documentation for the cartconvserv API{{if .ConcreteHeading}} - {{.ConcreteHeading}}{{end}}</title>
    <link rel="shortcut icon" href="../{{template "Back"}}/static/images/favicon.png" type="image/png"/> 
    <link rel="icon" href="../{{template "Back"}}/static/images/favicon.png" type="image/png"/>
  </head>
  <body>{{template "Payload"}}
  <nav>{{range .Navigation}}
    <a href="{{template "Back"}}{{.URL}}">{{.Documentation}}</a>{{end}}
  </nav>
  </body>
</html>{{end}}