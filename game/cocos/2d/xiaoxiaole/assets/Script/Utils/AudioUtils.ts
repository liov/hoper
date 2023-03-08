import { _decorator, Component, AudioClip,Node,AudioSource,assert,find,director } from 'cc';
import { AudioManager } from './AudioManger';

import {Toast} from '../Utils/Toast';

const { ccclass, property } = _decorator;

@ccclass('AudioUtils')
export class AudioUtils extends Component {
        
    @property(Node)
    public audioButton: Node = null;
    @property(AudioClip)
    public bgm: AudioClip = null;

    @property(AudioClip)
    public swap: AudioClip = null;
    @property(AudioClip)
    public click: AudioClip = null;
    @property(AudioClip)
    public eliminate: AudioClip[] = [];
    @property(AudioClip)
    public continuousMatch: AudioClip[] = [];

    onLoad () {
      
    }

    callback () {
        let state = AudioManager.instance.audioSource.state; 
        state === 1 ? AudioManager.instance.pause() : AudioManager.instance.play(this.bgm); 
        Toast(state === 1 ? '关闭背景音乐🎵' : '打开背景音乐🎵' ) 
    }

    onEnable () {
        AudioManager.instance.play(this.bgm);
    }
    
    start () {
        let audioButton = this.node.parent.getChildByName('audioButton');
        audioButton.on('click', this.callback, this);
    }

    playClick () {
        AudioManager.instance.playOneShot(this.click, 1); 
    }

    playSwap () {
        AudioManager.instance.playOneShot(this.swap, 1); 
    }

    playEliminate (step: any) {
        step = Math.min(this.eliminate.length - 1, step); 
        AudioManager.instance.playOneShot(this.eliminate[step], 1); 
    }

    playContinuousMatch (step: any) {
        console.log("step = ", step); 
        step = Math.min(step, 11); 
        if(step < 2){ 
                return  
        } 
        AudioManager.instance.playOneShot(this.continuousMatch[Math.floor(step/2) - 1], 1); 
    }

    playAudio () {
        AudioManager.instance.play(this.bgm);
    }
}


