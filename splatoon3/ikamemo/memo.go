package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type MemoGenerator struct {
	templatePath string
	outputDir    string
}

type MemoConfig struct {
	Weapon  string
	Rules   []string
	Stages  []string
	Results []string
}

func NewMemoGenerator(templatePath, outputDir string) *MemoGenerator {
	return &MemoGenerator{
		templatePath: templatePath,
		outputDir:    outputDir,
	}
}

func (mg *MemoGenerator) Generate(config MemoConfig) (string, error) {
	now := time.Now()
	filename := fmt.Sprintf("result_%s.md", now.Format("20060102150405"))

	content, err := os.ReadFile(mg.templatePath)
	if err != nil {
		return "", fmt.Errorf("テンプレートファイルが読み込めません: %v", err)
	}

	modifiedContent := string(content)

	modifiedContent = mg.insertDate(modifiedContent, now)
	modifiedContent = mg.insertWeapon(modifiedContent, config.Weapon)
	modifiedContent = mg.checkRules(modifiedContent, config.Rules)
	modifiedContent = mg.checkStages(modifiedContent, config.Stages)
	modifiedContent = mg.checkResults(modifiedContent, config.Results)

	if err := os.MkdirAll(mg.outputDir, 0755); err != nil {
		return "", fmt.Errorf("出力ディレクトリの作成に失敗: %v", err)
	}

	outputPath := filepath.Join(mg.outputDir, filename)
	if err := os.WriteFile(outputPath, []byte(modifiedContent), 0644); err != nil {
		return "", fmt.Errorf("ファイルの作成に失敗: %v", err)
	}

	return outputPath, nil
}

func (mg *MemoGenerator) insertDate(content string, now time.Time) string {
	dateStr := now.Format("2006-01-02 15:04:05")
	return strings.Replace(content, "日付: ", fmt.Sprintf("日付: %s", dateStr), 1)
}

func (mg *MemoGenerator) insertWeapon(content, weapon string) string {
	if weapon == "" {
		return content
	}

	content = strings.Replace(content, "- (自分)", fmt.Sprintf("- %s (自分)", weapon), 1)
	content = strings.Replace(content, "使用ブキ:", fmt.Sprintf("使用ブキ: %s", weapon), 1)

	return content
}

func (mg *MemoGenerator) checkRules(content string, rules []string) string {
	for _, rule := range rules {
		switch strings.ToLower(rule) {
		case "エリア", "area":
			content = strings.Replace(content, "  - [ ] ガチエリア", "  - [x] ガチエリア", 1)
		case "ヤグラ", "yagura":
			content = strings.Replace(content, "  - [ ] ガチヤグラ", "  - [x] ガチヤグラ", 1)
		case "ホコ", "hoko":
			content = strings.Replace(content, "  - [ ] ガチホコ", "  - [x] ガチホコ", 1)
		case "アサリ", "asari":
			content = strings.Replace(content, "  - [ ] ガチアサリ", "  - [x] ガチアサリ", 1)
		default:
			content = mg.findAndCheckRule(content, rule)
		}
	}
	return content
}

func (mg *MemoGenerator) findAndCheckRule(content, rule string) string {
	rulePattern := regexp.MustCompile(`  - \[ \] (ガチ.+)`)
	matches := rulePattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			fullRuleName := match[1]
			if strings.Contains(strings.ToLower(fullRuleName), strings.ToLower(rule)) {
				content = strings.Replace(content,
					fmt.Sprintf("  - [ ] %s", fullRuleName),
					fmt.Sprintf("  - [x] %s", fullRuleName), 1)
				break
			}
		}
	}

	return content
}

func (mg *MemoGenerator) checkStages(content string, stages []string) string {
	for _, stage := range stages {
		content = mg.findAndCheckStage(content, stage)
	}
	return content
}

func (mg *MemoGenerator) findAndCheckStage(content, stage string) string {
	stagePattern := regexp.MustCompile(`- ステージ:\n((?:\s+- \[ \] .+\n)+)`)
	stageMatch := stagePattern.FindStringSubmatch(content)

	if len(stageMatch) > 1 {
		stageSection := stageMatch[1]
		itemPattern := regexp.MustCompile(`\s+- \[ \] (.+)`)
		itemMatches := itemPattern.FindAllStringSubmatch(stageSection, -1)

		for _, match := range itemMatches {
			if len(match) > 1 {
				fullStageName := match[1]
				if strings.Contains(strings.ToLower(fullStageName), strings.ToLower(stage)) {
					content = strings.Replace(content,
						fmt.Sprintf("  - [ ] %s", fullStageName),
						fmt.Sprintf("  - [x] %s", fullStageName), 1)
					break
				}
			}
		}
	}

	return content
}

func (mg *MemoGenerator) checkResults(content string, results []string) string {
	for _, result := range results {
		switch strings.ToLower(result) {
		case "win", "勝利":
			content = strings.Replace(content, "  - [ ] Win", "  - [x] Win", 1)
		case "lose", "敗北":
			content = strings.Replace(content, "  - [ ] Lose", "  - [x] Lose", 1)
		}
	}
	return content
}
