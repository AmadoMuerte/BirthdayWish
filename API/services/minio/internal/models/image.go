package models

type ImageRecord struct {
	OriginalName string `json:"originalName"`
	StoragePath  string `json:"storagePath"`
	PublicURL    string `json:"publicURL"`
}
