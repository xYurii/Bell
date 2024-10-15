package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Rarity int

const (
	Common Rarity = iota
	Rare
	Epic
	Legendary
	Special
	Mythic
)

type CosmeticType int

const (
	Background CosmeticType = iota
	Badge
	Skin
)

func (rarity Rarity) String() string {
	return [...]string{"Comum", "Raro", "Epico", "Lendario", "Especial", "Mitico"}[rarity]
}

func (rarity Rarity) Price() int {
	return [...]int{30, 160, 500, 1200, 500, 3000}[rarity]
}

func (rarity Rarity) Color() int {
	return [...]int{13493247, 255, 9699539, 16748544, 16728128, 16777201}[rarity]
}

type Effect struct {
	Name   string `json:"name"`
	Class  int    `json:"class"`
	Type   int    `json:"type"`
	Self   bool   `json:"self"`
	Phrase string `json:"phrase"`
	Turns  int    `json:"turns"`
	Range  [2]int `json:"range"`
}

type Class struct {
	Name          string `json:"name"`
	Desc          string `json:"desc"`
	Disadvantages []int  `json:"disadvantages"`
	Rarity        Rarity `json:"rarity"`
}

type Skill struct {
	Name    string     `json:"name"`
	Damage  [2]int     `json:"damage"`
	Level   int        `json:"level"`
	Effect  [2]float64 `json:"effect"`
	Self    bool       `json:"self"`
	Evolved bool       `json:"evolved"`
}

type Cosmetic struct {
	Type         CosmeticType `json:"type"`
	Name         string       `json:"name"`
	Value        string       `json:"value"`
	ReverseValue string       `json:"reverseValue"`
	Rarity       Rarity       `json:"rarity"`
	Extra        int          `json:"extra"`
}

var Effects []*Effect
var Roosters []*Class
var Skills []([]*Skill)
var Sprites [][]string
var Cosmetics []*Cosmetic

var cache = struct {
	sync.Mutex
	data map[string][]*Class
}{
	data: make(map[string][]*Class),
}

var RoostersSprites = struct {
	sync.Mutex
	data map[string][][]string
}{
	data: make(map[string][][]string),
}

var skillsCache = struct {
	sync.Mutex
	data map[string][]*Skill
}{
	data: make(map[string][]*Skill),
}

func GetCosmetics() ([]*Cosmetic, error) {
	if len(Cosmetics) > 0 {
		return Cosmetics, nil
	}

	urls := []string{
		"https://raw.githubusercontent.com/Acnologla/asura/master/resources/galo/cosmetics.json",
		"https://raw.githubusercontent.com/Acnologla/asura/master/resources/galo/newCosmetics.json",
		"https://raw.githubusercontent.com/Acnologla/asura/master/resources/galo/skins.json",
	}

	var allCosmetics []*Cosmetic

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("erro ao fazer requisição: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("erro ao buscar cosméticos, status: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler o corpo da resposta: %v", err)
		}

		var cosmetics []*Cosmetic
		if err := json.Unmarshal(body, &cosmetics); err != nil {
			return nil, fmt.Errorf("erro ao decodificar o JSON: %v", err)
		}

		allCosmetics = append(allCosmetics, cosmetics...)
	}

	Cosmetics = allCosmetics
	fmt.Printf("Cosmetics Loaded: %d\n", len(Cosmetics))

	return Cosmetics, nil
}

func GetEffects() ([]*Effect, error) {
	if len(Effects) > 0 {
		return Effects, nil
	}

	url := "https://raw.githubusercontent.com/Acnologla/asura/master/resources/galo/effects.json"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao buscar efeitos, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o corpo da resposta: %v", err)
	}

	var effects []*Effect
	if err := json.Unmarshal(body, &effects); err != nil {
		return nil, fmt.Errorf("erro ao decodificar o JSON: %v", err)
	}

	Effects = effects

	return Effects, nil
}

func CalcDamage(min, max int, resets float64) (int, int) {
	if resets > 0 {
		min += int(float64(min) * 0.15 * resets)
		max += int(float64(max) * 0.15 * resets)
	}
	return min, max
}

func GetRoosterSkills(rooster *Class) ([]*Skill, error) {
	skillsCache.Lock()
	if cachedSkills, exists := skillsCache.data[rooster.Name]; exists {
		skillsCache.Unlock()
		return cachedSkills, nil
	}
	skillsCache.Unlock()

	url := fmt.Sprintf("https://raw.githubusercontent.com/Acnologla/asura/master/resources/galo/attacks/%s.json", rooster.Name)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	var skills []*Skill

	err = json.Unmarshal(body, &skills)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %w", err)
	}

	skillsCache.Lock()
	skillsCache.data[rooster.Name] = skills
	skillsCache.Unlock()

	return skills, nil
}

func GetRarityColor(rarity int) int {
	return [...]int{13493247, 255, 9699539, 16748544, 16728128, 16777201}[rarity]
}

func GetRoostersSprites(url string) ([][]string, error) {
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

	err = json.Unmarshal(body, &Sprites)
	if err != nil {
		fmt.Println("Error unmarshalling json:", err)
		return nil, err
	}

	RoostersSprites.data[url] = Sprites

	return Sprites, nil
}

func GetRoostersClasses(url string) (interface{}, error) {
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

	err = json.Unmarshal(body, &Roosters)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %w", err)
	}

	cache.data[url] = Roosters
	return Roosters, nil
}

func GetRoostersNames(data interface{}) []Class {
	roosters := make([]Class, 0)

	if classList, ok := data.([]*Class); ok {
		for _, class := range classList {
			roosters = append(roosters, Class{
				Name:          class.Name,
				Rarity:        class.Rarity,
				Disadvantages: class.Disadvantages,
			})
		}
	}

	return roosters
}
