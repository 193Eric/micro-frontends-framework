package wechat

import (  
	"strings"
	"encoding/json"
	"crypto/sha1"
	"encoding/hex"
	"github.com/astaxie/beego"
	"io/ioutil"
	"github.com/nanjishidu/gomini/gocrypto"
	"strconv"
	"errors"
	"fmt"
	"net/http"
	"encoding/base64"
	"crypto/md5"
	"encoding/xml"
    "sort"
	"bytes"
	"github.com/astaxie/beego/orm"
	"time"
	"github.com/astaxie/beego/logs"
)

type WeixinTestController struct {
	beego.Controller
}
type WXPayNotifyResp struct {
    Return_code string `xml:"return_code"`
    Return_msg  string `xml:"return_msg"`
}

type RefundNotify struct {
	Return_code    string `xml:"return_code"`
    Return_msg     string `xml:"return_msg"`
    Appid          string `xml:"appid"`
    Mch_id         string `xml:"mch_id"`
    Nonce          string `xml:"nonce_str"`
	Req_info	   string `xml:"req_info"`
	Out_refund_no string `xml:"out_refund_no"`
	Out_trade_no string `xml:"out_trade_no"`
	Refund_fee string `xml:"refund_fee"`
	Refund_status string `xml:"refund_status"`
	Success_time string `xml:"success_time"`
	Transaction_id string `xml:"transaction_id"`
}
type GetOrderDetail struct{
	Appid string `xml:"appid"`
	Mch_id string `xml:"mch_id"`
	Out_trade_no string `xml:"out_trade_no"`
	Noncestr string `xml:"nonce_str"`
	Sign string `xml:"sign"`
	Sign_type string `xml:"sign_type"`
}
type WXPayNotifyReq struct {
    Return_code    string `xml:"return_code"`
    Return_msg     string `xml:"return_msg"`
    Appid          string `xml:"appid"`
    Mch_id         string `xml:"mch_id"`
    Nonce          string `xml:"nonce_str"`
    Sign           string `xml:"sign"`
    Result_code    string `xml:"result_code"`
    Openid         string `xml:"openid"`
    Is_subscribe   string `xml:"is_subscribe"`
    Trade_type     string `xml:"trade_type"`
    Bank_type      string `xml:"bank_type"`
    Total_fee      int    `xml:"total_fee"`
    Fee_type       string `xml:"fee_type"`
    Cash_fee       int    `xml:"cash_fee"`
    Cash_fee_Type  string `xml:"cash_fee_type"`
    Transaction_id string `xml:"transaction_id"`
    Out_trade_no   string `xml:"out_trade_no"`
    Attach         string `xml:"attach"`
    Time_end       string `xml:"time_end"`
}

type UnifyOrderReq struct {
    Appid            string `xml:"appid"`            //公众账号ID
    Body             string `xml:"body"`             //商品描述
    Mch_id           string `xml:"mch_id"`           //商户号
    Nonce_str        string `xml:"nonce_str"`        //随机字符串
    Notify_url       string `xml:"notify_url"`       //通知地址
    Trade_type       string `xml:"trade_type"`       //交易类型
    Spbill_create_ip string `xml:"spbill_create_ip"` //支付提交用户端ip
    Total_fee        int    `xml:"total_fee"`        //总金额
    Out_trade_no     string `xml:"out_trade_no"`     //商户订单号
    Sign             string `xml:"sign"`             //签名
    Openid           string `xml:"openid"`           //购买商品的用户wxid
}

type UnifyOrderResp struct {
    Return_code string `xml:"return_code"`
    Return_msg  string `xml:"return_msg"`
    Appid       string `xml:"appid"`
    Mch_id      string `xml:"mch_id"`
    Nonce_str   string `xml:"nonce_str"`
    Sign        string `xml:"sign"`
    Result_code string `xml:"result_code"`
    Prepay_id   string `xml:"prepay_id"`
    Trade_type  string `xml:"trade_type"`
}

