import Helper from "../helper.js";

export default class Actions {
    collectionName = 'actions';

    constructor(pb, cells, users) {
        this.pb = pb;
        this.actions = null;
        this.cells = cells;
        this.users = users;

        this.page = 1;
        this.limit = 10;
        this.totalPages = 1;
        this.isLoading = false;

        this.actionContainer = document.querySelector('.actions .container');
        this.actionsSentinel = document.querySelector('.actions .sentinel');
        this.actionTemplate = document.getElementById('action-template');

        const observer = new IntersectionObserver(async (entries, observer) => {
            const entry = entries[0];

            if (entry.isIntersecting && !this.isLoading) {
                if (this.page > this.totalPages) {
                    observer.unobserve(this.actionsSentinel);
                    return;
                }

                await this.fetch();
            }
        });

        observer.observe(this.actionsSentinel);
    }

    async fetch() {
        this.isLoading = true;

        this.actions = await this.pb.collection(this.collectionName).getList(this.page, this.limit, {
            filter: '\'["roll", "reroll", "drop", "chooseResult", "chooseGame", "rollCell", "rollWheelPreset"]\' ~ type',
            sort: '-created',
        });

        this.totalPages = this.actions.totalPages;
        this.page++;

        for (const action of this.actions.items) {
            const actionNode = this.createActionNode(action);
            this.appendActionNode(actionNode);
        }

        this.isLoading = false;
    }

    createActionNode(action) {
        const actionNode = this.actionTemplate.content.cloneNode(true).firstElementChild;
        const actionParams = Helper.actions[action.type];
        actionNode.style.background = actionParams?.color;

        const actionDate = actionNode.querySelector('.action__date');
        actionDate.innerText = Helper.formatDateLocalized(action.created);

        const userAvatar = actionNode.querySelector('.action__user img');
        const userName = actionNode.querySelector('.action__user span');
        const user = this.users.getById(action.user);
        userAvatar.src = Helper.getFile('avatar', user);
        userAvatar.style.borderColor = user.color;
        userAvatar.loading = 'lazy';
        userName.innerText = user.name;

        const cellIcon = actionNode.querySelector('.action__cell img');
        const cellName = actionNode.querySelector('.action__cell span');
        const cell = this.cells.getById(action.cell);
        cellIcon.src = Helper.getFile('icon', cell, {'thumb': '250x0'});
        cellIcon.dataset.id = cell.id;
        cellIcon.dataset.description = cell.description;
        cellIcon.loading = 'lazy';
        cellName.innerText = `НА КЛЕТКЕ: ${cell.name}`;

        const actionText = actionNode.querySelector('.action__type span');
        let text = action.value;
        const textTemplate = actionParams?.template;
        if (textTemplate) {
            text = textTemplate.replace('{{VALUE}}', action.value);
        }
        actionText.innerText = text;

        if (action.comment) {
            const actionComment = actionNode.querySelector('.action__comment p');
            actionComment.innerText = action.comment;
            actionComment.classList.remove('hidden');
        }

        return actionNode;
    }

    appendActionNode(actionNode) {
        this.actionContainer.appendChild(actionNode);
    }

    prependActionNode(actionNode) {
        this.actionContainer.prepend(actionNode);
    }
}