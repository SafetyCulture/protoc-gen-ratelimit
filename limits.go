package genratelimit

import (
	"sort"
	"strings"
)

type Limit struct {
	Key   string
	Value *YamlRateLimit
}

type limits []*Limit

func (s limits) Len() int      { return len(s) }
func (s limits) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s limits) Less(i, j int) bool {
	iKeys := strings.Split(s[i].Key, delimiter)
	jKeys := strings.Split(s[j].Key, delimiter)

	for idx := range iKeys {
		if iKeys[idx] == jKeys[idx] {
			continue
		}

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

func sortLimits(limitsMap map[string]*Limit) limits {
	limitArr := make(limits, 0, len(limitsMap))

	for _, l := range limitsMap {
		limitArr = append(limitArr, l)
	}

	sort.Sort(limitArr)

	return limitArr
}

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
