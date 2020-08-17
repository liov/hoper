contentView = View()
contentView:bgColor(Color(34,66,121,1)):width(MeasurementType.MATCH_PARENT):height(88):marginTop(100)
window:addView(contentView)

imageView = ImageView():width(70):height(70):setGravity(Gravity.CENTER_VERTICAL):cornerRadius(35):marginLeft(10)
imageView:setImageUrl("http://www.qq22.com.cn/uploads/allimg/201609071046/ldxdoeundzf118556.jpg")
contentView:addView(imageView)