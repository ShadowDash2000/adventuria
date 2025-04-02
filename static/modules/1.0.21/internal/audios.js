import Helper from "../helper.js";

export default class Audios {
    audioCollectionName = 'audio';
    audioPresetCollectionName = 'audio_presets';

    constructor(pb) {
        this.pb = pb;
        this.audios = new Map();
        this.audioPresets = new Map();

        document.addEventListener('record.audio.create', async (e) => {
            this.addAudio(e.detail.record);
        });
        document.addEventListener('record.audio.update', async (e) => {
            this.addAudio(e.detail.record);
        });

        document.addEventListener('record.audio_presets.create', async (e) => {
            this.addAudioPreset(e.detail.record);
        });
        document.addEventListener('record.audio_presets.update', async (e) => {
            this.addAudioPreset(e.detail.record);
        });
    }

    async fetch() {
        for (const audio of await this.pb.collection(this.audioCollectionName).getFullList()) {
            this.addAudio(audio);
        }
        for (const audio of await this.pb.collection(this.audioPresetCollectionName).getFullList()) {
            this.addAudioPreset(audio);
        }
    }

    addAudio(item) {
        this.audios.set(item.id, item);
    }

    addAudioPreset(item) {
        this.audioPresets.set(item.id, item);
    }

    getRandomAudioFromCellByEvent(cell, event) {
        const cellAudioPresetsIds = cell.audioPresets;

        for (const audioPresetId of cellAudioPresetsIds) {
            const audioPreset = this.audioPresets.get(audioPresetId);
            const audioPresetEvent = audioPreset?.event;

            if (audioPresetEvent === event) {
                const audiosIds = audioPreset.audios;
                const randomAudioId = audiosIds[Math.floor(Math.random() * audiosIds.length)];

                return this.audios.get(randomAudioId);
            }
        }
    }

    getAudiosSrcByEvent(event) {
        let audioPreset;
        for (const preset of this.audioPresets) {
            if (preset[1].event === event) {
                audioPreset = preset[1];
                break;
            }
        }

        if (!audioPreset) return;

        const audiosIds = audioPreset.audios;
        let src = [];
        for (const id of audiosIds) {
            src.push(Helper.getFile('audio', this.audios.get(id)));
        }

        return src;
    }
}