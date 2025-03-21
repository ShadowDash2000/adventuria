import {app} from "../app.js";

const gameResultModal = document.querySelector('.graph-modal__content.game-result');
const gameResultForm = gameResultModal.querySelector('form');
const rerollButton = gameResultModal.querySelector('.button.reroll');
const dropButton = gameResultModal.querySelector('.button.drop');
const doneButton = gameResultModal.querySelector('.button.done');

gameResultForm.querySelector('input[type="file"]').addEventListener('change', updateFileName);
rerollButton.addEventListener('click', gameResultActions);
dropButton.addEventListener('click', gameResultActions);
doneButton.addEventListener('click', gameResultActions);

function updateFileName(e) {
    const fileName = gameResultForm.querySelector('.file-name');
    fileName.innerText = e.target.files[0]?.name;
}

function gameResultActions(e) {
    e.preventDefault();

    const action = e.currentTarget.dataset.action;
    let text = '';
    switch (action) {
        case 'reroll':
            text = 'Вы уверены, что хотите рерольнуть?';
            break;
        case 'drop':
            text = 'Вы уверены, что хотите дропнуть?';
            break;
        case 'done':
            text = 'Вы уверены, что хотите завершить прохождение / просмотр фильма?';
    }

    app.submit.open({
        text: text,
        onAccept: async () => {
            const formData = new FormData(gameResultForm);

            const res = await fetch('/api/' + action, {
                method: "POST",
                headers: {
                    "Authorization": app.getUserAuthToken(),
                },
                body: formData,
            });

            if (!res.ok) {
                app.modal.open('game-result');
                return;
            }

            gameResultForm.reset();
        },
        onDecline: () => {
            app.modal.open('game-result');
        },
    });
}

document.addEventListener('modal.open', async (e) => {
    const modalName = e.detail.modalName;

    if (modalName !== 'game-result') return;

    const res = await fetch('/api/get-last-action', {
        method: "GET",
        headers: {
            "Authorization": app.getUserAuthToken(),
        },
    });

    if (!res.ok) return;

    const json = await res.json();

    const currentCell = app.users.getUserCurrentCell(app.getUserId());

    const gameTitle = gameResultModal.querySelector('.game-title');
    gameTitle.innerText = json.title;

    const dropButton = gameResultModal.querySelector('.button.drop');
    if (currentCell.cantDrop || json.isInJail) {
        dropButton.classList.add('hidden');
    } else {
        dropButton.classList.remove('hidden');
    }

    const rerollButton = gameResultModal.querySelector('.button.reroll');
    if (currentCell.cantReroll) {
        rerollButton.classList.add('hidden');
    } else {
        rerollButton.classList.remove('hidden');
    }
});