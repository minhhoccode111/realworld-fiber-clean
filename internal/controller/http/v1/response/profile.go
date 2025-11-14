package response

import "github.com/minhhoccode111/realworld-fiber-clean/internal/entity"

type ProfilePreviewResponse struct {
	Profile *entity.ProfilePreview `json:"profile"`
}
