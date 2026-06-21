package types

type JSONResult struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

type RecordPayload struct {
	BabyID          uint                   `json:"baby_id"`
	BabyName        string                 `json:"baby_name"`
	ActionType      string                 `json:"action_type"`
	ActionTime      string                 `json:"action_time"`
	Summary         string                 `json:"summary"`
	Data            map[string]interface{} `json:"data"`
	Source          string                 `json:"source"`
	ClientRequestID string                 `json:"client_request_id"`
}
