{{define "DocSetup"}}<!DOCTYPE HTML>
<html>
<head>
  <title>Documentation for the cartconvserv API{{if .ConcreteHeading}} - {{.ConcreteHeading}}{{end}}</title>
  <link rel="shortcut icon" href="../{{template "Back"}}/static/images/favicon.png" type="image/png"/> 
  <link rel="icon" href="../{{template "Back"}}/static/images/favicon.png" type="image/png"/>
  <link rel="stylesheet" type="text/css" href="../{{template "Back"}}/static/css/styles.css"/> 
</head>
<body>{{template "Payload"}}
  <nav>
    <ul>{{range .Navigation}}
      <li><a href="{{template "Back"}}{{.URL}}">{{.Documentation}}</a></li>{{end}}
    </ul>
  </nav>
</body>
</html>{{end}}