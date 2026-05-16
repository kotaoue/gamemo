package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rangeWeapon1     string
	rangeWeapon2     string
	rangeVal1        float64
	rangeVal2        float64
	rangeMoverWeight string
	rangeFireRate    float64
)

var rangeCmd = &cobra.Command{
	Use:   "range",
	Short: "射程差から時間・発数アドバンテージを計算する",
	Long: `射程差を基に、長射程側が何秒・何発分有利かを計算します。

ブキ名で計算する例:
  ikamemo range -a バレル -b スシ

射程・重量・発射速度を直接指定して計算する例:
  ikamemo range --range1 4.1 --range2 2.9 --weight 軽 --fire-rate 15

利用可能なブキ名: バレル, スシ, シャーカー, プラコラ, 52, 96`,
	RunE: runRangeCommand,
}

func init() {
	rangeCmd.Flags().StringVarP(&rangeWeapon1, "weapon1", "a", "", "長射程ブキ名")
	rangeCmd.Flags().StringVarP(&rangeWeapon2, "weapon2", "b", "", "短射程ブキ名")
	rangeCmd.Flags().Float64Var(&rangeVal1, "range1", 0, "長射程値（ライン）")
	rangeCmd.Flags().Float64Var(&rangeVal2, "range2", 0, "短射程値（ライン）")
	rangeCmd.Flags().StringVarP(&rangeMoverWeight, "weight", "W", "", "短射程側の重量クラス（軽/中/重）")
	rangeCmd.Flags().Float64VarP(&rangeFireRate, "fire-rate", "f", 0, "長射程ブキの発射速度（発/秒）")

	rootCmd.AddCommand(rangeCmd)
}

func runRangeCommand(cmd *cobra.Command, args []string) error {
	// ブキ名で指定された場合
	if rangeWeapon1 != "" || rangeWeapon2 != "" {
		return runRangeByWeaponName()
	}

	// 直接パラメータ指定の場合
	if rangeVal1 > 0 || rangeVal2 > 0 {
		return runRangeByParams()
	}

	// どちらも指定がない場合はブキ一覧を表示
	return printWeaponList()
}

func runRangeByWeaponName() error {
	if rangeWeapon1 == "" {
		return fmt.Errorf("長射程ブキ名(-a)を指定してください")
	}
	if rangeWeapon2 == "" {
		return fmt.Errorf("短射程ブキ名(-b)を指定してください")
	}

	result, err := CalcRangeAdvantage(rangeWeapon1, rangeWeapon2)
	if err != nil {
		return err
	}

	fmt.Printf("=== 射程差アドバンテージ計算 ===\n")
	fmt.Printf("長射程: %s (射程 %.1f, 重量 %s)\n", result.LongRangeWeapon.Name, result.LongRangeWeapon.Range, result.LongRangeWeapon.Weight)
	fmt.Printf("短射程: %s (射程 %.1f, 重量 %s)\n", result.ShortRangeWeapon.Name, result.ShortRangeWeapon.Range, result.ShortRangeWeapon.Weight)
	fmt.Println("---")
	fmt.Println(result.FormatResult())
	return nil
}

func runRangeByParams() error {
	if rangeVal1 <= 0 {
		return fmt.Errorf("長射程値(--range1)を正の値で指定してください")
	}
	if rangeVal2 <= 0 {
		return fmt.Errorf("短射程値(--range2)を正の値で指定してください")
	}
	if rangeMoverWeight == "" {
		return fmt.Errorf("短射程側の重量クラス(-W)を指定してください（軽/中/重）")
	}
	if rangeFireRate <= 0 {
		return fmt.Errorf("長射程ブキの発射速度(-f)を正の値で指定してください")
	}

	weight := WeightClass(rangeMoverWeight)
	result, err := CalcRangeAdvantageByParams(rangeVal1, rangeVal2, weight, rangeFireRate)
	if err != nil {
		return err
	}

	fmt.Printf("=== 射程差アドバンテージ計算 ===\n")
	fmt.Printf("長射程: %.1f ライン (発射速度 %.1f 発/秒)\n", result.LongRangeWeapon.Range, result.LongRangeWeapon.FireRate)
	fmt.Printf("短射程: %.1f ライン (重量 %s)\n", result.ShortRangeWeapon.Range, result.ShortRangeWeapon.Weight)
	fmt.Println("---")
	fmt.Println(result.FormatResult())
	return nil
}

func printWeaponList() error {
	fmt.Println("=== 利用可能なブキ一覧 ===")
	fmt.Printf("%-12s %-20s %6s %8s %6s %10s\n", "キー", "名前", "射程", "ダメージ", "重量", "発射速度")
	fmt.Println(strings.Repeat("-", 70))

	// 表示順を固定するためにキーをソート
	keys := []string{"スシ", "シャーカー", "52", "プラコラ", "96", "バレル"}
	for _, k := range keys {
		w := Weapons[k]
		fmt.Printf("%-12s %-20s %6.1f %8.0f %6s %8.1f発/秒\n",
			k, w.Name, w.Range, w.Damage, w.Weight, w.FireRate)
	}
	fmt.Println()
	fmt.Printf("重量クラス別移動速度:\n")
	for _, wc := range []WeightClass{WeightLight, WeightMedium, WeightHeavy} {
		fmt.Printf("  %s: %.1f ライン/秒\n", wc, MovementSpeed[wc])
	}
	fmt.Println()
	fmt.Printf("使用例:\n")
	fmt.Printf("  ikamemo range -a バレル -b スシ\n")
	fmt.Printf("  ikamemo range --range1 4.1 --range2 2.9 -W 軽 -f 15\n")

	return nil
}

// formatFloat は小数点以下の余分なゼロを省いてフォーマットする
func formatFloat(f float64) string {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	return s
}
