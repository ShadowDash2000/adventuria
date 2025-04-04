import {app} from "../app.js";

const gameResultModal = document.querySelector('.graph-modal__content.game-result-update');
const gameResultForm = gameResultModal.querySelector('form');
const gameTitle = gameResultModal.querySelector('.game-title');
const comment = gameResultForm.querySelector('textarea[name="comment"]');
const actionIdInput = gameResultForm.querySelector('input[name="actionId"]');
const doneButton = gameResultModal.querySelector('.button.done');

const resultFileBlock = gameResultModal.querySelector('.result-file');
const resultFile = gameResultForm.querySelector('input[type="file"]');
const error = gameResultForm.querySelector('.error');

resultFile.addEventListener('change', onInputChange);
resultFileBlock.addEventListener('dragenter', onDragEnter);
resultFileBlock.addEventListener('dragleave', onDragLeave);
resultFileBlock.addEventListener('dragover', (e) => {e.preventDefault()});
resultFileBlock.addEventListener('drop', onDrop);

doneButton.addEventListener('click', update);

let actionId = 0;

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

function update(e) {
    e.preventDefault();

    app.submit.open({
        text: 'Вы уверены, что хотите сохранить результат?',
        onAccept: async () => {
            const formData = new FormData(gameResultForm);
            updateError('');

            const res = await fetch('/api/update-action', {
                method: "POST",
                headers: {
                    "Authorization": app.getUserAuthToken(),
                },
                body: formData,
            });

            if (!res.ok) {
                updateError(await res.text());
                app.modal.open('game-result-update');
                return;
            }

            gameResultForm.reset();
            updateFileName('');
        },
        onDecline: () => {
            app.modal.open('game-result-update');
        },
    });
}

export async function openActionUpdateModal(id) {
    actionId = id;

    const action = await app.pb.collection('actions').getOne(actionId);

    actionIdInput.value = actionId;
    gameTitle.innerText = action.value;
    comment.value = action.comment;
    updateFileName(action.icon)

    app.modal.open('game-result-update', {
        speed: 100,
        animation: 'fadeInUp',
    });
}