package models

import "time"

type WindowInfo struct {
	Title  string `json:"title"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type CursorInfo struct {
	ID string `json:"id"`
	X  int    `json:"x"`
	Y  int    `json:"y"`
}

type MessageInfo struct {
	ID    string `json:"id"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
	GridX int    `json:"gridX"`
	GridY int    `json:"gridY"`
}

type ClientInfo struct {
	ID          string
	ConnectedAt time.Time
}
