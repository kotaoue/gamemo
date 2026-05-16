package main

import (
	"fmt"
	"math"
)

// WeightClass は武器の重量クラスを表す
type WeightClass string

const (
	WeightLight  WeightClass = "軽"
	WeightMedium WeightClass = "中"
	WeightHeavy  WeightClass = "重"
)

// MovementSpeed は重量クラスごとの移動速度（ライン/秒）
var MovementSpeed = map[WeightClass]float64{
	WeightLight:  2.4,
	WeightMedium: 2.0,
	WeightHeavy:  1.6,
}

// Weapon はスプラトゥーン3のブキ情報を表す
type Weapon struct {
	Name     string
	Range    float64     // 射程（ライン単位）
	Damage   float64     // ダメージ
	Weight   WeightClass // 重量クラス
	FireRate float64     // 発射速度（発/秒）
}

// Weapons は主要なブキのデータ
var Weapons = map[string]Weapon{
	"バレル": {
		Name:     "バレルスピナー",
		Range:    4.1,
		Damage:   30,
		Weight:   WeightMedium,
		FireRate: 15.0, // 一周チャージ48F(0.8秒)で12発 = 15発/秒
	},
	"スシ": {
		Name:     "スプラシューター",
		Range:    2.9,
		Damage:   35,
		Weight:   WeightLight,
		FireRate: 10.0, // 6F間隔 = 10発/秒
	},
	"シャーカー": {
		Name:     "スプラッシュボム・シャーカー",
		Range:    2.4,
		Damage:   28,
		Weight:   WeightLight,
		FireRate: 11.0, // 約5-6F間隔
	},
	"プラコラ": {
		Name:     "フォルテスプラシューターコラボ",
		Range:    3.4,
		Damage:   42,
		Weight:   WeightMedium,
		FireRate: 6.0, // 約10F間隔
	},
	"52": {
		Name:     ".52ガロン",
		Range:    3.2,
		Damage:   52,
		Weight:   WeightMedium,
		FireRate: 5.7, // 約10-11F間隔
	},
	"96": {
		Name:     ".96ガロン",
		Range:    3.6,
		Damage:   62,
		Weight:   WeightHeavy,
		FireRate: 3.0, // 約20F間隔
	},
}

// RangeAdvantage は射程差から算出したアドバンテージ情報
type RangeAdvantage struct {
	LongRangeWeapon  Weapon
	ShortRangeWeapon Weapon
	RangeDiff        float64 // 射程差（ライン）
	TimeAdv          float64 // 時間有利（秒）
	ShotsAdv         float64 // 発数有利（発）
}

// CalcRangeAdvantage は2つのブキ名から射程差アドバンテージを計算する
// shortMoverWeight は射程の短い側（詰めてくる側）の重量クラス
func CalcRangeAdvantage(longName, shortName string) (RangeAdvantage, error) {
	long, ok := Weapons[longName]
	if !ok {
		return RangeAdvantage{}, fmt.Errorf("ブキが見つかりません: %s", longName)
	}
	short, ok := Weapons[shortName]
	if !ok {
		return RangeAdvantage{}, fmt.Errorf("ブキが見つかりません: %s", shortName)
	}

	return calcAdvantage(long, short)
}

// CalcRangeAdvantageByParams は射程・重量・発射速度を直接指定して計算する
func CalcRangeAdvantageByParams(longRange float64, shortRange float64, moverWeight WeightClass, longFireRate float64) (RangeAdvantage, error) {
	speed, ok := MovementSpeed[moverWeight]
	if !ok {
		return RangeAdvantage{}, fmt.Errorf("重量クラスが不正です: %s", moverWeight)
	}
	if longRange <= shortRange {
		return RangeAdvantage{}, fmt.Errorf("長射程(%v)は短射程(%v)より大きくなければなりません", longRange, shortRange)
	}
	if speed <= 0 {
		return RangeAdvantage{}, fmt.Errorf("移動速度は正の値でなければなりません")
	}

	rangeDiff := longRange - shortRange
	timeAdv := rangeDiff / speed
	shotsAdv := timeAdv * longFireRate

	return RangeAdvantage{
		LongRangeWeapon:  Weapon{Range: longRange, FireRate: longFireRate},
		ShortRangeWeapon: Weapon{Range: shortRange, Weight: moverWeight},
		RangeDiff:        math.Round(rangeDiff*100) / 100,
		TimeAdv:          math.Round(timeAdv*1000) / 1000,
		ShotsAdv:         math.Round(shotsAdv*10) / 10,
	}, nil
}

func calcAdvantage(long, short Weapon) (RangeAdvantage, error) {
	if long.Range <= short.Range {
		return RangeAdvantage{}, fmt.Errorf("%s(%v)の射程は%s(%v)以下です", long.Name, long.Range, short.Name, short.Range)
	}

	speed, ok := MovementSpeed[short.Weight]
	if !ok {
		return RangeAdvantage{}, fmt.Errorf("重量クラスが不正です: %s", short.Weight)
	}

	rangeDiff := long.Range - short.Range
	timeAdv := rangeDiff / speed
	shotsAdv := timeAdv * long.FireRate

	return RangeAdvantage{
		LongRangeWeapon:  long,
		ShortRangeWeapon: short,
		RangeDiff:        math.Round(rangeDiff*100) / 100,
		TimeAdv:          math.Round(timeAdv*1000) / 1000,
		ShotsAdv:         math.Round(shotsAdv*10) / 10,
	}, nil
}

// FormatResult はRangeAdvantageを人間が読みやすい形式にフォーマットする
func (r RangeAdvantage) FormatResult() string {
	shortLabel := r.ShortRangeWeapon.Name
	if shortLabel == "" {
		shortLabel = fmt.Sprintf("%.1fライン", r.ShortRangeWeapon.Range)
	}
	longLabel := r.LongRangeWeapon.Name
	if longLabel == "" {
		longLabel = fmt.Sprintf("%.1fライン", r.LongRangeWeapon.Range)
	}
	return fmt.Sprintf(
		"射程差: %.2f ライン\n移動時間有利: %.3f 秒 （%s側の移動速度 %.1f ライン/秒）\n発数有利: %.1f 発 （%sの発射速度 %.1f 発/秒）",
		r.RangeDiff,
		r.TimeAdv,
		shortLabel,
		MovementSpeed[r.ShortRangeWeapon.Weight],
		r.ShotsAdv,
		longLabel,
		r.LongRangeWeapon.FireRate,
	)
}
