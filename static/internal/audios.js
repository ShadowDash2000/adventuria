export default class Audios {
    collectionName = 'audio';

    constructor(pb) {
        this.pb = pb;
        this.audios = new Map();

        document.addEventListener('record.audio.create', async (e) => {
            this.addItem(e.detail.record);
        });
        document.addEventListener('record.audio.update', async (e) => {
            this.addItem(e.detail.record);
        });
    }

    async fetch() {
        for (const audio of await this.pb.collection(this.collectionName).getFullList()) {
            this.addItem(audio);
        }
    }

    addItem(item) {
        const audioType = this.audios.get(item.event) || new Map();
        audioType.set(item.id, item);
        this.audios.set(item.event, audioType);
    }

    getRandomAudio(type) {
        const audioItemsKeys = Array.from(this.audios.get(type)?.keys());
        const randomKey = audioItemsKeys[Math.floor(Math.random() * audioItemsKeys.length)];
        return this.audios.get(type)?.get(randomKey);
    }
}