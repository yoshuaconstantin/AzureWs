package models

type UserProfileDataModel struct {
	Nickname *string `json:"nickname,omitempty"`
	Age      *string `json:"age,omitempty"`
	Gender   *string `json:"gender,omitempty"`
	ImageUrl *string `json:"image_url,omitempty"`
}

type GetUserProfileDataModel struct {
	Nickname     *string `json:"nickname,omitempty"`
	Age          *string `json:"age,omitempty"`
	Gender       *string `json:"gender,omitempty"`
	ImageUrl     *string `json:"image_url,omitempty"`
	CreatedSince *string `json:"created_since,omitempty"`
}
