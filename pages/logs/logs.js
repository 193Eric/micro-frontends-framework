
Page({
  data: {
    logs: [
      {
        img:"http://img02.tooopen.com/images/20150928/tooopen_sy_143912755726.jpg",
        text:'好开心啊！！！！',
        icon:'http://img06.tooopen.com/images/20160818/tooopen_sy_175866434296.jpg',
        name:"193ric"
      },
      {
        img: "http://img02.tooopen.com/images/20150928/tooopen_sy_143912755726.jpg",
        text: '好开心啊！！！！',
        icon: 'http://img06.tooopen.com/images/20160818/tooopen_sy_175866434296.jpg',
        name: "193ric"
      },
      {
        img: "http://img02.tooopen.com/images/20150928/tooopen_sy_143912755726.jpg",
        text: '好开心啊！！！！',
        icon: 'http://img06.tooopen.com/images/20160818/tooopen_sy_175866434296.jpg',
        name: "193ric"
      }
    ],
    imgUrls: [
      'http://img02.tooopen.com/images/20150928/tooopen_sy_143912755726.jpg',
      'http://img06.tooopen.com/images/20160818/tooopen_sy_175866434296.jpg',
      'http://img06.tooopen.com/images/20160818/tooopen_sy_175833047715.jpg'
    ],
    indicatorDots: false,
    autoplay: true,
    interval: 5000,
    duration: 1000
  },
  onLoad: function () {
    const that = this;
    wx.request({
      url: 'http://119.29.140.135/getLogs',
      method: 'POST',
      header: {
        "Content-Type": "application/x-www-form-urlencoded"
      },
      data: {
        page:1
      },
      success: (res)=>{
        console.log(res.data.ms.logs)
        that.setData({
          imgUrls:res.data.ms.lunbo,
          logs: res.data.ms.logs,
        })
      }
    })

  },
  send_log(){
    wx.navigateTo({ url: "../send_log/send_log"})
  }
})
