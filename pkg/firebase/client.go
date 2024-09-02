package firebase

import (
	"context"
	"log"
	"os"

	"github.com/DevJHansen/headlines/internal"

	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func NewFirebaseApp(ctx context.Context) (*firebase.App, error) {
	credFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credFile == "" {
		credFile = "./firebase-sa.json"
	}

	opt := option.WithCredentialsFile(credFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
		return nil, err
	}
	return app, nil
}

func GetHeadline(app *firebase.App, ctx context.Context) ([]internal.Headline, error) {
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	headlines := make([]internal.Headline, 0)
	iter := client.Collection(internal.HEADLINE_COLLECTION).Documents(ctx)

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		var headline internal.Headline
		if err := doc.DataTo(&headline); err != nil {
			log.Printf("Error decoding document: %v", err)
			continue
		}

		headlines = append(headlines, headline)
	}

	return headlines, nil
}

func GetHeadlineByField(app *firebase.App, ctx context.Context, field string, value string) (internal.Headline, error) {
	client, err := app.Firestore(ctx)
	if err != nil {
		return internal.Headline{}, err
	}
	defer client.Close()

	headlines := make([]internal.Headline, 0)
	iter := client.Collection(internal.HEADLINE_COLLECTION).Where(field, "==", value).Limit(1).Documents(ctx)

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return internal.Headline{}, err
		}

		var headline internal.Headline
		if err := doc.DataTo(&headline); err != nil {
			log.Printf("Error decoding document: %v", err)
			continue
		}

		headlines = append(headlines, headline)
	}

	if len(headlines) == 0 {
		return internal.Headline{}, nil
	}

	return headlines[0], nil
}

func AddHeadline(app *firebase.App, ctx context.Context, headline internal.Headline) error {
	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	data := map[string]interface{}{
		"title":      headline.Title,
		"content":    headline.Content,
		"source":     headline.Source,
		"link":       headline.Link,
		"media":      headline.Media,
		"createdAt":  int64(headline.CreatedAt),
		"posted":     headline.Posted,
		"datePosted": int64(headline.DatePosted),
		"deleted":    headline.Deleted,
	}

	_, _, err = client.Collection(internal.HEADLINE_COLLECTION).Add(ctx, data)
	if err != nil {
		return err
	}

	return nil
}
