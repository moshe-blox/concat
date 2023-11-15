# concat

Walks through a directory, concatenating files based on specified extensions, and exclude certain files or directories.

`.gitignore` and `.dockerignore` files are respected.

## Usage

```bash
concat [directory] -x [.ext1, .ext2, ...] -e [exclude_pattern1, exclude_pattern2, ...]
```

- `directory`: The directory to walk through.
- `-x`: Specify file extensions to include. Leave blank to not filter by extension.
- `-e`: Exclude files or directories based on patterns.

## Installation

```bash
go get github.com/your-repo/directory-walker-concatenator
```

## Example

```bash
concat /path-to-dir --exclude .go --ext *.test
```

_Explanation: walk through `/path-to-dir`, concatenate all `.go` files, excluding files matching `*_test`._

Output:

````bash
# /path-to-dir/a.go
```go
package main

func a() {
    // ...
}
```

# /path-to-dir/b.go
```go
package main

func b() {
    // ...
}
```
````
