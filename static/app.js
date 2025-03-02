import GraphModal from "/graph-modal/graph-modal.js";
import PocketBase from "/pocketbase/pocketbase.es.js";
import Submit from "./modal/submit.js";
import Timer from './timer.js';

class App {
    constructor() {
        this.pb = new PocketBase('/');
        this.isAuthorized = !!this.pb.authStore.token;
        this.usersCells = new Map();
        this.usersList = new Map();
        this.audio = new Map();
        this.inventories = {};
        this.items = new Map();
        this.wheelItems = new Map();
        this.nextStepType = '';
        this.submit = new Submit();
        this.volume = localStorage.getItem('volume') ? localStorage.getItem('volume') : 30;
        this.audioPlayer = new Audio();
        this.audioPlayer.volume = this.volume / 100;

        const collections = [
            'users',
            'actions',
            'inventory',
            'audio',
            'items',
            'wheel_items',
        ];
        for (const collection of collections) {
            this.pb.collection(collection).subscribe('*', (e) => {
                document.dispatchEvent(new CustomEvent(`record.${collection}.${e.action}`, {
                    detail: {
                        'record': e.record,
                    },
                }));
            });
        }

        this.modal = new GraphModal({
            isOpen: (modal) => {
                const activeModal = modal.modalContainer;
                const modalName = activeModal.dataset.graphTarget;

                document.dispatchEvent(new CustomEvent("modal.open", {
                    detail: {
                        activeModal,
                        modalName,
                    },
                }));
            },
            isClose: (modal) => {
                const activeModal = modal.modalContainer;
                const modalName = activeModal.dataset.graphTarget;

                document.dispatchEvent(new CustomEvent("modal.close", {
                    detail: {
                        activeModal,
                        modalName,
                    },
                }));
            },
        });

        document.addEventListener('DOMContentLoaded', async () => {
            this.positions = {
                'left': document.querySelector('.left-row'),
                'top': document.querySelector('.top-row'),
                'right': document.querySelector('.right-row'),
                'bottom': document.querySelector('.bottom-row'),
                'special': {
                    'start': document.getElementById('start'),
                    'jail': document.getElementById('jail'),
                    'big-win': document.getElementById('big-win'),
                    'preset': document.getElementById('preset'),
                },
            };
            this.cellTemplate = document.getElementById('cell-template');
            this.cellTemplateRight = document.getElementById('cell-template-right');
            this.specialCellTemplate = document.getElementById('special-cell-template');
            this.usersTableItemTemplate = document.getElementById('users-table-item');

            const volumeSlider = document.getElementById('volume-slider');
            volumeSlider.value = this.volume;
            volumeSlider.addEventListener('change', () => {
                this.setVolume(volumeSlider.value);
            });

            document.addEventListener('modal.open', () => {
                volumeSlider.parentElement.classList.add('fixed');
            });
            document.addEventListener('modal.close', () => {
                volumeSlider.parentElement.classList.remove('fixed');
            });

            if (this.getUserAuthToken()) {
                const avatar = this.getFile('avatar', this.pb.authStore.record);
                const profile = document.querySelector('.profile');
                profile.classList.remove('hidden');
                profile.querySelector('img').src = avatar;

                const timerBlock = document.getElementById('timer');
                const timer = new Timer(this.getUserAuthToken(), timerBlock);
                const timerStopButton = document.querySelector('.timer button.red');
                const timerStartButton = document.querySelector('.timer button.green');
                timerStopButton.addEventListener('click', () => {
                    timer.stopTimer();
                });
                timerStartButton.addEventListener('click', () => {
                    timer.startTimer();
                });
            }

            await this.fetchCells();
            await this.fetchUsers();
            await this.fetchItems();
            await this.fetchInventories();
            await this.fetchWheelItems();
            await this.fetchAudio();

            this.updateCells();
            this.updateUsers();

            await this.updateInnerField();
        });

        document.addEventListener('record.users.update', (e) => {
            this.usersList.set(e.detail.record.id, e.detail.record)
            this.updateUsersFields();
            this.updateUsersTable();
        });

        document.addEventListener('record.actions.create', async (e) => {
            if (e.detail.record.user !== this.pb.authStore.record.id) return;
            setTimeout(async () => {
                await this.showActionButtons();
            }, 1000);
        });
        document.addEventListener('record.actions.delete', async (e) => {
            if (e.detail.record.user !== this.pb.authStore.record.id) return;
            setTimeout(async () => {
                await this.showActionButtons();
            }, 1000);
        });

        document.addEventListener('record.inventory.create', async (e) => {
            this.addInventoryItem(e.detail.record);
        });
        document.addEventListener('record.inventory.update', async (e) => {
            this.addInventoryItem(e.detail.record);
        });
        document.addEventListener('record.inventory.delete', async (e) => {
            this.deleteInventoryItem(e.detail.record);
        });

        document.addEventListener('record.audio.create', async (e) => {
            this.addAudioItem(e.detail.record);
        });
        document.addEventListener('record.audio.update', async (e) => {
            this.addAudioItem(e.detail.record);
        });

        document.addEventListener('record.items.create', async (e) => {
            this.items.set(e.detail.record.id, e.detail.record);
        });
        document.addEventListener('record.items.update', async (e) => {
            this.items.set(e.detail.record.id, e.detail.record);
        });

        document.addEventListener('record.wheel_items.create', async (e) => {
            this.addWheelItem(e.detail.record);
        });
        document.addEventListener('record.wheel_items.update', async (e) => {
            this.addWheelItem(e.detail.record);
        });
        document.addEventListener('record.wheel_items.delete', async (e) => {
            this.wheelItems[e.detail.record.type].delete(e.detail.record.id);
        });
    }

