import {app} from "../app.js";

document.addEventListener('DOMContentLoaded', () => {
    const gameResultModal = document.querySelector('.graph-modal__content.game-result');
    const gameResultComment = gameResultModal.querySelector('textarea');
    const rerollButton = gameResultModal.querySelector('.button.reroll');
    const dropButton = gameResultModal.querySelector('.button.drop');
    const doneButton = gameResultModal.querySelector('.button.done');

    rerollButton.addEventListener('click', gameResultActions);
    dropButton.addEventListener('click', gameResultActions);
    doneButton.addEventListener('click', gameResultActions);

    const submitModal = document.querySelector('.graph-modal__content.submit');
    const submitDeclineButton = submitModal.querySelector('.button.decline');
    const submitAcceptButton = submitModal.querySelector('.button.accept');

    submitDeclineButton.addEventListener('click', submitActions);
    submitAcceptButton.addEventListener('click', submitActions);

    document.addEventListener('modal.submit.accept', async (e) => {
        const action = e.detail.action;
        const previousModal = e.detail.previousModal;

        if (previousModal !== 'game-result') return;

        const res = await fetch('/api/' + action, {
            method: "POST",
            headers: {
                "Authorization": app.auth.token,
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                "comment": gameResultComment.value,
            }),
        });

        if (!res.ok) return;

        await app.updateInnerField();

        if (action === 'done') {
            app.modal.close();
            app.modal.open('dice');
        } else {
            app.modal.close();
            app.modal.open('game-picker');
        }

        gameResultComment.value = '';
    });

    async function gameResultActions(e) {
        e.preventDefault();

        const action = e.currentTarget.dataset.action;

        submitModal.dataset.action = action;
        submitModal.dataset.previousModal = 'game-result';

        const text = submitModal.querySelector('.text');

        switch (action) {
            case 'reroll':
                text.innerHTML = 'Вы уверены, что хотите рерольнуть игру?';
                break;
            case 'drop':
                text.innerHTML = 'Вы уверены, что хотите дропнуть игру?';
                break;
            case 'done':
                text.innerHTML = 'Вы уверены, что хотите завершить прохождение?';
        }

        app.modal.close();
        app.modal.open('submit');
    }

    async function submitActions(e) {
        e.preventDefault();

        const action = e.currentTarget.dataset.action;
        const modalAction = submitModal.dataset.action;
        const previousModal = submitModal.dataset.previousModal;

        switch (action) {
            case 'decline':
                app.modal.close();
                app.modal.open(previousModal);
                break;
            case 'accept':
                document.dispatchEvent(new CustomEvent("modal.submit.accept", {
                    detail: {
                        action: modalAction,
                        previousModal: previousModal,
                    },
                }));
                break;
        }
    }
});