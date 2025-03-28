export default class Inventories {
    collectionName = 'inventory';

    constructor(pb) {
        this.pb = pb;
        this.inventories = new Map();

        document.addEventListener('record.inventory.create', async (e) => {
            this.addInventoryItem(e.detail.record);
        });
        document.addEventListener('record.inventory.update', async (e) => {
            this.addInventoryItem(e.detail.record);
        });
        document.addEventListener('record.inventory.delete', async (e) => {
            this.deleteInventoryItem(e.detail.record);
        });
    }

    async fetch() {
        for (const inventoryItem of await this.pb.collection(this.collectionName).getFullList()) {
            this.addInventoryItem(inventoryItem);
        }
    }

    getByUserId(userId) {
        return this.inventories.get(userId);
    }

    addInventoryItem(inventoryItem) {
        const userInventory = this.getByUserId(inventoryItem.user) || new Map();
        userInventory.set(inventoryItem.id, inventoryItem);
        this.inventories.set(inventoryItem.user, userInventory);
    }

    deleteInventoryItem(inventoryItem) {
        const userInventory = this.getByUserId(inventoryItem.user);
        userInventory?.delete(inventoryItem.id);
    }
}