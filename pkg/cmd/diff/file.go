package diff

import (
	"bytes"
	"io"
	"os"

	"github.com/bmatcuk/doublestar/v4"
)

// matchPatterns returns true if path should be included
func matchPatterns(includePatterns []string, excludePatterns []string, path string) bool {
	// if any include pattern matches, continue
	includeMatched := false
	for _, pattern := range includePatterns {
		match, err := doublestar.Match(pattern, path)
		if err != nil {
			log.Error(err, "error matching file pattern", "file", path)
			return false
		}

		if match {
			includeMatched = true
			break
		}
	}
	if len(includePatterns) > 0 && !includeMatched {
		return false
	}

	// if any exclude pattern matches, return false
	for _, pattern := range excludePatterns {
		match, err := doublestar.Match(pattern, path)
		if err != nil {
			log.Error(err, "error matching file pattern", "file", path)
			return false
		}

		if match {
			return false
		}
	}

	return true
}

func compareModificationTime(a string, b string) (bool, error) {
	stA, err := os.Stat(a)
	if err != nil {
		return false, err
	}

	stB, err := os.Stat(b)
	if err != nil {
		return false, err
	}

	return stA.ModTime().Equal(stB.ModTime()), nil
}

func compareFileSize(a string, b string) (bool, error) {
	stA, err := os.Stat(a)
	if err != nil {
		return false, err
	}

	stB, err := os.Stat(b)
	if err != nil {
		return false, err
	}

	return stA.Size() == stB.Size(), nil
}

func compareFilesByteForByte(file1 string, file2 string, chunkSize int) (bool, error) {
	f1, err := os.Open(file1)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false, err
	}
	defer f2.Close()

	for {
		b1 := make([]byte, chunkSize)
		n1, err1 := f1.Read(b1)
		if err1 != nil && err1 != io.EOF {
			return false, err1
		}

		b2 := make([]byte, chunkSize)
		n2, err2 := f2.Read(b2)
		if err2 != nil && err2 != io.EOF {
			return false, err2
		}

		// exit if we're at EOF of a file
		if err1 == io.EOF && err2 == io.EOF {
			return true, nil
		} else if err1 == io.EOF || err2 == io.EOF {
			return false, nil
		}

		if n1 != n2 || !bytes.Equal(b1, b2) {
			return false, nil
		}
	}
}
