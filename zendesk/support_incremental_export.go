package zendesk

import "time"

type IncrementalExportResponse struct {
	EndTimeUnix int64 `json:"end_time"`
	EndOfStream bool  `json:"end_of_stream"`
}

func (response IncrementalExportResponse) EndTime() time.Time {
	return time.Unix(int64(response.EndTimeUnix), 0)
}
