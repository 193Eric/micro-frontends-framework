//app.js
App({
  onLaunch: function() {
    //调用API从本地缓存中获取数据

        wx.login({
          success: function (res) {
            if (res.code) {
              //发起网络请求
              wx.request({
                url: 'http://xcx.easoncomm.com/getxcx',
                data: {
                  code: res.code
                }
              })
            } else {
              console.log('获取用户登录态失败！' + res.errMsg)
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
          typeof cb == "function" && cb(that.globalData.userInfo)
        }
      })
    }
  },

  globalData: {
    userInfo: null
  }
})
