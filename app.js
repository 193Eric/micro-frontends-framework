//app.js
App({
  onLaunch: function() {
    //调用API从本地缓存中获取数据

        wx.login({
          success: function (res) { 
            if (res.code) {
              //发起网络请求
              wx.request({
                url: 'http://119.29.140.135/getxcx',
                data: {
                  code: res.code
                },
                success: function (res) {
                  var uid = res.data.ms.id;
                  wx.setStorage({
                    key: "uid",
                    data: uid
                  })
                }
              })
            } else {
              console.log('获取用户登录态失败！ ' + res.errMsg)
            }
          }
        }); 

  },

  getUserInfo: function(cb) {
    var that = this
    if (this.globalData.userInfo) {
      typeof cb == "function" && cb(this.globalData.userInfo)
    } else {
      //调用登录接口
      wx.getUserInfo({ 
        success: function(res) {
          that.globalData.userInfo = res.userInfo
          console.log(res)
          typeof cb == "function" && cb(that.globalData.userInfo)
          wx.request({
            url: 'http://119.29.140.135/setRunUser',
            method: 'POST',
            header: {
              "Content-Type": "application/x-www-form-urlencoded"
            },
            data: {
              "name": res.userInfo.nickName,
              "logo": res.userInfo.avatarUrl
            }
          })
        }
      })
    }
  },

  globalData: {
    userInfo: null
  }
})
