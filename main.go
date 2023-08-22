package main

import (
	"bufio"
	"os"
	"strings"
	"thewisepigeon/sql-to-ts/categorizer"
	"thewisepigeon/sql-to-ts/parser"
)

func main() {
	if len(os.Args) < 2 {
		println("Missing file path")
		return
	}
	sql_file_path := os.Args[1]
	if !strings.HasSuffix(sql_file_path, ".sql") {
		println("Only SQL files are supported")
		return
	}
	reader, err := os.Open(sql_file_path)
	if err != nil {
		println("Error reading SQL file")
		return
	}
	defer reader.Close()
	scanner := bufio.NewScanner(reader)
	var parsed_tokens [][]string
	context := ""
	previous_context := ""
	line_number := 0
	for scanner.Scan() {
		line_number += 1
		line := strings.ToLower(scanner.Text())
		line_category := categorizer.Categorize(strings.TrimSpace(line))
		if line_category == "MULTILINE_COMMENT_START" {
			previous_context = context
			context = "MULTILINE_COMMENT"
			continue
		}
		if line_category == "MULTILINE_COMMENT_END" {
			if context != "MULTILINE_COMMENT" {
				println("Invalid character error */")
				return
			}
			context = previous_context
		}
		if line_category == "NEXT" {
			context = ""
		}
		if line_category == "DELIMITER_START" {
			if context == "PARSING" || context == "START_PARSING" {
				println("Error parsing")
				return
			}
			context = "START_PARSING"
		}
		if line_category == "FIELD" {
			if context == "MULTILINE_COMMENT" {
				continue
			}
			if context != "START_PARSING" && context != "PARSING" {
				println("Error parsing")
				return
			}
			if context == "START_PARSING" {
				context = "PARSING"
			}
		}
		if line_category == "DELIMITER_END" {
			if context == "START_PARSING" {
				println("Error can not parse empty table")
				return
			}
			if context == "" {
				println("Error parsing")
				return
			}
		}
		err := parser.Parse(line, context, line_category, &parsed_tokens)
		if err != nil {
			println(
				"Parsing error at line ",
				line_number,
				": ",
				"`",
				line,
				"`",
			)
			return
		}
	}
}
