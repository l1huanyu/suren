package suren

import (
	"crypto/sha1"
	"fmt"
	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
	"io"
	"sort"
	"strings"
	"time"
)

type (
	Suren struct {
		appID       string
		secret      string
		token       string
		accessToken string
		expiresIn   int
		refreshed   chan byte
	}
)

var (
	accessTokenErr = [...]int{40001, 40014, 41001, 42001}
)

func New(appID, secret, token string) *Suren {
	s := &Suren{
		appID:     appID,
		secret:    secret,
		token:     token,
		expiresIn: 7200,
	}
	s.refreshed = make(chan byte)
	err := s.updateAccessToken()
	if err != nil {
		logrus.Warn(err.Error())
	}
	go s.refresher()
	return s
}

func (s *Suren) updateAccessToken() error {
	r, err := grequests.Get(fmt.Sprintf(getAccessTokenUrl, s.appID, s.secret), nil)
	if err != nil {
		return err
	}
	wxr := new(getAccessTokenResp)
	err = r.JSON(wxr)
	if err != nil {
		return err
	}
	if wxr.ErrCode != 0 {
		s.checkAccessToken(wxr.ErrCode)
		return fmt.Errorf(wxr.ErrMsg)
	}
	s.accessToken = wxr.AccessToken
	s.expiresIn = wxr.ExpiresIn
	//重置定时器
	s.refreshed <- 1
	return nil
}

func (s *Suren) refresher() {
	for {
		select {
		case <-time.After(time.Duration(s.expiresIn) * time.Second):
			err := s.updateAccessToken()
			if err != nil {
				logrus.Warn(err.Error())
			}
		case <-s.refreshed:
			continue
		}
	}
}

func (s *Suren) checkAccessToken(errcode int) {
	for _, code := range accessTokenErr {
		if code == errcode {
			time.Sleep(1 * time.Second)
			s.updateAccessToken()
		}
	}
}

func (s *Suren) AppID() string {
	return s.appID
}

func (s *Suren) Secret() string {
	return s.secret
}

func (s *Suren) Token() string {
	return s.token
}

func (s *Suren) AccessToken() string {
	return s.accessToken
}

func (s *Suren) ExpiresIn() int {
	return s.expiresIn
}

//微信参数校验
func (s *Suren) CheckSignature(sig *Signature) (bool, error) {
	sl := []string{s.token, sig.Timestamp, sig.Nonce}
	//升序排序
	sort.Strings(sl)
	//sha1加密
	sh1 := sha1.New()
	_, err := io.WriteString(sh1, strings.Join(sl, ""))
	if err != nil {
		return false, err
	}
	localSignature := fmt.Sprintf("%x", sh1.Sum(nil))
	if localSignature == sig.Signature {
		return true, nil
	}
	return false, nil
}

//获取微信服务器IP地址
func (s *Suren) GetCallBackIP() ([]string, error) {
	r, err := grequests.Get(fmt.Sprintf(getCallBackUrl, s.accessToken), nil)
	if err != nil {
		return nil, err
	}
	wxr := new(getCallBackIpResp)
	err = r.JSON(wxr)
	if err != nil {
		return nil, err
	}
	if wxr.ErrCode != 0 {
		s.checkAccessToken(wxr.ErrCode)
		return nil, fmt.Errorf(wxr.ErrMsg)
	}
	return wxr.IpList, nil
}
