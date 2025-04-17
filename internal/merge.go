package internal

func MergeMaps(a, b map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	// Copy all from a
	for key, value := range a {
		merged[key] = value
	}

	// Copy all from b (overwrites if key already exists)
	for key, value := range b {
		merged[key] = value
	}

	return merged
}
