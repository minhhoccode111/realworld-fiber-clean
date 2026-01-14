package entity

// ProfilePreview captures public-facing user profile details.
type ProfilePreview struct {
	Username       string `json:"username"`
	Bio            string `json:"bio"`
	Image          string `json:"image"`
	Following      bool   `json:"following"`
	FollowersCount int    `json:"followersCount"`
}
