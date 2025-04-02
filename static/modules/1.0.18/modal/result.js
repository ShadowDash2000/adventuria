import {app} from "../app.js";

const gameResultModal = document.querySelector('.graph-modal__content.game-result');
const gameResultForm = gameResultModal.querySelector('form');
const rerollButton = gameResultModal.querySelector('.button.reroll');
const dropButton = gameResultModal.querySelector('.button.drop');
const doneButton = gameResultModal.querySelector('.button.done');

const resultFileBlock = gameResultModal.querySelector('.result-file');
const resultFile = gameResultForm.querySelector('input[type="file"]');
const error = gameResultForm.querySelector('.error');

resultFile.addEventListener('change', onInputChange);
resultFileBlock.addEventListener('dragenter', onDragEnter);
resultFileBlock.addEventListener('dragleave', onDragLeave);
resultFileBlock.addEventListener('dragover', (e) => {e.preventDefault()});
resultFileBlock.addEventListener('drop', onDrop);

rerollButton.addEventListener('click', gameResultActions);
dropButton.addEventListener('click', gameResultActions);
doneButton.addEventListener('click', gameResultActions);

function updateFileName(name) {
    const fileName = gameResultForm.querySelector('.file-name');
    fileName.innerText = name;
}

function onInputChange(e) {
    updateFileName(e.target.files[0]?.name);
}

let dragTarget = null;

function onDragEnter(e) {
    e.preventDefault();
    if (dragTarget) return;
    dragTarget = e.currentTarget;
    resultFileBlock.classList.add('drag-n-drop');
}

function onDragLeave(e) {
    e.preventDefault();
    if (dragTarget !== e.target) return;
    dragTarget = null;
    resultFileBlock.classList.remove('drag-n-drop');
}

function onDrop(e) {
    e.preventDefault();
    updateFileName(e.dataTransfer.files[0]?.name);
    dragTarget = null;

    const newData = new DataTransfer();
    newData.items.add(e.dataTransfer.files[0]);
    resultFile.files = newData.files;
    resultFileBlock.classList.remove('drag-n-drop');
}

function updateError(text) {
    error.innerText = text;

    if (text.length === 0) error.classList.add('hidden');
    else error.classList.remove('hidden');
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
            updateError('');

            const res = await fetch('/api/' + action, {
                method: "POST",
                headers: {
                    "Authorization": app.getUserAuthToken(),
                },
                body: formData,
            });

            if (!res.ok) {
                updateError(await res.text());
                app.modal.open('game-result');
                return;
            }

            gameResultForm.reset();
            updateFileName('');
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