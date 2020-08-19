contentView = View()
contentView:bgColor(Color(235, 122, 119, 1)):width(MeasurementType.MATCH_PARENT):height(MeasurementType.MATCH_PARENT)
window:addView(contentView)

imageView = ImageView():width(MeasurementType.MATCH_PARENT):setGravity(Gravity.CENTER_VERTICAL)
imageView:setImageUrl("https://ss1.bdstatic.com/70cFvXSh_Q1YnxGkpoWK1HF6hhy/it/u=1944539833,1909000872&fm=26&gp=0.jpg")
contentView:addView(imageView)

baseLabelView = LinearLayout(LinearType.HORIZONTAL)
baseLabelView:bgColor(Color(95, 78, 79, 18)):width(MeasurementType.MATCH_PARENT):height(80):marginTop(10)
             :setGravity(Gravity.BOTTOM)
contentView:addView(baseLabelView)

openFlutter = Label()
openFlutter:bgColor(Color(172, 13, 56, 1)):width(100):height(50)
           :setGravity(Gravity.CENTER_VERTICAL)
           :text("open flutter"):fontSize(16):cornerRadius(8):textAlign(TextAlign.CENTER)
           :onClick(function()
    LuaBridge:openPage("flutterPage")
end)
baseLabelView:addView(openFlutter)

openNative = Label()
openNative:bgColor(Color(28, 80, 199, 18)):width(100):height(50)
          :setGravity(Gravity.CENTER_VERTICAL)
          :text("open native"):fontSize(16):cornerRadius(8):marginLeft(20):textAlign(TextAlign.CENTER)
          :onClick(function()
    LuaBridge:openPage("nativePage")
end)
baseLabelView:addView(openNative)

openLua = Label()
openLua:bgColor(Color(160, 103, 5, 18)):width(100):height(50)
          :setGravity(Gravity.CENTER_VERTICAL)
          :text("open lua"):fontSize(16):cornerRadius(8):marginLeft(20):textAlign(TextAlign.CENTER)
          :onClick(function()
    Navigator:gotoPage("https://hoper.xyz/static/lua/view.lua",Map(),AnimType.RightToLeft)
end)
baseLabelView:addView(openLua)