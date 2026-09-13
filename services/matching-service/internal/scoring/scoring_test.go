package scoring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateScore_WithMatches(t *testing.T) {
    description := "We use Docker and PostgreSQL extensively in our stack."
    keywords := []string{"Docker", "PostgreSQL", "Kafka"}

    score := CalculateScore(description, keywords)

    assert.Equal(t, int32(2), score)
}
