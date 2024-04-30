package main

import (
	"embed"
	"log"
	"net/http"

	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

//go:embed i18n/localizeToml/*
var fs embed.FS

func main() {
	// new gin engine
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// apply i18n middleware
	router.Use(ginI18n.Localize(ginI18n.WithBundle(&ginI18n.BundleCfg{
		DefaultLanguage:  language.English,
		FormatBundleFile: "toml",
		AcceptLanguage:   []language.Tag{language.English, language.German, language.French},
		RootPath:         "./i18n/localizeToml/",
		UnmarshalFunc:    toml.Unmarshal,
		// After commenting this line, use defaultLoader
		// it will be loaded from the file
		Loader: &ginI18n.EmbedLoader{
			FS: fs,
		},
	})))

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, ginI18n.MustGetMessage(ctx, "welcome"))
	})

	router.GET("/:name", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, ginI18n.MustGetMessage(
			ctx,
			&i18n.LocalizeConfig{
				MessageID: "welcomeWithName",
				TemplateData: map[string]string{
					"name": ctx.Param("name"),
				},
			}))
	})

	router.GET("/foo", func(ctx *gin.Context) {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": ginI18n.MustGetMessage(ctx, "foo.bar"),
			},
		)
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
