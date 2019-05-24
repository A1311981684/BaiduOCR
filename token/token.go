package token

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	en_decrypt "ocr/en-decrypt"
	"os"
	"strconv"
	"time"
)

//API key
var clientId = "" //Your API key

//Secret key
var secretKey = "" //Your Secret key of app

const tokenJson = `token.json`

//Get token url
var tokenRequestUrl = "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=" + clientId +
	"&client_secret=" + secretKey

type ResponseOfToken struct {
	RefreshToken  string `json:"refresh_token"`
	ExpiresIn     int    `json:"expires_in"`
	Scope         string `json:"scope"`
	SessionKey    string `json:"session_key"`
	AccessToken   string `json:"access_token"`
	SessionSecret string `json:"session_secret"`
}

type StorageToken struct {
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
}

var _token string

func GetToken() string {
	if !checkTokenAvailableInJson() {
		_token = ""
		err := UpdateToken()
		if err != nil {
			return ""
		}
	}

	return _token
}

//Refresh local token json
func UpdateToken() error {
	t, n, err := getToken()
	if err != nil {
		return err
	}
	f, err := os.Create(tokenJson)
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()
	teAfter, err := time.ParseDuration(strconv.Itoa(n) + "s")
	if err != nil {
		return err
	}
	te := time.Now().Add(teAfter).Format("2006-01-02 15:04:05")
	var st = StorageToken{
		RefreshToken: t,
		ExpiresIn:    te,
	}
	bytesData, err := json.MarshalIndent(&st, "", "	")
	if err != nil {
		return err
	}
	encrypted, err := en_decrypt.EncryptText(string(bytesData))
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(encrypted))
	if err != nil {
		panic(err)
	}
	_token = t
	return nil
}

//Check if token is out-date or not available
func checkTokenAvailableInJson() bool {
	f, err := os.Open(tokenJson)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()

	byteData, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	if len(byteData) == 0 {
		return false
	}

	decrypted := en_decrypt.DecryptText(string(byteData))
	var ct StorageToken
	err = json.Unmarshal([]byte(decrypted), &ct)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	if ct.ExpiresIn == "" {
		return false
	}

	te, err := time.Parse("2006-01-02 15:04:05", ct.ExpiresIn)
	if err != nil {
		return false
	}
	tn, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	if te.Sub(tn) > 0 {
		_token = ct.RefreshToken
		return true
	}
	return false
}

//get token
func getToken() (string, int, error) {
	if clientId == "" || secretKey == "" {
		return "", 0, errors.New("check your clientId or secretKey variable")
	}
	resp, err := http.Get(tokenRequestUrl)
	if err != nil {
		return "", 0, err
	}

	var t ResponseOfToken
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	if len(body) == 0 {
		return "", 0, errors.New("no response")
	}
	err = json.Unmarshal(body, &t)
	if err != nil {
		return "", 0, err
	}
	return t.AccessToken, t.ExpiresIn, nil
}
