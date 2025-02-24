import GraphModal from "/graph-modal/graph-modal.js";
import PocketBase from "/pocketbase/pocketbase.es.js";

class App {
    constructor() {
        this.pb = new PocketBase('');
        let auth = localStorage.getItem('pocketbase_auth');
        if (auth) this.auth = JSON.parse(auth);
        this.isAuthorized = !!auth;
        this.usersCells = new Map();
        this.usersList = new Map();
        this.nextStepType = '';

        this.pb.collection('users').subscribe('*', (e) => {
            document.dispatchEvent(new CustomEvent("record.users."+e.action, {
                detail: {
                    'record': e.record,
                },
            }));
        });

        this.pb.collection('actions').subscribe('*', (e) => {
            document.dispatchEvent(new CustomEvent("record.actions."+e.action, {
                detail: {
                    'record': e.record,
                },
            }));
        });

        this.modal = new GraphModal({
            isOpen: (modal) => {
                const activeModal = modal.modal.querySelector('.graph-modal-open');
                const modalName = activeModal.dataset.graphTarget;

                document.dispatchEvent(new CustomEvent("modal.open", {
                    detail: {
                        activeModal,
                        modalName,
                    },
                }));
            },
            isClose: () => {
                document.dispatchEvent(new CustomEvent("modal.close"));
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

            await this.fetchCells();
            await this.fetchUsers();
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
            if (e.detail.record.user !== this.auth.record.id) return;
            await this.showActionButtons();
        });
    }


    async fetchCells() {
        this.cellsList = await app.pb.collection('cells').getFullList({
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

            cellNode.querySelector('img').src = "/api/files/" + cell.collectionId + "/" + cell.id + "/" + cell.icon;
            cellNode.querySelector('.name').innerHTML = cell.name;

            this.cellsList[key]['cellElement'] = cellContainer.appendChild(cellNode.firstElementChild);
        }
    }

    async fetchUsers() {
        const usersList = await app.pb.collection('users').getFullList({
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

    updateUsersFields() {
        this.usersCells.forEach((userCell) => {
            userCell.remove();
        });

        this.usersCells.clear();
        this.usersList.forEach((user) => {
            const currentCellNum = user.cellsPassed % this.cellsList.length;
            const currentCell = this.cellsList[this.cellsList.length - currentCellNum - 1].cellElement;

            const userElement = document.createElement("img");
            userElement.src = "/api/files/" + user.collectionId + "/" + user.id + "/" + user.avatar;
            userElement.setAttribute('style', "border: 2px solid" + user.color);

            const usersNode = currentCell.querySelector('.users');
            usersNode.appendChild(userElement);

            this.usersCells.set(user.name, userElement);
        });
    }

    updateUsersTable() {
        const usersTable = document.querySelector('table.users tbody');
        usersTable.innerHTML = '';

        this.usersList.forEach((user) => {
            const userItemNode = this.usersTableItemTemplate.content.cloneNode(true);

            userItemNode.querySelector('.users__avatar img').src = "/api/files/" + user.collectionId + "/" + user.id + "/" + user.avatar;
            userItemNode.querySelector('.users__name').innerHTML = user.name;
            userItemNode.querySelector('.users__points').innerHTML = user.points;

            usersTable.appendChild(userItemNode.firstElementChild);
        });
    }

    async showActionButtons() {
        if (!this.isAuthorized) return;

        const res = await fetch('/api/get-next-step-type', {
            method: "GET",
            headers: {
                "Authorization": this.auth.token,
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
                button = actionsButtons.querySelector('button.wheel');
                break;
        }

        button.classList.remove('hidden');
    }

    async updateInnerField() {
        await this.showActionButtons();
    }
}

export const app = new App();