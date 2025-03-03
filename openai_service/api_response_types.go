package openaiservice

type CompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

type Choice struct {
	Index        int64   `json:"index"`
	Message      Message `json:"message,omitempty"`
	Delta        Delta   `json:"delta,omitempty"`
	FinishReason string  `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Roles struct {
	System    string
	Assistant string
	User      string
}

type Delta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens          int64               `json:"prompt_tokens"`
	CompletionTokens      int64               `json:"completion_tokens"`
	TotalTokens           int64               `json:"total_tokens"`
	PromptTokensDetails   PromptTokensDetails `json:"prompt_tokens_details"`
	PromptCacheHitTokens  int64               `json:"prompt_cache_hit_tokens"`
	PromptCacheMissTokens int64               `json:"prompt_cache_miss_tokens"`
}

type PromptTokensDetails struct {
	CachedTokens int64 `json:"cached_tokens"`
}
