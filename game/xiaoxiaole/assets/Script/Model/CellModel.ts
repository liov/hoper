/**
 * 注意：已把原脚本注释，由于脚本变动过大，转换的时候可能有遗落，需要自行手动转换
 */
import { CELL_TYPE, ANITIME, CELL_STATUS, GRID_HEIGHT } from "./ConstValue";
import {v2} from "cc";

export default class CellModel {
    type: number;
    status: string;
    x: number;
    y: number;
    startX: number;
    startY: number;
    cmd: any[];
    isDeath:boolean;
    objecCount:number;
  constructor() {
    this.type = null;
    this.status = CELL_STATUS.COMMON;
    this.x = 1;
    this.y = 1;
    this.startX = 1;
    this.startY = 1;
    this.cmd = [];
    this.isDeath = false;
    this.objecCount = Math.floor(Math.random() * 1000);
  }

  init(type) {
    this.type = type;
  }

  isEmpty() {
    return this.type == CELL_TYPE.EMPTY;
  }

  setEmpty() {
    this.type = CELL_TYPE.EMPTY;
  }
  setXY(x, y) {
    this.x = x;
    this.y = y;
  }

  setStartXY(x, y) {
    this.startX = x;
    this.startY = y;
  }

  setStatus(status) {
    this.status = status;
  }

  moveToAndBack(pos) {
    var srcPos = v2(this.x, this.y);
    this.cmd.push({
      action: "moveTo",
      keepTime: ANITIME.TOUCH_MOVE,
      playTime: 0,
      pos: pos
    });
    this.cmd.push({
      action: "moveTo",
      keepTime: ANITIME.TOUCH_MOVE,
      playTime: ANITIME.TOUCH_MOVE,
      pos: srcPos
    });
  }

  moveTo(pos, playTime) {
    var srcPos = v2(this.x, this.y);
    this.cmd.push({
      action: "moveTo",
      keepTime: ANITIME.TOUCH_MOVE,
      playTime: playTime,
      pos: pos
    });
    this.x = pos.x;
    this.y = pos.y;
  }

  toDie(playTime) {
    this.cmd.push({
      action: "toDie",
      playTime: playTime,
      keepTime: ANITIME.DIE
    });
    this.isDeath = true;
  }

  toShake(playTime) {
    this.cmd.push({
      action: "toShake",
      playTime: playTime,
      keepTime: ANITIME.DIE_SHAKE
    });
  }

  setVisible(playTime, isVisible) {
    this.cmd.push({
      action: "setVisible",
      playTime: playTime,
      keepTime: 0,
      isVisible: isVisible
    });
  }

  moveToAndDie(pos) {

  }

  isBird() {
    return this.type == CELL_TYPE.BIRD;
  }

}
