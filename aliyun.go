package gsms

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"
	"time"
)

import (
	"github.com/sanxia/glib"
)

/* ================================================================================
 * 阿里云短信发送
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */
type (
	aliyunSms struct {
		Geteway          string `form:"Geteway" json:"Geteway"`                   //网关
		Action           string `form:"Action" json:"Action"`                     //操作接口名，系统规定参数，取值：SingleSendSms
		SignName         string `form:"SignName" json:"SignName"`                 //短信签名
		TemplateCode     string `form:"TemplateCode" json:"TemplateCode"`         //短信模板的模板CODE（状态必须是验证通过）
		RecNum           string `form:"RecNum" json:"RecNum"`                     //目标手机号，多个手机号可以逗号分隔
		ParamString      string `form:"ParamString" json:"ParamString"`           //短信模板中的变量；数字需要转换为字符串
		RegionId         string `form:"RegionId" json:"RegionId"`                 //区域ID
		AccessKeyId      string `form:"AccessKeyId" json:"AccessKeyId"`           //access id
		AccessKeySecret  string `form:"AccessKeySecret" json:"AccessKeySecret"`   //私匙
		Signature        string `form:"Signature" json:"Signature"`               //签名结果串
		SignatureNonce   string `form:"SignatureNonce" json:"SignatureNonce"`     //唯一随机数，用于防止网络重放攻击。用户在不同请求间要使用不同的随机数值
		SignatureMethod  string `form:"SignatureMethod" json:"SignatureMethod"`   //签名方式，目前支持HMAC-SHA1
		SignatureVersion string `form:"SignatureVersion" json:"SignatureVersion"` //签名算法版本，目前版本是1.0
		Format           string `form:"Format" json:"Format"`                     //返回值的类型，支持JSON与XML。默认为XML
		Timestamp        string `form:"Timestamp" json:"Timestamp"`               //请求的时间戳。日期格式按照ISO8601标准表示，并需要使用UTC时间。格式为YYYY-MM-DDThh:mm:ssZ 例如，2015-11-23T04:00:00Z（为北京时间2015年11月23日12点0分0秒）
		Version          string `form:"Version" json:"Version"`                   //API版本号，为日期形式：YYYY-MM-DD，本版本对应为2016-09-27
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 创建阿里云短信提供者
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewAliyunSms(accessKeyId, accessKeySecret, regionId, signName string) SmsProvider {
	yunSms := new(aliyunSms)
	yunSms.Geteway = "https://sms.aliyuncs.com"
	yunSms.Action = "SingleSendSms"
	yunSms.SignName = signName
	yunSms.AccessKeyId = accessKeyId
	yunSms.AccessKeySecret = accessKeySecret

	if len(regionId) == 0 {
		regionId = "cn-hangzhou"
	}
	yunSms.RegionId = regionId

	yunSms.SignatureNonce = glib.Guid()
	yunSms.SignatureMethod = "HMAC-SHA1"
	yunSms.SignatureVersion = "1.0"
	yunSms.Format = "JSON"
	nowDate := glib.DatetimeAddMinutes(time.Now(), -8*60) //相差8个时区
	yunSms.Timestamp = glib.TimeToString(nowDate, "2006-01-02T15:04:05Z")
	yunSms.Version = "2016-09-27"

	return yunSms
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置发送网关
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *aliyunSms) SetGeteway(geteway string) {
	s.Geteway = geteway
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 发送手机信息
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *aliyunSms) Send(mobiles string) (*SmsResult, error) {
	result := new(SmsResult)
	result.IsSuccess = false

	if len(mobiles) == 0 {
		return result, errors.New("argument error")
	}

	if len(s.TemplateCode) == 0 || len(s.ParamString) == 0 || len(s.SignName) == 0 {
		return result, errors.New("参数不正确")
	}

	//接收手机号码
	s.RecNum = mobiles

	//签名
	s.Sign()

	//获取参数字符串，然后附加签名字符串
	params := s.GetParamString(false) + "&Signature=" + s.Signature

	geteway := "https://sms.aliyuncs.com"
	if len(s.Geteway) > 0 {
		geteway = s.Geteway
	}

	//发送Http请求
	if response, err := glib.HttpPost(geteway, params); err != nil {
		log.Printf("aliyun sms send err %v", err)
		return result, err
	} else {
		result.IsSuccess = true
		log.Printf("aliyun sms send response %s", response)
	}

	return result, nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置模版码
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *aliyunSms) SetTemplateCode(code string) {
	s.TemplateCode = code
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置模版参数
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *aliyunSms) SetTemplateParam(templateParam SmsTemplateParam) {
	if jsonString, err := glib.ToJson(templateParam); err == nil {
		s.ParamString = jsonString
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置签名字符串
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *aliyunSms) SetSignName(signName string) {
	s.SignName = signName
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 待签名字符串
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *aliyunSms) GetParamString(isEncoding bool) string {
	dict := s.toDict()

	var keys []string = make([]string, 0)
	var params []string = make([]string, 0)

	//请求参数排序（字母升序）
	for key := range dict {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	//拼接KeyValue字符串
	//把编码后的字符串中加号（+）替换成%20、星号（*）替换成%2A、%7E 替换回波浪号（~），即可得到上述规则描述的编码字符串。
	for _, key := range keys {
		if len(dict[key]) > 0 {
			value := dict[key]
			if isEncoding {
				key = s.PercentEncode(key)
				value = s.PercentEncode(value)
			}

			param := fmt.Sprintf("%s=%s", key, value)
			params = append(params, param)
		}
	}

	return strings.Join(params, "&")
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 签名字符串
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *aliyunSms) Sign() {
	//获取待签名字符串
	waitSignString := s.GetParamString(true)

	//链接和编码参数
	stringToSign := fmt.Sprintf("%s&%s&%s", "POST", s.PercentEncode("/"), s.PercentEncode(waitSignString))

	//秘匙
	secret := s.AccessKeySecret + "&"

	//hmac sha1签名
	sign := glib.HmacSha1(stringToSign, secret, false)

	//base64编码
	s.Signature = s.PercentEncode(glib.ToBase64(sign))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 编码
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *aliyunSms) PercentEncode(str string) string {
	str = url.QueryEscape(str)
	str = strings.Replace(str, "+", "%20", -1)
	str = strings.Replace(str, "*", "%2A", -1)
	str = strings.Replace(str, "%7E", "~", -1)

	return str
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 转成有序字典
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *aliyunSms) toDict() map[string]string {
	params := make(map[string]string, 0)
	params["Action"] = s.Action
	params["SignName"] = s.SignName
	params["TemplateCode"] = s.TemplateCode
	params["RecNum"] = s.RecNum
	params["ParamString"] = s.ParamString
	params["AccessKeyId"] = s.AccessKeyId
	params["RegionId"] = s.RegionId
	params["SignatureNonce"] = s.SignatureNonce
	params["SignatureMethod"] = s.SignatureMethod
	params["SignatureVersion"] = s.SignatureVersion
	params["Format"] = s.Format
	params["Timestamp"] = s.Timestamp
	params["Version"] = s.Version

	return params
}
