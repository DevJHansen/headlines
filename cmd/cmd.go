package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	firebasePkg "firebase.google.com/go"
	"github.com/DevJHansen/headlines/internal"

	"github.com/DevJHansen/headlines/pkg/firebase"
	"github.com/DevJHansen/headlines/pkg/scraper"
	"github.com/gocolly/colly"
)

func Handler(w http.ResponseWriter, r *http.Request, ctx context.Context, app *firebasePkg.App) {

	headlinesChannel := make(chan internal.Headline, 100)
	scrapedHeadlines := make([]internal.Headline, 0)

	var wg sync.WaitGroup

	go func() {
		for headline := range headlinesChannel {
			scrapedHeadlines = append(scrapedHeadlines, headline)
		}
	}()

	wg.Add(1)
	go scraper.ScrapeTheNamibian(colly.NewCollector(), headlinesChannel, &wg, app, ctx)

	wg.Add(1)
	go scraper.ScrapeTheBrief(colly.NewCollector(), headlinesChannel, &wg, app, ctx)

	wg.Add(1)
	go scraper.ScrapeFutureMedia(colly.NewCollector(), headlinesChannel, &wg, app, ctx)

	wg.Add(1)
	go scraper.ScrapeOilAndGas(colly.NewCollector(), headlinesChannel, &wg, app, ctx)

	wg.Add(1)
	go scraper.ScrapeNewEra(colly.NewCollector(), headlinesChannel, &wg, app, ctx)

	wg.Add(1)
	go scraper.ScrapeInformante(colly.NewCollector(), headlinesChannel, &wg, app, ctx)

	// wg.Add(1)
	// go scraper.ScrapeRepublikein(colly.NewCollector(), headlinesChannel, &wg, app, ctx)

	wg.Wait()
	close(headlinesChannel)

	jsonItems, _ := json.MarshalIndent(scrapedHeadlines, "", "    ")
	fmt.Println("Collected Items:")
	fmt.Println(string(jsonItems))

	if app != nil {
		for _, headline := range scrapedHeadlines {
			wg.Add(1)

			go func(sp internal.Headline) {
				defer wg.Done()

				err := firebase.AddHeadline(app, context.Background(), sp)
				if err != nil {
					fmt.Println(err)
				}
			}(headline)
		}
	}

	wg.Wait()

	fmt.Fprintln(w, "Scraper ran successfully")
}
