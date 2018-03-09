package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	pocketAPIEndpoint = "https://getpocket.com/v3"
)

var (
	httpClient *http.Client
)

type PocketClient struct {
	ConsumerKey string
	AccessToken string
}

type RetrieveReq struct {
	State       string `json:"state,omitempty"`
	Favorite    string `json:"favorite,omitempty"`
	Tag         string `json:"tag,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Sort        string `json:"sort,omitempty"`
	DetailType  string `json:"detailType,omitempty"`
	Search      string `json:"search,omitempty"`
	Domain      string `json:"domain,omitempty"`
	Since       int64  `json:"since,omitempty"`
	Count       int    `json:"count,omitempty"`
	Offset      int    `json:"offset,omitempty"`
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
}

type RetrieveResponse struct {
	List     map[string]Item
	Status   int
	Complete int
	Since    int
}

type Item struct {
	ItemID        int    `json:"item_id,string"`
	ResolvedId    int    `json:"resolved_id,string"`
	GivenURL      string `json:"given_url"`
	ResolvedURL   string `json:"resolved_url"`
	GivenTitle    string `json:"given_title"`
	ResolvedTitle string `json:"resolved_title"`
	Favorite      int    `json:",string"`
	Status        int    `json:",string"`
	IsArticle     int    `json:"is_article,string"`
	HasImage      int    `json:"has_image,string"`
	HasVideo      int    `json:"has_video,string"`
	WordCount     int    `json:"word_count,string"`
	SortId        int    `json:"sort_id"`
	TimeAdded     string `json:"time_added"`
	TimeUpdated   string `json:"time_updated"`
	TimeRead      string `json:"time_read"`
	TimeFavorited string `json:"time_favorited"`
}

func (i *Item) AddedAt() time.Time {
	unixtime, _ := strconv.ParseInt(i.TimeAdded, 10, 64)
	return time.Unix(unixtime, 0).In(location)
}

func (c *PocketClient) SetHttpClient(client *http.Client) {
	httpClient = client
}

func (c *PocketClient) Retrieve(retrieve *RetrieveReq, res *RetrieveResponse) error {

	url := pocketAPIEndpoint + "/get"

	retrieve.ConsumerKey = c.ConsumerKey
	retrieve.AccessToken = c.AccessToken

	body, err := json.Marshal(retrieve)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("http response %d; Error=%s", resp.StatusCode, resp.Header.Get("X-Error"))
	}

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(res)
}

func NewClient(consumerKey, accessToken string) *PocketClient {
	return &PocketClient{
		ConsumerKey: consumerKey,
		AccessToken: accessToken,
	}
}
