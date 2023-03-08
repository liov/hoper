import {color,director,Canvas,Node,Label,Graphics,v2,v3,UITransform, tween,dragonBones } from 'cc';
//一个简单的tost组件，用法：
// let Toast = reqire('Toast.js')
// Toast(text,{gravity,duration,bg_color})
//text:要显示的字符串
//gravity(可选):位置，String类型，可选值('CENTER','TOP','BOTTOM'),默认为'CENTER'
//duration(可选):时间，Number类型，单位为秒，默认1s
//bg_color(可选):颜色，cc.color类型，默认cc.color(102, 102, 102, 200)
export function Toast(
  text = "",
  {
    gravity = "CENTER",
    duration = 1,
    bg_color = color(102, 102, 102, 200)
  } = {}
) {
  // canvas
  let canvas = director.getScene().getComponentInChildren(Canvas);
  const uiTrans = canvas.node.getComponent(UITransform)!;
  let width = uiTrans.width;
  let height = uiTrans.height;

  let bgNode = new Node();

  // Lable文本格式设置
  let textNode = new Node();
  let textLabel = textNode.addComponent(Label);
  textLabel.horizontalAlign = Label.HorizontalAlign.CENTER;
  textLabel.verticalAlign = Label.VerticalAlign.CENTER;
  textLabel.fontSize = 30;
  textLabel.string = text;
  const textuiTrans = textNode.getComponent(UITransform)!;
  // 当文本宽度过长时，设置为自动换行格式
  if (text.length * textLabel.fontSize > (width * 3) / 5) {
    textuiTrans.width = (width * 3) / 5;
    textLabel.overflow = Label.Overflow.RESIZE_HEIGHT;
  } else {
    textuiTrans.width = text.length * textLabel.fontSize;
  }
  let lineCount =
    ~~((text.length * textLabel.fontSize) / ((width * 3) / 5)) + 1;
    textuiTrans.height = textLabel.fontSize * lineCount;

  // 背景设置
  let ctx = bgNode.addComponent(Graphics);
  ctx.arc(
    -textuiTrans.width / 2,
    0,
    textuiTrans.height / 2 + 20,
    0.5 * Math.PI,
    1.5 * Math.PI,
    true
  );
  ctx.lineTo(textuiTrans.width / 2, -(textuiTrans.height / 2 + 20));
  ctx.arc(
    textuiTrans.width / 2,
    0,
    textuiTrans.height / 2 + 20,
    1.5 * Math.PI,
    0.5 * Math.PI,
    true
  );
  ctx.lineTo(-textuiTrans.width / 2, textuiTrans.height / 2 + 20);
  ctx.fillColor = bg_color;
  ctx.fill();

  bgNode.addChild(textNode);

  // gravity 设置Toast显示的位置
  if (gravity === "CENTER") {
    bgNode.position = v3(0,0,1000);
  } else if (gravity === "TOP") {
    bgNode.position = v3(bgNode.position.x,bgNode.position.y + (height / 5) * 2,1000);
  } else if (gravity === "BOTTOM") {
    bgNode.position = v3(bgNode.position.x,bgNode.position.y - (height / 5) * 2,1000);
  }

  canvas.node.addChild(bgNode);


  tween(bgNode).by(duration,{position: v3(0,0)},{easing:"fade"}).call(function() {
    bgNode.destroy();
  }).start();
}
