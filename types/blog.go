package types

import "time"

type PostInfo struct {
	Title       string    `toml:"title"`
	Description string    `toml:"description"`
	Category    string    `toml:"category"`
	Date        time.Time `toml:"date"`
	Author      string    `toml:"author"`
	Email       string    `toml:"email"`
	FileName    string
}
