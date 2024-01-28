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
	Source string // path, not pointing to ZFS snapshot
	Dest   string
	PathA  string // current path on disk, could point to /tank/.zfs/snaphot/...
	PathB  string // current path on disk, could point to /tank/.zfs/snaphot/...
}
