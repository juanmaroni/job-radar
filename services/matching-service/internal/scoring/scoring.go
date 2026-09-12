package scoring

func CalculateScore(role string) int32 {
    return int32(len(role) % 100)
}
