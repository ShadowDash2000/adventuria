import {app} from "../app.js";
import Dice from "../dice.js";
import Helper from "../helper.js";

const dice = new Dice();

const graphModal = document.querySelector('.graph-modal');
const rollModal = document.querySelector('.graph-modal__container.dice');
const scene = rollModal.querySelector('.scene');
const rollButton = document.getElementById('roll');
const rollResult = rollModal.querySelector('.roll-result');
const cell = rollModal.querySelector('.cell');

document.addEventListener('modal.open', async (e) => {
    const modalName = e.detail.modalName;

    if (modalName !== 'dice') return;

    graphModal.classList.add('bg-black');

    scene.innerHTML = '';
    rollResult.classList.add('hidden');
    cell.classList.add('hidden');
    rollButton.classList.remove('hidden');

    setTimeout(async () => {
        dice.initDices(await fetchDices(), scene);
        rollButton.addEventListener('click', roll);
    }, 0);
});

document.addEventListener('modal.close', () => {
    graphModal.classList.remove('bg-black');
    rollButton.removeEventListener('click', roll);
});

async function fetchDices() {
    const res = await fetch('/api/get-roll-effects', {
        method: "GET",
        headers: {
            "Authorization": app.getUserAuthToken(),
        },
    });

    if (!res.ok) return;

    const json = await res.json();

    let dices = [];
    if (json.changeDices) {
        for (const dice of json.changeDices) {
            dices.push(dice.type);
        }
    } else dices = ['d6', 'd6'];

    return dices;
}

async function roll() {
    const res = await fetch('/api/roll', {
        method: "POST",
        headers: {
            "Authorization": app.getUserAuthToken(),
        },
    });

    if (!res.ok) return;

    const json = await res.json();

    app.modal.lockClose();
    rollButton.classList.add('hidden');

    const currentCell = app.users.getUserCurrentCell(app.getUserId());
    let rollInfo = app.audios.getRandomAudioFromCellByEvent(currentCell, 'roll');

    if (!rollInfo) {
        rollInfo = {duration: 20};
    }

    let duration = rollInfo.duration;
    const durations = [];
    for (let i = 0; i < json.diceRolls.length; i++) {
        duration -= 2;
        durations.unshift(duration);
    }

    dice.rollDice(json.diceRolls, durations);

    app.setAudioSrc(Helper.getFile('audio', rollInfo));
    app.audioPlayer.play();

    setTimeout(async () => {
        app.modal.unlockClose();
        showRollResult(json);
        app.updateInnerField();
    }, rollInfo.duration * 1000);

    setTimeout(() => {
        app.audioPlayer.pause();
    }, (rollInfo.duration + 1) * 1000);
}

function showRollResult(result) {
    rollResult.querySelector('.roll-result__number').innerHTML = result.roll;

    const newCell = app.cells.getById(result.cellId);

    cell.querySelector('img').src = Helper.getFile('icon', newCell, {'thumb': '250x0'});
    cell.querySelector('.cell-info__name').innerHTML = newCell.name;
    cell.querySelector('.cell-info__description').innerHTML = newCell.description;

    rollResult.classList.remove('hidden');
    cell.classList.remove('hidden');
}