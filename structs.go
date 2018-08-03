package suren

//request
type ()

//Response
type (
	errResp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	getCallBackIpResp struct {
		errResp
		IpList []string `json:"ip_list"`
	}

	getAccessTokenResp struct {
		errResp
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
)

//Data
type (
	Signature struct {
		Signature string
		Timestamp string
		Nonce     string
		Echostr   string
	}

	//接受普通消息-文本消息
	TextMsgRx struct {
		ToUserName   string `xml:"ToUserName"`   //开发者微信号
		FromUserName string `xml:"FromUserName"` //发送方账号（一个OpenID）
		CreateTime   int    `xml:"CreateTime"`   //消息创建时间
		MsgType      string `xml:"MsgType"`      //text
		Content      string `xml:"Content"`      //文本消息内容
		MsgId        int64  `xml:"MsgId"`        //消息id
	}

	//被动回复-文本消息
	TextMsgTx struct {
		ToUserName   string //接收方账号（收到的OpenID）
		FromUserName string //开发者微信号
		CreateTime   int    //消息创建时间
		MsgType      string //text
		Content      string //回复的消息内容（可换行）
	}
)
