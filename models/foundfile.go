package models

import (
	"database/sql"
	"log"
	"time"
)

// FoundFile ...
type FoundFile struct {
	Source      string
	Path        string
	Md5hash     string
	Name        string
	Extension   string
	Type        string
	Size        int64
	Category    string
	Label       string
	Modified    time.Time
	Discovered  time.Time
	LastChecked time.Time
}

// CreateFoundFileTable ...
func CreateFoundFileTable(db *sql.DB) {
	const sql = `
		CREATE TABLE if not exists found_files (
			source TEXT NOT NULL,
			path TEXT NOT NULL,
			md5hash TEXT NOT NULL,
			name TEXT NOT NULL,
			size int NOT NULL,
			modified TIMESTAMP NOT NULL,
			extension TEXT NOT NULL,
			type TEXT NOT NULL DEFAULT '',
			category TEXT NOT NULL DEFAULT '',
			label TEXT NOT NULL DEFAULT '',
			tags TEXT NOT NULL DEFAULT '',
			notes TEXT NOT NULL DEFAULT '',
			discovered TIMESTAMP NOT NULL,
			last_checked TIMESTAMP NOT NULL,
			unique(source, path, md5hash)
	    )`
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
}

// GetFoundFile ...
func GetFoundFile(db *sql.DB, source string, path string) *FoundFile {
	// FIXME: Not following go pattern, need to use interface
	const sql = `
		SELECT source, path, md5hash, name, size, modified, extension, type, category, label, discovered, last_checked
		FROM found_files WHERE source = ? and path = ?`
	rows, err := db.Query(sql, source, path)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var source, path, md5hash, name, extension, fileType, category, label string
		var modified, lastChecked, discovered time.Time
		var size int64
		err = rows.Scan(&source, &path, &md5hash, &name, &size, &modified, &extension, &fileType, &category, &label, &discovered, &lastChecked)
		if err != nil {
			log.Fatal(err)
		}
		return &FoundFile{Source: source, Path: path, Md5hash: md5hash, Name: name, Extension: extension, Type: fileType, Size: size, Modified: modified, Category: category, Label: label, Discovered: discovered, LastChecked: lastChecked}
	}
	return nil
}

// Save ...
func (ff *FoundFile) Save(db *sql.DB) {
	// If the file changes, it is considered a different file, even if it is in the same path.
	const sql = `
		INSERT INTO found_files (source, path, md5hash, name, extension, type, size, modified, discovered, last_checked, category, label)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (source, path, md5hash) DO UPDATE SET
			name=excluded.name,
			type=excluded.type,
			extension=excluded.extension,
			size=excluded.size,
			modified=excluded.modified,
			discovered=excluded.discovered,
			last_checked=excluded.last_checked,
			category=excluded.category,
			label=excluded.label`
	_, err := db.Exec(sql, ff.Source, ff.Path, ff.Md5hash, ff.Name, ff.Extension, ff.Type, ff.Size, ff.Modified, ff.Discovered, ff.LastChecked, ff.Category, ff.Label)
	if err != nil {
		log.Panic(err)
	}
}
