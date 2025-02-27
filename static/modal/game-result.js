import {app} from "../app.js";
import Submit from "./submit.js";

document.addEventListener('DOMContentLoaded', () => {
    const gameResultModal = document.querySelector('.graph-modal__content.game-result');
    const gameResultComment = gameResultModal.querySelector('textarea');
    const rerollButton = gameResultModal.querySelector('.button.reroll');
    const dropButton = gameResultModal.querySelector('.button.drop');
    const doneButton = gameResultModal.querySelector('.button.done');

    rerollButton.addEventListener('click', gameResultActions);
    dropButton.addEventListener('click', gameResultActions);
    doneButton.addEventListener('click', gameResultActions);

    function gameResultActions(e) {
        e.preventDefault();

        const action = e.currentTarget.dataset.action;
        let text = '';
        switch (action) {
            case 'reroll':
                text = 'Вы уверены, что хотите рерольнуть игру?';
                break;
            case 'drop':
                text = 'Вы уверены, что хотите дропнуть игру?';
                break;
            case 'done':
                text = 'Вы уверены, что хотите завершить прохождение?';
        }

        app.submit.open({
            text: text,
            onAccept: () => {
                fetch('/api/' + action, {
                    method: "POST",
                    headers: {
                        "Authorization": app.auth.token,
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({
                        "comment": gameResultComment.value,
                    }),
                });
                gameResultComment.value = '';
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
                "Authorization": app.auth.token,
            },
        });

        if (!res.ok) return;

        const json = await res.json();

        const gameTitle = gameResultModal.querySelector('.game-title');
        gameTitle.innerText = json.title;

        if (json.isInJail) {
            gameResultModal.querySelector('.button.drop').classList.add('hidden');
        }
    });
});