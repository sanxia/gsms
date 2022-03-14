# gsms
golang  aliyun, alidayu sms sdk

Aliyun, Alidayu Sms Api
==========================

--------------------------
Aliyun Sms Example:
--------------------------
```
import (
    "github.com/sanxia/gsms"
)

var smsProvider *gsms.SmsProvider

func init(){
    accessKeyId := "you aliyun accesskey id"
    accessKeySecret := "you aliyun accesskey secret"
    regionId := "cn-hangzhou"
    signName := "you sms sign name"
    smsProvider = gsms.NewAliyunSms(accessKeyId, accessKeySecret, regionId,     signName)
}

func AliyunSmsSend(mobiles string) (*gsms.SmsResult, error){
    smsProvider.SetTemplateCode = "sms_123456"
    smsProvider.SetTemplateParam = gsms.SmsTemplateParam{
        Code: "S-123",
        Product: "Test Validate Code 1",
    }
    return smsProvider.Send(mobiles)
}
```

--------------------------
Alidayun Sms Example:
--------------------------
```
import (
    "github.com/sanxia/gsms"
)

var smsProvider *gsms.SmsProvider

func init(){
    appKey := "you alidayu app key"
    appSecret := "you alidayu app secret"
    signName := "you sms sign name"
    smsProvider = gsms.NewAlidayunSms(appKey, appSecret, signName)
}

func AlidayuSmsSend(mobiles string) (*gsms.SmsResult, error){
    smsProvider.SetTemplateCode = "sms_123456"
    smsProvider.SetTemplateParam = gsms.SmsTemplateParam{
        Code: "S-123",
        Product: "Test Validate Code 2",
    }
    return smsProvider.Send(mobiles)
}
```

