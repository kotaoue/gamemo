package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (g *Generator) WriteToFile(content string) (string, error) {
	now := time.Now()
	filename := fmt.Sprintf("result_%s.md", now.Format("20060102150405"))

	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return "", fmt.Errorf("出力ディレクトリの作成に失敗: %v", err)
	}

	outputPath := filepath.Join(g.outputDir, filename)
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("ファイルの作成に失敗: %v", err)
	}

	return outputPath, nil
}