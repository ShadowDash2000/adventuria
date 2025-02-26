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
        switch (app.nextStepType) {
            case 'rollJailCell':
                const cells = await app.pb.collection('cells').getFullList({
                    sort: '-sort',
                    filter: 'type = "game"',
                });

                for (const cell of cells) {
                    wheelItems.push({
                        id: cell.id,
                        src: "/api/files/" + cell.collectionId + "/" + cell.id + "/" + cell.icon,
                        text: cell.name
                    });
                }
                break;
            case 'rollBigWin':
                const games = await app.pb.collection('wheel_items').getFullList({
                    filter: 'type = "legendaryGame"',
                });

                for (const game of games) {
                    wheelItems.push({
                        id: game.id,
                        src: "/api/files/" + game.collectionId + "/" + game.id + "/" + game.icon,
                        text: game.name
                    });
                }
                break;
            case 'rollMovie':
                const movies = await app.pb.collection('wheel_items').getFullList({
                    filter: 'type = "movie"',
                });

                for (const movie of movies) {
                    wheelItems.push({
                        id: movie.id,
                        src: "/api/files/" + movie.collectionId + "/" + movie.id + "/" + movie.icon,
                        text: movie.name
                    });
                }
                break;
            case 'rollItem':
                const items = await app.pb.collection('items').getFullList({
                    filter: 'isRollable = true',
                });

                for (const item of items) {
                    wheelItems.push({
                        id: item.id,
                        src: "/api/files/" + item.collectionId + "/" + item.id + "/" + item.icon,
                        text: item.name
                    });
                }
                break;
        }

        if (wheelItems) {
            wheel.createWheel(wheelItems);
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

    const audio = new Audio("/api/files/" + rollInfo.collectionId + "/" + rollInfo.id + "/" + rollInfo.audio,);
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