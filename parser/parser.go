package parser

import (
	"errors"
	"strings"
	"thewisepigeon/sql-to-ts/purifier"
)

func Parse(line string, context string, category string, parsed_tokens *[][]string) error {
  if category == "NEXT"{
    return nil
  }
	if category == "DELIMITER_START" {
    if context == "" {
      return errors.New("Parsing error")
    }
    result := strings.Split(strings.TrimSpace(purifier.Purify(line, category)), " ")
    if len(result)>1 {
      return errors.New("Parsing error")
    }
    table_name := result[0]
    println(table_name)
	}
	return nil
}
