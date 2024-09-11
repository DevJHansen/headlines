package scraper

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"time"

	firebaseSDK "firebase.google.com/go"
	"github.com/DevJHansen/headlines/internal"
	"github.com/PuerkitoBio/goquery"

	firebaseUtils "github.com/DevJHansen/headlines/pkg/firebase"
	"github.com/DevJHansen/headlines/pkg/utils"
	"github.com/gocolly/colly"
)

func ScrapeTheNamibian(c *colly.Collector, headlineChan chan<- internal.Headline, wg *sync.WaitGroup, app *firebaseSDK.App, ctx context.Context) {
	defer wg.Done()
	c.OnHTML("div.ps-container.ps-mx-auto.wp-block-ps-main-story", func(e *colly.HTMLElement) {

		articleLinkElement := e.DOM.Find("a").First()
		linkToArticle, exists := articleLinkElement.Attr("href")

		if !exists {
			return
		}

		articleCollector := c.Clone()

		// Open article link
		articleCollector.OnHTML("main.wp-block-group", func(e *colly.HTMLElement) {
			source := "The Namibian"
			currentTime := time.Now()
			createdAt := currentTime.Unix()

			fbHeadline, _ := firebaseUtils.GetHeadlineByField(app, ctx, "link", linkToArticle)

			if fbHeadline.Link == linkToArticle {
				return
			}

			// Done after firestore check for efficiency
			mediaContainer := e.DOM.Find("figure.wp-block-ps-post-featured-image")
			mediaElement := mediaContainer.Find("img.attachment-post-thumbnail").First()
			mediaLink, _ := mediaElement.Attr("src")

			title := e.ChildText("h1.nam_title")
			content := ""

			e.DOM.Find("div.entry-content.post_content.wp-block-post-content.is-layout-flow.wp-block-post-content-is-layout-flow p").Each(func(_ int, s *goquery.Selection) {
				content += s.Text()
			})
			content = strings.TrimSpace(content)

			headlineChan <- internal.Headline{
				Media:      mediaLink,
				Title:      title,
				Content:    content,
				CreatedAt:  createdAt,
				Source:     source,
				Link:       linkToArticle,
				Posted:     false,
				DatePosted: 0,
				Deleted:    false,
			}
		})

		articleCollector.Visit(e.Request.AbsoluteURL(linkToArticle))
	})

	c.OnHTML("main.wp-block-group.nam-main-wrap", func(e *colly.HTMLElement) {
		e.ForEach("div.ps-container.ps-mx-auto.ps-py-6.ps-px-4.wp-block-ps-post-category", func(_ int, el *colly.HTMLElement) {
			func(el *colly.HTMLElement) {

				sectionTitleEl := el.DOM.Find("h2").First()
				fmtSectionTitle := utils.TrimWhiteSpace(sectionTitleEl.Text())

				fmt.Println(fmtSectionTitle)

				if fmtSectionTitle != "More Top Stories" && fmtSectionTitle != "Politics" && fmtSectionTitle != "Business" && fmtSectionTitle != "Energy Centre" {
					return
				}

				el.ForEach("article", func(_ int, articleEl *colly.HTMLElement) {
					linkToArticle := articleEl.ChildAttr("a", "href")

					fmt.Println(linkToArticle)

					if linkToArticle == "" {
						return
					}

					articleCollector := c.Clone()

					// Open article link
					articleCollector.OnHTML("main.wp-block-group", func(e *colly.HTMLElement) {
						source := "The Namibian"
						currentTime := time.Now()
						createdAt := currentTime.Unix()

						fbHeadline, _ := firebaseUtils.GetHeadlineByField(app, ctx, "link", linkToArticle)

						if fbHeadline.Link == linkToArticle {
							return
						}

						mediaContainer := e.DOM.Find("figure.wp-block-ps-post-featured-image")
						if mediaContainer.Length() > 0 {
							mediaElement := mediaContainer.Find("img.attachment-post-thumbnail").First()
							mediaLink, _ := mediaElement.Attr("src")

							title := e.ChildText("h1.nam_title")
							content := ""

							e.DOM.Find("div.entry-content.post_content.wp-block-post-content.is-layout-flow.wp-block-post-content-is-layout-flow p").Each(func(_ int, s *goquery.Selection) {
								content += s.Text()
							})
							content = strings.TrimSpace(content)

							headlineChan <- internal.Headline{
								Media:      mediaLink,
								Title:      title,
								Content:    content,
								CreatedAt:  createdAt,
								Source:     source,
								Link:       linkToArticle,
								Posted:     false,
								DatePosted: 0,
								Deleted:    false,
							}
						}
					})

					articleCollector.Visit(e.Request.AbsoluteURL(linkToArticle))
				})
			}(el)
		})

	})

	c.OnScraped(func(_ *colly.Response) {
		fmt.Println("Finished scraping The Namibian")
	})

	c.Visit("https://www.namibian.com.na/")
}

