package zendesk

type CustomFieldOptionID uint64

type CustomFieldOption struct {
	ID       CustomFieldOptionID `json:"id"`
	Name     string              `json:"name"`
	Position uint64              `json:"position"`
	RawName  string              `json:"raw_name"`
	URL      string              `json:"url"`
	Value    string              `json:"value"`
}
