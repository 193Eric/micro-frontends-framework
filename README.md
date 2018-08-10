# go-wechat
golang 微信公众号接口库，包括微信支付统一下单、获取access_token、获取用户信息、生成前端js-sdk授权信息、微信支付回调、微信退款回调、订单查询等等.....

### golang对接公众号和微信商户平台开发

> 首先我们要从公众号和商户平台获取一些参数  

- appid // 大家都懂
- appsecret  // 大家都懂
- noncestr //随机字符串，好像是小于32位就行了，做sign签名用的
- key //微信商户key，商户后台生成的，用来做sign签名的
- mch_id //微信商户id


> 介绍一下接口微信授权部分func  

- GetWeixinToken(code string, appid string, secret string)(*WXBody,error)  
  code参数是前端通过微信授权回调的code码,这个接口成功主要会返回access_token 和openid
  
- GetWeixinInfo(AccessToken string , Openid string)(*WXUser,error)  
  这个接口是接上一个接口调用，可以获取到微信的用户信息（头像，城市等等）
  
- GetWeixinAccessToken(appid string , appsecret string)(string,error)  
  返回access_token,因为access_token有2小时的限制，如果过期了又不想重新去取用户code，可以用这个接口直接获取access_token
  

> 接口微信JSSDK签名的func  

- GetWeixinSDK(appid string , appsecret string,url string,noncestr string)(*SDK,error)  
  url参数是当前要调用js-sdk的页面url（注意如果是单页面前端项目的话，ios和安卓要做兼容处理哦）,这个接口会返回前端调用js-sdk wx.config()需要的全部参数
  
> 接口微信下单支付部分的func  

- SetOrder(appid string,body string,mch_id string,nonce_str string,spbill_create_ip string,total_fee int,out_trade_no string,notify_url string,trade_type string,openid string,key string)( *UnifyOrderResp , error)  
  微信统一下单接口，参数一堆，在微信文档里面参数都又介绍，这里只介绍主要参数。total_fee（注意单位是分）,notify_url 是微信支付回调地址（下面会提到支付回调接口，这里主要填那个接口的地址）, key 是上面提到的商户平台key  
  这里返回的主要是看Prepay_id 这个是前端唤起支付的主要参数。  
  
- WxpayCallback(w http.ResponseWriter, r *http.Request, f func(string,string,string),key string)(string,string)  
  微信支付回调接口，f参数是支付成功的回调函数，可以自己定义，支付成功后需要做什么事情。key还是商户平台的key。  
  成功会返回(out_trade_no，""),失败会返回("","FAIL")  --- out_trade_no是微信订单号,我们统一下单传过去的。  
  
- WxRefundCallback(w http.ResponseWriter, r *http.Request, f func(string),key string)(string,string)  
  微信退款回调接口，微信退款成功后，访问的接口，***这里注意的是有两种退款，一种是微信后台退款，需要在后台设置退款通知url，设置成这个接口。一种是接口退款，需要传参数notice_url为这个接口地址***  
  参数和返回和支付回调一样。  

- GetWeixinOrderInfo(appid string , mch_id string,out_trade_no string,noncestr string,sign_type string,key string)(*WXPayNotifyReq,error)  
  查询微信订单接口，为什么要有这个接口呢，因为微信支付回调可能不准，微信文档自己也说了，不能保证完成通知到。所以我要通过定时任务。通过这个接口去拿支付信息  
  out_trade_no（统一下单传输的订单号），_type(签名加密方式，跟统一下单填一致)，这个接口会返回微信订单信息。  
  
> 接口微信公众号模版推送  

- SendWeixinMs(accessToken string,data interface{}) (string,error)  
  公众号模版推送接口，data是推送的内容（具体看微信文档），成功会返回string格式的推送信息提示  

  
