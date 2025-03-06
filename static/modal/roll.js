import {app} from "../app.js";
import Dice from "../dice.js";
import Helper from "../helper.js";

const dice = new Dice();

let rollModal, scene, rollButton, rollResult, cell;

document.addEventListener('DOMContentLoaded', () => {
    rollModal = document.querySelector('.graph-modal__container.dice');
    scene = rollModal.querySelector('.scene');
    rollButton = document.getElementById('roll');
    rollResult = rollModal.querySelector('.roll-result');
    cell = rollModal.querySelector('.cell');

    document.addEventListener('modal.open', async (e) => {
        const modalName = e.detail.modalName;

        if (modalName !== 'dice') return;

        scene.innerHTML = '';
        rollResult.classList.add('hidden');
        cell.classList.add('hidden');
        rollButton.classList.remove('hidden');

        const res = await fetch('/api/get-roll-effects', {
            method: "GET",
            headers: {
                "Authorization": app.getUserAuthToken(),
            },
        });

        if (!res.ok) return;

        const json = await res.json();

        let dices = [];
        if (json.dices) {
            for (const dice of json.dices) {
                dices.push(dice.type);
            }
        } else  dices = ['d6', 'd6'];

        dice.initDices(dices, scene);

        rollButton.addEventListener('click', roll);
    });

    document.addEventListener('modal.close', () => {
        rollButton.removeEventListener('click', roll);
    });
});

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

        rollResult.querySelector('.roll-result__number').innerHTML = json.roll;

        const newCell = app.cells.getById(json.cellId);

        cell.querySelector('img').src = Helper.getFile('icon', newCell, {'thumb': '250x0'});
        cell.querySelector('.cell-info__name').innerHTML = newCell.name;
        cell.querySelector('.cell-info__description').innerHTML = newCell.description;

        rollResult.classList.remove('hidden');
        cell.classList.remove('hidden');

        await app.updateInnerField();

    }, rollInfo.duration * 1000);

    setTimeout(() => {
        app.audioPlayer.pause();
    }, (rollInfo.duration + 1) * 1000);
}