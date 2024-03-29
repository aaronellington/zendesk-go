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

type BusinessRuleConditions struct {
	All []BusinessRuleCondition `json:"all"`
	Any []BusinessRuleCondition `json:"any"`
}

type BusinessRuleCondition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

type BusinessRuleAction struct {
	Field string `json:"field"`
	Value any    `json:"value"`
}

type JobStatusResponse struct {
	JobStatus JobStatus `json:"job_status"`
}

type JobStatus struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Total    uint64 `json:"total"`
	Progress uint64 `json:"progress"`
	Status   string `json:"status"`
	Message  string `json:"message"`
	Results  any    `json:"results"`
}
