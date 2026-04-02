package internal


type Blob []byte

// [mode] [type] [encoded hash] [filename]
type IndexEntry struct {
	Mode int64 // string representations of octal numbers, and also denotes what the type is
	Hash string
	Mtime int64
	Size int64
	Filename string
}

type Tree struct {
	Entries []IndexEntry
}
