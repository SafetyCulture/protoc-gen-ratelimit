package genratelimit

import (
	"fmt"

	gendoc "github.com/pseudomuto/protoc-gen-doc"
)

func getServicePath(file *gendoc.File, service *gendoc.Service) string {
	return fmt.Sprintf("%s.%s", file.Package, service.Name)
}

func getDefaultMethodPath(file *gendoc.File, service *gendoc.Service, method *gendoc.ServiceMethod) string {
	return fmt.Sprintf("/%s/%s", getServicePath(file, service), method.Name)
}

func newDescriptorTuple(clientClass string, orgID string, userID string, descriptors []yamlDescriptor) yamlDescriptor {
	return yamlDescriptor{
		Key:   "client_class",
		Value: clientClass,
		Descriptors: []yamlDescriptor{
			{
				Key:   "org_id",
				Value: orgID,
				Descriptors: []yamlDescriptor{
					{
						Key:         "user_id",
						Value:       userID,
						Descriptors: descriptors,
					},
				},
			},
		},
	}
}
