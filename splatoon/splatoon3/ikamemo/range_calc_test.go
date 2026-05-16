package main

import (
	"math"
	"testing"
)

// setupWeapons はテスト用にフォールバックデータで Weapons を初期化する
func setupWeapons() {
	Weapons = fallbackWeapons()
}

func TestCalcRangeAdvantage_BarrelVsSplattershot(t *testing.T) {
	setupWeapons()
	result, err := CalcRangeAdvantage("バレル", "スシ")
	if err != nil {
		t.Fatalf("CalcRangeAdvantage failed: %v", err)
	}

	// 射程差: 4.1 - 2.9 = 1.2 ライン
	wantRangeDiff := 0.12e1 // 1.2
	if math.Abs(result.RangeDiff-wantRangeDiff) > 0.001 {
		t.Errorf("RangeDiff = %.3f, want %.3f", result.RangeDiff, wantRangeDiff)
	}

	// 時間有利: 1.2 / 2.4 = 0.5 秒 （スシは軽: 2.4ライン/秒）
	wantTimeAdv := 0.5
	if math.Abs(result.TimeAdv-wantTimeAdv) > 0.001 {
		t.Errorf("TimeAdv = %.3f, want %.3f", result.TimeAdv, wantTimeAdv)
	}

	// 発数有利: 0.5秒 × 15発/秒 = 7.5発
	wantShotsAdv := 7.5
	if math.Abs(result.ShotsAdv-wantShotsAdv) > 0.1 {
		t.Errorf("ShotsAdv = %.1f, want %.1f", result.ShotsAdv, wantShotsAdv)
	}
}

func TestCalcRangeAdvantage_UnknownWeapon(t *testing.T) {
	setupWeapons()
	_, err := CalcRangeAdvantage("存在しないブキ", "スシ")
	if err == nil {
		t.Error("expected error for unknown weapon, got nil")
	}

	_, err = CalcRangeAdvantage("バレル", "存在しないブキ")
	if err == nil {
		t.Error("expected error for unknown weapon, got nil")
	}
}

func TestCalcRangeAdvantage_LongRangeMustBeGreater(t *testing.T) {
	setupWeapons()
	// スシの射程 > バレルの射程 という間違った引数
	_, err := CalcRangeAdvantage("スシ", "バレル")
	if err == nil {
		t.Error("expected error when long range is shorter than short range, got nil")
	}
}

func TestCalcRangeAdvantageByParams_Basic(t *testing.T) {
	// バレル(4.1) vs スシ(2.9), 軽重量, 15発/秒
	result, err := CalcRangeAdvantageByParams(4.1, 2.9, WeightLight, 15.0)
	if err != nil {
		t.Fatalf("CalcRangeAdvantageByParams failed: %v", err)
	}

	if math.Abs(result.RangeDiff-1.2) > 0.001 {
		t.Errorf("RangeDiff = %.3f, want 1.200", result.RangeDiff)
	}

	wantTimeAdv := 1.2 / 2.4 // = 0.5
	if math.Abs(result.TimeAdv-wantTimeAdv) > 0.001 {
		t.Errorf("TimeAdv = %.3f, want %.3f", result.TimeAdv, wantTimeAdv)
	}

	wantShotsAdv := wantTimeAdv * 15.0 // = 7.5
	if math.Abs(result.ShotsAdv-wantShotsAdv) > 0.1 {
		t.Errorf("ShotsAdv = %.1f, want %.1f", result.ShotsAdv, wantShotsAdv)
	}
}

func TestCalcRangeAdvantageByParams_Errors(t *testing.T) {
	tests := []struct {
		name      string
		longRange float64
		shortRng  float64
		weight    WeightClass
		fireRate  float64
	}{
		{"long <= short", 2.9, 4.1, WeightLight, 15.0},
		{"equal range", 4.1, 4.1, WeightLight, 15.0},
		{"invalid weight", 4.1, 2.9, WeightClass("超重"), 15.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CalcRangeAdvantageByParams(tt.longRange, tt.shortRng, tt.weight, tt.fireRate)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestMovementSpeed(t *testing.T) {
	// 軽 > 中 > 重 であることを確認
	if MovementSpeed[WeightLight] <= MovementSpeed[WeightMedium] {
		t.Errorf("Light speed (%.1f) should be greater than Medium speed (%.1f)", MovementSpeed[WeightLight], MovementSpeed[WeightMedium])
	}
	if MovementSpeed[WeightMedium] <= MovementSpeed[WeightHeavy] {
		t.Errorf("Medium speed (%.1f) should be greater than Heavy speed (%.1f)", MovementSpeed[WeightMedium], MovementSpeed[WeightHeavy])
	}
}

func TestFallbackWeaponsData(t *testing.T) {
	fb := fallbackWeapons()
	// バレルとスシが登録されていること
	for _, key := range []string{"バレル", "スシ"} {
		w, ok := fb[key]
		if !ok {
			t.Errorf("weapon %q not found in fallback weapons", key)
			continue
		}
		if w.Range <= 0 {
			t.Errorf("weapon %q has invalid range: %.1f", key, w.Range)
		}
		if w.FireRate <= 0 {
			t.Errorf("weapon %q has invalid fire rate: %.1f", key, w.FireRate)
		}
	}
}

func TestFormatResult(t *testing.T) {
	setupWeapons()
	result, err := CalcRangeAdvantage("バレル", "スシ")
	if err != nil {
		t.Fatalf("CalcRangeAdvantage failed: %v", err)
	}

	formatted := result.FormatResult()
	if formatted == "" {
		t.Error("FormatResult returned empty string")
	}
}

func TestFormatResult_NoFireRate(t *testing.T) {
	// FireRate が 0 の場合、発数有利の行が出力されないこと
	result := RangeAdvantage{
		LongRangeWeapon:  Weapon{Name: "テスト長", Range: 4.0, Weight: WeightMedium, FireRate: 0},
		ShortRangeWeapon: Weapon{Name: "テスト短", Range: 2.0, Weight: WeightLight},
		RangeDiff:        2.0,
		TimeAdv:          0.833,
		ShotsAdv:         0,
	}
	formatted := result.FormatResult()
	if formatted == "" {
		t.Error("FormatResult returned empty string")
	}
	// 発数行が含まれていないこと
	if contains(formatted, "発数有利") {
		t.Error("FormatResult should not include shots advantage when FireRate is 0")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && func() bool {
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	}()
}

