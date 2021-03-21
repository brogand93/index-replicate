package index

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	dataSource = source{
		url: "https://svcga.com/sc/index",
	}
)

func NewClient() *Client {
	return &Client{
		Http: &http.Client{},
	}
}

// Get returns information for a given index from a configured data source
func (c *Client) Get(index string) (*Index, error) {
	indexUrl := strings.Join([]string{dataSource.url, index}, "/")
	request, err := http.NewRequest(http.MethodGet, indexUrl, nil)
	if err != nil {
		return &Index{}, err
	}

	response, err := c.Http.Do(request)
	if err != nil {
		return &Index{}, err
	}

	if response.StatusCode == http.StatusNotFound {
		return &Index{}, NotFound{
			Message: fmt.Sprintf("%s is not a valid index for this data source", index),
		}
	}

	if response.StatusCode != http.StatusOK {
		return &Index{}, DataSourceUnavailable{
			Message: fmt.Sprintf("%s data source is either unavailable or misconfigured", dataSource.url),
		}
	}

	defer response.Body.Close()

	sourceData, _ := ioutil.ReadAll(response.Body)
	return parseIndex(sourceData)
}

func parseIndex(sourceData []byte) (*Index, error) {
	var index Index
	if err := json.Unmarshal(sourceData, &index); err != nil {
		return &Index{}, DataSourceChanged{
			Message: "it looks like the data source format has changed",
		}
	}
	return &index, nil
}

func (e DataSourceChanged) Error() string {
	return e.Message
}

func (e NotFound) Error() string {
	return e.Message
}

func (e DataSourceUnavailable) Error() string {
	return e.Message
}
