import {app} from "../app.js";
import Submit from "./submit.js";

document.addEventListener('DOMContentLoaded', () => {
    const resultModal = document.querySelector('.graph-modal__content.movie-result');
    const comment = resultModal.querySelector('textarea');
    const doneButton = resultModal.querySelector('.button.done');

    doneButton.addEventListener('click', gameResultActions);

    function gameResultActions(e) {
        e.preventDefault();

        const submit = new Submit({
            text: 'Вы уверены, что хотите завершить просмотр?',
            onAccept: () => {
                fetch('/api/movie-done', {
                    method: "POST",
                    headers: {
                        "Authorization": app.auth.token,
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({
                        "comment": comment.value,
                    }),
                });
            },
            onDecline: () => {
                app.modal.open('movie-result', {
                    speed: 100,
                    animation: 'fadeInUp',
                });
            },
        });

        submit.open();
    }

    document.addEventListener('modal.open', async (e) => {
        const modalName = e.detail.modalName;

        if (modalName !== 'movie-result') return;

        const res = await fetch('/api/get-last-action', {
            method: "GET",
            headers: {
                "Authorization": app.auth.token,
            },
        });

        if (!res.ok) return;

        const json = await res.json();

        const title = resultModal.querySelector('.game-title');
        title.innerText = json.title;
    });
});