func ScrapeTheBrief(c *colly.Collector, headlineChan chan<- internal.Headline, wg *sync.WaitGroup, app *firebaseSDK.App, ctx context.Context) {
	defer wg.Done()
	c.OnHTML("div.jeg_news_ticker_items", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {

			func(el *colly.HTMLElement) {

				linkToArticle := el.Attr("href")

				if linkToArticle == "" {
					return
				}

				articleCollector := c.Clone()

				// Open article link
				articleCollector.OnHTML("div.jeg_content.jeg_singlepage", func(e *colly.HTMLElement) {
					source := "The Brief"
					currentTime := time.Now()
					createdAt := currentTime.Unix()

					fbHeadline, _ := firebaseUtils.GetHeadlineByField(app, ctx, "link", linkToArticle)

					if fbHeadline.Link == linkToArticle {
						return
					}

					// Done after firestore check for efficiency
					mediaContainer := e.DOM.Find("figure.wp-block-image")
					mediaElement := mediaContainer.Find("img").First()
					mediaLink, _ := mediaElement.Attr("src")

					title := e.ChildText("h1.jeg_post_title")
					content := ""

					e.DOM.Find("div.entry-content.no-share p").Each(func(_ int, s *goquery.Selection) {
						content += s.Text() + " "
					})
					content = strings.TrimSpace(content)

					headlineChan <- internal.Headline{
						Media:      mediaLink,
						Title:      title,
						Content:    content,
						CreatedAt:  createdAt,
						Source:     source,
						Link:       linkToArticle,
						Posted:     false,
						DatePosted: 0,
						Deleted:    false,
					}
				})

				articleCollector.Visit(e.Request.AbsoluteURL(linkToArticle))
			}(el)
		})
	})

	c.OnScraped(func(_ *colly.Response) {
		fmt.Println("Finished scraping The Brief")
	})

	c.Visit("https://thebrief.com.na/")
}

func ScrapeFutureMedia(c *colly.Collector, headlineChan chan<- internal.Headline, wg *sync.WaitGroup, app *firebaseSDK.App, ctx context.Context) {
	defer wg.Done()
	c.OnHTML("div.proradio-col.proradio-s12.proradio-m12.proradio-l8", func(e *colly.HTMLElement) {
		e.ForEach("article", func(_ int, el *colly.HTMLElement) {

			func(e *colly.HTMLElement) {
				linkParentEl := el.DOM.Find("h3.proradio-post__title.proradio-h2").First()
				linkEl := linkParentEl.Find("a")
				linkToArticle, _ := linkEl.Attr("href")

				if linkToArticle == "" {
					return
				}

				articleCollector := c.Clone()

				// Open article link
				articleCollector.OnHTML("div#proradio-pagecontent", func(e *colly.HTMLElement) {
					source := "Future Media News"
					currentTime := time.Now()
					createdAt := currentTime.Unix()

					fbHeadline, _ := firebaseUtils.GetHeadlineByField(app, ctx, "link", linkToArticle)

					if fbHeadline.Link == linkToArticle {
						return
					}

					// Done after firestore check for efficiency
					mediaContainer := e.DOM.Find("div.proradio-entrycontent")
					mediaElement := mediaContainer.Find("img").First()
					mediaLink, _ := mediaElement.Attr("src")

					title := e.ChildText("h1.proradio-pagecaption.proradio-glitchtxt")
					content := ""

					e.DOM.Find("div.proradio-entrycontent p").Each(func(_ int, s *goquery.Selection) {
						content += s.Text() + " "
					})
					content = strings.TrimSpace(content)

					headlineChan <- internal.Headline{
						Media:      mediaLink,
						Title:      title,
						Content:    content,
						CreatedAt:  createdAt,
						Source:     source,
						Link:       linkToArticle,
						Posted:     false,
						DatePosted: 0,
						Deleted:    false,
					}
				})

				articleCollector.Visit(e.Request.AbsoluteURL(linkToArticle))
			}(el) // This is to immediately invoke the function and passing el as a param
		})
	})

	c.OnScraped(func(_ *colly.Response) {
		fmt.Println("Finished scraping Future Media")
	})

	c.Visit("https://futuremedianews.com.na/category/namibia/")
	c.Visit("https://futuremedianews.com.na/category/business-economics/")
}

