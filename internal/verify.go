package internal

import "os"

// FileExists checks if a file or directory exists and returns an error if it doesn't or if another error occurred.
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err // return other unexpected errors (e.g., permission denied)
}
