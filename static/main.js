document.addEventListener('DOMContentLoaded', async () => {
    const pb = new PocketBase('');

    //const resultList = await pb.collection('actions').getFullList({});

    /*
    const authData = await pb.collection('users').authWithPassword(
        'scout0501@gmail.com',
        'kFWo4d_MDZRACSH',
    );
     */

    const modal = new GraphModal({
        isOpen: async (modal) => {
            const activeModal = modal.modal.querySelector('.graph-modal-open');
            const modalName = activeModal.dataset.graphTarget;

            if (modalName === 'game-result') {
                const res = await fetch('/api/game-result', {
                    method: "GET",
                    headers: {
                        "Authorization": auth.token,
                    },
                });

                if (!res.ok) return;

                const json = await res.json();

                const gameTitle = activeModal.querySelector('.game-title');
                gameTitle.innerHTML = json.game;

                if (!json.canDrop) {
                    activeModal.querySelector('.button.drop').classList.add('hidden');
                }
            }
        },
    });

    const authForm = document.getElementById('auth');
    authForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const formData = new FormData(authForm);

        const authResult = await pb.collection('users').authWithPassword(
            formData.get('login'),
            formData.get('password'),
        );

        if (authResult.token) {
            window.location.reload();
        }
    });

    const positions = {
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

    const cellsList = await pb.collection('cells').getFullList({
        sort: '-sort',
    });

    for (const key in cellsList) {
        let cell = cellsList[key];
        let cellHTML = cell.htmlTemplate.replace("{{NAME}}", cell.name);
        cellHTML = cellHTML.replace("{{IMG}}", "/api/files/" + cell.collectionId + "/" + cell.id + "/" + cell.icon);

        const cellElement = document.createElement("div");
        cellElement.innerHTML = cellHTML;

        let container = cellElement.querySelector('.container');
        const usersDiv = document.createElement("div");
        usersDiv.classList.add('users');
        container.appendChild(usersDiv);

        let cellContainer;
        if (cell.position === 'special') {
            cellContainer = positions[cell.position][cell.code].appendChild(cellElement.firstElementChild);
        } else {
            cellContainer = positions[cell.position].appendChild(cellElement.firstElementChild);
        }

        cellsList[key]['cellElement'] = cellContainer;
    }

    await showUsers(cellsList);

    async function showUsers(cellsList) {
        const usersList = await pb.collection('users').getFullList({});

        for (const user of usersList) {
            const currentCellNum = user.cellsPassed % cellsList.length;
            const currentCell = cellsList[cellsList.length - currentCellNum - 1].cellElement;

            const userElement = document.createElement("img");
            userElement.src = "/api/files/" + user.collectionId + "/" + user.id + "/" + user.avatar;
            userElement.setAttribute('style', "border: 2px solid" + user.color);

            currentCell.querySelector('.users')?.appendChild(userElement);
        }
    }

    const rollButton = document.getElementById('roll');
    let auth = localStorage.getItem('pocketbase_auth');
    auth = JSON.parse(auth);

    rollButton.addEventListener('click', async () => {
        const res = await fetch('/api/roll', {
            method: "POST",
            headers: {
                "Authorization": auth.token,
            },
        });

        if (!res.ok) {
            return;
        }

        const json = await res.json()

        rollDice(json.roll, 4);

        setTimeout(() => {
            const modal = document.querySelector('.graph-modal__container.dice');
            const rollResult = modal.querySelector('.roll-result');
            const cell = modal.querySelector('.cell');

            rollResult.querySelector('.roll-result__number').innerHTML = json.roll;

            cell.querySelector('img').src = json.cell.icon;
            cell.querySelector('.cell-info__name').innerHTML = json.cell.name;
            cell.querySelector('.cell-info__description').innerHTML = json.cell.description;

            rollResult.classList.remove('hidden');
            cell.classList.remove('hidden');

            rollButton.classList.add('hidden');
            modal.querySelector('.choose-game').classList.remove('hidden');

        }, 4000);
    });

    const chooseGameButton = document.getElementById('choose-game');
    const gamePicker = document.getElementById('game-picker');
    chooseGameButton.addEventListener('click', async (e) => {
        e.preventDefault();

        await fetch('/api/choose-game', {
            method: "POST",
            headers: {
                "Authorization": auth.token,
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                "game": gamePicker.querySelector('input[name="game"]').value,
            }),
        });
    });

    async function showActionButtons() {
        const res = await fetch('/api/get-last-action', {
            method: "GET",
            headers: {
                "Authorization": auth.token,
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
        if (json.canRoll) {
            button = actionsButtons.querySelector('button.game-roll');
        }

        if (json.cellType !== 'game') {
            button = 'test';
        }

        if (!button) {
            switch (json.status) {
                case 'inProgress':
                    button = actionsButtons.querySelector('button.game-result');
                    break;
                case 'reroll':
                case 'drop':
                case 'notChosen':
                    button = actionsButtons.querySelector('button.game-picker');
                    break;
                case 'done':
                default:
                    button = actionsButtons.querySelector('button.game-roll');
            }
        }

        button.classList.remove('hidden');
    }

    await showActionButtons();

    const gameResultModal = document.querySelector('.graph-modal__content.game-result');
    const gameResultComment = gameResultModal.querySelector('textarea');
    const rerollButton = gameResultModal.querySelector('.button.reroll');
    const dropButton = gameResultModal.querySelector('.button.drop');
    const doneButton = gameResultModal.querySelector('.button.done');

    rerollButton.addEventListener('click', actions);
    dropButton.addEventListener('click', actions);
    doneButton.addEventListener('click', actions);

    async function actions(e) {
        e.preventDefault();

        const action = e.currentTarget.dataset.action;

        const res = await fetch('/api/' + action, {
            method: "POST",
            headers: {
                "Authorization": auth.token,
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                "comment": gameResultComment.value,
            }),
        });

        if (!res.ok) return;

        if (action === 'done') {
            modal.close();
            modal.open('dice');
        } else {
            modal.close();
            modal.open('game-picker');
        }
    }
});