package response

import "github.com/minhhoccode111/realworld-fiber-clean/internal/entity"

// ProfilePreviewResponse wraps a profile preview.
type ProfilePreviewResponse struct {
	Profile *entity.ProfilePreview `json:"profile"`
}