    setVolume(v) {
        this.volume = v;
        this.audioPlayer.volume = v / 100;
        localStorage.setItem('volume', v);
    }

    setAudioSrc(src) {
        this.audioPlayer.src = src;
    }

    getRandomAudio(type) {
        if (!this.audio[type]) return;

        const audioItemsKeys = Array.from(this.audio[type].keys());
        const randomKey = audioItemsKeys[Math.floor(Math.random() * audioItemsKeys.length)];
        return this.audio[type].get(randomKey);
    }

    async fetchCells() {
        this.cellsList = await this.pb.collection('cells').getFullList({
            sort: '-sort',
        });
    }

    updateCells() {
        for (const key in this.cellsList) {
            let cell = this.cellsList[key];

            let cellContainer, cellNode;
            switch (cell.position) {
                case 'special':
                    cellContainer = this.positions[cell.position][cell.code];
                    cellNode = this.specialCellTemplate.content.cloneNode(true);
                    break;
                case 'right':
                case 'bottom':
                    cellContainer = this.positions[cell.position];
                    cellNode = this.cellTemplateRight.content.cloneNode(true);
                    break;
                default:
                    cellContainer = this.positions[cell.position];
                    cellNode = this.cellTemplate.content.cloneNode(true);
            }

            const colorBar = cellNode.querySelector('.color-bar');
            if (colorBar) {
                colorBar.style.background = cell.color;
            }

            cellNode.querySelector('img').src = "/api/files/" + cell.collectionId + "/" + cell.id + "/" + cell.icon;
            const name = cellNode.querySelector('.name');
            name.innerHTML = cell.name;
            name.dataset.id = cell.id;
            name.dataset.type = 'cell';

            this.cellsList[key]['cellElement'] = cellContainer.appendChild(cellNode.firstElementChild);
        }
    }

    async fetchInventories() {
        const inventories = await this.pb.collection('inventory').getFullList({});

        for (const inventoryItem of inventories) {
            this.addInventoryItem(inventoryItem);
        }
    }

    addInventoryItem(item) {
        if (!this.inventories[item.user]) {
            this.inventories[item.user] = new Map();
        }
        this.inventories[item.user].set(item.id, item);
    }

    deleteInventoryItem(item) {
        this.inventories[item.user].delete(item.id);
    }

    async fetchUsers() {
        const usersList = await this.pb.collection('users').getFullList({
            sort: '-points',
        });

        for (const user of usersList) {
            this.usersList.set(user.id, user);
        }
    }

    updateUsers() {
        this.updateUsersFields();
        this.updateUsersTable();
    }

    getUserId() {
        if (this.pb.authStore) {
            return this.pb.authStore.record.id;
        }
    }

    getUserAuthToken() {
        if (this.pb.authStore) {
            return this.pb.authStore.token;
        }
    }

