import {app} from "../app.js";

const gameStatsButton = document.getElementById('game-stats');
const gameStatsModal = document.getElementById('game-stats-modal');

document.addEventListener('modal.open', (e) => {
    const modalName = e.detail.modalName;
    if (modalName !== 'game-stats') return;


});
