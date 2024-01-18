package diff

type Change string

const (
	Modified Change = "modified"
	Created  Change = "created"
	Removed  Change = "removed"
	Renamed  Change = "renamed" // can also be modified
)

type Diff struct {
	Change Change
	Source string
	Dest   string
	PathA  string
	PathB  string
}
