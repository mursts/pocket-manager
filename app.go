package main

import (
	"bytes"
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

func postToSlack(item Item, postUrl string) error {

	param := struct {
		Text string `json:"text"`
	}{
		Text: fmt.Sprintf("- %s (%d)\n    %s\n\n", item.ResolvedTitle, item.ItemID, item.ResolvedURL),
	}
	paramBytes, err := json.Marshal(param)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", postUrl, bytes.NewReader(paramBytes))
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

func postSlackHandle(w http.ResponseWriter, r *http.Request) {

	config, err := NewConfig()
	if err != nil {
		log.Printf("failed to load config. %v", err)
		return
	}

	baseDate := time.Now().UTC().In(location).AddDate(0, 0, -1)
	since := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 0, 0, 0, 0, location)
	until := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 23, 59, 59, 0, location)

	client := NewClient(config.Pocket.ConsumerKey, config.Pocket.AccessToken)
	client.SetHttpClient(http.DefaultClient)

	retrieve := &RetrieveReq{
		State: "unread",
		Sort:  "newest",
		Since: since.Unix(),
	}

	var response RetrieveResponse
	if err := client.Retrieve(retrieve, &response); err != nil {
		log.Printf("failed to request to pocket. %v", err)
		return
	}

	var wg sync.WaitGroup

	for _, item := range response.List {

		added := item.AddedAt()

		log.Printf("%d\t%s\t%s\t%s", item.ItemID, item.ResolvedTitle, item.ResolvedURL, added)

		if since.After(added) || until.Before(added) {
			continue
		}

		wg.Add(1)
		go func(item Item, url string) {
			defer wg.Done()
			postToSlack(item, url)
		}(item, config.Slack.PostUrl)
	}

	wg.Wait()
}

func main() {
	http.HandleFunc("/post", postSlackHandle)
	port := os.Getenv("port")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
