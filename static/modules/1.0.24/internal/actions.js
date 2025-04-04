import Helper from "../helper.js";
import {app} from "../app.js";
import {openActionUpdateModal} from "../modal/result-update.js";

export default class Actions {
    collectionName = 'actions';
    editableActions = [
        'chooseResult',
        'drop',
        'reroll',
    ]

    constructor(pb, cells, users) {
        this.pb = pb;
        this.cells = cells;
        this.users = users;

        this.page = 1;
        this.limit = 10;
        this.totalPages = 1;
        this.isLoading = false;

        this.actionContainer = document.querySelector('.actions .container');
        this.actionsSentinel = document.querySelector('.actions .sentinel');
        this.actionTemplate = document.getElementById('action-template');
        this.actionsTypes = document.querySelectorAll(' .actions .actions__actions-type .actions-type__item');
        this.actionType = 'all';

        this.actionsTypes.forEach(el => {
            el.addEventListener('click', async (e) => {
                await this.changeActionType(e);
            });
        });

         this.observer = new IntersectionObserver(async (entries, observer) => {
            const entry = entries[0];

            if (entry.isIntersecting && !this.isLoading) {
                if (this.page > this.totalPages) {
                    observer.unobserve(this.actionsSentinel);
                    return;
                }

                await this.fetchNextActions();
            }
        });

        this.observer.observe(this.actionsSentinel);
    }

    async fetchNextActions() {
        this.isLoading = true;
        const actions = await this.fetch(this.page, this.limit, this.actionType);
        this.totalPages = actions.totalPages;
        this.page++;

        for (const action of actions.items) {
            const actionNode = this.createActionNode(action);
            this.appendActionNode(actionNode);
        }

        this.isLoading = false;
    }

    async fetch(page, limit, action = '') {
        let actions = ["roll", "reroll", "drop", "chooseResult", "chooseGame", "rollCell", "rollWheelPreset"];
        if (action.length > 0 && action !== 'all') {
            actions = [action];
        }

        const actionsFilter = actions.map(action => `type="${action}"`).join("||");

        return await this.pb.collection(this.collectionName).getList(page, limit, {
            filter: `${actionsFilter}`,
            sort: '-created',
        });
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

        const actionType = actionNode.querySelector('.action__type')
        const actionIcon = actionType.querySelector('img');
        const actionText = actionType.querySelector('span');

        if (action.icon) {
            actionIcon.src = Helper.getFile('icon', action);
        } else if (actionParams.icon) {
            actionIcon.src = actionParams.icon;
        }

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

        const currentUser = app.getUserRecord();
        if (currentUser?.id === user.id && this.editableActions.includes(action.type)) {
            const actionEdit = actionNode.querySelector('.action__edit');
            actionEdit.dataset.actionId = action.id;
            actionEdit.classList.remove('hidden');
            actionEdit.addEventListener('click', async () => {
                await openActionUpdateModal(action.id);
            });
        }

        return actionNode;
    }

    appendActionNode(actionNode) {
        this.actionContainer.appendChild(actionNode);
    }

    prependActionNode(actionNode) {
        this.actionContainer.prepend(actionNode);
    }

    resetActions() {
        this.page = 1;
        this.totalPages = 1;
        this.observer.unobserve(this.actionsSentinel);
    }

    async changeActionType(e) {
        if (this.isLoading) return;

        const target = e.currentTarget;

        if (target.classList.contains('active')) return;

        this.actionsTypes.forEach(el => {
            el.classList.remove('active');
        });

        target.classList.add('active');
        this.actionType = target.dataset.type;

        this.resetActions();
        this.actionContainer.innerHTML = '';
        await this.fetchNextActions();
        this.observer.observe(this.actionsSentinel);
    }
}