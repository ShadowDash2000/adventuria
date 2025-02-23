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
        for (const dice of json.dices) {
            dices.push(dice.type);
        }

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

    const duration = 4;

    dice.rollDice(json.diceRolls, duration);

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

    }, duration * 1000);
}