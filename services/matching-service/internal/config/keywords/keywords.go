package keywords

import (
	"errors"
	"jobradar/matching-service/internal/config/environment"
	"strings"
)

func Load() ([]string, error) {
	keywordsEnv, err := environment.GetEnvVar("MATCHING_SERVICE_PROFILE_KEYWORDS")

	if err != nil {
		return nil, err
	}

	if keywordsEnv == "" {
		return nil, errors.New("no keywords configured")
	}

	keywords := strings.Split(keywordsEnv, ",")
	var keywordsTrimmed []string

	for _, k := range(keywords) {
		keywordsTrimmed = append(keywordsTrimmed, strings.TrimSpace(k))
	}


	return keywordsTrimmed, nil
}
