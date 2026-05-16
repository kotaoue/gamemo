package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	// WikiBaseURL はスプラ3 wikiwiki の基本URL
	WikiBaseURL = "https://wikiwiki.jp/splatoon3mix"

	// UserAgent はスクレイピング時のUser-Agent
	UserAgent = "ikamemo/1.0 (https://github.com/kotaoue/gamemo)"

	// scrapeTimeout はHTTPリクエストのタイムアウト
	scrapeTimeout = 10 * time.Second
)

// weaponCategoryURLs は重量クラスごとのブキカテゴリURL
var weaponCategoryURLs = []struct {
	url    string
	weight WeightClass
}{
	{WikiBaseURL + "/ブキ/シューター", WeightLight},
	{WikiBaseURL + "/ブキ/マニューバー", WeightLight},
	{WikiBaseURL + "/ブキ/スピナー", WeightMedium},
	{WikiBaseURL + "/ブキ/リールガン", WeightMedium},
	{WikiBaseURL + "/ブキ/チャージャー", WeightHeavy},
	{WikiBaseURL + "/ブキ/スロッシャー", WeightMedium},
	{WikiBaseURL + "/ブキ/フデ", WeightLight},
	{WikiBaseURL + "/ブキ/ローラー", WeightHeavy},
	{WikiBaseURL + "/ブキ/ブラスター", WeightMedium},
	{WikiBaseURL + "/ブキ/弓", WeightMedium},
	{WikiBaseURL + "/ブキ/ワイパー", WeightLight},
}

// ScrapeWeapons はwikiwiki.jpからブキデータをスクレイピングする
// ネットワークエラーやパースエラーが起きた場合はerrorを返す
func ScrapeWeapons() (map[string]Weapon, error) {
	client := &http.Client{Timeout: scrapeTimeout}
	result := make(map[string]Weapon)

	for _, cat := range weaponCategoryURLs {
		weapons, err := scrapeCategoryPage(client, cat.url, cat.weight)
		if err != nil {
			return nil, fmt.Errorf("%s のスクレイピングに失敗: %w", cat.url, err)
		}
		for k, v := range weapons {
			result[k] = v
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("ブキデータが1件も取得できませんでした")
	}

	return result, nil
}

// scrapeCategoryPage は1つのカテゴリページからブキデータを取得する
func scrapeCategoryPage(client *http.Client, url string, defaultWeight WeightClass) (map[string]Weapon, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, url)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseWeaponTable(doc, defaultWeight), nil
}

// parseWeaponTable はHTMLドキュメントからブキのステータステーブルをパースする
// wikiwiki.jpのテーブル形式: 列ヘッダが 名前/射程/ダメージ/重量 を含む
func parseWeaponTable(doc *goquery.Document, defaultWeight WeightClass) map[string]Weapon {
	result := make(map[string]Weapon)

	// wikiwiki.jpのtableを順番に探して、射程列を持つものを使う
	doc.Find("table").Each(func(_ int, table *goquery.Selection) {
		headers, colIndex := extractHeaders(table)
		if colIndex["射程"] < 0 {
			return // 射程列がないテーブルはスキップ
		}

		table.Find("tr").Each(func(rowIdx int, row *goquery.Selection) {
			if rowIdx == 0 {
				return // ヘッダ行をスキップ
			}
			cells := row.Find("td")
			if cells.Length() == 0 {
				return
			}

			w := parseWeaponRow(cells, headers, colIndex, defaultWeight)
			if w.Name != "" && w.Range > 0 {
				result[w.Name] = w
			}
		})
	})

	return result
}

// extractHeaders はテーブルのヘッダ行からカラム名と位置を返す
func extractHeaders(table *goquery.Selection) ([]string, map[string]int) {
	colIndex := map[string]int{
		"名前": -1, "射程": -1, "ダメージ": -1, "重量": -1,
	}
	var headers []string

	table.Find("tr").First().Find("th, td").Each(func(i int, cell *goquery.Selection) {
		text := strings.TrimSpace(cell.Text())
		headers = append(headers, text)
		for key := range colIndex {
			if strings.Contains(text, key) {
				colIndex[key] = i
			}
		}
	})

	return headers, colIndex
}

