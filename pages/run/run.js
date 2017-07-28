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
      color: "#66db38AA",
      width: 4,
    }],
    index:1,
    length: 0,
    _time:"00:00:00",
    km:"0.00",
    m_km:"--",
    ka:"0"
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

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  getLocation: function () {//获取当前地址
    var that = this;
    wx.getLocation({
      type: 'wgs84',
      success: function (res) {
        var site = {
          latitude: res.latitude,
          longitude: res.longitude
        }
        that.data.polyline[0].points.push(site);
        var polyline = [{
          points: that.data.polyline[0].points,
          color: "#66db38AA",
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
        var km = getDistance(first.latitude, first.longitude, last.latitude, last.longitude);
        var m_km = Math.floor(that.data.km/time/60)!=0&&time!=0?Math.floor(that.data.km/time/60):'--';
        if(that.data.option == 1){
        var ka = m_km!="--"?60*time/3600*(30/m_km*0.4)/1000:0;
        }else{
          var ka = m_km != "--"?m_km*60*60*1.05*time/3600/1000:0;
        }
        console.log(Math.floor(that.data.km / time / 60))
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
      this.getLocation();
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
  continue: function () {
    this.setData({
      index: 2
    })
    timeout = true;
    timeout = setTimeout(this.judgeTime, 1000)
  }
})