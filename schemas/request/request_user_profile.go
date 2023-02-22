package request

type RequestUploadImageData struct {
	Token string `json:"token"`
	Data  []byte `json:"data"`
}

type RequestUpdateProfileImageData struct {
	Token       string `json:"token"`
	OldImageUrl string `json:"oldImgUrl"`
	Data        []byte `json:"data"`
}

type RequestDeleteProfileImageData struct {
	Token       string `json:"token"`
	OldImageUrl string `json:"oldImgUrl"`
}

type RequestTokenOnlyData struct {
	Token string `json:"token"`
}

type RequestInsertProfileData struct {
	Token string `json:"token"`
	Data  []struct {
		Nickname string `json:"nickname,omitempty"`
		Age      string `json:"age,omitempty"`
		Gender   string `json:"gender,omitempty"`
		Nation	 string `json:"nation,omitempty"`
		ImageUrl string `json:"image_url,omitempty"`
	} `json:"data"`
}