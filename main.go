package main

import (
	"net/http"

	awsnews "github.com/circa10a/go-aws-news/news"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func newsStatement(n []CarouselItem) string {
	if len(n) == 0 {
		return "No cloud computing news yet."
	}
	return "Here's the latest cloud computing news."
}

func defaultNewsItem() []CarouselItem {
	var newsItems []CarouselItem
	// Needs at minimium 2 items
	for i := 1; i <= 2; i++ {
		newsItems = append(newsItems, CarouselItem{
			Title:       "No recent news",
			Description: "Check back soon",
			OpenURLAction: OpenURLAction{
				URL: "https://aws.amazon.com/new",
			},
		})
	}
	return newsItems
}

func newsListItems() []CarouselItem {
	newsItems := make([]CarouselItem, 0)
	news, _ := awsnews.ThisMonth()

	if len(news) == 0 {
		return defaultNewsItem()
	}

	for _, newsItem := range news {
		itemInfo := CarouselItem{
			Title:       newsItem.Title,
			Description: newsItem.PostDate,
			OpenURLAction: OpenURLAction{
				URL: newsItem.Link,
			},
		}
		newsItems = append(newsItems, itemInfo)
	}
	return newsItems
}

func buildFulfillment() *Response {
	news := newsListItems()

	return &Response{
		Payload{
			Google{
				ExpectUserResponse: false,
				RichResponse: RichResponse{
					Items: []Item{
						{
							SimpleResponse: &SimpleResponse{
								TextToSpeech: newsStatement(news),
							},
						},
						{
							CarouselBrowse: &CarouselBrowse{
								Items: news,
							},
						},
					},
				},
			},
		},
	}
}

func handleWebhook(c *gin.Context) {
	c.JSON(http.StatusOK, buildFulfillment())
}

func main() {
	var err error

	r := gin.Default()
	r.POST("/webhook", handleWebhook)

	if err = r.Run(); err != nil {
		log.WithError(err).Fatal("Couldn't start server")
	}
}