    updateUsersFields() {
        this.usersCells.forEach((userCell) => {
            userCell.remove();
        });

        this.usersCells.clear();
        this.usersList.forEach((user, userId) => {
            const currentCellNum = user.cellsPassed % this.cellsList.length;
            const currentCell = this.cellsList[this.cellsList.length - currentCellNum - 1];
            const currentCellElement = currentCell.cellElement;

            user.currentCell = currentCell.id;

            const userElement = document.createElement("img");
            userElement.src = "/api/files/" + user.collectionId + "/" + user.id + "/" + user.avatar;
            userElement.setAttribute('style', "border: 2px solid" + user.color);

            const usersNode = currentCellElement.querySelector('.users');
            usersNode.appendChild(userElement);

            this.usersCells.set(user.name, userElement);
        });
    }

    getUserCurrentCell(userId) {
        const cellId = this.usersList.get(userId).currentCell;
        return this.getCellById(cellId);
    }

    getCellById(cellId) {
        for (const cell of this.cellsList) {
            if (cell.id === cellId) return cell;
        }
    }

    updateUsersTable() {
        const usersTable = document.querySelector('table.users tbody');
        usersTable.innerHTML = '';

        this.usersList.forEach((user) => {
            const userItemNode = this.usersTableItemTemplate.content.cloneNode(true);

            userItemNode.querySelector('.users__avatar img').src = "/api/files/" + user.collectionId + "/" + user.id + "/" + user.avatar;
            userItemNode.querySelector('.users__name').innerHTML = user.name;
            userItemNode.querySelector('.users__points').innerHTML = user.points;
            userItemNode.querySelector('.users__cells-passed').innerHTML = user.cellsPassed;

            const inventoryButton = userItemNode.querySelector('.users__inventory button');
            inventoryButton.dataset.inventory = user.id;

            inventoryButton.addEventListener('click', () => {
                document.dispatchEvent(new CustomEvent('inventory.open', {
                    detail: {
                        userId: user.id,
                    },
                }));
            });

            usersTable.appendChild(userItemNode.firstElementChild);
        });
    }

    async fetchAudio() {
        for (const audio of await this.pb.collection('audio').getFullList()) {
            this.addAudioItem(audio);
        }
    }

    addAudioItem(item) {
        if (!this.audio[item.event]) {
            this.audio[item.event] = new Map();
        }
        this.audio[item.event].set(item.id, item);
    }

    async fetchItems() {
        for (const item of await this.pb.collection('items').getFullList()) {
            this.items.set(item.id, item);
        }
    }

    getFile(key, item) {
        return "/api/files/" + item.collectionId + "/" + item.id + "/" + item[key];
    }

    async fetchWheelItems() {
        for (const item of await this.pb.collection('wheel_items').getFullList()) {
            this.addWheelItem(item);
        }
    }

    addWheelItem(item) {
        if (!this.wheelItems[item.type]) {
            this.wheelItems[item.type] = new Map();
        }
        this.wheelItems[item.type].set(item.id, item);
    }

    async showActionButtons() {
        if (!this.isAuthorized) return;

        const res = await fetch('/api/get-next-step-type', {
            method: "GET",
            headers: {
                "Authorization": this.getUserAuthToken(),
            },
        });

        if (!res.ok) return;

        const json = await res.json();
        this.nextStepType = json.nextStepType;

        const actionsButtons = document.querySelector('.actions-buttons');
        const buttons = actionsButtons.querySelectorAll('button');
        for (const button of buttons) {
            button.classList.add('hidden');
        }

        let button;

        switch (json.nextStepType) {
            case 'chooseResult':
                button = actionsButtons.querySelector('button.game-result');
                break;
            case 'chooseGame':
                button = actionsButtons.querySelector('button.game-picker');
                break;
            case 'roll':
                button = actionsButtons.querySelector('button.game-roll');
                break;
            case 'rollJailCell':
            case 'rollBigWin':
            case 'rollMovie':
            case 'rollPreset':
            case 'rollItem':
            case 'rollDeveloper':
                button = actionsButtons.querySelector('button.wheel');
                break;
            case 'movieResult':
                button = actionsButtons.querySelector('button.movie-result');
                break;
        }

        button.classList.remove('hidden');
    }

    async updateInnerField() {
        await this.showActionButtons();
    }
}

export const app = new App();