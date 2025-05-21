package genratelimit

import (
	"fmt"
	"strings"
)

func formatKey(key, bucketName string, count int) (string, error) {
	updatedKey := key

	keys := strings.Count(key, delimiter)
	if keys == count && (string(key[len(key)-1]) != delimiter && bucketName != "") {
		return "", fmt.Errorf("key %s has too tuples, last one should be reserved for bucket", key)
	}
	if keys > count-1 {
		return "", fmt.Errorf("key %s has too many delimiters", key)
	}

	for i := 0; i < (count - keys - 1); i++ {
		updatedKey = updatedKey + delimiter
	}

	return updatedKey + bucketName, nil
}
