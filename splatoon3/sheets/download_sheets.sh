#!/bin/zsh

echo "Please enter your API key: "
read API_KEY
export API_KEY

echo "Please enter your Spreadsheet ID: "
read SPREADSHEET_ID
export SPREADSHEET_ID

echo "Getting sheet list..."
SHEETS_RESPONSE=$(curl -s "https://sheets.googleapis.com/v4/spreadsheets/${SPREADSHEET_ID}?fields=sheets.properties.title&key=${API_KEY}")

if [[ $SHEETS_RESPONSE == *"error"* ]]; then
  echo "An error occurred:"
  echo $SHEETS_RESPONSE | jq .
  exit 1
fi

SHEET_NAMES=$(echo $SHEETS_RESPONSE | jq -r '.sheets[].properties.title')

echo "Downloading sheet data..."
for SHEET_NAME in $SHEET_NAMES; do
  ENCODED_SHEET_NAME=$(echo $SHEET_NAME | sed 's/ /%20/g')

  echo "Getting data for sheet \"${SHEET_NAME}\"..."

  curl -s "https://sheets.googleapis.com/v4/spreadsheets/${SPREADSHEET_ID}/values/${ENCODED_SHEET_NAME}?key=${API_KEY}" \
       > "./${SHEET_NAME}.json"

  if [[ -s "./${SHEET_NAME}.json" ]]; then
    echo "\"${SHEET_NAME}.json\" has been saved successfully"
  else
    echo "Failed to get data for \"${SHEET_NAME}\""
  fi
done

echo "All sheet data has been downloaded successfully"
