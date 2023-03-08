import { _decorator, Component, ProgressBar, Button, AudioClip,resources,Prefab,director,AssetManager  } from 'cc';
const { ccclass, property } = _decorator;
const { Task } = AssetManager;

import {AudioManager} from "../Utils/AudioManger";
@ccclass('LoginController')
export class LoginController extends Component {
    @property(ProgressBar)
    public loadingBar: ProgressBar = null;
    @property(Button)
    public loginButton: Button = null;
    @property(AudioClip)
    public worldSceneBGM: AudioClip = null;

    last:number;

    onLoad () {
      console.log("加载LoginController");
      AudioManager.instance.play(this.worldSceneBGM);
    }

    onLogin () {
        this.last = 0;
        this.loadingBar.node.active = true;
        this.loginButton.node.active = false;
        this.loadingBar.progress = 0;
        this.loadingBar.barSprite.fillRange = 0;
        console.time("load");
        director.preloadScene("Game",
         (count, amount, item)=> {
          let progress = Number((count / amount).toFixed(2));
          if (progress > this.loadingBar.barSprite.fillRange) {
            this.loadingBar.barSprite.fillRange = count / amount;
          }
        }, ()=> {
          this.loadingBar.node.active = false;
          this.loginButton.node.active = false;
          director.loadScene("Game");
        });
    }

    onDestroy () {
      console.log("销毁");
      console.timeEnd("load");
      AudioManager.instance.stop();
    }

}


