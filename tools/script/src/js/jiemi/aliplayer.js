const handler=setInterval(function(){
  console.clear();
  const _0xbfb1c1=new Date();
  debugger;
  const _0x3cabea=new Date();
  const _0x3cf64b=_0x3cabea.getTime()-_0xbfb1c1.getTime();
  if(_0x3cf64b>100){}
},1);
document.oncontextmenu=function(_0x33fe06){
  if(window.event){
    _0x33fe06=window.event;
  }
  try{
    var _0x4fd181=_0x33fe06.srcElement;
    if(!(_0x4fd181.tagName=='INPUT'&&_0x4fd181.type.toLowerCase()=='text'||_0x4fd181.tagName=='TEXTAREA')){
      return false;
    }
    return true;
  }catch(_0x3ed964){
    return false;
  }
};
document.onpaste=function(_0x5aa150){
  if(window.event){
    _0x5aa150=window.event;
  }
  try{
    var _0x13a96a=_0x5aa150.srcElement;
    if(!(_0x13a96a.tagName=='INPUT'&&_0x13a96a.type.toLowerCase()=='text'||_0x13a96a.tagName=='TEXTAREA')){
      return false;
    }
    return true;
  }catch(_0x5734a0){
    return false;
  }
};
document.oncut=function(_0x12fb23){
  if(window.event){
    _0x12fb23=window.event;
  }
  try{
    var _0x544cb1=_0x12fb23.srcElement;
    if(!(_0x544cb1.tagName=='INPUT'&&_0x544cb1.type.toLowerCase()=='text'||_0x544cb1.tagName=='TEXTAREA')){
      return false;
    }
    return true;
  }catch(_0x25bcc8){
    return false;
  }
};
function fuckyou(){
  window.open('/','_blank');
  window.close();
  window.location='about:blank';
}
var arr=[123,17,18];
document.oncontextmenu=new Function('event.returnValue=false;'),window.onkeydown=function(_0x32473c){
  var _0x203d1b=_0x32473c.keyCode||_0x32473c.which||_0x32473c.charCode;
  var _0x384a4c=_0x32473c.ctrlKey||_0x32473c.metaKey;
  console.log(_0x203d1b+'--'+_0x203d1b);
  if(_0x384a4c&&_0x203d1b==85){
    _0x32473c.preventDefault();
  }
  if(arr.indexOf(_0x203d1b)>-1){
    _0x32473c.preventDefault();
  }
};
function ck(){
  console.profile();
  console.profileEnd();
  if(console.clear){
    console.clear();
  };
  if(typeof console.profiles=='object'){
    return console.profiles.length>0;
  }
}
function hehe(){
  if(window.console&&(console.firebug||console.table&&/firebug/i.test(console.table()))||typeof opera=='object'&&typeof opera.postError=='function'&&console.profile.length>0){
    fuckyou();
  }
  if(typeof console.profiles=='object'&&console.profiles.length>0){
    fuckyou();
  }
}
hehe();
window.onresize=function(){
  if(window.outerHeight-window.innerHeight>100)fuckyou();
};
document.onkeydown=function(_0x377b5f){
  if(_0x377b5f.keyCode==112||_0x377b5f.keyCode==113||_0x377b5f.keyCode==114||_0x377b5f.keyCode==115||_0x377b5f.keyCode==117||_0x377b5f.keyCode==118||_0x377b5f.keyCode==119||_0x377b5f.keyCode==120||_0x377b5f.keyCode==121||_0x377b5f.keyCode==122||_0x377b5f.keyCode==123){
    return false;
  }
};
window.onhelp=function(){
  return false;
};

function bb(_0x32ec5f){
  var _0x404fa0=new RegExp('(^|&)'+_0x32ec5f+'=([^&]*)(&|$)');
  var _0x170084=window.location.search.substr(1).match(_0x404fa0);
  if(_0x170084!=null)return decodeURI(_0x170084[2]);
  return null;
}
var Params=bb('/');
var t=bb('//');
if(!Params){
  document.getElementsByTagName('img')[0].style.display='block';
}else{
  document.getElementsByTagName('title')[0].innerText=t;
  var player=new Aliplayer({'id':'playerCnt','source':'http://45.127.129.36/rms/'+Params+'?download=*.m3u8','width':'100%','height':'100%','autoplay':true,'preload':true,'rePlay':false,'playsinline':false,'useH5Prism':false,'isLive':false,'components':[{'name':'MemoryPlayComponent','type':AliPlayerComponent.MemoryPlayComponent,'args':[true]}]},function(_0x412aa2){
    console.log('The player is created');
  });
};