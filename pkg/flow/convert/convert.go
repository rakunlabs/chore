package convert

import "strings"

func GetBoolean(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return v == "true"
	default:
		return false
	}
}

func GetList(value interface{}) []string {
	switch v := value.(type) {
	case []string:
		return v
	case string:
		return strings.Fields(strings.ReplaceAll(v, ",", " "))
	default:
		return nil
	}
}

func IsTagsEnabled(tags []string, enabledTags map[string]struct{}) bool {
	if len(tags) == 0 {
		return true
	}

	for _, tag := range tags {
		if _, ok := enabledTags[tag]; ok {
			return true
		}
	}

	return false
}

func SliceToMap(list []string) map[string]struct{} {
	m := make(map[string]struct{}, len(list))
	for _, v := range list {
		m[v] = struct{}{}
	}

	return m
}
