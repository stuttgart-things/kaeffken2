package internal

import (
	"fmt"
	"math/rand"
	"time"
)

func GetRandomStringFromMap(m map[string]interface{}, key string) (string, error) {
	val, ok := m[key]
	if !ok {
		return "", fmt.Errorf("key %q not found", key)
	}

	strSlice, ok := val.([]interface{})
	if !ok {
		return "", fmt.Errorf("value at key %q is not a slice", key)
	}

	if len(strSlice) == 0 {
		return "", fmt.Errorf("slice at key %q is empty", key)
	}

	// Convert []interface{} to []string
	strings := make([]string, 0, len(strSlice))
	for _, item := range strSlice {
		if s, ok := item.(string); ok {
			strings = append(strings, s)
		}
	}

	if len(strings) == 0 {
		return "", fmt.Errorf("no string elements found at key %q", key)
	}

	// Use a local rand.Rand instance
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strings[r.Intn(len(strings))], nil
}
