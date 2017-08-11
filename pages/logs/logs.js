
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
    this.setData({

    })
  },
  send_log(){
    wx.navigateTo({ url: "../send_log/send_log"})
  }
})
