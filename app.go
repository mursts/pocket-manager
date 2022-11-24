package pocket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	location, _ = time.LoadLocation("Asia/Tokyo")
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func postToSlack(item Item, postURL string) error {

	param := struct {
		Text string `json:"text"`
	}{
		Text: fmt.Sprintf("- %s (%d)\n    %s\n\n", item.ResolvedTitle, item.ItemID, item.ResolvedURL),
	}
	paramBytes, err := json.Marshal(param)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", postURL, bytes.NewReader(paramBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("%d: %s", item.ItemID, string(body))

	return nil
}

func Run(ctx context.Context, m PubSubMessage) error {
	consumerKey := os.Getenv("POCKET_CONSUMER_KEY")
	accessToken := os.Getenv("POCKET_ACCESS_TOKEN")
	slackPostURL := os.Getenv("SLACK_POST_URL")

	baseDate := time.Now().UTC().In(location).AddDate(0, 0, -1)
	since := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 0, 0, 0, 0, location)
	until := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 23, 59, 59, 0, location)

	client := NewClient(consumerKey, accessToken)
	client.SetHttpClient(http.DefaultClient)

	retrieve := &RetrieveReq{
		State: "unread",
		Sort:  "newest",
		Since: since.Unix(),
	}

	var response RetrieveResponse
	if err := client.Retrieve(retrieve, &response); err != nil {
		log.Printf("failed to request to pocket. %v", err)
		return err
	}

	var wg sync.WaitGroup

	for _, item := range response.List {

		added := item.AddedAt()

		log.Printf("%d\t%s\t%s\t%s", item.ItemID, item.ResolvedTitle, item.ResolvedURL, added)

		if since.After(added) || until.Before(added) {
			continue
		}

		if slackPostURL != "" {
			wg.Add(1)
			go func(item Item, url string) {
				defer wg.Done()
				postToSlack(item, url)
			}(item, slackPostURL)
		}
	}

	wg.Wait()

	return nil
}
