package gsms

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

import (
	"github.com/sanxia/glib"
)

/* ================================================================================
 * 野狗短信发送
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	getewayUrl struct {
		Code   string
		Notify string
		Check  string
	}

	yegouSms struct {
		Geteway    string     `form:"geteway" json:"geteway"`
		Url        getewayUrl `form:"geteway_url" json:"geteway_url"`
		AppKey     string     `form:"app_key" json:"app_key"`
		AppSecret  string     `form:"app_secret" json:"app_secret"`
		Mobiles    []string   `form:"mobiles" json:"mobiles"`
		TemplateId string     `form:"template_id" json:"template_id"`
		Params     []string   `form:"params" json:"params"`
		Type       string     `form:"type" json:"type"`
		Timestamp  string     `form:"timestamp" json:"timestamp"`
	}

	yegouErrorResult struct {
		Errcode int    `form:"errcode" json:"errcode"`
		Message string `form:"message" json:"message"`
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 创建野狗短信提供者
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewYeGouSms(appKey, appSecret string) SmsProvider {
	sms := new(yegouSms)
	sms.AppKey = appKey
	sms.AppSecret = appSecret
	sms.Geteway = "https://sms.wilddog.com/api/v1/"
	sms.Url = getewayUrl{
		Code:   "/code/send",
		Notify: "/notify/send",
		Check:  "/code/check",
	}
	sms.Type = "code"

	timespame := glib.UnixTimestamp() * int64(1000)
	sms.Timestamp = fmt.Sprintf("%d", timespame)

	return sms
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 发送手机信息
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *yegouSms) Send(mobile string) (*SmsResult, error) {
	result := new(SmsResult)

	if len(mobile) == 0 {
		return result, errors.New("手机号不能为空")
	}

	if len(s.TemplateId) == 0 {
		return nil, errors.New("参数不正确")
	}

	//接收手机号码
	s.Mobiles = []string{mobile}

	//签名请求参数
	requestString := s.GetRequestString()

	url := s.Url.Code
	if strings.ToLower(s.Type) == "notify" {
		url = s.Url.Notify
	}
	url = s.Geteway + s.AppKey + url

	//发起Http请求
	if response, err := glib.HttpPost(url, requestString); err != nil {
		result.Message = err.Error()
		return result, err
	} else {
		result.Message = response

		//错误处理
		//{"status" : "ok","data":{"rrid":"bdd977d825084bd0ad7a00597dbd0f69"}}
		//{"errcode": 79998,"message": "request error ,error is null"}
		var errorResult *yegouErrorResult
		glib.FromJson(response, &errorResult)
		if len(errorResult.Message) > 0 {
			result.IsSuccess = false
			return result, errors.New(errorResult.Message)
		} else {
			result.IsSuccess = true
		}
	}

	return result, nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置模版Id
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *yegouSms) SetGeteway(geteway string) {
	s.Geteway = geteway
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置模版Id
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *yegouSms) SetTemplateCode(templateCode string) {
	s.TemplateId = templateCode
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置模版参数
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *yegouSms) SetTemplateParam(templateParam SmsTemplateParam) {
	s.Params = []string{templateParam.Code}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置模版参数字符串
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *yegouSms) SetTemplateString(templateString string) {
	s.Params = []string{templateString}
}

func (s *yegouSms) SetSignName(signName string) {

}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取请求字符串
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *yegouSms) GetRequestString() string {
	params := s.toDict()

	//Md5签名串
	sign := s.Sign(params)

	//参数值url编码
	var options []string = make([]string, 0)
	for k, v := range params {
		item := fmt.Sprintf("%s=%s", k, url.QueryEscape(v))
		options = append(options, item)
	}

	//把签名拼接到参数
	options = append(options, fmt.Sprintf("%s=%s", "signature", url.QueryEscape(sign)))

	//用&链接请求参数
	return strings.Join(options, "&")
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 签名算法
 * params里的每个Value都需要进行url编码
 * fmt.Sprintf("%s=%s", key, url.QueryEscape(value))
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *yegouSms) Sign(params map[string]string) string {
	var keys []string = make([]string, 0)
	var values []string = make([]string, 0)

	//请求参数排序（字母升序）
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	//拼接KeyValue字符串
	for _, key := range keys {
		if len(params[key]) > 0 {
			keyValue := fmt.Sprintf("%s=%s", key, params[key])
			values = append(values, keyValue)
		}
	}
	paramString := strings.Join(values, "&")

	//Sha256签名
	signString := fmt.Sprintf("%s&%s", paramString, s.AppSecret)

	return glib.Sha256(signString)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取参数字典
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *yegouSms) toDict() map[string]string {
	var params map[string]string = make(map[string]string, 0)
	params["templateId"] = s.TemplateId

	if strings.ToLower(s.Type) == "notify" {
		params["mobiles"] = glib.StringSliceToString(s.Mobiles)
	} else {
		params["mobile"] = s.Mobiles[0]
	}

	if len(s.Params) > 0 {
		paramsJson, _ := glib.ToJson(s.Params)
		//params["params"] = glib.StringSliceToString(s.Params)
		params["params"] = paramsJson
	}

	params["timestamp"] = s.Timestamp

	return params
}
