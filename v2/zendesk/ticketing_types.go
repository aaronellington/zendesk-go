package zendesk

const (
	PriorityUrgent = "urgent"
	PriorityHigh   = "high"
	PriorityNormal = "normal"
	PriorityLow    = "low"
)

type CustomFieldOptionID uint64

type CustomFieldOption struct {
	ID       CustomFieldOptionID `json:"id"`
	Name     string              `json:"name"`
	Position uint64              `json:"position"`
	RawName  string              `json:"raw_name"`
	URL      string              `json:"url"`
	Value    string              `json:"value"`
}

type Tags []Tag

func (tags Tags) HasTag(targetTag Tag) bool {
	for _, tag := range tags {
		if tag == targetTag {
			return true
		}
	}

	return false
}
