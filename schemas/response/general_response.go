package response

type GeneralResponseNoData struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}