import {app} from "../app.js";
import Wheel from "../wheel.js";
import Helper from "../helper.js";

const wheel = new Wheel();

const wheelModal = document.querySelector('.wheel-modal');
const startButton = wheelModal.querySelector('.start-btn');

document.addEventListener('modal.open', async (e) => {
    const modalName = e.detail.modalName;

    if (modalName !== 'wheel') return;

    let wheelItems = [];
    const currentCell = app.users.getUserCurrentCell(app.getUserId());
    const cellPresetId = currentCell.preset;
    switch (app.nextStepType) {
        case 'rollCell':
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
        case 'rollWheelPreset':
            app.wheelItems.getByPreset(cellPresetId).forEach((item) => {
                wheelItems.push({
                    id: item.id,
                    src: Helper.getFile('icon', item),
                    text: item.name,
                    type: 'wheelItem',
                });
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

async function startSpin() {
    let url = '';
    switch (app.nextStepType) {
        case 'rollCell':
            url = '/api/roll-cell';
            break;
        case 'rollWheelPreset':
            url = '/api/roll-wheel-preset';
            break;
        case 'rollItem':
            url = '/api/roll-item';
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

    const currentCell = app.users.getUserCurrentCell(app.getUserId());
    let rollInfo = app.audios.getRandomAudioFromCellByEvent(currentCell, app.nextStepType);

    if (!rollInfo) {
        rollInfo = {duration: 20};
    }

    wheel.startSpin(json.itemId, rollInfo.duration);

    app.setAudioSrc(Helper.getFile('audio', rollInfo));
    app.audioPlayer.play();

    setTimeout(() => {
        app.modal.unlockClose();
        app.audioPlayer.pause();
    }, (rollInfo.duration + 1) * 1000);
}