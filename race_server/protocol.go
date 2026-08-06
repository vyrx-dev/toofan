package main

import "encoding/json"

type ServerMsg struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type JoinMsg struct {
	Room       string   `json:"room"`
	Players    []string `json:"players"`
	Online     int      `json:"online"`
	Host       string   `json:"host"`
	AutoStart  bool     `json:"auto_start"`
	CanStart   bool     `json:"can_start"`
	Difficulty string   `json:"difficulty"`
	Mode       string   `json:"mode"`
	Lang       string   `json:"lang"`
	Duration   int      `json:"duration"`
	IsPrivate  bool     `json:"is_private"`
	State      string   `json:"state"`
	Text       string   `json:"text,omitempty"`
	TimeLeft   int      `json:"time_left,omitempty"`
}

type CountdownMsg struct {
	Seconds int `json:"seconds"`
}

type StartMsg struct {
	Text       string `json:"text"`
	Difficulty string `json:"difficulty"`
	Mode       string `json:"mode"`
	Lang       string `json:"lang"`
	Duration   int    `json:"duration"`
}

type PlayerProgress struct {
	Name     string  `json:"name"`
	Progress float64 `json:"progress"`
	WPM      float64 `json:"wpm"`
	Finished bool    `json:"finished"`
}

type ProgressMsg struct {
	Players []PlayerProgress `json:"players"`
}

type FinishMsg struct {
	Placements []PlayerProgress `json:"placements"`
}

type OnlineMsg struct {
	Count int `json:"count"`
}

type ProgressUpdate struct {
	Name     string  `json:"name"`
	Room     string  `json:"room"`
	Progress float64 `json:"progress"`
	WPM      float64 `json:"wpm"`
}

type StartRequest struct {
	Name string `json:"name"`
	Room string `json:"room"`
}

func marshal(msg ServerMsg) []byte {
	data, _ := json.Marshal(msg)
	return data
}
