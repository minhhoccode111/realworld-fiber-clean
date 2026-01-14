package response

// Error represents an error message returned to the client.
type Error struct {
	Error string `json:"error" example:"message"`
}
