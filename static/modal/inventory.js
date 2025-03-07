import {app} from "../app.js";
import Helper from "../helper.js";


const inventoryModal = document.getElementById('inventory-modal');
const inventoryItems = inventoryModal.querySelector('.inventory__items');
const inventorySideEffects = inventoryModal.querySelector('.inventory__side-effects');

const inventoryItemTemplate = document.getElementById('inventory-item');
const inventorySideEffectTemplate = document.getElementById('inventory-side-effect');

document.addEventListener('inventory.open', (e) => {
    openInventory(e.detail.userId);
});

function openInventory(userId) {
    inventoryItems.innerHTML = '';
    inventorySideEffects.innerHTML = '';

    const user = app.users.getById(userId);
    const inventory = app.inventories.getByUserId(userId);
    if (!inventory) {
        inventoryModal.querySelector('h2').innerHTML = `В ИНВЕНТАРЕ ${user.name} ПУСТО`;

        app.modal.open('inventory', {
            speed: 100,
            animation: 'fadeInUp',
        });

        return;
    }

    inventoryModal.querySelector('h2').innerHTML = `ИНВЕНТАРЬ ${user.name}`;

    setTimeout(() => {putInventoryToModal(userId, inventory)}, 0);

    app.modal.open('inventory', {
        speed: 100,
        animation: 'fadeInUp',
    });
}

function putInventoryToModal(userId, inventory) {
    inventory.forEach((inventoryItem) => {
        const itemId = inventoryItem.item;
        const item = app.items.getById(itemId);

        if (item.isUsingSlot) {
            const itemNode = inventoryItemTemplate.content.cloneNode(true).firstElementChild;

            const img = itemNode.querySelector('img');
            img.src = Helper.getFile('icon', item);
            img.dataset.id = itemId;
            img.dataset.type = 'item';

            itemNode.querySelector('span').innerText = item.name;
            itemNode.dataset.inventoryItemId = inventoryItem.id;
            itemNode.dataset.itemId = item.id;

            if (userId === app.getUserId()) {
                itemNode.querySelector('.inventory__item-actions').classList.remove('hidden');

                if (!item.canDrop) {
                    itemNode.querySelector('button.drop').classList.add('disabled');
                } else {
                    itemNode.querySelector('button.drop').addEventListener('click', (e) => {
                        submitBeforeItemDrop(userId, e);
                    });
                }

                if (inventoryItem.isActive) {
                    itemNode.querySelector('button.use').classList.add('disabled');
                } else {
                    itemNode.querySelector('button.use').addEventListener('click', useItem);
                }
            }

            inventoryItems.appendChild(itemNode);
        } else {
            const itemNode = inventorySideEffectTemplate.content.cloneNode(true);

            itemNode.querySelector('img').src = Helper.getFile('icon', item);

            inventorySideEffects.appendChild(itemNode);
        }
    });
}

async function useItem(e) {
    const inventoryItemId = e.target.closest('.inventory__item').dataset.inventoryItemId;

    const res = await fetch('/api/use-item', {
        method: "POST",
        headers: {
            "Authorization": app.getUserAuthToken(),
            "Content-type": 'application/json',
        },
        body: JSON.stringify({
            "itemId": inventoryItemId,
        }),
    });

    if (!res.ok) return;

    e.target.classList.add('disabled');
}

function submitBeforeItemDrop(userId, e) {
    const inventoryItem = e.target.closest('.inventory__item');
    const inventoryItemId = inventoryItem.dataset.inventoryItemId;
    const itemId = inventoryItem.dataset.itemId;
    const item = app.items.getById(itemId);

    app.submit.open({
        text: `Вы уверены, что хотите выбросить предмет ${item.name}?`,
        onAccept: async () => {
            await dropItem(inventoryItemId);
            openInventory(userId);
        },
        onDecline: () => {
            openInventory(userId);
        },
    });
}

async function dropItem(itemId) {
    await fetch('/api/drop-item', {
        method: "POST",
        headers: {
            "Authorization": app.getUserAuthToken(),
            "Content-type": 'application/json',
        },
        body: JSON.stringify({
            "itemId": itemId,
        }),
    });
}