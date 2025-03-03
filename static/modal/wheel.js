import {app} from "../app.js";
import Wheel from "../wheel.js";
import Helper from "../helper.js";

const wheel = new Wheel();

document.addEventListener('DOMContentLoaded', () => {
    const wheelModal = document.querySelector('.wheel-modal');
    const startButton = wheelModal.querySelector('.start-btn');

    document.addEventListener('modal.open', async (e) => {
        const modalName = e.detail.modalName;

        if (modalName !== 'wheel') return;

        let wheelItems = [];
        const currentCell = app.users.getUserCurrentCell(app.getUserId());
        switch (app.nextStepType) {
            case 'rollJailCell':
                for (const cell of app.cells.getAll()) {
                    if (cell.type === 'game') {
                        wheelItems.push({
                            id: cell.id,
                            src: Helper.getFile('icon', cell),
                            text: cell.name,
                            type: 'cell',
                        });
                    }
                }
                break;
            case 'rollBigWin':
                app.wheelItems.getByType('legendaryGame').forEach((game) => {
                    wheelItems.push({
                        id: game.id,
                        src: Helper.getFile('icon', game),
                        text: game.name
                    });
                });
                break;
            case 'rollMovie':
                app.wheelItems.getByType('movie').forEach((movie) => {
                    if (movie.preset === currentCell.preset) {
                        wheelItems.push({
                            id: movie.id,
                            src: Helper.getFile('icon', movie),
                            text: movie.name
                        });
                    }
                });
                break;
            case 'rollItem':
                app.items.getAll().forEach((item) => {
                    if (item.isRollable) {
                        wheelItems.push({
                            id: item.id,
                            src: Helper.getFile('icon', item),
                            text: item.name,
                            type: 'item',
                        });
                    }
                })
                break;
            case 'rollDeveloper':
                app.wheelItems.getByType('developer').forEach((game) => {
                    if (game.preset === currentCell.preset) {
                        wheelItems.push({
                            id: game.id,
                            src: Helper.getFile('icon', game),
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
            "Authorization": app.getUserAuthToken(),
        },
    });

    if (!res.ok) return;

    app.modal.lockClose();

    const json = await res.json();

    const rollInfo = app.audios.getRandomAudio(app.nextStepType);

    wheel.startSpin(json.itemId, rollInfo.duration);

    const wheelContainer = document.querySelector('.graph-modal__content.wheel-modal');
    const wheelTitle = wheelContainer.querySelector('h2');

    app.setAudioSrc(Helper.getFile('audio', rollInfo));
    app.audioPlayer.play();

    const interval = setInterval(() => {
        wheelTitle.innerText = wheel.getCurrentWinner().text;
    }, 100);

    setTimeout(() => {
        app.modal.unlockClose();
        app.audioPlayer.pause();
        clearTimeout(interval);
    }, (rollInfo.duration + 1) * 1000);
}