// parseWeaponRow は1行のtd要素からWeaponを組み立てる
func parseWeaponRow(cells *goquery.Selection, _ []string, colIndex map[string]int, defaultWeight WeightClass) Weapon {
	getCellText := func(idx int) string {
		if idx < 0 || idx >= cells.Length() {
			return ""
		}
		return strings.TrimSpace(cells.Eq(idx).Text())
	}

	name := getCellText(colIndex["名前"])
	if name == "" && colIndex["名前"] < 0 {
		// 名前列がない場合は先頭セルを使う
		name = strings.TrimSpace(cells.Eq(0).Text())
	}
	// リンクテキストやアイコン代替テキストが混入している場合はImg alt属性を優先
	if cells.Eq(0).Find("img").Length() > 0 {
		alt := cells.Eq(0).Find("img").AttrOr("alt", "")
		if alt != "" {
			name = strings.TrimSpace(alt)
		}
	}
	if name == "" {
		return Weapon{}
	}

	rangeStr := getCellText(colIndex["射程"])
	rangeVal, err := strconv.ParseFloat(rangeStr, 64)
	if err != nil {
		return Weapon{}
	}

	damageVal := 0.0
	if colIndex["ダメージ"] >= 0 {
		damageStr := getCellText(colIndex["ダメージ"])
		// "35~42" のような形式から最初の数値を取る
		damageStr = strings.Split(damageStr, "~")[0]
		damageStr = strings.Split(damageStr, "*")[0]
		damageStr = strings.TrimSpace(damageStr)
		if v, err := strconv.ParseFloat(damageStr, 64); err == nil {
			damageVal = v
		}
	}

	weight := defaultWeight
	if colIndex["重量"] >= 0 {
		wText := getCellText(colIndex["重量"])
		switch {
		case strings.Contains(wText, "軽"):
			weight = WeightLight
		case strings.Contains(wText, "重"):
			weight = WeightHeavy
		default:
			weight = WeightMedium
		}
	}

	return Weapon{
		Name:   name,
		Range:  rangeVal,
		Damage: damageVal,
		Weight: weight,
	}
}

// LoadWeaponsWithFallback はスクレイピングを試み、失敗した場合はフォールバックデータを返す
func LoadWeaponsWithFallback() (map[string]Weapon, bool) {
	scraped, err := ScrapeWeapons()
	if err == nil && len(scraped) > 0 {
		return scraped, true
	}
	return fallbackWeapons(), false
}

// fallbackWeapons は固定値のフォールバックデータを返す
func fallbackWeapons() map[string]Weapon {
	return map[string]Weapon{
		"バレル": {
			Name:     "バレルスピナー",
			Range:    4.1,
			Damage:   30,
			Weight:   WeightMedium,
			FireRate: 15.0,
		},
		"スシ": {
			Name:     "スプラシューター",
			Range:    2.9,
			Damage:   35,
			Weight:   WeightLight,
			FireRate: 10.0,
		},
		"シャーカー": {
			Name:     "スプラッシュボム・シャーカー",
			Range:    2.4,
			Damage:   28,
			Weight:   WeightLight,
			FireRate: 11.0,
		},
		"プラコラ": {
			Name:     "フォルテスプラシューターコラボ",
			Range:    3.4,
			Damage:   42,
			Weight:   WeightMedium,
			FireRate: 6.0,
		},
		"52": {
			Name:     ".52ガロン",
			Range:    3.2,
			Damage:   52,
			Weight:   WeightMedium,
			FireRate: 5.7,
		},
		"96": {
			Name:     ".96ガロン",
			Range:    3.6,
			Damage:   62,
			Weight:   WeightHeavy,
			FireRate: 3.0,
		},
	}
}
