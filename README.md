# zgit

Explore ZFS Snapshots with a git command line interface, where a ZFS Snapshot represents a git commit:

```
Usage:
  zgit [command]

Available Commands:
  diff        show diff between two snapshots (if only given one argument, show diff between snapshot and working copy)
  log         list all non empty snapshots
  show        show changes of a snapshot
  status      show changes in working copy

Flags:
      --difftool string      use a custom diff program (default "git --no-pager diff --no-index --color=always")
  -i, --ignore string        ignore files matching this pattern in the diff
      --max-diff-count int   maximum number of changes per snapshot (default 50)
  -v, --verbose count        verbosity
```
