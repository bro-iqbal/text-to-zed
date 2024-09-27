package domain

type ParsedContent struct {
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	Values   []Values `json:"values"`
	Comments []string `json:"comments"`
}

type Values struct {
	Comments []string `json:"comments"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Value    string   `json:"value"`
}
