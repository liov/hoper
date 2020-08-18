contentView = View()
contentView:bgColor(Color(235,122,119,1)):width(MeasurementType.MATCH_PARENT):height(MeasurementType.MATCH_PARENT)
window:addView(contentView)

imageView = ImageView():width(MeasurementType.MATCH_PARENT):setGravity(Gravity.CENTER_VERTICAL)
imageView:setImageUrl("https://ss1.bdstatic.com/70cFvXSh_Q1YnxGkpoWK1HF6hhy/it/u=1944539833,1909000872&fm=26&gp=0.jpg")
contentView:addView(imageView)