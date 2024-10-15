package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Class struct {
	Name          string `json:"name"`
	Rarity        int    `json:"rarity"`
	Disadvantages []int  `json:"disadvantages,omitempty"`
}

var cache = struct {
	sync.Mutex
	data map[string]interface{}
}{
	data: make(map[string]interface{}),
}

func GetJSON(url string) (interface{}, error) {
	cache.Lock()
	defer cache.Unlock()

	if content, exists := cache.data[url]; exists {
		return content, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	var jsonData interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %w", err)
	}

	cache.data[url] = jsonData
	return jsonData, nil
}

func GetRoostersNames(data interface{}) []string {
	roosters := make([]string, 0)

	for _, item := range data.([]interface{}) {
		if classMap, ok := item.(map[string]interface{}); ok {
			name := classMap["name"].(string)
			roosters = append(roosters, name)
		}
	}

	return roosters
}
