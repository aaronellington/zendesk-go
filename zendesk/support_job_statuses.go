package zendesk

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
