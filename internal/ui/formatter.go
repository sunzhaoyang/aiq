package ui

import (
	"regexp"
	"strings"
)

var (
	// SQL keywords for syntax highlighting
	sqlKeywords = []string{
		"SELECT", "FROM", "WHERE", "INSERT", "UPDATE", "DELETE",
		"CREATE", "DROP", "ALTER", "TABLE", "INDEX", "VIEW",
		"JOIN", "INNER", "LEFT", "RIGHT", "OUTER", "ON",
		"GROUP", "BY", "ORDER", "HAVING", "LIMIT", "OFFSET",
		"AS", "AND", "OR", "NOT", "IN", "LIKE", "BETWEEN",
		"IS", "NULL", "DISTINCT", "COUNT", "SUM", "AVG", "MAX", "MIN",
		"UNION", "ALL", "EXISTS", "CASE", "WHEN", "THEN", "ELSE", "END",
		"ASC", "DESC", "PRIMARY", "KEY", "FOREIGN", "REFERENCES",
		"CONSTRAINT", "DEFAULT", "CHECK", "UNIQUE",
	}
	
	keywordPattern *regexp.Regexp
)

func init() {
	// Create regex pattern for SQL keywords (case-insensitive, word boundaries)
	pattern := "\\b(?i)(" + strings.Join(sqlKeywords, "|") + ")\\b"
	keywordPattern = regexp.MustCompile(pattern)
}

// HighlightSQL highlights SQL keywords in a query string
func HighlightSQL(sql string) string {
	// First, preserve the original case
	result := keywordPattern.ReplaceAllStringFunc(sql, func(match string) string {
		return SQLKeyword.Render(strings.ToUpper(match))
	})
	return result
}

// FormatSQL formats SQL query with proper indentation (basic)
func FormatSQL(sql string) string {
	// Basic formatting: uppercase keywords, add newlines after major clauses
	sql = strings.TrimSpace(sql)
	
	// Replace common patterns with formatted versions
	replacements := map[string]string{
		"SELECT": "\nSELECT",
		"FROM":   "\nFROM",
		"WHERE":  "\nWHERE",
		"JOIN":   "\nJOIN",
		"LEFT JOIN": "\nLEFT JOIN",
		"RIGHT JOIN": "\nRIGHT JOIN",
		"INNER JOIN": "\nINNER JOIN",
		"GROUP BY": "\nGROUP BY",
		"ORDER BY": "\nORDER BY",
		"HAVING":  "\nHAVING",
		"LIMIT":   "\nLIMIT",
		"UNION":   "\nUNION",
	}
	
	formatted := sql
	for old, new := range replacements {
		formatted = regexp.MustCompile("(?i)\\b"+old+"\\b").ReplaceAllString(formatted, new)
	}
	
	return strings.TrimSpace(formatted)
}

// MaskPassword masks password in connection strings
func MaskPassword(text string) string {
	// Simple masking - replace password=xxx with password=***
	re := regexp.MustCompile(`(?i)password\s*=\s*[^\s&]+`)
	return re.ReplaceAllString(text, "password=***")
}
