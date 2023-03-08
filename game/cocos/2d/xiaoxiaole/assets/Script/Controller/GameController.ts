import { _decorator, Component, Node, director,find, Vec2 } from 'cc';
const { ccclass, property } = _decorator;

import GameModel from "../Model/GameModel";
import {GridView}  from "../View/GridView";

@ccclass('GameController')
export class GameController extends Component {
    @property(Node)
    public grid: Node = null;

    gameModel:GameModel;

    onLoad () {
        console.log("加载GameController");
        this.gameModel = new GameModel();
        this.gameModel.init(6);
        var gridScript = this.grid.getComponent(GridView);
        gridScript.setController(this);
        gridScript.initWithCellModels(this.gameModel.getCells());
    }


    selectCell (pos: Vec2) {
         return this.gameModel.selectCell(pos);
    }

    cleanCmd () {
         this.gameModel.cleanCmd();
    }

}