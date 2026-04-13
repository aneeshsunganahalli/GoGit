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
	Hash []byte
	Mode int
	Exists bool
	IsDirty bool
	IsFile bool
}

type CommitObject struct {
	TreeHash string
	ParentHash string
	Author string
	Committer string
	Timestamp int64
	Timezone string
	Message string
}

type LogData struct {
		hash       string
		authorLine string
		message    []string
	}
	