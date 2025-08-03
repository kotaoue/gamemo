package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	weapon   string
	rules    []string
	stages   []string
	outcomes []string
)

var rootCmd = &cobra.Command{
	Use:   "ikamemo",
	Short: "Splatoon 3 battle memo generator",
	Long:  "Generate battle memo from template for Splatoon 3 analysis",
	RunE:  runCommand,
}

func init() {
	rootCmd.Flags().StringVarP(&weapon, "weapon", "w", "", "武器名")
	rootCmd.Flags().StringSliceVarP(&rules, "rule", "r", []string{}, "ルール (エリア/ヤグラ/ホコ/アサリ)")
	rootCmd.Flags().StringSliceVarP(&stages, "stage", "s", []string{}, "ステージ名")
	rootCmd.Flags().StringSliceVarP(&outcomes, "outcome", "o", []string{}, "結果 (win/lose)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runCommand(cmd *cobra.Command, args []string) error {
	gen := NewGenerator("../results/template.md", "../results")
	
	config := Config{
		Weapon:  weapon,
		Rules:   rules,
		Stages:  stages,
		Results: outcomes,
	}
	
	outputPath, err := gen.Generate(config)
	if err != nil {
		return err
	}
	
	fmt.Printf("メモファイルを作成しました: %s\n", outputPath)
	return nil
}