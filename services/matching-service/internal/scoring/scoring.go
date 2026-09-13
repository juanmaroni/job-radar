package scoring

import "strings"

func CalculateScore(description string, keywords []string) int32 {
	if len(description) == 0 {
		return 0
	}

	description = strings.ToLower(description)
	var count int32 = 0

	for _, k := range keywords {
		k = strings.ToLower(k)

		if strings.Contains(description, k) {
			count += 1
		}
	}

    return count
}
