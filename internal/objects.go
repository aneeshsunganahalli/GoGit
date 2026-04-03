package internal

type Blob []byte

// [mode] [type] [encoded hash] [filename]
type IndexEntry struct {
	Mode int // string representations of octal numbers, and also denotes what the type is
	Hash string
	Mtime int64
	Size int64
	Filename string
}

// Temporary Trie Structure Node
type TrieNode struct {
	Children map[string]*TrieNode // Only for directories, since files can't have children
	Hash string
	Mode int
	IsFile bool
}
