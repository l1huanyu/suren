package suren

const (
	getAccessTokenUrl = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	getCallBackUrl    = "https://api.weixin.qq.com/cgi-bin/getcallbackip?access_token=%s"
)

const (
	TEXT = "text"
)