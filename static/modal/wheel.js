import {app} from "../app.js";
import Wheel from "../wheel.js";

const wheel = new Wheel();

document.addEventListener('DOMContentLoaded', () => {
    const wheelModal = document.querySelector('.wheel-modal');
    const startButton = wheelModal.querySelector('.start-btn');

    document.addEventListener('modal.open', async (e) => {
        const modalName = e.detail.modalName;

        if (modalName !== 'wheel') return;

        let wheelItems = [];
        const currentCell = app.getUserCurrentCell(app.auth.record.id);
        switch (app.nextStepType) {
            case 'rollJailCell':
                for (const cell of app.cellsList) {
                    if (cell.type === 'game') {
                        wheelItems.push({
                            id: cell.id,
                            src: app.getFile('icon', cell),
                            text: cell.name
                        });
                    }
                }
                break;
            case 'rollBigWin':
                app.wheelItems['legendaryGame'].forEach((game) => {
                    wheelItems.push({
                        id: game.id,
                        src: app.getFile('icon', game),
                        text: game.name
                    });
                });
                break;
            case 'rollMovie':
                app.wheelItems['movie'].forEach((movie) => {
                    if (movie.preset === currentCell.preset) {
                        wheelItems.push({
                            id: movie.id,
                            src: app.getFile('icon', movie),
                            text: movie.name
                        });
                    }
                });
                break;
            case 'rollItem':
                app.items.forEach((item) => {
                    if (item.isRollable) {
                        wheelItems.push({
                            id: item.id,
                            src: app.getFile('icon', item),
                            text: item.name
                        });
                    }
                })
                break;
            case 'rollDeveloper':
                app.wheelItems['developer'].forEach((game) => {
                    if (game.preset === currentCell.preset) {
                        wheelItems.push({
                            id: game.id,
                            src: app.getFile('icon', game),
                            text: game.name
                        });
                    }
                });
                break;
        }

        if (wheelItems) {
            wheel.createWheel(wheelItems);
            wheel.rotate();
            startButton.addEventListener('click', startSpin);
        }
    });

    document.addEventListener('modal.close', () => {
        startButton.removeEventListener('click', startSpin);
        wheel.clearWheel();
    });
});

async function startSpin() {
    let url = '';
    switch (app.nextStepType) {
        case 'rollJailCell':
            url = '/api/roll-cell';
            break;
        case 'rollBigWin':
            url = '/api/roll-big-win';
            break;
        case 'rollMovie':
            url = '/api/roll-movie';
            break;
        case 'rollItem':
            url = '/api/roll-item';
            break;
        case 'rollDeveloper':
            url = '/api/roll-developer';
            break;
    }

    const res = await fetch(url, {
        method: "POST",
        headers: {
            "Authorization": app.auth.token,
        },
    });

    if (!res.ok) return;

    const json = await res.json();
    let itemId = json.itemId;

    const audioItemsKeys = Array.from(app.audio[app.nextStepType].keys());
    const randomKey = audioItemsKeys[Math.floor(Math.random() * audioItemsKeys.length)];
    const rollInfo = app.audio[app.nextStepType].get(randomKey);

    if (itemId) wheel.startSpin(itemId, rollInfo.duration);

    const wheelContainer = document.querySelector('.graph-modal__content.wheel-modal');
    const wheelTitle = wheelContainer.querySelector('h2');

    const audio = new Audio("/api/files/" + rollInfo.collectionId + "/" + rollInfo.id + "/" + rollInfo.audio);
    audio.volume = app.volume / 100;
    audio.play();

    const interval = setInterval(() => {
        wheelTitle.innerText = wheel.getCurrentWinner().text;
    }, 100);

    setTimeout(() => {
        audio.pause();
        clearTimeout(interval);
    }, (rollInfo.duration + 1) * 1000);
}