//index.js
//获取应用实例
var app = getApp()
Page({
  data: {
    motto: 'Hello World',
    userInfo: {},
    step:0,
    all_time:0,
    week_time:0,
    index:0,
  },
  gopage: function (event){
    const id = event.currentTarget.dataset.select;
    this.setData({
      index:id
    })
  },
  bindChange:function(event){
    const id = event.detail.current;
    this.setData({
      index: id
    })
  },
  onLoad: function () {
    const that = this;
    let uid;
    //调用应用实例的方法获取全局数据
    app.getUserInfo((userInfo)=>{
      //更新数据
      that.setData({
        userInfo:userInfo
      })
    })
    //获取uid
    wx.getStorage({
      key: 'uid',
      success: (res)=>{
        uid = res.data;
        //获取运动详情
        wx.request({
          url: 'http://119.29.140.135/getRunDetail',
          method: 'POST',
          header: {
            "Content-Type": "application/x-www-form-urlencoded"
          },
          data: {
            "uid": uid,
          },
          success:(res)=>{
            that.setData({
              all_time:(res.data.ms.total/60).toFixed(1),
              week_time:(res.data.ms.week/60).toFixed(1)
            })
          }
        })
      }
    })
    //获取微信今日步数
    wx.getWeRunData({
      success(res) {
        wx.request({
          url: 'http://119.29.140.135/getvalue',
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
  },
  goRun:function(event){
    var run = event.currentTarget.dataset.run;
    wx.navigateTo({ url: "../run/run?index="+run });
  },
  getLocation:function(){//获取当前地址
    var that = this;
    wx.getLocation({
      type: 'wgs84',
      success: function (res) {
        that.setData({
          latitude: res.latitude,
          longitude: res.longitude,
        })

      }
    })
  }
})
