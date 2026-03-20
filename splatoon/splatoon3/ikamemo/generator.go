package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type Generator struct {
	templatePath string
	outputDir    string
}

type Config struct {
	Weapon  string
	Rules   []string
	Stages  []string
	Results []string
}

func NewGenerator(templatePath, outputDir string) *Generator {
	return &Generator{
		templatePath: templatePath,
		outputDir:    outputDir,
	}
}

func (g *Generator) Generate(config Config) (string, error) {
	content, err := g.GenerateContent(config)
	if err != nil {
		return "", err
	}

	return g.WriteToFile(content)
}

func (g *Generator) GenerateContent(config Config) (string, error) {
	now := time.Now()

	content, err := os.ReadFile(g.templatePath)
	if err != nil {
		return "", fmt.Errorf("テンプレートファイルが読み込めません: %v", err)
	}

	modifiedContent := string(content)

	modifiedContent = g.insertDate(modifiedContent, now)
	modifiedContent = g.insertWeapon(modifiedContent, config.Weapon)
	modifiedContent = g.checkRules(modifiedContent, config.Rules)
	modifiedContent = g.checkStages(modifiedContent, config.Stages)
	modifiedContent = g.checkResults(modifiedContent, config.Results)

	return modifiedContent, nil
}

func (g *Generator) insertDate(content string, now time.Time) string {
	dateStr := now.Format("2006-01-02 15:04:05")
	return strings.Replace(content, "日付: ", fmt.Sprintf("日付: %s", dateStr), 1)
}

func (g *Generator) insertWeapon(content, weapon string) string {
	if weapon == "" {
		return content
	}

	content = strings.Replace(content, "- (自分)", fmt.Sprintf("- %s (自分)", weapon), 1)
	content = strings.Replace(content, "使用ブキ:", fmt.Sprintf("使用ブキ: %s", weapon), 1)

	return content
}

func (g *Generator) checkRules(content string, rules []string) string {
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
			content = g.findAndCheckRule(content, rule)
		}
	}
	return content
}

func (g *Generator) findAndCheckRule(content, rule string) string {
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

func (g *Generator) checkStages(content string, stages []string) string {
	for _, stage := range stages {
		content = g.findAndCheckStage(content, stage)
	}
	return content
}

func (g *Generator) findAndCheckStage(content, stage string) string {
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

func (g *Generator) checkResults(content string, results []string) string {
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