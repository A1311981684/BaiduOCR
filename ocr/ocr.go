package ocr

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"ocr/token"
	"strings"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

//百度高精度OCR接口URL
//Baidu high accuracy OCR API url
const ocrURL = "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic?access_token="

//请求返回的Response结构类型
//Response type
type Response struct {
	OgId      uint64 `json:"og_id"`
	Direction int32  `json:"direction"`

	//WordsResult是我们真正想要的识别结果
	//切片类型（slice，[]）是因为返回的结果是按照图片中识别出来的每一行文字返回的
	//map类型（key = string， value = interface{}）是因为我会要求返回每一行的置信度。获取该行的执行度用probability作为key从map中获取，
	//获取识别结果文字就用words作为key。

	//WordsResult is what we truly care about
	//The slice type exists cause the result returned contains multiple text lines as the image indicates
	//The map element exists cause I require Probability of each line; If you are interested in Probability, retrieve it
	// from the map by querying with key 'probability'; And result is by key 'words'.
	WordsResult    []map[string]interface{} `json:"words_result"`
	WordsResultNum uint32                   `json:"words_result_num"`
	Words          string                   `json:"+words"`
	Probability    float32                  `json:"probability"`
}

func OCR(imagePath string) ([]map[string]interface{}, error) {
	//百度OCR要求使用base64字符串表示图片
	//Baidu OCR requires base64 string type images
	base64Image, err := OpenImageFileToBase64(imagePath)
	if err != nil {
		return nil, err
	}
	//发送识别请求
	//Send request
	body := strings.NewReader(base64Image)
	if token.GetToken() == "" {
		return  nil, errors.New("get token error, check clientId or secretKey")
	}
	r, err := http.Post(ocrURL+token.GetToken(), "application/x-www-form-urlencoded", body)
	if err != nil {
		return nil, err
	}
	//Get response
	var resp Response
	byteRes, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(byteRes, &resp)
	if err != nil {
		return nil, err
	}

	return resp.WordsResult, nil
}

func OpenImageFileToBase64(imagePath string) (string, error) {
	fileContents, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return "", err
	}
	//实现过程中最坑的部分，官方文档说要将图片base64化，然后对base64结果进行UrlEncode。我以为base64自带的URLEncode可以搞，但是转出来的结果格
	// 式不正确，使用url.Values{}.Encode()可以解决URLEncode的问题
	bs64 := base64.StdEncoding.EncodeToString(fileContents) //Make image base64
	u := url.Values{}
	u.Set("image", bs64)
	u.Set("detect_direction", "true")
	u.Set("probability", "true")
	//Construct URLEncoded string
	return u.Encode(), nil
}
