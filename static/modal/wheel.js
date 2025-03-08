import {app} from "../app.js";
import Wheel from "../wheel.js";
import Helper from "../helper.js";

const wheel = new Wheel();

const wheelModal = document.querySelector('.wheel-modal');
const startButton = wheelModal.querySelector('.start-btn');

document.addEventListener('modal.open', async (e) => {
    const modalName = e.detail.modalName;

    if (modalName !== 'wheel') return;

    const wheelItems = getWheelItems();

    if (wheelItems) {
        setTimeout(() => {
            wheel.createWheel(wheelItems);
            wheel.rotate();
            startButton.addEventListener('click', requestWheelSpinResult);
        }, 0);
    }
});

document.addEventListener('modal.close', () => {
    startButton.removeEventListener('click', startSpin);
    wheel.clearWheel();
});

function getWheelItems() {
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
                        description: cell.description,
                        type: 'cell',
                    });
                }
            }
            break;
        case 'rollWheelPreset':
            app.wheelItems.getByPreset(cellPresetId).forEach((item) => {
                wheelItems.push({
                    id: item.id,
                    src: Helper.getFile('icon', item, {'thumb': '600x0'}),
                    text: item.name,
                    description: item.name,
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
                        description: item.description,
                        type: 'item',
                    });
                }
            })
            break;
    }

    return wheelItems;
}

function requestWheelSpinResult() {
    app.requestsWorker.postMessage({
        method: 'fetchWheelSpinResult',
        payload: app.nextStepType,
    });
}

document.addEventListener('worker.fetchWheelSpinResult', async (e) => {
    await startSpin(e.detail.result);
});

async function startSpin(itemId) {
    if (wheel.isSpinning()) return;

    app.modal.lockClose();

    const currentCell = app.users.getUserCurrentCell(app.getUserId());
    let rollInfo = app.audios.getRandomAudioFromCellByEvent(currentCell, app.nextStepType);

    if (!rollInfo) {
        rollInfo = {duration: 20};
    }

    wheel.startSpin(itemId, rollInfo.duration);

    app.setAudioSrc(Helper.getFile('audio', rollInfo));
    app.audioPlayer.play();

    setTimeout(() => {
        app.modal.unlockClose();
        app.audioPlayer.pause();
    }, (rollInfo.duration + 1) * 1000);
}