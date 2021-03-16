package index

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type source struct {
	url string
}

var (
	dataSource = source{
		url: "https://svcga.com/sc/index",
	}
)

// Get returns information for a given index from a configured data source
func Get(index string) (*Index, error) {
	indexUrl := strings.Join([]string{dataSource.url, index}, "/")
	request, err := http.NewRequest(http.MethodGet, indexUrl, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, IndexNotFound{
			Message: fmt.Sprintf("%s is not a valid index for this data source", index),
		}
	}

	if response.StatusCode != http.StatusOK {
		return nil, DataSourceUnavailable{
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

func (e IndexNotFound) Error() string {
	return e.Message
}

func (e DataSourceUnavailable) Error() string {
	return e.Message
}
