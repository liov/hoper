/**
* 智障的引擎设计，一群SB
*/
import { _decorator, Component, SpriteFrame, Animation,tween,v3,quat,Sprite} from 'cc';
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
        var actionArray = [];
        var curTime = 0;
        for(var i in cmd){
            if( cmd[i].playTime > curTime){
                var delay = tween(this.node).delay(cmd[i].playTime - curTime);
                actionArray.push(delay);
            }
            if(cmd[i].action == "moveTo"){
                var x = (cmd[i].pos.x - 0.5) * CELL_WIDTH;
                var y = (cmd[i].pos.y - 0.5) * CELL_HEIGHT;
                var move = tween(this.node).to(ANITIME.TOUCH_MOVE,{position:v3(x,y)});
                actionArray.push(move);
            }
            else if(cmd[i].action == "toDie"){
                let action = tween(this.node).call(function(){ 
                    this.node.destroy();
                }.bind(this));

                if(this.status == CELL_STATUS.BIRD){
                    let animation = this.node.getComponent(Animation);
                    animation.play("effect");
                    action = action.delay(ANITIME.BOMB_BIRD_DELAY);
                }
                actionArray.push(action);
            }
            else if(cmd[i].action == "setVisible"){
                let isVisible = cmd[i].isVisible;
                actionArray.push(tween(this.node).call(function(){
                    if(isVisible){
                        this.node.opacity = 255;
                    }
                    else{
                        this.node.opacity = 0;
                    }
                }.bind(this)));
            }
            else if(cmd[i].action == "toShake"){
                let action = tween(this.node).by(ANITIME.DIE_SHAKE,{rotation:quat(0.06,30)}).by(ANITIME.DIE_SHAKE,{rotation:quat(0.12, -60)}).repeat(2);

                actionArray.push(action);
            }
            curTime = cmd[i].playTime + cmd[i].keepTime;
        }
        if(actionArray.length == 1){
            actionArray[0].start();
        }
        else{
            tween(this.node).sequence(...actionArray).start();
        }
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

