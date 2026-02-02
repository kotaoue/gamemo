package main

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Character struct {
	GroupName     string
	GroupNameKana string
	Name          string
	NameKana      string
	MCName        string
	Punchline     string
	Birthday      string
	Age           string
	Height        string
	Weight        string
	BloodType     string
	Occupation    string
}

func ScrapeCharacters(url string) ([]Character, error) {
	// User-Agentを設定してリクエストを作成
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "HypnosismicWikiScraper/1.0 (Educational purpose)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var characters []Character

	// h3とdtを順に処理
	doc.Find("h3").Each(func(i int, h3 *goquery.Selection) {
		headingText := h3.Text()
		groupRegex := regexp.MustCompile(`^(.+?)（(.+?)）`)
		if matches := groupRegex.FindStringSubmatch(headingText); len(matches) >= 3 {
			currentGroup := strings.TrimSpace(matches[1])
			currentGroupKana := strings.TrimSpace(matches[2])

			// このh3の後から次のh3の前までのdt要素を探す
			h3.Parent().NextUntil("div.mw-heading").Each(func(j int, sibling *goquery.Selection) {
				sibling.Find("dt").Each(func(k int, dt *goquery.Selection) {
					dtText := dt.Text()
					if !strings.Contains(dtText, "【") {
						return
					}

					// dt の後の dd 要素を集める
					var ddTexts []string
					for next := dt.Next(); next.Length() > 0 && goquery.NodeName(next) == "dd"; next = next.Next() {
						ddTexts = append(ddTexts, next.Text())
					}

					char := parseCharacterFromDt(dtText, ddTexts, currentGroup, currentGroupKana)
					if char.Name != "" {
						characters = append(characters, char)
					}
				})
			})
		}
	})

	return characters, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func parseCharacterFromDt(dtText string, ddTexts []string, groupName, groupKana string) Character {
	char := Character{
		GroupName:     groupName,
		GroupNameKana: groupKana,
	}

	// 名前とMC名の抽出（例: "山田 一郎（やまだ いちろう）【MC.B.B/エムシービッグブラザー】"）
	nameRegex := regexp.MustCompile(`^(.+?)（(.+?)）【(.+?)】`)
	if matches := nameRegex.FindStringSubmatch(dtText); len(matches) > 3 {
		char.Name = strings.TrimSpace(matches[1])
		char.NameKana = strings.TrimSpace(matches[2])
		char.MCName = strings.TrimSpace(matches[3])
	} else {
		return char
	}

	// dd要素からパンチラインとプロフィールを抽出
	for _, ddText := range ddTexts {
		// パンチライン（名言）の抽出（引用符で囲まれた部分）
		if strings.Contains(ddText, "「") && char.Punchline == "" {
			punchlineRegex := regexp.MustCompile(`「(.+?)」`)
			if matches := punchlineRegex.FindStringSubmatch(ddText); len(matches) > 1 {
				char.Punchline = strings.TrimSpace(matches[1])
			}
		}

		// プロフィール情報の抽出
		if strings.Contains(ddText, "誕生日：") {
			profileRegex := regexp.MustCompile(`誕生日：(.+?)\s*/\s*年齢：(.+?)\s*/\s*身長：(.+?)\s*/\s*体重：(.+?)\s*/\s*血液型.{0,5}:\s*(.+?)\s*/\s*職業：(.+?)$`)
			if matches := profileRegex.FindStringSubmatch(ddText); len(matches) > 6 {
				char.Birthday = strings.TrimSpace(matches[1])
				char.Age = strings.TrimSpace(matches[2])
				char.Height = strings.TrimSpace(matches[3])
				char.Weight = strings.TrimSpace(matches[4])
				char.BloodType = strings.TrimSpace(matches[5])
				char.Occupation = strings.TrimSpace(matches[6])
			}
		}
	}

	return char
}
