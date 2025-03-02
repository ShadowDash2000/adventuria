import {app} from "../app.js";
import Submit from "./submit.js";

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

        inventoryModal.querySelector('h2').innerHTML = `ИНВЕНТАРЬ ${user.name}`;

        app.inventories[userId].forEach((inventoryItem) => {
            const itemId = inventoryItem.item;
            const item = app.items.get(itemId);

            if (item.isUsingSlot) {
                const itemNode = inventoryItemTemplate.content.cloneNode(true);

                const img = itemNode.querySelector('img');
                img.src = app.getFile('icon', item);
                img.dataset.id = itemId;
                img.dataset.type = 'item';

                itemNode.querySelector('span').innerText = item.name;
                itemNode.firstElementChild.dataset.id = inventoryItem.id;

                if (userId === app.getUserId()) {
                    itemNode.querySelector('.inventory__item-actions').classList.remove('hidden');
                }

                if (!item.canDrop) {
                    itemNode.querySelector('button.drop').classList.add('disabled');
                } else {
                    itemNode.querySelector('button.drop').addEventListener('click', () => {
                        app.submit.open({
                            text: `Вы уверены, что хотите выбросить предмет ${item.name}?`,
                            onAccept: async () => {
                                await dropItem(inventoryItem.id);
                                openInventory(userId);
                            },
                            onDecline: () => {
                                openInventory(userId);
                            },
                        });
                    });
                }

                if (inventoryItem.isActive) {
                    itemNode.querySelector('button.use').classList.add('disabled');
                } else {
                    itemNode.querySelector('button.use').addEventListener('click', useItem);
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

    async function useItem(e) {
        const itemId = e.target.closest('.inventory__item').dataset.id;

        const res = await fetch('/api/use-item', {
            method: "POST",
            headers: {
                "Authorization": app.getUserAuthToken(),
                "Content-type": 'application/json',
            },
            body: JSON.stringify({
                "itemId": itemId,
            }),
        });

        if (!res.ok) return;

        e.target.classList.add('disabled');
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
});