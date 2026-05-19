package entities

type APIResponse struct {
	Message string      `json:"message" description:"Success or error message"`
	Data    interface{} `json:"data,omitempty" description:"The actual response payload"`
}
