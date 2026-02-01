package main

import (
	"encoding/csv"
	"io"
	"strings"
)

func WriteCSV(w io.Writer, characters []Character) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// ヘッダー行
	header := []string{
		"グループ名",
		"グループよみ",
		"名前",
		"よみ",
		"MC名",
		"パンチライン",
		"誕生日",
		"年齢",
		"身長",
		"体重",
		"血液型",
		"職業",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// データ行
	for _, char := range characters {
		record := []string{
			strings.TrimSpace(char.GroupName),
			strings.TrimSpace(char.GroupNameKana),
			strings.TrimSpace(char.Name),
			strings.TrimSpace(char.NameKana),
			strings.TrimSpace(char.MCName),
			strings.TrimSpace(char.Punchline),
			strings.TrimSpace(char.Birthday),
			strings.TrimSpace(char.Age),
			strings.TrimSpace(char.Height),
			strings.TrimSpace(char.Weight),
			strings.TrimSpace(char.BloodType),
			strings.TrimSpace(char.Occupation),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
