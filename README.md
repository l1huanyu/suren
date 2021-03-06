# suren
golang微信公众号开发框架，苏伦（Suren），取自《文理双修》女主角 苏伦·曼彻斯特·斯托克。

### Install
```text
go get -u -v github.com/l1huanyu/suren
```

### Example
```go
package main

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"suren"
	"time"
)

func main() {
	e := echo.New()
	s := suren.New("APP_ID", "SECRET", "TOKEN")
	e.GET("/wechat", func(c echo.Context) error {
		echostr := c.QueryParam("echostr")
		if ok, err := s.CheckSignature(&suren.Signature{
			Signature: c.QueryParam("signature"),
			Timestamp: c.QueryParam("timestamp"),
			Nonce:     c.QueryParam("nonce"),
			Echostr:   echostr,
		}); ok && err != nil {
			return c.String(http.StatusOK, echostr)
		}
		return c.NoContent(http.StatusAccepted)
	})
	e.POST("/wechat", func(c echo.Context) error {
		msgRx := new(suren.TextMsgRx)
		err := c.Bind(msgRx)
		if err != nil {
			return err
		}
		msgTx := &suren.TextMsgTx{
			ToUserName:   msgRx.FromUserName,
			FromUserName: msgRx.ToUserName,
			CreateTime:   int(time.Now().Unix()),
			MsgType:      suren.TEXT,
			Content:      fmt.Sprintf("收到消息\"%s\"。", msgRx.Content),
		}
		return c.XML(http.StatusOK, msgTx)
	})
	e.Start(":8823")
}
```
