--1
ui_views = {}
require("packet/BindMeta")
local window_1 = VStack()
ui_views.window_1 = window_1
window_1:widthPercent(100)
window_1:bgColor(Color(235, 122, 119, 1))
window_1:width(MeasurementType.MATCH_PARENT)
window_1:height(MeasurementType.MATCH_PARENT)
local window_1_1 = ImageView()
ui_views.window_1_1 = window_1_1
window_1_1:width(MeasurementType.MATCH_PARENT)
window_1_1:setGravity(Gravity.CENTER_VERTICAL)
window_1_1:setImageUrl("https://ss1.bdstatic.com/70cFvXSh_Q1YnxGkpoWK1HF6hhy/it/u=1944539833,1909000872&fm=26&gp=0.jpg")
local window_1_2 = HStack()
ui_views.window_1_2 = window_1_2
window_1_2:bgColor(Color(95, 78, 79, 18))
window_1_2:width(MeasurementType.MATCH_PARENT)
window_1_2:height(80)
window_1_2:marginTop(10)
local kvar1="open flutter"
local kvar2=Color(172, 13, 56, 1)
local kvar3=function()
 LuaBridge:openPage("flutterPage")
end
local window_1_2_1 = Label()
ui_views.window_1_2_1 = window_1_2_1
window_1_2_1:bgColor(kvar2)
window_1_2_1:width(100)
window_1_2_1:height(50)
window_1_2_1:setGravity(Gravity.CENTER_VERTICAL)
window_1_2_1:text(kvar1)
window_1_2_1:fontSize(16)
window_1_2_1:cornerRadius(8)
window_1_2_1:marginLeft(20)
window_1_2_1:textAlign(TextAlign.CENTER)
window_1_2_1:onClick(kvar3)
local kvar4="open native"
local kvar5=Color(28, 80, 199, 18)
local kvar6=function()
 LuaBridge:openPage("nativePage")
end
local window_1_2_2 = Label()
ui_views.window_1_2_2 = window_1_2_2
window_1_2_2:bgColor(kvar5)
window_1_2_2:width(100)
window_1_2_2:height(50)
window_1_2_2:setGravity(Gravity.CENTER_VERTICAL)
window_1_2_2:text(kvar4)
window_1_2_2:fontSize(16)
window_1_2_2:cornerRadius(8)
window_1_2_2:marginLeft(20)
window_1_2_2:textAlign(TextAlign.CENTER)
window_1_2_2:onClick(kvar6)
local kvar7="open lua"
local kvar8=Color(160, 103, 5, 18)
local kvar9=function()
 Link:link("HotUpdate", Map(), AnimType.RightToLeft)
end
local window_1_2_3 = Label()
ui_views.window_1_2_3 = window_1_2_3
window_1_2_3:bgColor(kvar8)
window_1_2_3:width(100)
window_1_2_3:height(50)
window_1_2_3:setGravity(Gravity.CENTER_VERTICAL)
window_1_2_3:text(kvar7)
window_1_2_3:fontSize(16)
window_1_2_3:cornerRadius(8)
window_1_2_3:marginLeft(20)
window_1_2_3:textAlign(TextAlign.CENTER)
window_1_2_3:onClick(kvar9)
window_1_2:children({window_1_2_1, window_1_2_2, window_1_2_3})
window_1:children({window_1_1, window_1_2})
window:addView(window_1)
