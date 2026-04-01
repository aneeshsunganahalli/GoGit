package internal

type Blob []byte

// [mode] [type] [encoded hash] [filename]
type TreeEntry struct {
	Mode string // string representations of octal numbers
	Type string
	Hash string
	Filename string
}

type Tree struct {
	Entries []TreeEntry
}
