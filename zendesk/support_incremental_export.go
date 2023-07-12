package zendesk

type IncrementalExportResponse struct {
	EndTime     uint64 `json:"end_time"`
	EndOfStream bool   `json:"end_of_stream"`
}
