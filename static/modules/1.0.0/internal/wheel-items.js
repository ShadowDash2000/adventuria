export default class WheelItems {
    collectionName = 'wheel_items';

    constructor(pb) {
        this.pb = pb;
        this.wheelItems = new Map();
        this.wheelItemsList = new Map();

        document.addEventListener('record.wheel_items.create', async (e) => {
            this.addItem(e.detail.record);
            this.addToList(e.detail.record);
        });
        document.addEventListener('record.wheel_items.update', async (e) => {
            this.addItem(e.detail.record);
            this.addToList(e.detail.record);
        });
        document.addEventListener('record.wheel_items.delete', async (e) => {
            this.wheelItems[e.detail.record.preset].delete(e.detail.record.id);
            this.wheelItemsList.delete(e.detail.record.id);
        });
    }


    async fetch() {
        for (const item of await this.pb.collection(this.collectionName).getFullList()) {
            this.addItem(item);
            this.addToList(item);
        }
    }

    addItem(item) {
        item.presets?.forEach(preset => {
            const presetItems = this.wheelItems.get(preset) || new Map();

            presetItems.set(item.id, item);

            this.wheelItems.set(preset, presetItems);
        });
    }

    addToList(item) {
        this.wheelItemsList.set(item.id, item);
    }

    getByPreset(preset) {
        return this.wheelItems.get(preset);
    }

    getById(id) {
        return this.wheelItemsList.get(id);
    }
}