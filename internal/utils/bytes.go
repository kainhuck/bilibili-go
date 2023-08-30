package utils

const (
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
)

// SplitBytes 将 total 按 chunk 分区
func SplitBytes(total []byte, chunk int) [][]byte {
	if chunk <= 0 {
		return nil
	}

	var result [][]byte

	for len(total) > 0 {
		currentChunk := chunk
		if len(total) < chunk {
			currentChunk = len(total)
		}
		result = append(result, total[:currentChunk])
		total = total[currentChunk:]
	}

	return result
}