type WXUser struct {
	Id         int    `orm:"column(id);pk;auto"`
	Name  string  `orm:"column(name)"`
	CreateTime int64  `orm:"column(create_time)"`
	Openid  string `orm:"column(open_id)"`
	City  string `orm:"column(city)"`
	Country  string `orm:"column(country)"`
	Province  string `orm:"column(province)"`
	HeadimgUrl  string `orm:"column(headimg_url)"`
}

type WXBody struct { 
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid string `json:"openid"`
	Scope string `json:"scope"`
} 

type WXInfo struct { 
	Openid string `json:"openid"`
	Nickname interface{} `json:"nickname"`
	City  interface{} `json:"city"`
	Country  interface{} `json:"country"`
	Province  interface{} `json:"province"`
	HeadimgUrl  interface{} `json:"headimgurl"`
} 

type Ticket struct { 
	Errcode int `json:"errcode"`
	Errmsg string `json:"errmsg"`
	Ticket  string `json:"ticket"`
} 

type SDK struct {
	Timestamp int64 `json:"timestamp"`
	Signature string `json:"signature"`
	Noncestr string `json:"noncestr"`
}

func  GetWeixinToken(code, appid, secret string) (*WXBody,error){

	//获取openid
	requestLine := strings.Join([]string{"https://api.weixin.qq.com/sns/oauth2/access_token",
        "?appid=", appid,
        "&secret=", secret,
        "&code=", code,
        "&grant_type=authorization_code"}, "")

	resp, err := http.Get(requestLine)	
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	logs.Info(err)
	if bytes.Contains(body, []byte("access_token")) {
		
	}else{
		return nil,errors.New("get msg fail")
	}

	atr := WXBody{}
	err = json.Unmarshal(body, &atr)
	if err !=nil {
		return nil,err
	}else{
		return &atr,nil
	}
	
}

func  GetWeixinInfo(AccessToken string , Openid string) (*WXUser,error){
	//获取用户信息
	userInfoUrl := strings.Join([]string{"https://api.weixin.qq.com/sns/userinfo",
        "?access_token=", AccessToken,
        "&openid=", Openid,
		"&lang=zh_CN"}, "")
	
	infoBody, err := http.Get(userInfoUrl)
	
	defer infoBody.Body.Close()
	
	info, _ := ioutil.ReadAll(infoBody.Body)
	updateUser := WXInfo{}
	json.Unmarshal(info, &updateUser)
	user := WXUser{Openid: updateUser.Openid}
	o := orm.NewOrm()
	err = o.Read(&user, "Openid")
	if err!=nil{
		return &user,err
	}else{
		user.City = updateUser.City.(string)
		user.Name = updateUser.Nickname.(string)
		user.Country = updateUser.Country.(string)
		user.Province = updateUser.Province.(string)
		user.HeadimgUrl = updateUser.HeadimgUrl.(string)
		if num, err := o.Update(&user); err == nil {
			logs.Info(num)
		}
		return &user,err
	}
}

func SendWeixinMs(accessToken string,data interface{}) (string,error){
	client := &http.Client{}
	content,err :=json.Marshal(data)
	postReq, err := http.NewRequest("POST",
		strings.Join([]string{`https://api.weixin.qq.com/cgi-bin/message/template/send`, "?access_token=", accessToken}, ""),
		bytes.NewReader(content))
	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(postReq)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
       return "",err
    }
	defer postReq.Body.Close()
	return string(body),nil
}

func GetWeixinAccessToken(appid string , appsecret string)(string,error){
	Url := strings.Join([]string{"https://api.weixin.qq.com/cgi-bin/token",
        "?grant_type=", "client_credential",
        "&appid=",appid,
		"&secret=",appsecret,}, "")
	infoBody, err := http.Get(Url)
	if err != nil{
		return "",err
	}
	defer infoBody.Body.Close()
	info, _ := ioutil.ReadAll(infoBody.Body)
	atr := WXBody{}
	err = json.Unmarshal(info, &atr)
	if err != nil{
		return "",err
	}else{
		return  atr.AccessToken,nil
	}
}

