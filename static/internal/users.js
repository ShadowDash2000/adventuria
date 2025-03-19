import Helper from "../helper.js";

export default class Users {
    collectionName = 'users';

    constructor(pb, cells) {
        this.pb = pb;
        this.cells = cells;
        this.users = new Map();
        this.usersElements = [];
        this.usersTableItemTemplate = document.getElementById('users-table-item');

        document.addEventListener('record.users.update', (e) => {
            this.users.set(e.detail.record.id, e.detail.record);
            this.sort();
            this.refreshCells();
            this.refreshTable();
        });
    }

    sort() {
        this.users = new Map([...this.users.entries()].sort((a, b) => b[1].points - a[1].points));
    }

    async fetch() {
        const users = await this.pb.collection(this.collectionName).getFullList({
            sort: '-points',
        });

        for (const user of users) {
            this.users.set(user.id, user);
        }
    }

    refreshCells() {
        this.usersElements.forEach((userCell) => {
            userCell.remove();
        });

        this.users.forEach((user) => {
            const currentCell = this.cells.getByPassed(user.cellsPassed);
            const currentCellElement = currentCell.cellElement;

            user.currentCellId = currentCell.id;

            const userElement = document.createElement("img");
            userElement.src = Helper.getFile('avatar', user);
            userElement.setAttribute('style', `border-color: ${user.color}`);

            const usersNode = currentCellElement.querySelector('.users');
            usersNode.appendChild(userElement);

            this.usersElements.push(userElement);
        });
    }

    refreshTable() {
        const table = document.querySelector('table.users');
        const tbody = document.querySelector(' tbody');
        table.classList.remove('hidden');
        tbody.innerHTML = '';

        this.users.forEach((user) => {
            const userItemNode = this.usersTableItemTemplate.content.cloneNode(true);

            const avatar = userItemNode.querySelector('.users__avatar img');
            avatar.src = Helper.getFile('avatar', user);
            avatar.style.borderColor = user.color;

            avatar.addEventListener('click', () => {
                document.dispatchEvent(new CustomEvent('profile.open', {
                    detail: {
                        userId: user.id,
                    }
                }));
            });

            userItemNode.querySelector('.users__name').innerHTML = user.name;
            userItemNode.querySelector('.users__points').innerHTML = user.points;
            userItemNode.querySelector('.users__cells-passed').innerHTML = user.cellsPassed;
            userItemNode.querySelector('.users__drops').innerHTML = user.stats?.drops || 0;
            userItemNode.querySelector('.users__finished').innerHTML = user.stats?.finished || 0;

            const inventoryButton = userItemNode.querySelector('.users__inventory button');
            inventoryButton.dataset.inventory = user.id;

            inventoryButton.addEventListener('click', () => {
                document.dispatchEvent(new CustomEvent('inventory.open', {
                    detail: {
                        userId: user.id,
                    },
                }));
            });

            tbody.appendChild(userItemNode.firstElementChild);
        });
    }

    getUserCurrentCell(userId) {
        const user = this.users.get(userId);
        return user?.currentCellId ? this.cells.getById(user.currentCellId) : null;
    }

    getById(userId) {
        return this.users.get(userId);
    }
}