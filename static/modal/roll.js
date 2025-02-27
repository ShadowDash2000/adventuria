import {app} from "../app.js";
import Dice from "../dice.js";

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

        const res = await fetch('/api/get-roll-effects', {
            method: "GET",
            headers: {
                "Authorization": app.auth.token,
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
            "Authorization": app.auth.token,
        },
    });

    if (!res.ok) {
        return;
    }

    const json = await res.json();

    const audioItemsKeys = Array.from(app.audio[app.nextStepType].keys());
    const randomKey = audioItemsKeys[Math.floor(Math.random() * audioItemsKeys.length)];
    const rollInfo = app.audio[app.nextStepType].get(randomKey);

    let duration = rollInfo.duration;
    const durations = [];
    for (let i = 0; i < json.diceRolls.length; i++) {
        duration -= 2;
        durations.unshift(duration);
    }

    dice.rollDice(json.diceRolls, durations);

    const audio = new Audio(app.getFile('audio', rollInfo));
    audio.volume = app.volume / 100;
    audio.play();

    setTimeout(async () => {
        rollResult.querySelector('.roll-result__number').innerHTML = json.roll;

        cell.querySelector('img').src = json.cell.icon;
        cell.querySelector('.cell-info__name').innerHTML = json.cell.name;
        cell.querySelector('.cell-info__description').innerHTML = json.cell.description;

        rollResult.classList.remove('hidden');
        cell.classList.remove('hidden');

        rollButton.classList.add('hidden');
        rollModal.querySelector('.choose-game').classList.remove('hidden');

        await app.updateInnerField();

    }, rollInfo.duration * 1000);

    setTimeout(() => {
        audio.pause();
    }, (rollInfo.duration + 1) * 1000);
}