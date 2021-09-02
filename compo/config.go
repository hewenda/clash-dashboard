package compo

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type Config struct {
	app.Compo
	fileName string
	conf     string
}

func (c *Config) OnNav(ctx app.Context) {
	ctx.Page.SetTitle("Config - Clash Dashboard")
}

var urlReg = regexp.MustCompile(`/config/(\d+)\.yaml`)

func (c *Config) OnMount(ctx app.Context) {
	params := urlReg.FindStringSubmatch(ctx.Page.URL().Path)
	fileName := params[1]

	c.fileName = fileName
	c.Update()

	ctx.Async(func() {
		resp, err := client.R().
			SetHeader("Accept", "application/x-yaml").
			Get(fmt.Sprintf("api/auth/config/%s", c.fileName))

		if err != nil {
			fmt.Println(err)
			return
		}

		if resp.StatusCode() == http.StatusOK {
			c.conf = resp.String()
			c.Update()
			c.editorInit()
		}
	})
}

func (c *Config) Render() app.UI {
	child := app.Div().Class("h-full").Body(
		app.Script().Src("/web/codemirror/codemirror.js"),
		app.Script().Src("/web/codemirror/yaml.js"),
		app.Script().Src("/web/codemirror/lint.js"),
		app.Script().Src("/web/codemirror/js-yaml.min.js"),
		app.Script().Src("/web/codemirror/yaml-lint.js"),
		app.Link().Type("text/css").Rel("stylesheet").Href("/web/codemirror/codemirror.css"),
		app.Link().Type("text/css").Rel("stylesheet").Href("/web/codemirror/lint.css"),
		app.Style().Body(
			app.Text(`
				.CodeMirror {
					height: auto;
					min-height: 300px;
				}
			`),
		),

		app.Div().Class("flex h-full flex-col").Body(
			app.Textarea().
				ID("editor").
				Class("w-full h-full flex-1").
				Title("clash-config").
				Body(
					app.Text(c.conf),
				),
			app.Div().Class("flex justify-end mt-2").Body(
				app.Button().Class("ui button").Body(app.Text("Save")).OnClick(c.onSave()),
			),
		),
	)

	return newLayout().Content(
		child,
	)
}

func (c *Config) editorInit() {
	onChange := app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		value := args[0].Call("getValue").String()
		c.conf = value

		return nil
	})

	c.Defer(func(ctx app.Context) {
		textarea := app.Window().
			Get("document").
			Call("querySelector", "#editor")

		editor := app.Window().Get("CodeMirror").Call("fromTextArea", textarea.JSValue(), map[string]interface{}{
			"mode":        "text/yaml",
			"lineNumbers": true,
			"tabSize":     2,
			"gutters":     []interface{}{`CodeMirror-lint-markers`},
			"lint": map[string]interface{}{
				"highlightLines": true,
			},
			"onChange":  onChange,
			"extraKeys": map[string]interface{}{},
		}).JSValue()

		editor.Call("on", "change", onChange)
	})
}

func (c *Config) onSave() app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		ctx.Async(func() {

			resp, err := client.R().
				SetHeader("Accept", "application/json").
				SetBody(map[string]interface{}{
					"Config": c.conf,
				}).
				Post(fmt.Sprintf("/api/auth/config/%s", c.fileName))

			if err != nil {
				fmt.Println(err)
				return
			}

			if resp.StatusCode() == http.StatusOK {
				ctx.Navigate("/config")
			}
		})
	}
}
