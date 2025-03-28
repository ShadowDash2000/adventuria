export default class Items {
    collectionName = 'items';

    constructor(pb) {
        this.pb = pb;
        this.items = new Map();

        document.addEventListener('record.items.create', async (e) => {
            this.items.set(e.detail.record.id, e.detail.record);
        });
        document.addEventListener('record.items.update', async (e) => {
            this.items.set(e.detail.record.id, e.detail.record);
        });
    }

    async fetch() {
        for (const item of await this.pb.collection(this.collectionName).getFullList()) {
            this.items.set(item.id, item);
        }
    }

    getById(itemId) {
        return this.items.get(itemId);
    }

    getAll() {
        return this.items;
    }
}