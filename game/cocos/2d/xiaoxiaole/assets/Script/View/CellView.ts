/**
* 智障的引擎设计，一群SB
*/
import { _decorator, Component, SpriteFrame, Animation,tween,v3,quat,Sprite, Tween} from 'cc';
const { ccclass, property } = _decorator;

import {CELL_STATUS, CELL_WIDTH, CELL_HEIGHT, ANITIME} from '../Model/ConstValue';
import CellModel from '../Model/CellModel';
@ccclass('CellView')
export class CellView extends Component {
    @property(SpriteFrame)
    public defaultFrame: SpriteFrame = null;

    isSelect:boolean;
    model:CellModel;
    status:string;

    onLoad () {
        this.isSelect = false;
    }

    initWithModel (model: CellModel) {
        this.model = model;
        var x = model.startX;
        var y = model.startY;
        this.node.setPosition(CELL_WIDTH * (x - 0.5),CELL_HEIGHT * (y - 0.5))
        var animation  = this.node.getComponent(Animation);
        if (model.status == CELL_STATUS.COMMON){
            animation.stop();
        }
        else{
            animation.play(model.status);
        }
    }

    updateView() {
        var cmd = this.model.cmd;
        if(cmd.length <= 0){
            return ;
        }

        let action = tween(this.node);
        var curTime = 0;
        for(var i in cmd){
            console.log(cmd[i]);
            if( cmd[i].playTime > curTime){
                action = action.delay(cmd[i].playTime - curTime);
            }
            if(cmd[i].action == "moveTo"){
                var x = (cmd[i].pos.x - 0.5) * CELL_WIDTH;
                var y = (cmd[i].pos.y - 0.5) * CELL_HEIGHT;
                action = action.to(ANITIME.TOUCH_MOVE,{position:v3(x,y)},{
                    onComplete:(target:Node)=>{
                        console.log("moveTo");
                    }
                });

            }
            else if(cmd[i].action == "toDie"){

                if(this.status == CELL_STATUS.BIRD){
                    let animation = this.node.getComponent(Animation);
                    animation.play("effect");
                    action = action.delay(ANITIME.BOMB_BIRD_DELAY);
                }
                action = action.removeSelf();

            }
            else if(cmd[i].action == "setVisible"){
                let isVisible = cmd[i].isVisible;
                console.log("setVisible");
                if(isVisible){
                    action = action.show();
                }else{
                    action = action.hide();
                }
            }
            else if(cmd[i].action == "toShake"){
                action = action.by(0.06,{angle:30}).by(0.12,{angle:-60}).by(0.06,{angle:30}).repeat(2);
            }
            curTime = cmd[i].playTime + cmd[i].keepTime;
        }
        action.start();
    }

    setSelect (flag: any) {
        var animation = this.node.getComponent(Animation);
        var bg = this.node.getChildByName("select");
        if(flag == false && this.isSelect && this.model.status == CELL_STATUS.COMMON){
            animation.stop();
            this.node.getComponent(Sprite).spriteFrame = this.defaultFrame;
        }
        else if(flag && this.model.status == CELL_STATUS.COMMON){
            animation.play(CELL_STATUS.CLICK);
        }
        else if(flag && this.model.status == CELL_STATUS.BIRD){
            animation.play(CELL_STATUS.CLICK);
        }
        bg.active = flag;
        this.isSelect = flag;
    }

}

