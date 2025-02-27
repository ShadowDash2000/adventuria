import {app} from "../app.js";

document.addEventListener('DOMContentLoaded', () => {
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

        const user = app.usersList.get(userId);
        if (!app.inventories[userId]) {
            inventoryModal.querySelector('h2').innerHTML = `В ИНВЕНТАРЕ ${user.name} ПУСТО`;

            app.modal.open('inventory', {
                speed: 100,
                animation: 'fadeInUp',
            });

            return;
        }

        app.inventories[userId].forEach((inventoryItem) => {
            const itemId = inventoryItem.item;
            const item = app.items.get(itemId);

            inventoryModal.querySelector('h2').innerHTML = `ИНВЕНТАРЬ ${user.name}`;

            if (item.isUsingSlot) {
                const itemNode = inventoryItemTemplate.content.cloneNode(true);

                itemNode.querySelector('img').src = app.getFile('icon', item);
                itemNode.querySelector('span').innerText = item.name;
                itemNode.firstElementChild.dataset.id = item.id;

                if (userId === app.auth.record.id) {
                    itemNode.querySelector('.inventory__item-actions').classList.remove('hidden');
                }

                if (!item.canDrop) {
                    itemNode.querySelector('button.drop').classList.add('disabled');
                }

                inventoryItems.appendChild(itemNode);
            } else {
                const itemNode = inventorySideEffectTemplate.content.cloneNode(true);

                itemNode.querySelector('img').src = app.getFile('icon', item);

                inventoryItems.appendChild(itemNode);
            }
        });

        app.modal.open('inventory', {
            speed: 100,
            animation: 'fadeInUp',
        });
    }
});