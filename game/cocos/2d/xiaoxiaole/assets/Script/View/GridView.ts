/**
* 根据cell的model返回对应的view
*/
import { _decorator, Component, Node, Prefab,Vec2,v2,v3,instantiate, tween, UITransform,EventTouch} from 'cc';
import { GameController } from '../Controller/GameController';
import {EffectLayer} from './EffectLayer';
import {CellView} from './CellView';
const { ccclass, property } = _decorator;

import {CELL_WIDTH, CELL_HEIGHT, GRID_PIXEL_WIDTH, GRID_PIXEL_HEIGHT, ANITIME} from '../Model/ConstValue';
import {AudioUtils} from "../Utils/AudioUtils";
import CellModel from '../Model/CellModel';
@ccclass('GridView')
export class GridView extends Component {
    @property(Prefab)
    public aniPre: Prefab[] = [];
    @property(Node)
    public effectLayer: Node = null;
    @property(AudioUtils)
    public audioUtils: AudioUtils = null;

    lastTouchPos:Vec2;
    isCanMove:boolean;
    isInPlayAni:boolean;
    cellViews:Node[][];
    controller:GameController;

    onLoad () {
        this.setListener();
        this.lastTouchPos = v2(-1, -1);
        this.isCanMove = true;
        this.isInPlayAni = false; // 是否在播放中
    }

    setController (controller: any) {
           this.controller = controller;
    }

    initWithCellModels (cellsModels: any) {
        this.cellViews = [];
        for(var i = 1;i<=9;i++){
            this.cellViews[i] = [];
            for(var j = 1;j<=9;j++){
                var type = cellsModels[i][j].type;
                var aniView = instantiate(this.aniPre[type]);
                aniView.parent = this.node;
                var cellViewScript = aniView.getComponent(CellView);
                cellViewScript.initWithModel(cellsModels[i][j]);
                this.cellViews[i][j] = aniView;
            }
        }
    }

    setListener() {
        this.node.on(Node.EventType.TOUCH_START, function(eventTouch:EventTouch){
            if(this.isInPlayAni){//播放动画中，不允许点击
                return true;
            }
            var touchPos = eventTouch.getUILocation();
            var cellPos = this.convertTouchPosToCell(touchPos);
            if(cellPos){
                console.log("2");
                var changeModels = this.selectCell(cellPos);
                this.isCanMove = changeModels.length < 3;
            } else{
                console.log("3");
                this.isCanMove = false;
            }
            return true;
        }, this);
            // 滑动操作逻辑
        this.node.on(Node.EventType.TOUCH_MOVE, function(eventTouch:EventTouch){
            if(this.isCanMove){
                var startTouchPos = eventTouch.getUIStartLocation();
                var startCellPos = this.convertTouchPosToCell(startTouchPos);
                var touchPos = eventTouch.getUILocation();
                var cellPos = this.convertTouchPosToCell(touchPos);
                if(startCellPos.x != cellPos.x || startCellPos.y != cellPos.y){
                    this.isCanMove = false;
                    var changeModels = this.selectCell(cellPos);
                }
            }
        }, this);
        this.node.on(Node.EventType.TOUCH_END, function(eventTouch:EventTouch){
        }, this);
        this.node.on(Node.EventType.TOUCH_CANCEL, function(eventTouch:EventTouch){
        }, this);
    }
    // 根据点击的像素位置，转换成网格中的位置
    convertTouchPosToCell (posv2: Vec2) {
        let uiTransform = this.node.getComponent(UITransform);
        const posv3 = uiTransform.convertToNodeSpaceAR(v3(posv2.x,posv2.y));
        if(posv3.x < 0 || posv3.x >= GRID_PIXEL_WIDTH || posv3.y < 0 || posv3.y >= GRID_PIXEL_HEIGHT){
            return false;
        }
        var x = Math.floor(posv3.x / CELL_WIDTH) + 1;
        var y = Math.floor(posv3.y / CELL_HEIGHT) + 1;
        return v2(x, y);
    }
    // 移动格子
    updateView(changeModels: CellModel[]) {
        let newCellViewInfo = [];
        for(var i in changeModels){
            var model = changeModels[i];
            var viewInfo = this.findViewByModel(model);
            var view = null;
                // 如果原来的cell不存在，则新建
            if(!viewInfo){
                var type = model.type;
                var aniView = instantiate(this.aniPre[type]);
                aniView.parent = this.node;
                var cellViewScript = aniView.getComponent(CellView);
                cellViewScript.initWithModel(model);
                view = aniView;
            } // 如果已经存在
            else{
                view = viewInfo.view;
                this.cellViews[viewInfo.y][viewInfo.x] = null;
            }
            var cellScript = view.getComponent(CellView);
            cellScript.updateView();// 执行移动动作
            if (!model.isDeath) {
                newCellViewInfo.push({
                    model: model,
                    view: view
                });
            }
        }
        // 重新标记this.cellviews的信息
        newCellViewInfo.forEach(function(ele){
            let model = ele.model;
            this.cellViews[model.y][model.x] = ele.view;
        },this);
    }

