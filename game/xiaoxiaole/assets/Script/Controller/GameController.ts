import { _decorator, Component, Node, AudioSource,find } from 'cc';
const { ccclass, property } = _decorator;

import GameModel from "../Model/GameModel";
import {GridView}  from "../View/GridView";

@ccclass('GameController')
export class GameController extends Component {
    @property(Node)
    public grid: Node = null;

    gameModel:GameModel;

    onLoad () {
        this.gameModel = new GameModel(); 
        this.gameModel.init(4); 
        var gridScript = this.grid.getComponent(GridView); 
        gridScript.setController(this); 
        gridScript.initWithCellModels(this.gameModel.getCells()); 
    }


    selectCell (pos: any) {
         return this.gameModel.selectCell(pos); 
    }

    cleanCmd () {
         this.gameModel.cleanCmd(); 
    }

}