func ScrapeOilAndGas(c *colly.Collector, headlineChan chan<- internal.Headline, wg *sync.WaitGroup, app *firebaseSDK.App, ctx context.Context) {
	defer wg.Done()
	c.OnHTML("div#tan-main-banner-latest-trending-popular-popular", func(e *colly.HTMLElement) {
		e.ForEach("div.small-post", func(_ int, el *colly.HTMLElement) {

			func(el *colly.HTMLElement) {
				linkParentEl := el.DOM.Find("h5.title").First()
				linkEl := linkParentEl.Find("a")
				linkToArticle, _ := linkEl.Attr("href")

				if linkToArticle == "" {
					return
				}

				articleCollector := c.Clone()

				// Open article link
				articleCollector.OnHTML("div.mg-blog-post-box", func(e *colly.HTMLElement) {
					source := "Namibia Oil and Gas"
					currentTime := time.Now()
					createdAt := currentTime.Unix()

					fbHeadline, _ := firebaseUtils.GetHeadlineByField(app, ctx, "link", linkToArticle)

					if fbHeadline.Link == linkToArticle {
						return
					}

					// Done after firestore check for efficiency
					mediaElement := e.DOM.Find("img.img-fluid.wp-post-image").First()
					mediaLink, _ := mediaElement.Attr("src")

					title := e.ChildText("h1.title")
					content := ""

					e.DOM.Find("article.small p").Each(func(_ int, s *goquery.Selection) {
						content += s.Text() + " "
					})
					content = strings.TrimSpace(content)

					headlineChan <- internal.Headline{
						Media:      mediaLink,
						Title:      title,
						Content:    content,
						CreatedAt:  createdAt,
						Source:     source,
						Link:       linkToArticle,
						Posted:     false,
						DatePosted: 0,
						Deleted:    false,
					}
				})

				articleCollector.Visit(e.Request.AbsoluteURL(linkToArticle))
			}(el) // This is to immediately invoke the function and passing el as a param
		})
	})

	c.OnHTML("div#tan-main-banner-latest-trending-popular-recent", func(e *colly.HTMLElement) {
		e.ForEach("div.small-post", func(_ int, el *colly.HTMLElement) {

			func(el *colly.HTMLElement) {
				linkParentEl := el.DOM.Find("h5.title").First()
				linkEl := linkParentEl.Find("a")
				linkToArticle, _ := linkEl.Attr("href")

				if linkToArticle == "" {
					return
				}

				articleCollector := c.Clone()

				// Open article link
				articleCollector.OnHTML("div.mg-blog-post-box", func(e *colly.HTMLElement) {
					source := "Namibia Oil and Gas"
					currentTime := time.Now()
					createdAt := currentTime.Unix()

					fbHeadline, _ := firebaseUtils.GetHeadlineByField(app, ctx, "link", linkToArticle)

					if fbHeadline.Link == linkToArticle {
						return
					}

					// Done after firestore check for efficiency
					mediaElement := e.DOM.Find("img.img-fluid.wp-post-image").First()
					mediaLink, _ := mediaElement.Attr("src")

					title := e.ChildText("h1.title")
					content := ""

					e.DOM.Find("article.small p").Each(func(_ int, s *goquery.Selection) {
						content += s.Text() + " "
					})
					content = strings.TrimSpace(content)

					headlineChan <- internal.Headline{
						Media:      mediaLink,
						Title:      title,
						Content:    content,
						CreatedAt:  createdAt,
						Source:     source,
						Link:       linkToArticle,
						Posted:     false,
						DatePosted: 0,
						Deleted:    false,
					}
				})

				articleCollector.Visit(e.Request.AbsoluteURL(linkToArticle))
			}(el) // This is to immediately invoke the function and passing el as a param
		})
	})

	c.OnHTML("div#tan-main-banner-latest-trending-popular-categorised", func(e *colly.HTMLElement) {
		e.ForEach("div.small-post", func(_ int, el *colly.HTMLElement) {

			func(el *colly.HTMLElement) {
				linkParentEl := el.DOM.Find("h5.title").First()
				linkEl := linkParentEl.Find("a")
				linkToArticle, _ := linkEl.Attr("href")

				if linkToArticle == "" {
					return
				}

				articleCollector := c.Clone()

				// Open article link
				articleCollector.OnHTML("div.mg-blog-post-box", func(e *colly.HTMLElement) {
					source := "Namibia Oil and Gas"
					currentTime := time.Now()
					createdAt := currentTime.Unix()

					fbHeadline, _ := firebaseUtils.GetHeadlineByField(app, ctx, "link", linkToArticle)

					if fbHeadline.Link == linkToArticle {
						return
					}

					// Done after firestore check for efficiency
					mediaElement := e.DOM.Find("img.img-fluid.wp-post-image").First()
					mediaLink, _ := mediaElement.Attr("src")

					title := e.ChildText("h1.title")
					content := ""

					e.DOM.Find("article.small p").Each(func(_ int, s *goquery.Selection) {
						content += s.Text() + " "
					})
					content = strings.TrimSpace(content)

					headlineChan <- internal.Headline{
						Media:      mediaLink,
						Title:      title,
						Content:    content,
						CreatedAt:  createdAt,
						Source:     source,
						Link:       linkToArticle,
						Posted:     false,
						DatePosted: 0,
						Deleted:    false,
					}
				})

				articleCollector.Visit(e.Request.AbsoluteURL(linkToArticle))
			}(el) // This is to immediately invoke the function and passing el as a param
		})
	})

	c.OnScraped(func(_ *colly.Response) {
		fmt.Println("Finished scraping Namibia oil and gas")
	})

	c.Visit("https://namibiaoilandgas.com/")
}

