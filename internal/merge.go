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

func CleanMap(input map[string]interface{}) map[string]interface{} {
	cleaned := make(map[string]interface{})
	for k, v := range input {
		switch val := v.(type) {
		case nil:
			// skip
		case string:
			if val != "" {
				cleaned[k] = val
			}
		case []interface{}:
			if len(val) > 0 {
				cleaned[k] = val
			}
		case map[string]interface{}:
			// recursively clean nested maps
			nested := CleanMap(val)
			if len(nested) > 0 {
				cleaned[k] = nested
			}
		default:
			cleaned[k] = val
		}
	}
	return cleaned
}
