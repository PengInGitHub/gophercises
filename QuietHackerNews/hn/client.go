//Package hn implements a Hacker News client
package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiBase = "https://hacker-news.firebaseio.com/v0"
)

//Client is an API client to interact with Hacker News API
type Client struct {
	apiBase string
}

//making the Client zero value useful without forcing users to do
//something like 'client := NewClient(apiBase string)'
//in this way zero value (var client Client) has apiBase of Hacker News
func (c *Client) defaultify() {
	if c.apiBase == "" {
		c.apiBase = apiBase //global var
	}
}

//GetTopItemsID returns roughly 450 top items IDs in decreasing order
func (c *Client) GetTopItemsID() ([]int, error) {
	c.defaultify()
	resp, err := http.Get(fmt.Sprintf("%s/topstories.json", c.apiBase))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var IDs []int
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&IDs)
	if err != nil {
		return nil, err
	}
	return IDs, nil
}

//GetItems returns the item defined by the provided ID
func (c *Client) GetItems(id int) (Item, error) {
	c.defaultify()
	var item Item
	resp, err := http.Get(fmt.Sprintf("%s/item/%d.json", c.apiBase, id))
	if err != nil {
		return item, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&item)
	if err != nil {
		return item, err
	}
	return item, nil
}

//Item represents a single item returned by the HN API
//This can have a type of 'story', 'comment' or 'job'
type Item struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`

	// Only one of these should exist
	Text string `json:"text"`
	URL  string `json:"url"`
}