func ScrapeNewEra(c *colly.Collector, headlineChan chan<- internal.Headline, wg *sync.WaitGroup, app *firebaseSDK.App, ctx context.Context) {
	defer wg.Done()
	c.OnHTML("div#cmsmasters_column_ec66bbce2e", func(e *colly.HTMLElement) {
		e.ForEach("div.cmsmasters_post_cont", func(_ int, el *colly.HTMLElement) {

			func(el *colly.HTMLElement) {
				linkParentEl := el.DOM.Find("h3.entry-title").First()
				linkEl := linkParentEl.Find("a")
				linkToArticle, _ := linkEl.Attr("href")

				if linkToArticle == "" {
					return
				}

				articleCollector := c.Clone()

				// Open article link
				articleCollector.OnHTML("article.cmsmasters_open_post", func(e *colly.HTMLElement) {
					source := "New Era"
					currentTime := time.Now()
					createdAt := currentTime.Unix()

					fbHeadline, _ := firebaseUtils.GetHeadlineByField(app, ctx, "link", linkToArticle)

					if fbHeadline.Link == linkToArticle {
						return
					}

					// Done after firestore check for efficiency
					mediaParentElement := e.DOM.Find("figure.cmsmasters_img_wrap").First()
					mediaElement := mediaParentElement.Find("a")
					mediaLink, _ := mediaElement.Attr("href")

					title := e.ChildText("h2.entry-title")
					content := ""

					e.DOM.Find("div.cmsmasters_post_content.entry-content p").Each(func(_ int, s *goquery.Selection) {
						content += s.Text() + " "
					})
					content = strings.TrimSpace(content)

					headlineChan <- internal.Headline{
						Media:      mediaLink,
						Title:      title,
						Content:    content,
						CreatedAt:  createdAt,
						Source:     source,
						Link:       linkToArticle,
						Posted:     false,
						DatePosted: 0,
						Deleted:    false,
					}
				})

				articleCollector.Visit(e.Request.AbsoluteURL(linkToArticle))
			}(el) // This is to immediately invoke the function and passing el as a param
		})
	})
	c.OnHTML("div#blog_eb4f7a9570", func(e *colly.HTMLElement) {
		e.ForEach("article", func(_ int, el *colly.HTMLElement) {

			func(el *colly.HTMLElement) {
				linkParentEl := el.DOM.Find("h3.entry-title").First()
				linkEl := linkParentEl.Find("a")
				linkToArticle, _ := linkEl.Attr("href")

				if linkToArticle == "" {
					return
				}

				articleCollector := c.Clone()

				// Open article link
				articleCollector.OnHTML("article.cmsmasters_open_post", func(e *colly.HTMLElement) {
					source := "New Era"
					currentTime := time.Now()
					createdAt := currentTime.Unix()

					fbHeadline, _ := firebaseUtils.GetHeadlineByField(app, ctx, "link", linkToArticle)

					if fbHeadline.Link == linkToArticle {
						return
					}

					// Done after firestore check for efficiency
					mediaParentElement := e.DOM.Find("figure.cmsmasters_img_wrap").First()
					mediaElement := mediaParentElement.Find("a")
					mediaLink, _ := mediaElement.Attr("href")

					title := e.ChildText("h2.entry-title")
					content := ""

					e.DOM.Find("div.cmsmasters_post_content.entry-content p").Each(func(_ int, s *goquery.Selection) {
						content += s.Text() + " "
					})
					content = strings.TrimSpace(content)

					headlineChan <- internal.Headline{
						Media:      mediaLink,
						Title:      title,
						Content:    content,
						CreatedAt:  createdAt,
						Source:     source,
						Link:       linkToArticle,
						Posted:     false,
						DatePosted: 0,
						Deleted:    false,
					}
				})

				articleCollector.Visit(e.Request.AbsoluteURL(linkToArticle))
			}(el) // This is to immediately invoke the function and passing el as a param
		})
	})
	c.OnHTML("div#blog_3b7nl6sigc", func(e *colly.HTMLElement) {
		e.ForEach("article", func(_ int, el *colly.HTMLElement) {

			func(el *colly.HTMLElement) {
				linkParentEl := el.DOM.Find("h3.entry-title").First()
				linkEl := linkParentEl.Find("a")
				linkToArticle, _ := linkEl.Attr("href")

				if linkToArticle == "" {
					return
				}

				articleCollector := c.Clone()

				// Open article link
				articleCollector.OnHTML("article.cmsmasters_open_post", func(e *colly.HTMLElement) {
					source := "New Era"
					currentTime := time.Now()
					createdAt := currentTime.Unix()

					fbHeadline, _ := firebaseUtils.GetHeadlineByField(app, ctx, "link", linkToArticle)

					if fbHeadline.Link == linkToArticle {
						return
					}

					// Done after firestore check for efficiency
					mediaParentElement := e.DOM.Find("figure.cmsmasters_img_wrap").First()
					mediaElement := mediaParentElement.Find("a")
					mediaLink, _ := mediaElement.Attr("href")

					title := e.ChildText("h2.entry-title")
					content := ""

					e.DOM.Find("div.cmsmasters_post_content.entry-content p").Each(func(_ int, s *goquery.Selection) {
						content += s.Text() + " "
					})
					content = strings.TrimSpace(content)

					headlineChan <- internal.Headline{
						Media:      mediaLink,
						Title:      title,
						Content:    content,
						CreatedAt:  createdAt,
						Source:     source,
						Link:       linkToArticle,
						Posted:     false,
						DatePosted: 0,
						Deleted:    false,
					}
				})

				articleCollector.Visit(e.Request.AbsoluteURL(linkToArticle))
			}(el) // This is to immediately invoke the function and passing el as a param
		})
	})

	c.OnScraped(func(_ *colly.Response) {
		fmt.Println("Finished scraping New Era")
	})

	c.Visit("https://neweralive.na/")
}

