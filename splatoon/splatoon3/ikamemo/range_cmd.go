package main

import (
	"fmt"
	"sort"
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
	rangeOffline     bool
)

var rangeCmd = &cobra.Command{
	Use:   "range",
	Short: "射程差から時間・発数アドバンテージを計算する",
	Long: `射程差を基に、長射程側が何秒・何発分有利かを計算します。
ブキデータはwikiwiki.jp/splatoon3mixから自動取得します（--offlineで省略可）。

ブキ名で計算する例:
  ikamemo range -a バレル -b スシ

射程・重量・発射速度を直接指定して計算する例:
  ikamemo range --range1 4.1 --range2 2.9 --weight 軽 --fire-rate 15

ブキ一覧を表示（データ取得を確認）:
  ikamemo range`,
	RunE: runRangeCommand,
}

func init() {
	rangeCmd.Flags().StringVarP(&rangeWeapon1, "weapon1", "a", "", "長射程ブキ名")
	rangeCmd.Flags().StringVarP(&rangeWeapon2, "weapon2", "b", "", "短射程ブキ名")
	rangeCmd.Flags().Float64Var(&rangeVal1, "range1", 0, "長射程値（ライン）")
	rangeCmd.Flags().Float64Var(&rangeVal2, "range2", 0, "短射程値（ライン）")
	rangeCmd.Flags().StringVarP(&rangeMoverWeight, "weight", "W", "", "短射程側の重量クラス（軽/中/重）")
	rangeCmd.Flags().Float64VarP(&rangeFireRate, "fire-rate", "f", 0, "長射程ブキの発射速度（発/秒）。省略するとブキ名検索時の登録値を使用")
	rangeCmd.Flags().BoolVar(&rangeOffline, "offline", false, "スクレイピングをスキップしてフォールバックデータを使う")

	rootCmd.AddCommand(rangeCmd)
}

// initWeapons はコマンド実行前にブキデータを初期化する
func initWeapons() {
	if Weapons != nil {
		return
	}
	if rangeOffline {
		Weapons = fallbackWeapons()
		return
	}
	weapons, fromWeb := LoadWeaponsWithFallback()
	if !fromWeb {
		fmt.Println("※ ブキデータの取得に失敗したため、フォールバックデータを使用しています")
	}
	Weapons = weapons
}

func runRangeCommand(cmd *cobra.Command, args []string) error {
	initWeapons()

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

	// --fire-rate が指定された場合はブキのFireRateを上書きする
	if rangeFireRate > 0 {
		if w, ok := Weapons[rangeWeapon1]; ok {
			w.FireRate = rangeFireRate
			Weapons[rangeWeapon1] = w
		}
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
	fmt.Println("=== ブキ一覧 ===")
	fmt.Printf("%-20s %6s %8s %6s %10s\n", "名前", "射程", "ダメージ", "重量", "発射速度")
	fmt.Println(strings.Repeat("-", 56))

	keys := make([]string, 0, len(Weapons))
	for k := range Weapons {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		w := Weapons[k]
		fireRateStr := "-"
		if w.FireRate > 0 {
			fireRateStr = fmt.Sprintf("%.1f発/秒", w.FireRate)
		}
		fmt.Printf("%-20s %6.1f %8.0f %6s %10s\n",
			w.Name, w.Range, w.Damage, w.Weight, fireRateStr)
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

