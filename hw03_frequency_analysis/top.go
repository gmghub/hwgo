package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(s string) []string {
	normalizedMap := make(map[string]int, len(s)/2)
	for _, v := range strings.Fields(s) {
		vn := strings.ToLower(v)
		vn = strings.Trim(vn, ",.!?()'")
		if vn == "-" {
			continue
		}
		normalizedMap[vn] += 1
	}

	strFreqSlice := make([]struct {
		str  string
		freq int
	}, len(normalizedMap))
	i := 0
	for k, v := range normalizedMap {
		strFreqSlice[i].str = k
		strFreqSlice[i].freq = v
		i++
	}

	sort.Slice(strFreqSlice, func(i, j int) bool {
		if strFreqSlice[i].freq > strFreqSlice[j].freq {
			return true
		} else if strFreqSlice[i].freq < strFreqSlice[j].freq {
			return false
		}
		if strings.Compare(strFreqSlice[i].str, strFreqSlice[j].str) < 0 {
			return true
		}
		return false
	})

	res := []string{}
	for i := 0; i < 10; i++ {
		if i > len(strFreqSlice)-1 {
			break
		}
		res = append(res, strFreqSlice[i].str)
	}
	return res
}
