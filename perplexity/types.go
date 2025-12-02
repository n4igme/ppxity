package perplexity

type Response struct {
	Output         string  `json:"output"`
	Final          bool    `json:"final"`
	ElapsedTime    float64 `json:"elapsed_time"`
	TokensStreamed int     `json:"tokens_streamed"`
	Status         string  `json:"status"`
}
type Message struct {
	Role     string `json:"role"`
	Content  string `json:"content"`
	Priority int    `json:"priority,omitempty"` // Only used in the old socket.io API
}
