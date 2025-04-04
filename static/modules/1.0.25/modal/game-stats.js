import {app} from "../app.js";
import Helper from "../helper.js";

const gameStatsModal = document.getElementById('game-stats-modal');

document.addEventListener('modal.open.game-stats', putStatsIntoModal, {once: true});

function putStatsIntoModal() {
    const stats = calculateStats();

    for (const key in stats) {
        const stat = stats[key];

        const statEl = document.createElement('div');
        statEl.classList.add('game-stats-modal__stat');

        const title = document.createElement('h3');
        const description = document.createElement('p');
        description.classList.add('stat__description');

        title.innerText = stat.title;
        description.innerText = stat.description;

        const statDetailEl = document.createElement('div');
        statDetailEl.classList.add('stat__detail');

        const avatar = document.createElement('img');
        const userName = document.createElement('span');
        const value = document.createElement('span');

        avatar.src = Helper.getFile('avatar', stat.user);
        userName.innerText = stat.user.name;
        value.innerText = stat.value;

        statDetailEl.append(avatar, userName, value);
        statEl.append(title, description, statDetailEl);
        gameStatsModal.appendChild(statEl);
    }
}

function calculateStats() {
    let stats = {
        winner: {
            title: 'Победитель!',
            description: '',
            user: null,
            value: null,
        },
        cellsPassed: {
            title: 'Бегунок',
            description: '',
            user: null,
            value: null,
        },
        finished: {
            title: 'Гроза игр',
            description: '',
            user: null,
            value: null,
        },
        drops: {
            title: 'Убегунчик',
            description: '',
            user: null,
            value: null,
        },
        rerolls: {
            title: 'Любитель мошны',
            description: '',
            user: null,
            value: null,
        },
        wasInJail: {
            title: 'Матёрый зек',
            description: '',
            user: null,
            value: null,
        },
        diceRolls: {
            title: 'Азартный',
            description: '',
            user: null,
            value: null,
        },
        maxDiceRoll: {
            title: 'Шарлатан',
            description: '',
            user: null,
            value: null,
        },
        itemsUsed: {
            title: '.....',
            description: '',
            user: null,
            value: null,
        },
        wheelRolled: {
            title: 'Крутилкин',
            description: '',
            user: null,
            value: null,
        },
    }

    const users = app.users.users;
    users.forEach(user => {
        if (stats.winner.user === null ||
            user.points > stats.winner.user?.points) {
            stats.winner.user = user;
            stats.winner.value = user.points;
        }

        if (stats.cellsPassed.user === null ||
            user.cellsPassed > stats.cellsPassed.user?.cellsPassed) {
            stats.cellsPassed.user = user;
            stats.cellsPassed.value = user.cellsPassed;
        }

        for (const statKey in user.stats) {
            const stat = user.stats[statKey];

            if (stats[statKey]?.user === null ||
                stat > stats[statKey]?.user[statKey]) {
                stats[statKey].user = user;
                stats[statKey].value = stat;
            }
        }
    });

    return stats;
}