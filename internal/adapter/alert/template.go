package alert

type ErrorMessage struct {
	Title        string
	Description  string
	Request      interface{}
	ErrorMessage string
}

type alertRequest struct {
	Embeds []embeds `json:"embeds"`
}

type embeds struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Color       int          `json:"color"`
	Fields      []embedField `json:"fields"`
}

type embedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
