package scraper

import (
	"context"
	"fmt"
	"log"
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
				Project:    "headlines.com.na",
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
								Project:    "headlines.com.na",
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
						Project:    "headlines.com.na",
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

					var mediaLink string

					// Done after firestore check for efficiency
					mediaContainer := e.DOM.Find("div.proradio-entrycontent")
					imgElement := mediaContainer.Find("img").First()
					imgLink, imgExists := imgElement.Attr("src")

					if imgExists {
						mediaLink = imgLink
					}

					videoElement := mediaContainer.Find("video").First()
					videoLink, videoExists := videoElement.Attr("src")

					if videoExists {
						mediaLink = videoLink
					}

					title := e.ChildText("h1.proradio-pagecaption.proradio-glitchtxt")
					content := ""

					e.DOM.Find("div.proradio-entrycontent p").Each(func(_ int, s *goquery.Selection) {
						content += s.Text() + " "
					})
					content = strings.TrimSpace(content)

					headlineChan <- internal.Headline{
						Project:    "headlines.com.na",
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
						Project:    "headlines.com.na",
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
						Project:    "headlines.com.na",
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
						Project:    "headlines.com.na",
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
						Project:    "headlines.com.na",
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
						Project:    "headlines.com.na",
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
						Project:    "headlines.com.na",
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
						Project:    "headlines.com.na",
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
						Project:    "headlines.com.na",
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

func ScrapeNbc(c *colly.Collector, headlineChan chan<- internal.Headline, wg *sync.WaitGroup, app *firebaseSDK.App, ctx context.Context) {
	defer wg.Done()

	c.OnHTML(`div#views_slideshow_cycle_teaser_section_-block_6`, func(e *colly.HTMLElement) {

		e.ForEach("div.views_slideshow_cycle_slide.views_slideshow_slide", func(_ int, el *colly.HTMLElement) {
			log.Print("Found article")
			func(el *colly.HTMLElement) {
				linkEl := el.DOM.Find("div.views-field-title a").First()
				linkToArticle, _ := linkEl.Attr("href")

				if linkToArticle == "" {
					log.Println("Warning: Empty article link found")
					return
				}

				linkToArticle = "https://nbcnews.na" + linkToArticle

				articleCollector := c.Clone()

				articleCollector.OnHTML("body", func(e *colly.HTMLElement) {
					log.Println("Found article content")

					source := "NBC"
					currentTime := time.Now()
					createdAt := currentTime.Unix()

					fbHeadline, err := firebaseUtils.GetHeadlineByField(app, ctx, "link", linkToArticle)
					if err != nil {
						log.Printf("Error getting headline from Firebase: %v", err)
					}

					if fbHeadline.Link == linkToArticle {
						log.Printf("Article already in database: %s", linkToArticle)
						return
					}

					mediaElement := e.DOM.Find("div.image-preview img").First()
					mediaLink, _ := mediaElement.Attr("src")

					fullLink := ""

					if mediaLink != "" {
						fullLink = "https://nbcnews.na" + mediaLink
					}

					title := e.DOM.Find("nav.breadcrumb ol li").Last().Text()

					title = strings.Trim(title, " ")

					if title == "" {
						title = e.ChildText("div.node-content h2")
					}

					var contentJoined string
					e.DOM.Find("div.field--name-body p").Each(func(_ int, el *goquery.Selection) {
						contentJoined += el.Text() + " "
					})

					headline := internal.Headline{
						Media:      fullLink,
						Title:      title,
						Content:    contentJoined,
						CreatedAt:  createdAt,
						Source:     source,
						Link:       linkToArticle,
						Posted:     false,
						DatePosted: 0,
						Deleted:    false,
					}

					headlineChan <- headline
					log.Printf("Headline sent to channel: %s", title)
				})

				log.Println(e.Request.AbsoluteURL(linkToArticle))
				articleCollector.Visit(e.Request.AbsoluteURL(linkToArticle))
			}(el)
		})
	})

	c.OnScraped(func(_ *colly.Response) {
		log.Println("Finished scraping NBC") // Log finish
	})

	log.Println("Visiting NBC main page") // Log main page visit
	c.Visit("https://nbcnews.na/")
}

func ScrapeBusinessExpress(c *colly.Collector, headlineChan chan<- internal.Headline, wg *sync.WaitGroup, app *firebaseSDK.App, ctx context.Context) {
	defer wg.Done()

	c.OnHTML(`div#frontpage-area_c`, func(e *colly.HTMLElement) {

		e.ForEach("div.hk-gridunit.hcolumn-1-4.hk-gridunit-size1", func(_ int, el *colly.HTMLElement) {
			log.Print("Found article")
			func(el *colly.HTMLElement) {
				linkEl := el.DOM.Find("div.hk-gridunit-bg a").First()
				linkToArticle, _ := linkEl.Attr("href")

				if linkToArticle == "" {
					log.Println("Warning: Empty article link found")
					return
				}

				articleCollector := c.Clone()

				articleCollector.OnHTML("body", func(e *colly.HTMLElement) {
					log.Println("Found article content")

					source := "Business Express"
					currentTime := time.Now()
					createdAt := currentTime.Unix()

					fbHeadline, err := firebaseUtils.GetHeadlineByField(app, ctx, "link", linkToArticle)
					if err != nil {
						log.Printf("Error getting headline from Firebase: %v", err)
					}

					if fbHeadline.Link == linkToArticle {
						log.Printf("Article already in database: %s", linkToArticle)
						return
					}

					// Modified selector to target the div with the data-image-src attribute
					mediaElement := e.DOM.Find("div.entry-featured-img-headerwrap").First()
					mediaLink, _ := mediaElement.Attr("data-image-src")

					title := e.DOM.Find("h1.loop-title.entry-title").Text()

					title = strings.Trim(title, " ")

					var contentJoined string
					e.DOM.Find("div.entry-the-content p").Each(func(_ int, el *goquery.Selection) {
						contentJoined += el.Text() + " "
					})

					headline := internal.Headline{
						Media:      mediaLink,
						Title:      title,
						Content:    contentJoined,
						CreatedAt:  createdAt,
						Source:     source,
						Link:       linkToArticle,
						Posted:     false,
						DatePosted: 0,
						Deleted:    false,
					}

					headlineChan <- headline
					log.Printf("Headline sent to channel: %s", title)
				})

				log.Println(e.Request.AbsoluteURL(linkToArticle))
				articleCollector.Visit(e.Request.AbsoluteURL(linkToArticle))
			}(el)
		})
	})

	c.OnScraped(func(_ *colly.Response) {
		log.Println("Finished scraping Business Express")
	})

	log.Println("Visiting Business Express main page")
	c.Visit("https://nambusinessexpress.com/")
}

func ScrapeNamibianSun(c *colly.Collector, headlineChan chan<- internal.Headline, wg *sync.WaitGroup, app *firebaseSDK.App, ctx context.Context) {
	defer wg.Done()

	c.OnHTML(`div.row[data-row_id="1"]`, func(e *colly.HTMLElement) {
		e.ForEach("div.col-md-8.col-xs-12", func(i int, el *colly.HTMLElement) {
			// We are interested in the first instance of this div which contains the top stories
			if i != 0 {
				return
			}

			log.Print("Found top stories section")
			storiesContainer := el.DOM.Find("div.focus-tabs div.tab-content div#tab-9201.tab-pane.fade.in.active").First()

			storiesContainer.Find("article.thumb-article").Each(func(_ int, article *goquery.Selection) {
				linkEl := article.Find("a").First()
				linkToArticle, _ := linkEl.Attr("href")

				if linkToArticle == "" {
					log.Println("Warning: Empty article link found")
					return
				}

				absoluteURL := e.Request.AbsoluteURL(linkToArticle)
				articleCollector := c.Clone()

				articleCollector.OnHTML("body", func(e *colly.HTMLElement) {
					log.Println("Found article content for:", absoluteURL)

					source := "Namibian Sun"
					currentTime := time.Now()
					createdAt := currentTime.Unix()

					fbHeadline, err := firebaseUtils.GetHeadlineByField(app, ctx, "link", absoluteURL)
					if err != nil {
						log.Printf("Error getting headline from Firebase: %v", err)
					}

					if fbHeadline.Link == absoluteURL {
						log.Printf("Article already in database: %s", absoluteURL)
						return
					}

					// Extracting the main image
					mediaElement := e.DOM.Find("img").First()
					mediaLink, _ := mediaElement.Attr("src")

					if mediaLink == "" {
						log.Println("Warning: Could not find main image for:", absoluteURL)
					}

					// Extracting the title
					title := e.DOM.Find("h1.article-title").Text()
					title = strings.TrimSpace(title)

					// Extracting the article content
					var contentJoined string
					e.DOM.Find("div.article-body div.articleBody").Each(func(_ int, el *goquery.Selection) {
						el.Find("div.articleAd").Remove() // Remove potential ads within the content
						contentJoined += el.Text() + " "
					})
					contentJoined = strings.TrimSpace(contentJoined)

					headline := internal.Headline{
						Project:    "headlines.com.na",
						Media:      mediaLink,
						Title:      title,
						Content:    contentJoined,
						CreatedAt:  createdAt,
						Source:     source,
						Link:       absoluteURL,
						Posted:     false,
						DatePosted: 0,
						Deleted:    false,
					}

					headlineChan <- headline
					log.Printf("Headline sent to channel: %s", title)
				})

				log.Println("Visiting article:", absoluteURL)
				articleCollector.Visit(absoluteURL)
			})
		})
	})

	c.OnScraped(func(_ *colly.Response) {
		log.Println("Finished scraping Namibian Sun")
	})

	log.Println("Visiting Namibian Sun main page")
	c.Visit("https://www.namibiansun.com/")
}

func ScrapeMiningAndEnergy(c *colly.Collector, headlineChan chan<- internal.Headline, wg *sync.WaitGroup, app *firebaseSDK.App, ctx context.Context) {
	defer wg.Done()

	c.OnHTML("div.jnews_module_4507_8_68af50dfed8df", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(i int, el *colly.HTMLElement) {
			func(el *colly.HTMLElement) {
				linkToArticle := el.Attr("href")

				// Skip category links and other non-article links
				if strings.Contains(linkToArticle, "/category/") || linkToArticle == "" {
					return
				}

				articleCollector := c.Clone()

				articleCollector.OnHTML("body", func(e *colly.HTMLElement) {
					source := "Mining and Energy"
					currentTime := time.Now()
					createdAt := currentTime.Unix()

					fbHeadline, _ := firebaseUtils.GetHeadlineByField(app, ctx, "link", linkToArticle)

					if fbHeadline.Link == linkToArticle {
						return
					}

					var mediaLink string

					fmt.Printf("DEBUG: Looking for images in article: %s\n", linkToArticle)

					// Try 1: jeg_featured img
					mediaElement := e.DOM.Find(".jeg_featured img").First()
					if mediaElement.Length() > 0 {
						if src, exists := mediaElement.Attr("src"); exists && src != "" {
							mediaLink = src
							fmt.Printf("DEBUG: Found image via .jeg_featured img: %s\n", mediaLink)
						} else {
							fmt.Println("DEBUG: .jeg_featured img exists but no src attribute")
						}
					} else {
						fmt.Println("DEBUG: No .jeg_featured img found")
					}

					// Try 2: jeg_post_thumbnail img
					if mediaLink == "" {
						mediaElement = e.DOM.Find(".jeg_post_thumbnail img").First()
						if mediaElement.Length() > 0 {
							if src, exists := mediaElement.Attr("src"); exists && src != "" {
								mediaLink = src
								fmt.Printf("DEBUG: Found image via .jeg_post_thumbnail img: %s\n", mediaLink)
							} else {
								fmt.Println("DEBUG: .jeg_post_thumbnail img exists but no src attribute")
							}
						} else {
							fmt.Println("DEBUG: No .jeg_post_thumbnail img found")
						}
					}

					// Try 3: article img
					if mediaLink == "" {
						mediaElement = e.DOM.Find("article img").First()
						if mediaElement.Length() > 0 {
							if src, exists := mediaElement.Attr("src"); exists && src != "" {
								mediaLink = src
								fmt.Printf("DEBUG: Found image via article img: %s\n", mediaLink)
							} else {
								fmt.Println("DEBUG: article img exists but no src attribute")
							}
						} else {
							fmt.Println("DEBUG: No article img found")
						}
					}

					// Try 4: wp-post-image
					if mediaLink == "" {
						mediaElement = e.DOM.Find(".wp-post-image").First()
						if mediaElement.Length() > 0 {
							if src, exists := mediaElement.Attr("src"); exists && src != "" {
								mediaLink = src
								fmt.Printf("DEBUG: Found image via .wp-post-image: %s\n", mediaLink)
							} else {
								fmt.Println("DEBUG: .wp-post-image exists but no src attribute")
							}
						} else {
							fmt.Println("DEBUG: No .wp-post-image found")
						}
					}

					// Try 5: any img in content areas
					if mediaLink == "" {
						mediaElement = e.DOM.Find(".entry-content img, .post-content img, .content-inner img").First()
						if mediaElement.Length() > 0 {
							if src, exists := mediaElement.Attr("src"); exists && src != "" {
								mediaLink = src
								fmt.Printf("DEBUG: Found image via content img: %s\n", mediaLink)
							} else {
								fmt.Println("DEBUG: content img exists but no src attribute")
							}
						} else {
							fmt.Println("DEBUG: No content img found")
						}
					}

					// Try 6: any img tag at all
					if mediaLink == "" {
						mediaElement = e.DOM.Find("img").First()
						if mediaElement.Length() > 0 {
							if src, exists := mediaElement.Attr("src"); exists && src != "" {
								mediaLink = src
								fmt.Printf("DEBUG: Found image via any img: %s\n", mediaLink)
							} else {
								fmt.Println("DEBUG: img exists but no src attribute")
							}
						} else {
							fmt.Println("DEBUG: No img found at all")
						}
					}

					fmt.Printf("DEBUG: Final mediaLink: '%s'\n", mediaLink)

					// Try multiple selectors for the title
					var title string
					titleElement := e.DOM.Find("h1.jeg_post_title").First()
					if titleElement.Length() > 0 {
						title = strings.TrimSpace(titleElement.Text())
						fmt.Printf("DEBUG: Found title via h1.jeg_post_title: %s\n", title)
					} else {
						titleElement = e.DOM.Find(".jeg_post_title").First()
						if titleElement.Length() > 0 {
							title = strings.TrimSpace(titleElement.Text())
							fmt.Printf("DEBUG: Found title via .jeg_post_title: %s\n", title)
						} else {
							titleElement = e.DOM.Find("h1").First()
							if titleElement.Length() > 0 {
								title = strings.TrimSpace(titleElement.Text())
								fmt.Printf("DEBUG: Found title via h1: %s\n", title)
							} else {
								title = strings.TrimSpace(e.ChildText("title"))
								fmt.Printf("DEBUG: Found title via title tag: %s\n", title)
							}
						}
					}

					content := ""
					e.DOM.Find(".content-inner p").Each(func(_ int, s *goquery.Selection) {
						content += s.Text() + " "
					})
					content = strings.TrimSpace(content)

					headline := internal.Headline{
						Project:    "headlines.com.na",
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

					headlineChan <- headline
				})

				articleCollector.Visit(e.Request.AbsoluteURL(linkToArticle))
			}(el)
		})
	})

	c.OnScraped(func(_ *colly.Response) {
		fmt.Println("Finished scraping Mining and Energy")
	})

	c.Visit("https://miningandenergy.com.na/")
}
