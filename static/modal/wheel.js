import {app} from "../app.js";
import Wheel from "../wheel.js";

const wheel = new Wheel();

document.addEventListener('DOMContentLoaded', () => {
    const wheelModal = document.querySelector('.wheel-modal');
    const startButton = wheelModal.querySelector('.start-btn');

    document.addEventListener('modal.open', async (e) => {
        const modalName = e.detail.modalName;

        if (modalName !== 'wheel') return;

        let items = [];
        switch (app.nextStepType) {
            case 'rollJailCell':
                const cells = await app.pb.collection('cells').getFullList({
                    sort: '-sort',
                    filter: 'type = "game"',
                });

                for (const cell of cells) {
                    items.push({
                        id: cell.id,
                        src: "/api/files/" + cell.collectionId + "/" + cell.id + "/" + cell.icon,
                        text: cell.name
                    });
                }
                break;
            case 'rollBigWin':
                const games = await app.pb.collection('wheel_items').getFullList({
                    filter: 'type = "bigWin"',
                });

                for (const game of games) {
                    items.push({
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
                    items.push({
                        id: movie.id,
                        src: "/api/files/" + movie.collectionId + "/" + movie.id + "/" + movie.icon,
                        text: movie.name
                    });
                }
                break;
        }

        if (items) {
            wheel.createWheel(items);
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

    if (itemId) wheel.startSpin(itemId, 6);
}