    // 显示选中的格子背景
    updateSelect(pos: Vec2) {
        for(var i = 1;i <=9 ;i++){
        for(var j = 1 ;j <=9 ;j ++){
            if(this.cellViews[i][j]){
                var cellScript = this.cellViews[i][j].getComponent(CellView);
                if(pos.x == j && pos.y ==i){
                    cellScript.setSelect(true);
                }
                else{
                    cellScript.setSelect(false);
                }
            }
        }
    }
    }

     /**
     * 根据cell的model返回对应的view
     */
    findViewByModel (model: CellModel) {
        for(var i = 1;i <=9 ;i++){
            for(var j = 1 ;j <=9 ;j ++){
                if(this.cellViews[i][j] && this.cellViews[i][j].getComponent(CellView).model == model){
                    return {view:this.cellViews[i][j],x:j, y:i};
                }
            }
        }
        return null;
    }

    getPlayAniTime (changeModels: any) {
        if(!changeModels){
            return 0;
        }
        var maxTime = 0;
        changeModels.forEach(function(ele){
            ele.cmd.forEach(function(cmd){
                if(maxTime < cmd.playTime + cmd.keepTime){
                    maxTime = cmd.playTime + cmd.keepTime;
                }
            },this)
        },this);
        return maxTime;
    }

    // 获得爆炸次数， 同一个时间算一个
    getStep (effectsQueue: any) {
        if(!effectsQueue){
            return 0;
        }
        return effectsQueue.reduce(function(maxValue, efffectCmd){
            return Math.max(maxValue, efffectCmd.step || 0);
        }, 0);
    }

    //一段时间内禁止操作
    disableTouch(time: number, step: number) {
        if(time <= 0){
            return ;
        }
        console.log("time:",time);
        this.isInPlayAni = true;
        tween(this.node).delay(time).call(function(){
            this.isInPlayAni = false;
            this.audioUtils.playContinuousMatch(step);
        }.bind(this)).start();
    }

    // 正常击中格子后的操作
    selectCell (cellPos: Vec2) {
        console.log(cellPos);
        var result = this.controller.selectCell(cellPos); // 直接先丢给model处理数据逻辑
        var changeModels = result[0]; // 有改变的cell，包含新生成的cell和生成马上摧毁的格子
        var effectsQueue = result[1]; //各种特效
        this.playEffect(effectsQueue);
        this.disableTouch(this.getPlayAniTime(changeModels), this.getStep(effectsQueue));
        this.updateView(changeModels);
        this.controller.cleanCmd();
        if(changeModels.length >= 2){
            this.updateSelect(v2(-1,-1));
            this.audioUtils.playSwap();
        }
        else{
            this.updateSelect(cellPos);
            this.audioUtils.playClick();
        }
        return changeModels;
    }

    playEffect (effectsQueue: any) {
         this.effectLayer.getComponent(EffectLayer).playEffects(effectsQueue);
    }

}
