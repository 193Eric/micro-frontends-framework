//index.js
//获取应用实例
var app = getApp()
Page({
  data: {
    motto: 'Hello World',
    userInfo: {},
    step:0
  },
  onLoad: function () {
    console.log('onLoad')
    var that = this
    //调用应用实例的方法获取全局数据
    app.getUserInfo(function(userInfo){
      //更新数据
      that.setData({
        userInfo:userInfo
      })
    })
    wx.getWeRunData({
      success(res) {
        wx.request({
          url: 'http://xcx.easoncomm.com/getvalue',
          method:'POST',
          header: {
            "Content-Type": "application/x-www-form-urlencoded"
          },
          data: {
            "encryptedData": res.encryptedData,
            "iv": res.iv
          },
          success: function (res) {
            that.setData({
              step: res.data.data.stepInfoList[30]['step']
            })
          }
        })
      }
    })
  }
})
