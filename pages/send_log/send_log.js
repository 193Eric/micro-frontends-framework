// send_log.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    img:[],
    text:''
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
  
  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady: function () {

  },
  textChange: function (e){
    this.setData({
      text:e.detail.value
    })
  },
  send_log:function(){
    var that = this;
    wx.getStorage({
      key: 'uid',
      success: function (res) {
        var uid = res.data;
        //获取运动详情
        wx.request({
          url: 'http://119.29.140.135/setlog',
          method: 'POST',
          header: {
            "Content-Type": "application/x-www-form-urlencoded"
          },
          data: {
            "uid": uid,
            "text":that.data.text
          },
          success: function (res) {
            console.log(res.data.ms.insertId);
            for(var i =0,len=that.data.img.length;i<len;i++){
              that.upload(res.data.ms.insertId, that.data.img[i]);
            }     
            wx.switchTab({ url: "../logs/logs" });
          }
        })
      }
    })
  },
  send_img:function(id){
    var that = this;
    wx.chooseImage({
      count: 1, // 默认9
      sizeType: ['original', 'compressed'], // 可以指定是原图还是压缩图，默认二者都有
      sourceType: ['album', 'camera'], // 可以指定来源是相册还是相机，默认二者都有
      success: function (res) {
        // 返回选定照片的本地文件路径列表，tempFilePath可以作为img标签的src属性显示图片
        var tempFilePaths = res.tempFilePaths[0];
        that.data.img.push(tempFilePaths);
        that.setData({
          img: that.data.img
        })
      }
    })
  },
  upload:function(id,img){
    wx.uploadFile({
      url: 'http://119.29.140.135/setRunImg', //仅为示例，非真实的接口地址
      filePath: img,
      name: 'file',
      formData: {
        'id': id
      },
      success: function (res) {
        var data = res.data
        //do something
      }
    })
  },
  /**
   * 生命周期函数--监听页面显示
   */
  onShow: function () {
  
  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide: function () {
  
  },

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload: function () {
  
  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh: function () {
  
  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom: function () {
  
  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage: function () {
  
  }
})