func ScrapeInformante(c *colly.Collector, headlineChan chan<- internal.Headline, wg *sync.WaitGroup, app *firebaseSDK.App, ctx context.Context) {
	defer wg.Done()
	c.OnHTML("ul.ultp-news-ticker", func(e *colly.HTMLElement) {
		e.ForEach("div.ultp-list-box", func(_ int, el *colly.HTMLElement) {

			func(el *colly.HTMLElement) {
				linkEl := el.DOM.Find("a").First()
				linkToArticle, _ := linkEl.Attr("href")

				if linkToArticle == "" {
					return
				}

				articleCollector := c.Clone()

				// Open article link
				articleCollector.OnHTML("body.post-template-default ", func(e *colly.HTMLElement) {
					source := "Informante"
					currentTime := time.Now()
					createdAt := currentTime.Unix()

					fbHeadline, _ := firebaseUtils.GetHeadlineByField(app, ctx, "link", linkToArticle)

					if fbHeadline.Link == linkToArticle {
						return
					}

					mediaElement := e.DOM.Find("div.post-layout1").First()
					mediaLink, _ := mediaElement.Attr("style")
					parsedMediaLink := utils.GetImgUrlFromStyleAtr(mediaLink)

					title := e.ChildText("h1.entry-title")
					content := ""

					e.DOM.Find("div.entry-content p").Each(func(_ int, s *goquery.Selection) {
						content += s.Text() + " "
					})
					content = strings.TrimSpace(content)

					headlineChan <- internal.Headline{
						Media:      parsedMediaLink,
						Title:      title,
						Content:    content,
						CreatedAt:  createdAt,
						Source:     source,
						Link:       linkToArticle,
						Posted:     false,
						DatePosted: 0,
						Deleted:    false,
					}
				})

				articleCollector.Visit(e.Request.AbsoluteURL(linkToArticle))
			}(el) // This is to immediately invoke the function and passing el as a param
		})
	})

	c.OnScraped(func(_ *colly.Response) {
		fmt.Println("Finished scraping Informante")
	})

	c.Visit("https://informante.web.na/")
}

