package util

import (
	"bufio"
	"os"
)

// ReadLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func ChunkSlice[T any](d []T, size int) [][]T {
	var chunks [][]T
	for i := 0; i < len(d); i += size {
		end := i + size

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(d) {
			end = len(d)
		}

		chunks = append(chunks, d[i:end])
	}

	return chunks
}
