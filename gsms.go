package gsms

/* ================================================================================
 * 短信接口
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	SmsProvider interface {
		Send(mobiles string) (*SmsResult, error)
		SetTemplateCode(code string)
		SetTemplateParam(templateParam SmsTemplateParam)
		SetTemplateString(templateString string)
		SetSignName(signName string)
		SetGeteway(geteway string)
	}

	SmsTemplateParam struct {
		Code string `form:"code" json:"code"`
	}

	SmsResult struct {
		Code      int32  `form:"err_code" json:"err_code"`
		Message   string `form:"msg" json:"msg"`
		Model     string `form:"model" json:"model"`
		RequestId string `form:"request_id" json:"request_id"`
		IsSuccess bool   `form:"is_success" json:"is_success"`
	}
)
