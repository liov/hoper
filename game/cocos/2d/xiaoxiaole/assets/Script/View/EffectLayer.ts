import { _decorator, Component, Prefab, instantiate, tween, Node, Animation, UITransform } from 'cc';
const { ccclass, property } = _decorator;

import { CELL_WIDTH } from '../Model/ConstValue';
import { AudioUtils } from "../Utils/AudioUtils";
@ccclass('EffectLayer')
export class EffectLayer extends Component {
    @property(Prefab)
    public bombWhite: Prefab = null;

    @property(Prefab)
    public crushEffect: Prefab = null;

    @property(AudioUtils)
    public audioUtils: AudioUtils = null;

    onLoad() {
    }

    playEffects(effectQueue: any[]) {
        if (!effectQueue || effectQueue.length <= 0) {
            return;
        }
        let soundMap = {}; //某一时刻，某一种声音是否播放过的标记，防止重复播放

        effectQueue.forEach(function (cmd) {
            tween(this.node).delay(cmd.playTime)
                .call(function () {
                    let instantEffect: Node = null;
                    let animationName = "";
                    if (cmd.action == "crush") {
                        instantEffect = instantiate(this.crushEffect);
                        animationName = "effect";
                        !soundMap["crush" + cmd.playTime] && this.audioUtils.playEliminate(cmd.step);
                        soundMap["crush" + cmd.playTime] = true;

                    }
                    else if (cmd.action == "rowBomb") {
                        instantEffect = instantiate(this.bombWhite);
                        animationName = "effect_line";

                    }
                    else if (cmd.action == "colBomb") {
                        instantEffect = instantiate(this.bombWhite);
                        animationName = "effect_col";
                    }

                    instantEffect.setPosition(CELL_WIDTH * (cmd.pos.x - 0.5), CELL_WIDTH * (cmd.pos.y - 0.5))
                    instantEffect.parent = this.node;
                    // 3.x挂载节点之后，动画的_nameState才会初始化，play和on才有效
                    const animation  = instantEffect.getComponent(Animation);
                    animation.play(animationName);
                    animation.on(Animation.EventType.FINISHED, function () {
                        instantEffect.destroy();
                    }, this);

                }.bind(this)).start()
        }, this);
    }

}

