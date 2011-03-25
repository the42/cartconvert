package httptempltest

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
	items []string
}

type Footer struct {
	posted string
}

type Page struct {
	title   string
	content Content
	footer  Footer
}

var templateNames = []string{
	"layout.tpl",
	"content.tpl",
	"header.tpl",
	"footer.tpl",
}

var templates = make(map[string]*template.Template)

func evalTemplate(wr io.Writer, formatter string, data ...interface{}) {
	templates[formatter].Execute(wr, data)
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
	templates["main.html"].Execute(os.Stdout, page)
}
