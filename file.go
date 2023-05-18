package gut

import (
	"fmt"
	"strings"
)

func createHeader(s Settings) string {
	sb := strings.Builder{}

	// Append first line if exists
	if s.FirstLine != "" {
		sb.WriteString(s.FirstLine)
	}

	// Check for uuid type
	if s.UuidType == "" {
		sb.WriteString(fmt.Sprintf("export type UuidType = %s\n", "string"))
	} else {
		sb.WriteString(fmt.Sprintf("export type UuidType = %s\n", s.UuidType))
	}

	// Check for int64/uint64 types
	if s.BigIntType == "" {
		sb.WriteString(fmt.Sprintf("export type BigIntType = %s\n", "BigInt"))
	} else {
		sb.WriteString(fmt.Sprintf("export type BigIntType = %s\n", s.BigIntType))
	}

	// Check for time.Time type
	if s.DateType == "" {
		sb.WriteString(fmt.Sprintf("export type DateType = %s\n", "Date"))
	} else {
		sb.WriteString(fmt.Sprintf("export type DateType = %s\n", s.DateType))
	}

	// seperate type definitions from the generated interfaces
	sb.WriteString("\n\n")

	return sb.String()
}
