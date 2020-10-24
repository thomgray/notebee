package model

// SearchResultItem ...
type SearchResultItem struct {
	File  *File
	Path  FilePath
	Score int
}
