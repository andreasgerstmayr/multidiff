package cli

type Options struct {
	Verbosity           int
	ShowPathOnly        bool
	NewFilesAsEmpty     bool
	ShowMetadataChanges bool
	CompareByteForByte  bool
	IncludePatterns     []string
	ExcludePatterns     []string

	Conv struct {
		Exif bool
	}
}

type OptionsKey struct{}