func GetWeixinTicket(appid string , appsecret string)(string,error){
	access_token,err := GetWeixinAccessToken(appid,appsecret)
	Url := strings.Join([]string{"https://api.weixin.qq.com/cgi-bin/ticket/getticket",
        "?access_token=", access_token,
        "&type=jsapi",}, "")
	infoBody, err := http.Get(Url)
	if err != nil{
		return "",err
	}
	defer infoBody.Body.Close()
	info, _ := ioutil.ReadAll(infoBody.Body)
	atr := Ticket{}
	err = json.Unmarshal(info, &atr)
	if err != nil{
		return atr.Errmsg,err
	}else{
		return  atr.Ticket,nil
	}
}

func InitSDK(ticket string,url string,noncestr string)(*SDK,error){
	timestamp := time.Now().Unix()
	sel := strings.Join([]string{"jsapi_ticket=",ticket,
        "&noncestr=", noncestr,
		"&timestamp=",strconv.FormatInt(timestamp,10),"&url=",url}, "")
	logs.Info(sel)
	signature := Sha1(sel)
	atr := SDK{
		Noncestr:noncestr,
		Timestamp:timestamp,
		Signature:signature,
	}

	return &atr,nil
}

func SetOrder(appid string,body string,mch_id string,nonce_str string,spbill_create_ip string,total_fee int,out_trade_no string,notify_url string,trade_type string,openid string,key string)( *UnifyOrderResp , error){
	sendData := UnifyOrderReq{
		Appid:appid,
		Body:body,
		Mch_id:mch_id,
		Nonce_str:nonce_str,
		Spbill_create_ip:spbill_create_ip,
		Total_fee:total_fee,
		Out_trade_no:out_trade_no,
		Notify_url:notify_url,
		Trade_type:trade_type,
		Openid:openid,
	}
	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = sendData.Appid
	m["body"] = sendData.Body
	m["mch_id"] = sendData.Mch_id
	m["notify_url"] = sendData.Notify_url
	m["trade_type"] = sendData.Trade_type
	m["spbill_create_ip"] = sendData.Spbill_create_ip
	m["total_fee"] = sendData.Total_fee
	m["out_trade_no"] = sendData.Out_trade_no
	m["nonce_str"] = sendData.Nonce_str
	m["openid"] = sendData.Openid
	sendData.Sign = wxpayCalcSign(m, key)
	bytes_req, err := xml.Marshal(sendData)     
	str_req := strings.Replace(string(bytes_req), "UnifyOrderReq", "xml", -1)
	//fmt.Println("转换为xml--------", str_req)
	bytes_req = []byte(str_req)

     //发送unified order请求.
     req, err := http.NewRequest("POST", "https://api.mch.weixin.qq.com/pay/unifiedorder", bytes.NewReader(bytes_req))
     if err != nil {
         fmt.Println("New Http Request发生错误，原因:", err)
         return nil,errors.New("Http Request发生错误")

	 }
	req.Header.Set("Accept", "application/xml")
     //这里的http header的设置是必须设置的.
     req.Header.Set("Content-Type", "application/xml;charset=utf-8")

     client := http.Client{}
     resp, _err := client.Do(req)
     if _err != nil {
         fmt.Println("请求微信支付统一下单接口发送错误, 原因:", _err)
         return nil,errors.New("请求微信支付统一下单接口发送错误")
	 }
	respBytes, err := ioutil.ReadAll(resp.Body)
     if err != nil {
         fmt.Println("解析返回body错误", err)
		 return nil,errors.New("解析返回body错误")
     }
     xmlResp := UnifyOrderResp{}
     _err = xml.Unmarshal(respBytes, &xmlResp)
     //处理return code.
     if xmlResp.Return_code == "FAIL" {
         fmt.Println("微信支付统一下单不成功，原因:", xmlResp.Return_msg, " str_req-->", str_req)
         return nil,errors.New("统一下单失败原因:"+xmlResp.Return_msg)
     }else{
		return &xmlResp,nil
	 }
	
}

