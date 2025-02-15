import GraphModal from "/graph-modal/graph-modal.js";
import PocketBase from "/pocketbase/pocketbase.es.js";

class App {
    constructor() {
        this.pb = new PocketBase('');
        let auth = localStorage.getItem('pocketbase_auth');
        if (auth) this.auth = JSON.parse(auth);
        this.isAuthorized = !!auth;
        this.usersCells = new Map();

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
            this.specialCellTemplate = document.getElementById('special-cell-template');

            await this.updateCells();
            await this.updateUsers();
            await this.updateInnerField();

            setTimeout(() => {
                this.updateUsers();
            }, 2000)
        });
    }

    async updateCells() {
        this.cellsList = await app.pb.collection('cells').getFullList({
            sort: '-sort',
        });

        for (const key in this.cellsList) {
            let cell = this.cellsList[key];

            let cellContainer, cellNode;
            if (cell.position === 'special') {
                cellContainer = this.positions[cell.position][cell.code];
                cellNode = this.specialCellTemplate.content.cloneNode(true);
            } else {
                cellContainer = this.positions[cell.position];
                cellNode = this.cellTemplate.content.cloneNode(true);
            }

            cellNode.querySelector('img').src = "/api/files/" + cell.collectionId + "/" + cell.id + "/" + cell.icon;
            cellNode.querySelector('.name').innerHTML = cell.name;

            this.cellsList[key]['cellElement'] = cellContainer.appendChild(cellNode.firstElementChild);
        }
    }

    async updateUsers() {
        this.usersList = await app.pb.collection('users').getFullList({
            sort: '-points',
        });

        this.updateUsersFields();
        this.updateUsersTable();
    }

    updateUsersFields() {
        this.usersCells.forEach((userCell) => {
            userCell.remove();
        });

        this.usersCells.clear();
        for (const user of this.usersList) {
            const currentCellNum = user.cellsPassed % this.cellsList.length;
            const currentCell = this.cellsList[this.cellsList.length - currentCellNum - 1].cellElement;

            const userElement = document.createElement("img");
            userElement.src = "/api/files/" + user.collectionId + "/" + user.id + "/" + user.avatar;
            userElement.setAttribute('style', "border: 2px solid" + user.color);

            const usersNode = currentCell.querySelector('.users');
            usersNode.appendChild(userElement);

            this.usersCells.set(user.name, userElement);
        }
    }

    updateUsersTable() {
        for (const user of this.usersList) {

        }
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
            default:
                button = actionsButtons.querySelector('button.game-roll');
        }

        button.classList.remove('hidden');
    }

    async updateInnerField() {
        await this.showActionButtons();
    }
}

export const app = new App();