func ScrapeRepublikein(c *colly.Collector, headlineChan chan<- internal.Headline, wg *sync.WaitGroup, app *firebaseSDK.App, ctx context.Context) {
	defer wg.Done()
	c.OnHTML(`[data-widget-id="9094"]`, func(e *colly.HTMLElement) {
		e.ForEach("h4.article-title", func(_ int, el *colly.HTMLElement) {

			func(el *colly.HTMLElement) {
				linkEl := el.DOM.Find("a").First()
				linkToArticle, _ := linkEl.Attr("href")

				if linkToArticle == "" {
					return
				}

				linkToArticle = "https://www.republikein.com.na" + linkToArticle

				articleCollector := c.Clone()

				// Open article link
				articleCollector.OnHTML("article.article.article-post ", func(e *colly.HTMLElement) {
					source := "Republikein"
					currentTime := time.Now()
					createdAt := currentTime.Unix()

					fbHeadline, _ := firebaseUtils.GetHeadlineByField(app, ctx, "link", linkToArticle)

					if fbHeadline.Link == linkToArticle {
						return
					}

					mediaElement := e.DOM.Find("a.fancybox").First()
					mediaLink, _ := mediaElement.Attr("href")

					title := e.ChildText("h1.article-title")
					content := e.ChildText("div.articleBody")

					headlineChan <- internal.Headline{
						Media:      mediaLink,
						Title:      title,
						Content:    content,
						CreatedAt:  createdAt,
						Source:     source,
						Link:       linkToArticle,
						Posted:     false,
						DatePosted: 0,
						Deleted:    false,
					}
				})

				fmt.Println(e.Request.AbsoluteURL(linkToArticle))
				articleCollector.Visit(e.Request.AbsoluteURL(linkToArticle))
			}(el) // This is to immediately invoke the function and passing el as a param
		})
	})

	c.OnScraped(func(_ *colly.Response) {
		fmt.Println("Finished scraping Republikein")
	})

	c.Visit("https://www.republikein.com.na/")
}
