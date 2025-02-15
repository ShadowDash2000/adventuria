import {app} from "../app.js";

document.addEventListener('modal.open', async (e) => {
    const activeModal = e.detail.activeModal;
    const modalName = e.detail.modalName;

    if (modalName !== 'game-result') return;

    const res = await fetch('/api/game-result', {
        method: "GET",
        headers: {
            "Authorization": app.auth.token,
        },
    });

    if (!res.ok) return;

    const json = await res.json();

    const gameTitle = activeModal.querySelector('.game-title');
    gameTitle.innerHTML = json.game;

    if (!json.canDrop) {
        activeModal.querySelector('.button.drop').classList.add('hidden');
    }
});