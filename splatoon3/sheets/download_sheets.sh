#!/bin/zsh

# 設定ファイルの読み込み
CONFIG_FILE="./config.env"
if [[ -f "$CONFIG_FILE" ]]; then
  source "$CONFIG_FILE"
else
  echo "Error: Configuration file not found: $CONFIG_FILE"
  exit 1
fi

# 必須変数の確認
if [[ -z "$API_KEY" || -z "$SPREADSHEET_ID" ]]; then
  echo "Error: API_KEY and SPREADSHEET_ID must be set in $CONFIG_FILE"
  exit 1
fi

echo "Getting sheet list..."
SHEETS_RESPONSE=$(curl -s "https://sheets.googleapis.com/v4/spreadsheets/${SPREADSHEET_ID}?fields=sheets.properties.title&key=${API_KEY}")

if [[ $SHEETS_RESPONSE == *"error"* ]]; then
  echo "An error occurred:"
  echo $SHEETS_RESPONSE | jq .
  exit 1
fi

echo "Downloading sheet data..."
while IFS= read -r SHEET_NAME; do
  ENCODED_SHEET_NAME=$(printf '%s' "$SHEET_NAME" | jq -sRr @uri)

  echo "Getting data for sheet \"${SHEET_NAME}\"..."

  curl -s "https://sheets.googleapis.com/v4/spreadsheets/${SPREADSHEET_ID}/values/${ENCODED_SHEET_NAME}?key=${API_KEY}" \
       > "./${SHEET_NAME}.json"

  if [[ -s "./${SHEET_NAME}.json" ]]; then
    echo "\"${SHEET_NAME}.json\" has been saved successfully"
  else
    echo "Failed to get data for \"${SHEET_NAME}\""
  fi
done < <(echo "$SHEETS_RESPONSE" | jq -r '.sheets[].properties.title')

echo "All sheet data has been downloaded successfully"
