package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// color: 1 red, 2 green, 3 yello, 4 blue, 5 purple, 6 blue
func p(color int, sep string, str ...any) {
	newStr := []any{}
	for index, v := range str {
		if index == 0 {
			newStr = append(newStr, v)
		} else {
			newStr = append(newStr, sep, v)
		}
	}

	suffixColor := "\033[3" + strconv.Itoa(color) + "m"
	fmt.Printf("%s%s%s", suffixColor, fmt.Sprint(newStr...), "\033[0m\n")
}

// write to files
func writeJson(data any, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("error marchal %s: %w", filename, err)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error write file %s: %w", filename, err)
	}

	return nil
}

// loading files
func loadJson(filename string) ([]ItemParse, error) {
	plan, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var data []ItemParse
	err = json.Unmarshal(plan, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// append json
func appendJson(data []ItemParse, filename string) {
	dataFiles, err := loadJson(filename)
	if err != nil {
		writeJson(data, filename)
	} else if len(dataFiles) > 0 {
		dataFiles = append(dataFiles, data...)
		writeJson(dataFiles, filename)
	}
}
