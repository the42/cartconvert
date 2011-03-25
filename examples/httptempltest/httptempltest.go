package main

/*
Consider the templates:
—-main.html —-

{@|header.html}

{content|content.html}
{footer|footer.html}

—header.html—-
{title}

—content.html—-

{.repeated section items}
{@}
{.end}

—footer.html—-
Posted: {posted}
*/


import (
	"template"
	"os"
	"io"
)

type Content struct {
	Items []string
}

type Footer struct {
	Posted string
}

type Page struct {
	Title   string
	Content Content
	Footer  Footer
}

var templateNames = []string{
	"layout1.tpl",
	"content1.tpl",
	"header1.tpl",
	"footer1.tpl",
}

var templates = make(map[string]*template.Template)

func evalTemplate(wr io.Writer, formatter string, data ...interface{}) {
	err := templates[formatter].Execute(wr, data)
	if err != nil {
	  print(err.String())
	}
}

func main() {
	fmap := template.FormatterMap{}

	for _, name := range templateNames {
		fmap[name] = evalTemplate
	}

	for _, name := range templateNames {
		templates[name] = template.MustParseFile(name, fmap)
	}

	page := Page{"test page", Content{[]string{"a", "b"}}, Footer{"today"}}
	err := templates["layout1.tpl"].Execute(os.Stdout, page)
	if err != nil {
	  print(err.String())
	}
}
