package cyoa

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHanlderTmpl))
}

var defaultHanlderTmpl = `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8">
        <title>Choose Your Own Adventure</title>
        <style>
            body {
                background: #f5ecd7;
                font-family: 'Georgia', serif;
                color: #3e2c1c;
                margin: 0;
                padding: 0;
            }
            .book-container {
                max-width: 700px;
                margin: 40px auto;
                background: #fffbe9;
                box-shadow: 0 8px 32px rgba(60,40,20,0.15), 0 1.5px 0 #e2d3b1;
                border-radius: 12px;
                padding: 48px 36px 36px 36px;
                border: 1.5px solid #e2d3b1;
                min-height: 500px;
            }
            h1 {
                font-family: 'Merriweather', 'Georgia', serif;
                font-size: 2.5em;
                margin-bottom: 0.5em;
                text-align: center;
                letter-spacing: 1px;
                color: #7c5c36;
                text-shadow: 0 1px 0 #e2d3b1;
            }
            p {
                font-size: 1.2em;
                line-height: 1.7;
                margin: 1.2em 0;
                text-indent: 2em;
            }
            ul {
                list-style: none;
                padding: 0;
                margin-top: 2.5em;
                text-align: center;
            }
            li {
                display: inline-block;
                margin: 0 1em;
            }
            a {
                display: inline-block;
                background: #e2d3b1;
                color: #7c5c36;
                text-decoration: none;
                font-weight: bold;
                padding: 0.6em 1.2em;
                border-radius: 6px;
                box-shadow: 0 2px 8px rgba(60,40,20,0.07);
                transition: background 0.2s, color 0.2s, box-shadow 0.2s;
                font-size: 1.1em;
            }
            a:hover {
                background: #7c5c36;
                color: #fffbe9;
                box-shadow: 0 4px 16px rgba(60,40,20,0.13);
            }
        </style>
    </head>
    <body>
        <div class="book-container">
            <h1>{{.Title}}</h1>
            {{range .Paragraphs}}
            <p>{{.}}</p>
            {{end}}
            <ul>
                {{range .Options}}
                <li><a href="/story/{{.Arc}}">{{.Text}}</a></li>
                {{end}}
            </ul>
        </div>
    </body>
</html>
`

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}

}

func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultPathFn}

	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/" || path == "" {
		path = "/intro"
	}

	return path[1:] // Remove leading slash

}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)
	if arc, ok := h.s[path]; ok {

		err := h.t.Execute(w, arc)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "something went wrong...", http.StatusInternalServerError)
			return
		}
		return

	}
	http.Error(w, "Arc not found", http.StatusNotFound)

}

func Jsonstory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		return nil, err
	}
	return story, nil
}

type Story map[string]Arc

type Arc struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}
