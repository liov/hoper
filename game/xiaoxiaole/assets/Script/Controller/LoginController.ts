import { _decorator, Component, ProgressBar, Button, AudioClip,loader,assetManager,director,AssetManager  } from 'cc';
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
      AudioManager.instance.play(this.worldSceneBGM);
    }

    onLogin () {
        this.last = 0; 
        this.loadingBar.node.active = true; 
        this.loginButton.node.active = false; 
        this.loadingBar.progress = 0; 
        this.loadingBar.barSprite.fillRange = 0; 
        director.preloadScene("Game",
        function  (count, amount, item) { 
          let progress = Number((count / amount).toFixed(2)); 
          if (progress > this.loadingBar.barSprite.fillRange) { 
            this.loadingBar.barSprite.fillRange = count / amount; 
          } 
        }.bind(this), function () { 
          this.loadingBar.node.active = false; 
          this.loginButton.node.active = false; 
          director.loadScene("Game"); 
        }.bind(this)); 
    }

    onDestroy () {
      AudioManager.instance.stop(); 
    }

}


