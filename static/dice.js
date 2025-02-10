function rollDice(forcedValue = null, duration = 4) {
    const dice = document.getElementById('dice');
    const rolls = [
        {x: 0, y: 0, value: 1},
        {x: 90, y: 0, value: 2},
        {x: 0, y: -90, value: 3},
        {x: 0, y: 90, value: 4},
        {x: -90, y: 0, value: 5},
        {x: 180, y: 0, value: 6}
    ];

    let roll = forcedValue ? rolls.find(r => r.value === forcedValue) : rolls[Math.floor(Math.random() * 6)];

    let randomX = Math.floor(Math.random() * 4 + 4) * 360;
    let randomY = Math.floor(Math.random() * 4 + 4) * 360;

    dice.style.transition = `transform ${duration}s ease-in-out`;
    dice.style.transform = `rotateX(${randomX + roll.x}deg) rotateY(${randomY + roll.y}deg)`;
}