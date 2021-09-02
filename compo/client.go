package compo

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

var client = resty.New().SetHostURL(
	fmt.Sprintf("%s://%s", app.Window().URL().Scheme, app.Window().URL().Host),
)
