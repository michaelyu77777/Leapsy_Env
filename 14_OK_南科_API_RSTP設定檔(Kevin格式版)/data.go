package main

import "time"

type Data struct {
	Time  time.Time `json:"time"`
	Score int       `json:"score"`
}
