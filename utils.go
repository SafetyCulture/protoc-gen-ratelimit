package genratelimit

import (
	"fmt"
	"strings"

	gendoc "github.com/pseudomuto/protoc-gen-doc"
)

func getServicePath(file *gendoc.File, service *gendoc.Service) string {
	return fmt.Sprintf("%s.%s", file.Package, service.Name)
}

func getDefaultMethodPath(file *gendoc.File, service *gendoc.Service, method *gendoc.ServiceMethod) string {
	return fmt.Sprintf("/%s/%s", getServicePath(file, service), method.Name)
}

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
