package entity

type Dump struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Id   int64  `json:"id"`
	// Messages []Messages `json:"messages"`
}

type Message struct {
	Id           int64   `json:"id"`
	Type         string  `json:"type"`
	DateUnixTime string  `json:"date_unixtime"`
	From         string  `json:"from"`
	FromId       string  `json:"from_id"`
	FullText     string  `json:"full_text"`
	Photo        *string `json:"photo"`
	MediaType    *string `json:"media_type"`
}
