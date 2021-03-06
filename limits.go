package genratelimit

import (
	"strings"
)

// Limit is the limit applied to specific descriptors
type Limit struct {
	Key   string
	Value *YamlRateLimit
}

type limits []*Limit

// Len returns the length of the limits
func (s limits) Len() int { return len(s) }

// Swap swaps the two limits
func (s limits) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less returns true if the limit at index i is less than the limit at index j
func (s limits) Less(i, j int) bool {
	iKeys := strings.Split(s[i].Key, delimiter)
	jKeys := strings.Split(s[j].Key, delimiter)

	for idx := range iKeys {
		if iKeys[idx] == jKeys[idx] {
			continue
		}

		// "" should be last, this ensure that the most specific tuple set is used
		if iKeys[idx] == "" {
			return false
		}
		if jKeys[idx] == "" {
			return true
		}

		return iKeys[idx] < jKeys[idx]
	}

	return false
}

// Descriptors returns envoyproxy/ratelimit formatted descriptors for the supplied limits
func (s limits) Descriptors(names []string) []*YamlDescriptor {
	descriptors := make([]*YamlDescriptor, 0, len(s))
	descriptorsMap := make(map[string]*YamlDescriptor)

	for _, l := range s {
		keys := strings.Split((l.Key), delimiter)
		aggregateKey := ""
		for i, key := range keys {
			prevKey := aggregateKey
			aggregateKey = aggregateKey + delimiter + key
			var desc *YamlDescriptor
			var ok bool
			if desc, ok = descriptorsMap[aggregateKey]; !ok {
				desc = &YamlDescriptor{
					Key:   names[i],
					Value: key,
				}
				descriptorsMap[aggregateKey] = desc
				if i == 0 {
					descriptors = append(descriptors, desc)
				} else {
					descriptorsMap[prevKey].Descriptors = append(descriptorsMap[prevKey].Descriptors, desc)
				}
			}

			if i == len(keys)-1 {
				desc.RateLimit = l.Value
			}
		}
	}

	return descriptors
}
