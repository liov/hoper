import { AudioClip, AudioSource, assert, warn, clamp01, resources,assetManager,Node,director } from "cc";
export class AudioManager {

    private static _instance: AudioManager;
    private  _audioSource?: AudioSource;
    private  _cachedAudioClipMap: Record<string, AudioClip> = {};

    constructor() {
        //@en create a node as audioMgr
        //@zh 创建一个节点作为 audioMgr
        let audioMgr = new Node();
        audioMgr.name = '__audioMgr__';

        //@en add to the scene.
        //@zh 添加节点到场景
        director.getScene().addChild(audioMgr);

        //@en make it as a persistent node, so it won't be destroied when scene change.
        //@zh 标记为常驻节点，这样场景切换的时候就不会被销毁了
        director.addPersistRootNode(audioMgr);

        //@en add AudioSource componrnt to play audios.
        //@zh 添加 AudioSource 组件，用于播放音频。
        this._audioSource = audioMgr.addComponent(AudioSource);
    }
    
    static get instance () {
        if (this._instance) {
            return this._instance;
        }

        this._instance = new AudioManager();
        return this._instance;
    }


    public get audioSource() {
        return this._audioSource;
    }

    // 设置音乐音量
    setMusicVolume (flag: number) {
        const audioSource = this._audioSource!;
        assert(audioSource, 'AudioManager not inited!');

        flag = clamp01(flag);
        audioSource.volume = flag;
    }

       /**
     * @en
     * play short audio, such as strikes,explosions
     * @zh
     * 播放短音频,比如 打击音效，爆炸音效等
     * @param sound clip or url for the audio
     * @param volume 
     */
    playOneShot(sound: AudioClip | string, volume: number = 1.0) {
        if (sound instanceof AudioClip) {
            this._audioSource.playOneShot(sound, volume);
        }
        else {
            let cachedAudioClip = this._cachedAudioClipMap[sound];
            if (cachedAudioClip) {
                this._audioSource.playOneShot(cachedAudioClip, 1);
            } else {    
                resources.load(`audio/sound/${sound}`, (err, clip: AudioClip) => {
                    if (err) {
                        console.log(err);
                    }
                    else {
                        this._cachedAudioClipMap[sound] = clip;
                        this._audioSource.playOneShot(clip, volume);
                    }
                });
            }
        }
    }

    /**
     * @en
     * play long audio, such as the bg music
     * @zh
     * 播放长音频，比如 背景音乐
     * @param sound clip or url for the sound
     * @param volume 
     */
    play(sound: AudioClip | string, volume: number = 1.0) {
        if (sound instanceof AudioClip) {
            this._audioSource.clip = sound;
            this._audioSource.play();
            this._audioSource.volume = volume;
        }
        else {
            let cachedAudioClip = this._cachedAudioClipMap[sound];
            if (cachedAudioClip) {
                this._audioSource.playOneShot(cachedAudioClip, 1);
            } else {
                resources.load(`audio/sound/${sound}`, (err, clip: AudioClip) => {
                    if (err) {
                        console.log(err);
                    }
                    else {
                        this._cachedAudioClipMap[sound] = clip;
                        this._audioSource.clip = clip;
                        this._audioSource.play();
                        this._audioSource.volume = volume;
                    }
                });
            }
        }
    }

    /**
     * stop the audio play
     */
    stop() {
        this._audioSource.stop();
    }

    /**
     * pause the audio play
     */
    pause() {
        this._audioSource.pause();
    }

    /**
     * resume the audio play
     */
    resume(){
        this._audioSource.play();
    }

}