//微信支付 下单签名
func wxpayCalcSign(mReq map[string]interface{}, key string) string {

    //fmt.Println("========STEP 1, 对key进行升序排序.========")
    //fmt.Println("微信支付签名计算, API KEY:", key)
    //STEP 1, 对key进行升序排序.
    sorted_keys := make([]string, 0)
    for k, _ := range mReq {
        sorted_keys = append(sorted_keys, k)
    }

    sort.Strings(sorted_keys)

    //fmt.Println("========STEP2, 对key=value的键值对用&连接起来，略过空值========")
    //STEP2, 对key=value的键值对用&连接起来，略过空值
    var signStrings string
    for _, k := range sorted_keys {
        //fmt.Printf("k=%v, v=%v\n", k, mReq[k])
        value := fmt.Sprintf("%v", mReq[k])
        if value != "" {
            signStrings = signStrings + k + "=" + value + "&"
        }
    }

    //fmt.Println("========STEP3, 在键值对的最后加上key=API_KEY========")
    //STEP3, 在键值对的最后加上key=API_KEY
    if key != "" {
        signStrings = signStrings + "key=" + key
    }

    //fmt.Println("========STEP4, 进行MD5签名并且将所有字符转为大写.========")
    //STEP4, 进行MD5签名并且将所有字符转为大写.
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(signStrings))
    cipherStr := md5Ctx.Sum(nil)
    upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))

    return upperSign
}
func Sha1(data string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(data))
	return hex.EncodeToString(sha1.Sum([]byte("")))
}

func GetWeixinSDK(appid string , appsecret string,url string,noncestr string)(*SDK,error){
	res,_ := GetWeixinTicket(appid,appsecret)
	sdk,err := InitSDK(res,url,noncestr)
	if err != nil{
		return nil,err
	}else{
		return sdk,nil
	}
}	

//微信退款回调
func WxRefundCallback(w http.ResponseWriter, r *http.Request, f func(string),key string)(string,string) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Println("读取http body失败，原因!", err)
        http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return "","FAIL"
    }
    defer r.Body.Close()

    fmt.Println("微信退款异步通知，HTTP Body:", "成功")
    var mr RefundNotify
	err = xml.Unmarshal(body, &mr)
	
	b, err := base64.StdEncoding.DecodeString(mr.Req_info)
	if err != nil {
		fmt.Println("解析HTTP Body格式到xml失败，原因!", err)
		return "","FAIL"
	}
	gocrypto.SetAesKey(strings.ToLower(gocrypto.Md5(key)))
	plaintext, err := gocrypto.AesECBDecrypt(b)
	if err != nil {
		fmt.Println(err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return "","FAIL"
	}
	var mnr RefundNotify
	err = xml.Unmarshal(plaintext, &mnr)
	if mnr.Refund_status == "SUCCESS"{
		f(mnr.Out_trade_no)
		return mnr.Out_trade_no,"SUCCESS"
	}else{
		return "","SUCCESS"
	}
}
//具体的微信支付回调函数的范例
func WxpayCallback(w http.ResponseWriter, r *http.Request, f func(string,string,string),key string)(string,string) {
    // body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Println("读取http body失败，原因!", err)
        http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return "","FAIL"
    }
    defer r.Body.Close()

    fmt.Println("微信支付异步通知，HTTP Body:", "成功")
    var mr WXPayNotifyReq
    err = xml.Unmarshal(body, &mr)
    if err != nil {
        fmt.Println("解析HTTP Body格式到xml失败，原因!", err)
        http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return "","FAIL"
    }

    var reqMap map[string]interface{}
    reqMap = make(map[string]interface{}, 0)

    reqMap["return_code"] = mr.Return_code
    reqMap["return_msg"] = mr.Return_msg
    reqMap["appid"] = mr.Appid
    reqMap["mch_id"] = mr.Mch_id
    reqMap["nonce_str"] = mr.Nonce
    reqMap["result_code"] = mr.Result_code
    reqMap["openid"] = mr.Openid
    reqMap["is_subscribe"] = mr.Is_subscribe
    reqMap["trade_type"] = mr.Trade_type
    reqMap["bank_type"] = mr.Bank_type
    reqMap["total_fee"] = mr.Total_fee
    reqMap["fee_type"] = mr.Fee_type
    reqMap["cash_fee"] = mr.Cash_fee
    reqMap["cash_fee_type"] = mr.Cash_fee_Type
    reqMap["transaction_id"] = mr.Transaction_id
    reqMap["out_trade_no"] = mr.Out_trade_no
    reqMap["attach"] = mr.Attach
    reqMap["time_end"] = mr.Time_end

    var resp WXPayNotifyResp
    //进行签名校验
    if wxpayVerifySign(reqMap, mr.Sign,key) {
		f(mr.Out_trade_no,mr.Transaction_id,mr.Result_code)
        //这里就可以更新我们的后台数据库了，其他业务逻辑同理。
        resp.Return_code = "SUCCESS"
        resp.Return_msg = "OK"
    } else {
        resp.Return_code = "FAIL"
        resp.Return_msg = "failed to verify sign, please retry!"
    }

    //结果返回，微信要求如果成功需要返回return_code "SUCCESS"
    bytes, _err := xml.Marshal(resp)
    strResp := strings.Replace(string(bytes), "WXPayNotifyResp", "xml", -1)
    if _err != nil {
        fmt.Println("xml编码失败，原因：", _err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return "","FAIL"
    }
    w.(http.ResponseWriter).WriteHeader(http.StatusOK)
	fmt.Fprint(w.(http.ResponseWriter), strResp)
	return mr.Out_trade_no,mr.Result_code
}

