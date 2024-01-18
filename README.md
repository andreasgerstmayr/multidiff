# multidiff
multidiff compares two sources (files/directories/ZFS snapshots) in a human readable way (e.g. comparing metadata of images)

```
Usage:
  multidiff [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  diff        show diff between two sources
  help        Help about any command
  zgit        zgit allows exploring ZFS snapshots like git repositories

Flags:
  -b, --byte                  compare byte for byte (default true)
      --conv.exif             compare EXIF metadata of images
  -e, --exclude stringArray   exclude files matching this pattern
  -h, --help                  help for multidiff
  -i, --include stringArray   include files matching this pattern
  -N, --new-file              show absent files as empty
      --path                  show paths only
  -m, --show-meta             include file metadata modifications in diff
  -v, --verbose count         set verbosity (can be used multiple times)
```

## zgit subcommand
zgit allows exploring ZFS snapshots like git repositories

```
Usage:
  multidiff zgit [command]

Available Commands:
  diff        show diff between two snapshots
  log         list all non empty snapshots
  show        show changes of a snapshot
  status      show changes in working copy

Flags:
  -h, --help   help for zgit

Global Flags:
  -b, --byte                  compare byte for byte (default true)
      --conv.exif             compare EXIF metadata of images
  -e, --exclude stringArray   exclude files matching this pattern
  -i, --include stringArray   include files matching this pattern
  -N, --new-file              show absent files as empty
      --path                  show paths only
  -m, --show-meta             include file metadata modifications in diff
  -v, --verbose count         set verbosity (can be used multiple times)
```
