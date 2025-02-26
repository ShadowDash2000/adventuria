const rolls = {
    'd4': [
        {x: 0, y: 0, value: 1},
        {x: 0, y: 120, value: 2},
        {x: 0, y: 240, value: 3},
        {x: 90, y: 180, value: 4},
    ],
    'd6': [
        {x: 0, y: 0, value: 1},
        {x: 90, y: 0, value: 2},
        {x: 0, y: -90, value: 3},
        {x: 0, y: 90, value: 4},
        {x: -90, y: 0, value: 5},
        {x: 180, y: 0, value: 6},
    ],
}

export default class Dice {
    constructor() {
        this.dicesTemplates = {
            'd4': document.getElementById('d4-template'),
            'd6': document.getElementById('d6-template'),
        }
        this.dices = [];
    }

    initDices(dices, container) {
        this.dices = [];
        container.innerHTML = '';

        dices.forEach((dice) => {
            if (!this.dicesTemplates[dice]) return;

            const diceTemplate = this.dicesTemplates[dice].content.cloneNode(true);
            const newDice = container.appendChild(diceTemplate.firstElementChild);

            this.dices.push({
                'type': dice,
                'element': newDice,
            });
        });
    }

    rollDice(values = null, duration = 4) {
        for (const key in this.dices) {
            const dice = this.dices[key];

            const roll = values[key] ?
                rolls[dice.type].find(r => r.value === values[key]) :
                rolls[dice.type][Math.floor(Math.random() * rolls[dice.type].length)];

            let randomX = Math.floor(Math.random() * 4 + 4) * 360;
            let randomY = Math.floor(Math.random() * 4 + 4) * 360;

            let rotate = dice.element.querySelector('.rotate');
            if (!rotate) rotate = dice.element;
            
            rotate.style.transition = `transform ${duration}s ease-in-out`;
            rotate.style.transform = `rotateX(${randomX + roll.x}deg) rotateY(${randomY + roll.y}deg)`;
        }
    }
}