//微信支付签名验证函数
func wxpayVerifySign(needVerifyM map[string]interface{}, sign string,key string) bool {
    signCalc := wxpayCalcSign(needVerifyM , key)
    if sign == signCalc {
        fmt.Println("签名校验通过!")
        return true
    }

    fmt.Println("签名校验失败!")
    return false
}

//查询订单
func GetWeixinOrderInfo(appid string , mch_id string,out_trade_no string,noncestr string,_type string,key string)(*WXPayNotifyReq,error){
	sendData := GetOrderDetail{
		Appid:appid,
		Mch_id:mch_id,
		Out_trade_no:out_trade_no,
		Noncestr:noncestr,
		Sign_type:_type,
	}
	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = sendData.Appid
	m["mch_id"] = sendData.Mch_id
	m["nonce_str"] = sendData.Noncestr
	m["out_trade_no"] = sendData.Out_trade_no
	m["sign_type"] = sendData.Sign_type
	sendData.Sign = wxpayCalcSign(m, key)
	bytes_req, err := xml.Marshal(sendData)     
	str_req := strings.Replace(string(bytes_req), "UnifyOrderReq", "xml", -1)
	//fmt.Println("转换为xml--------", str_req)
	bytes_req = []byte(str_req)

     //发送unified order请求.
     req, err := http.NewRequest("POST", "https://api.mch.weixin.qq.com/pay/orderquery", bytes.NewReader(bytes_req))
     if err != nil {
         fmt.Println("New Http Request发生错误，原因:", err)
         return nil,errors.New("Http Request发生错误")

	 }
	 req.Header.Set("Accept", "application/xml")
     //这里的http header的设置是必须设置的.
     req.Header.Set("Content-Type", "application/xml;charset=utf-8")

     client := http.Client{}
     resp, _err := client.Do(req)
     if _err != nil {
         fmt.Println("请求微信支付统一下单接口发送错误, 原因:", _err)
         return nil,errors.New("请求微信支付统一下单接口发送错误")
	 }
	respBytes, err := ioutil.ReadAll(resp.Body)
     if err != nil {
         fmt.Println("解析返回body错误", err)
		 return nil,errors.New("解析返回body错误")
     }
     xmlResp := WXPayNotifyReq{}
     _err = xml.Unmarshal(respBytes, &xmlResp)
     //处理return code.
     if xmlResp.Return_code == "FAIL" {
         fmt.Println("微信支付查询订单不成功，原因:", xmlResp.Return_msg, " str_req-->", str_req)
         return nil,errors.New("不成功失败原因:"+xmlResp.Return_msg)
     }else{
		return &xmlResp,nil
	 }
}	