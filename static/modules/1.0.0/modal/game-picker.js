import {app} from "../app.js";

const chooseGameButton = document.getElementById('choose-game');
const gamePicker = document.getElementById('game-picker');
const game = gamePicker.querySelector('input[name="game"]');

chooseGameButton.addEventListener('click', async (e) => {
    e.preventDefault();

    await fetch('/api/choose-game', {
        method: "POST",
        headers: {
            "Authorization": app.getUserAuthToken(),
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            "game": game.value,
        }),
    });

    app.updateInnerField();

    app.modal.close();

    game.value = '';
});