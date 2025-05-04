package utils

import "strings"

type SortField struct {
	Field     string
	Direction string
}

func ParseSortBy(sortBy string) []SortField {
	var sortFields []SortField

	items := strings.Split(sortBy, ",")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		var field string
		var direction string

		if strings.Contains(item, ":") {
			parts := strings.SplitN(item, ":", 2)
			field = parts[0]
			if strings.ToLower(parts[1]) == "desc" {
				direction = "DESC"
			} else {
				direction = "ASC"
			}
		} else {
			field = item
			direction = "ASC"
		}

		sortFields = append(sortFields, SortField{
			Field:     field,
			Direction: direction,
		})
	}

	return sortFields
}
