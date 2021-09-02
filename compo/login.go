package compo

import (
	"fmt"
	"net/http"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type User struct {
	Name     string
	Password string
}

type Login struct {
	app.Compo
	User
}

func (c *Login) OnNav(ctx app.Context) {
	ctx.Page.SetTitle("Login - Clash Dashboard")
}

func (c *Login) OnMount(ctx app.Context) {
	c.Name = ""
	c.Password = ""
	c.Update()
}

func (c *Login) Render() app.UI {
	login := app.Form().
		Class("ui form w-96").
		Body(
			app.Div().Class("field").
				Body(
					app.Div().Class("field").Body(
						app.Label().Body(app.Text("User")),
						app.Input().Type("text").Name("user").Value(c.Name).OnChange(c.ValueTo(&c.User.Name)),
					),
				),
			app.Div().Class("field").
				Body(
					app.Div().Class("field").Body(
						app.Label().Body(app.Text("Password")),
						app.Input().Type("password").Name("password").Value(c.Password).OnChange(c.ValueTo(&c.User.Password)),
					),
				),
			app.Div().Class("flex justify-end mt-2").Body(
				app.Button().Class("ui sumbmit button").Type("submit").Body(app.Text("Submit")).OnClick(c.onLogin()),
			),
		)

	return newLayout().Content(
		app.Main().
			Class("container mx-auto min-h-screen flex items-center justify-center").
			Body(
				login,
			),
	)
}

type AuthInfo struct {
	Code   int    `json:"code"`
	Expire string `json:"expire"`
	Token  string `json:"token"`
}

func (c *Login) onLogin() app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		e.PreventDefault()

		ctx.Async(func() {
			var auth AuthInfo

			resp, err := client.R().
				SetHeader("Accept", "application/json").
				SetBody(User{Name: c.User.Name, Password: c.User.Password}).
				SetResult(&auth).
				Post("/api/auth")

			if err != nil {
				fmt.Println(err)
				return
			}

			if resp.StatusCode() == http.StatusOK {
				ctx.LocalStorage().Set("token", auth.Token)
				ctx.Navigate("/config")
			}
		})
	}
}
