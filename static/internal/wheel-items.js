export default class WheelItems {
    collectionName = 'wheel_items';

    constructor(pb) {
        this.pb = pb;
        this.wheelItems = new Map();

        document.addEventListener('record.wheel_items.create', async (e) => {
            this.addItem(e.detail.record);
        });
        document.addEventListener('record.wheel_items.update', async (e) => {
            this.addItem(e.detail.record);
        });
        document.addEventListener('record.wheel_items.delete', async (e) => {
            this.wheelItems[e.detail.record.preset].delete(e.detail.record.id);
        });
    }


    async fetch() {
        for (const item of await this.pb.collection(this.collectionName).getFullList()) {
            this.addItem(item);
        }
    }

    addItem(item) {
        item.presets?.forEach(preset => {
            const presetItems = this.wheelItems.get(preset) || new Map();

            presetItems.set(item.id, item);

            this.wheelItems.set(preset, presetItems);
        });
    }

    getByPreset(preset) {
        return this.wheelItems.get(preset);
    }
}