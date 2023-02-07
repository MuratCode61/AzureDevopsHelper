package htmlparser

import (
	"azure_devops_helper/models"
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func ParseHtmlText(htmlText string, imageFileNamePrefix string) models.HtmlContent {
	tokenizer := html.NewTokenizer(strings.NewReader(htmlText))
	token := tokenizer.Token()

	var textBuffer bytes.Buffer
	var images []models.Image
	imageId := 0

parseLoop:
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.TextToken:
			textBuffer.Write(tokenizer.Text())
		case html.StartTagToken:
			token = tokenizer.Token()
			tag := token.Data
			if tag == "img" {
				for _, attr := range token.Attr {
					if attr.Key == "src" {
						imageUrl := attr.Val
						imageId++
						fileName := fmt.Sprintf("%s_%d.png", imageFileNamePrefix, imageId)
						textBuffer.WriteString(fmt.Sprintf("look at the image file named: %s \n", fileName))
						images = append(images, models.Image{Name: fileName, Url: imageUrl})
					}
				}
			}

			if tag == "br" {
				textBuffer.WriteString("\n")
			}

		case html.ErrorToken:
			break parseLoop
		}
	}

	return models.HtmlContent{Text: textBuffer.String(), Images: images}
}
