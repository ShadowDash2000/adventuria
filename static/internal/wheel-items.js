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
            this.wheelItems[e.detail.record.type].delete(e.detail.record.id);
        });
    }


    async fetch() {
        for (const item of await this.pb.collection(this.collectionName).getFullList()) {
            this.addItem(item);
        }
    }

    addItem(item) {
        const wheelItemsType = this.wheelItems.get(item.type) || new Map();
        wheelItemsType.set(item.id, item);
        this.wheelItems.set(item.type, wheelItemsType);
    }

    getByType(type) {
        return this.wheelItems.get(type);
    }
}