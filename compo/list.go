package compo

import (
	"fmt"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type List struct {
	app.Compo
	fileList []string
}

func (c *List) OnNav(ctx app.Context) {
	ctx.Page.SetTitle("Config - Clash Dashboard")
}

func (c *List) OnMount(ctx app.Context) {
	ctx.Async(func() {
		type Result struct {
			Message []string
		}

		var result Result

		resp, err := client.R().
			SetHeader("Accept", "application/x-yaml").
			SetResult(&result).
			Get("api/auth/config")

		if err != nil {
			fmt.Println(err)
			return
		}

		if resp.StatusCode() == http.StatusOK {
			c.fileList = result.Message
			c.Update()
		}
	})
}

func (c *List) Render() app.UI {
	child := app.Div().Class("h-full").Body(
		app.Div().Class("flex justify-start space-x-2").Body(
			app.A().Class("ui secondary button").Href(fmt.Sprintf("/config/%d.yaml", time.Now().Unix())).Body(app.Text("New")),
		),
		app.Div().Class("flex flex-col divide-y").Body(
			app.Range(c.fileList).Slice(func(i int) app.UI {
				return app.Div().Class("flex items-center py-3").Body(
					app.Div().Class("flex-1 h-10 flex items-center").Body(
						app.Span().Body(app.Text(c.fileList[i])),
					),
					app.Div().Class("flex items-center").Body(
						app.Button().Class("ui button").Body(app.Text("Copy")).OnClick(c.onCopy(c.fileList[i])),
						app.A().Href(fmt.Sprintf("/config/%s", c.fileList[i])).Class("ui button").Body(app.Text("Edit")),
					),
				)
			}),
		),
	)

	return newLayout().Content(
		child,
	)
}

func (c *List) onCopy(fileName string) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		var token string
		ctx.LocalStorage().Get("token", &token)
		uri := fmt.Sprintf("%s://%s/api/auth/config/%s?token=%s", ctx.Page.URL().Scheme, ctx.Page.URL().Host, fileName, token)

		app.Window().Get("navigator").Get("clipboard").Call("writeText", uri)
	}
}
