package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	readability "github.com/go-shiori/go-readability"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/html"
)

type Article struct {
	Title       string
	Byline      string
	Content     string
	TextContent string
	Node        *html.Node `json:"-"`
	Length      int
	Excerpt     string
	SiteName    string
	Image       string
	Favicon     string
}

func main() {
	app := &cli.App{
		Name:  "get",
		Usage: "transform web article into simple text",
		Action: func(ctx *cli.Context) error {
			rawUrl := ctx.Args().Get(0)
			url, err := url.Parse(rawUrl)
			if err != nil {
				return err
			}
			resp, err := http.Get(rawUrl)
			if err != nil {
				log.Fatalf("failed to download %s: %v\n", rawUrl, err)
			}
			defer resp.Body.Close()

			fullArticle, err := readability.FromReader(resp.Body, url)
			if err != nil {
				log.Fatalf("failed to parse %s: %v\n", rawUrl, err)
			}
			printableArticle := Article{
				Title:       fullArticle.Title,
				Byline:      fullArticle.Byline,
				Content:     fullArticle.Content,
				TextContent: fullArticle.TextContent,
				Length:      fullArticle.Length,
				Excerpt:     fullArticle.Excerpt,
				SiteName:    fullArticle.SiteName,
				Image:       fullArticle.Image,
				Favicon:     fullArticle.Favicon,
			}
			bytes, err := json.Marshal(printableArticle)
			if err != nil {
				return err
			}
			fmt.Println(string(bytes))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
