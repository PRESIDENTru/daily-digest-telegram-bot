package models

type Translation struct {
	TextData ResponseData `json:"responseData"`
}

type ResponseData struct {
	Text string `json:"translatedText"`
}
