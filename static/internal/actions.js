import Helper from "../helper.js";

export default class Actions {
    collectionName = 'actions';

    constructor(pb, cells, users) {
        this.pb = pb;
        this.actions = null;
        this.cells = cells;
        this.users = users;

        this.actionContainer = document.querySelector('.actions .container');
        this.actionTemplate = document.getElementById('action-template');
    }

    async fetch(page) {
        this.actions = await this.pb.collection(this.collectionName).getList(page, 10, {
            filter: '\'["roll", "reroll", "drop", "chooseResult", "rollCell", "rollWheelPreset"]\' ~ type',
            sort: '-created',
        });
    }

    refresh() {
        for (const action of this.actions.items) {
            const actionNode = this.createActionNode(action);
            this.appendActionNode(actionNode);
        }
    }

    createActionNode(action) {
        const actionNode = this.actionTemplate.content.cloneNode(true);
        const userAvatar = actionNode.querySelector('.action__user img');
        const userName = actionNode.querySelector('.action__user span');
        const cellIcon = actionNode.querySelector('.action__cell img');
        const cellName = actionNode.querySelector('.action__cell span');
        const actionIcon = actionNode.querySelector('.action__type img');
        const actionText = actionNode.querySelector('.action__type span');

        const user = this.users.getById(action.user);
        userAvatar.src = Helper.getFile('avatar', user);
        userAvatar.style.borderColor = user.color;
        userName.innerText = user.name;

        const cell = this.cells.getById(action.cell);
        cellIcon.src = Helper.getFile('icon', cell, {'thumb': '250x0'});
        cellIcon.dataset.id = cell.id;
        cellIcon.dataset.type = 'cell';
        cellName.innerText = `НА КЛЕТКЕ ${cell.name}`;

        let text = action.value;
        const textTemplate = Helper.actions[action.type]?.template;
        if (textTemplate) {
            text = textTemplate.replace('{{VALUE}}', action.value);
        }

        actionText.innerText = text;

        return actionNode.firstElementChild;
    }

    appendActionNode(actionNode) {
        this.actionContainer.appendChild(actionNode);
    }

    prependActionNode(actionNode) {
        this.actionContainer.prepend(actionNode);
    }
}