package common_types

// SuccessResponse : struct that represents success response of all requests
type SuccessResponse struct {
	StatusCode    int    `json:"status_code"`
	Msg           string `json:"msg"`
	FnName        string `json:"fn_name"`
	TrustVerified bool   `json:"trust_verified,omitempty"`
}

// ErrorResponse : struct that represents error response of all requests
type ErrorResponse struct {
	StatusCode    int    `json:"status_code"`
	ErrorMsg      string `json:"error_msg"`
	FnName        string `json:"fn_name,omitempty"`
	TrustVerified *bool  `json:"trust_verified,omitempty"`
}
