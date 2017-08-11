// run.js
var time = 0;
var timeout = true;
//计算距离
function getDistance(lat1, lng1, lat2, lng2) {
  var dis = 0;
  var radLat1 = toRadians(lat1);
  var radLat2 = toRadians(lat2);
  var deltaLat = radLat1 - radLat2;
  var deltaLng = toRadians(lng1) - toRadians(lng2);
  var dis = 2 * Math.asin(Math.sqrt(Math.pow(Math.sin(deltaLat / 2), 2) + Math.cos(radLat1) * Math.cos(radLat2) * Math.pow(Math.sin(deltaLng / 2), 2)));
  return (dis * 6378137).toFixed(2);
  function toRadians(d) { return d * Math.PI / 180; }
}
//获取计算的时间
function getTime(micro_second) {
  // 秒数
  var second = micro_second;
  // 小时位
  var hr = fill_zero_prefix( Math.floor(second / 3600));
  // 分钟位
  var min = fill_zero_prefix(Math.floor((second - hr * 3600) / 60));
  // 秒位
  var sec = fill_zero_prefix((second - hr * 3600 - min * 60));// equal to => var sec = second % 60;
  return hr + ":" + min + ":" + sec + " ";
}
function fill_zero_prefix(num) {
  return num < 10 ? "0" + num : num
}
Page({
  /**
   * 页面的初始数据
   */
  data: {
    polyline: [{
      points: [
      ],
      color: "#0000FFAA",
      width: 4,
    }],
    index:1,
    length: 0,
    _time:"00:00:00",
    km:"0.00",
    m_km:"--",
    ka:"0",
    animation:"",
    show:false,
    show_num:1
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    this.getLocation();
    this.setData({
      option:options.index
    })
  },
  onReady:function(){
    var that = this;
    this.show(1);
    setTimeout(function(){
      that.show(2);
    },1200)
    setTimeout(function () {
      that.setData({
        show: true
      })
    }, 2400)
  },
  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  show:function(num){

    var that = this;
    var animation = wx.createAnimation({
      duration: 600,
      transformOrigin: "50% 50%",
      timingFunction:'ease'
    })
    this.animation = animation
    animation.scale(0.01,0.01).step()
    this.setData({
      animation: animation.export()
    })
    setTimeout(function () {
      animation.scale(1, 1).step();
      that.setData({
        show_num: num+1,
        animation: animation.export()
      })
    }, 600)
  },
  getLocation: function () {//获取当前地址
    var that = this;
    wx.getLocation({
      type: 'gcj02',
      success: function (res) {
        var site = {
          latitude: res.latitude,
          longitude: res.longitude
        }
        that.data.polyline[0].points.push(site);
        var polyline = [{
          points: that.data.polyline[0].points,
          color: "#0000FFAA",
          width: 4,
        }]
        that.data.length++;
        that.setData({
          latitude: res.latitude,
          longitude: res.longitude,
          polyline: polyline,
          length: that.data.length
        })
        var last = that.data.polyline[0].points[that.data.length-1];
        var first = that.data.polyline[0].points[0];
        var km = getDistance(first.latitude, first.longitude, last.latitude, last.longitude)/1000;
        if (km < 0.0015) {
          km = 0.0;
        }
        km = km.toFixed(1);
        var m_km =that.data.km!=0&&time!=0?(that.data.km/(time/10/60)).toFixed(1):'--';
        if(that.data.option == 1){
        var ka = m_km!="--"?(60*time/10/3600*(30/m_km*0.4)/1000).toFixed(1):0;
        }else{
          var ka = m_km != "--"?(m_km*60*60*1.05*time/10/3600/1000).toFixed(1):0;
        }
        that.setData({
            km:km,
            m_km:m_km,
            ka: ka
        })
      }
    })
  },

  judgeTime:function(){
    if (timeout){
      time++;
      var num_time = getTime(time);
      this.setData({
        _time: num_time
      })
      if(time%10 == 0){
        this.getLocation();
      }
      setTimeout(this.judgeTime,1000);
    }else{
      return;
    }
  },
  begin:function(){
    this.setData({
      index:2
    })
    timeout = setTimeout(this.judgeTime, 1000)
  },
  stop: function () {
    this.setData({
      index: 3
    })
    timeout = false;
  },
  over:function(){
    var uid;
    var that =this;
    wx.getStorage({
      key: 'uid',
      success: function (res) {
        uid = res.data;
        //获取运动详情
        wx.request({
          url: 'http://119.29.140.135/setRunList',
          method: 'POST',
          header: {
            "Content-Type": "application/x-www-form-urlencoded"
          },
          data: {
            "uid": uid,
            "run_time":time,
            "end_time":(new Date().getTime()/1000).toFixed(0),
            "type": that.data.option,
            "km": that.data.km
          },
          success:function(res){
            wx.switchTab({ url: "../index/index"});
          }
        })
      }
    })
  },
  continue: function () {
    this.setData({
      index: 2
    })
    timeout = true;
    timeout = setTimeout(this.judgeTime, 1000)
  }
})