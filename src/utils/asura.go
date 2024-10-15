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

var RoostersSprites = struct {
	sync.Mutex
	data map[string][][]interface{}
}{
	data: make(map[string][][]interface{}),
}

func GetRarityColor(rarity int) int {
	return [...]int{13493247, 255, 9699539, 16748544, 16728128, 16777201}[rarity]
}

func GetRoostersSprites(url string) ([][]interface{}, error) {
	RoostersSprites.Lock()
	defer RoostersSprites.Unlock()

	if cachedData, exists := RoostersSprites.data[url]; exists {
		return cachedData, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error getting", url, ":", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return nil, err
	}

	var data [][]interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshalling json:", err)
		return nil, err
	}

	RoostersSprites.data[url] = data

	return data, nil
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

func GetRoostersNames(data interface{}) []Class {
	roosters := make([]Class, 0)

	for _, item := range data.([]interface{}) {
		if classMap, ok := item.(map[string]interface{}); ok {
			name := classMap["name"].(string)
			rarity := int(classMap["rarity"].(float64))
			var disadvantages []int
			if disadvantagesRaw, ok := classMap["disadvantages"].([]interface{}); ok {
				for _, d := range disadvantagesRaw {
					disadvantages = append(disadvantages, int(d.(float64)))
				}
			}
			roosters = append(roosters, Class{
				Name:          name,
				Rarity:        rarity,
				Disadvantages: disadvantages,
			})
		}
	}

	return roosters
}
