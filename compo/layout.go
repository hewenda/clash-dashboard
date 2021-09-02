package compo

import (
	"net/http"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type Layout struct {
	app.Compo
	IContent []app.UI
}

func (c *Layout) OnMount(ctx app.Context) {
	var token string
	ctx.LocalStorage().Get("token", &token)

	if token != "" {
		client.SetAuthToken(token)

		ctx.Async(func() {
			resp, _ := client.R().Get("/api/auth/ping")

			url := ctx.Page.URL()

			if ctx.Page.URL().Path == "/" && resp.StatusCode() == http.StatusOK {
				url.Path = "/config"
				ctx.NavigateTo(url)
			} else if ctx.Page.URL().Path != "/" && resp.StatusCode() == http.StatusUnauthorized {
				url.Path = "/"
				ctx.NavigateTo(url)
			}

		})
	}
}

func (c *Layout) Content(v ...app.UI) *Layout {
	c.IContent = app.FilterUIElems(v...)

	return c
}

func (c *Layout) Render() app.UI {
	return app.Main().
		Class("ui container mx-auto h-screen py-5").
		Body(
			app.Range(c.IContent).Slice(func(i int) app.UI {
				return c.IContent[i]
			}),
		)
}

func newLayout() *Layout {
	return &Layout{}
}
