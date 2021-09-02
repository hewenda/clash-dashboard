package main

import (
	"clash/compo"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"gopkg.in/yaml.v2"
)

func init() {
	env := os.Getenv("RUN_ENV")
	if env == "" {
		env = "local"
	}
	godotenv.Load(".env." + env)
}

func main() {
	appHandler := makeAppHandler()

	router := gin.Default()
	makeGroupApi(router)
	router.Use(gin.WrapH(&appHandler))

	if err := router.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}

func makeGroupApi(r *gin.Engine) {

	authMiddleware := makeAuthMiddleware()
	if errInit := authMiddleware.MiddlewareInit(); errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	api := r.Group("api")
	{
		api.POST("/auth", authMiddleware.LoginHandler)

		auth := api.Group("auth")
		auth.Use(authMiddleware.MiddlewareFunc())
		{
			auth.GET("/ping", func(c *gin.Context) {
				c.String(http.StatusOK, "pong")
			})
			auth.GET("/config", getConfigList)
			auth.GET("/config/:file", getConfig)
			auth.POST("/config/:file", setConfig)
		}
	}
}

func makeAppHandler() app.Handler {
	app.Route("/", &compo.Login{})
	app.Route("/config", &compo.List{})
	app.RouteWithRegexp(`^/config/\d+\.yaml\\?$`, &compo.Config{})
	app.RunWhenOnBrowser()

	handler := app.Handler{
		Name:  "Clash",
		Title: "Clash Dashboard",
		Styles: []string{
			"/web/tailwind.min.css",
			"/web/semantic.min.css",
		},
		Scripts: []string{
			"/web/jquery.min.js",
			"/web/semantic.min.js",
		},
	}

	return handler
}

func getConfigList(c *gin.Context) {
	var fileList []string
	var fileNameReg = regexp.MustCompile(`\d+\.yaml`)

	files, _ := ioutil.ReadDir("configs")

	for _, f := range files {
		if fileNameReg.MatchString(f.Name()) {
			fileList = append(fileList, f.Name())
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fileList,
	})
}
func getConfig(c *gin.Context) {
	fileName := c.Param("file")

	var file []byte
	var conf interface{}
	var err error

	if ExistsFile(fileName) {
		file, err = ioutil.ReadFile(fmt.Sprintf("configs/%s.yaml", fileName))
	} else {
		file, err = ioutil.ReadFile("configs/default.yaml")
	}

	if err != nil {
		log.Println(err)
		c.String(http.StatusTeapot, err.Error())
	}

	err = yaml.Unmarshal(file, &conf)

	if err != nil {
		log.Println(err)
		c.String(http.StatusTeapot, err.Error())
	}

	c.YAML(http.StatusOK, conf)
}

func setConfig(c *gin.Context) {
	fileName := c.Param("file")
	var conf interface{}

	type Body struct {
		Config string `json:"config" binding:"required"`
	}

	var jsonData Body

	if err := c.ShouldBind(&jsonData); err != nil {
		c.JSON(http.StatusTeapot, gin.H{
			"message": err.Error(),
		})
		return
	}

	fileBytes := []byte(jsonData.Config)

	if err := yaml.Unmarshal(fileBytes, &conf); err != nil {
		c.JSON(http.StatusTeapot, gin.H{
			"message": err.Error(),
		})
	} else {

		if err := ioutil.WriteFile(fmt.Sprintf("configs/%s.yaml", fileName), fileBytes, 0755); err != nil {
			c.JSON(http.StatusTeapot, gin.H{
				"message": err.Error(),
			})
		}

		c.YAML(http.StatusOK, conf)
	}
}
