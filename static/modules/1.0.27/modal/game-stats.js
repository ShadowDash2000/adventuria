import {app} from "../app.js";
import Helper from "../helper.js";

const gameStatsModal = document.getElementById('game-stats-modal');
const container = gameStatsModal.querySelector('.graph-modal__content .container');

document.addEventListener('modal.open.game-stats', putStatsIntoModal, {once: true});

function putStatsIntoModal() {
    const stats = calculateStats();
    const textShadow = ['text-shadow', 'text-soft-shadow'];

    for (const key in stats) {
        const stat = stats[key];

        for (const key in stat.values) {
            const statValue = stat.values[key];

            const statEl = document.createElement('div');
            statEl.classList.add('game-stats-modal__stat');

            const title = document.createElement('h3');
            title.classList.add(...textShadow);
            const description = document.createElement('p');
            description.classList.add('stat__description');
            description.classList.add(...textShadow);

            title.innerText = statValue.title;
            description.innerText = statValue.description;

            const statDetailEl = document.createElement('div');
            statDetailEl.classList.add('stat__detail');

            const avatar = document.createElement('img');
            avatar.style.borderColor = statValue.user.color;
            const userName = document.createElement('span');
            userName.classList.add(...textShadow);
            const value = document.createElement('span');
            value.classList.add(...textShadow);

            avatar.src = Helper.getFile('avatar', statValue.user);
            userName.innerText = statValue.user.name;
            value.innerText = statValue.value;

            statDetailEl.append(avatar, userName, value);
            statEl.append(title, description, statDetailEl);
            container.appendChild(statEl);
        }
    }
}

function calculateStats() {
    let stats = {
        winner: {
            values: {
                max: {
                    title: 'Победитель!',
                    description: 'Набрал больше всего очков.',
                    user: null,
                    value: null,
                },
            },
        },
        cellsPassed: {
            values: {
                max: {
                    title: 'Бегунок',
                    description: 'Прошел больше всего клеток.',
                    user: null,
                    value: null,
                },
                min: {
                    title: 'Улиткус',
                    description: 'Прошел меньше всего клеток.',
                    user: null,
                    value: null,
                },
            },
        },
        finished: {
            values: {
                max: {
                    title: 'Гроза игр',
                    description: 'Завершено больше всего игр / фильмов.',
                    user: null,
                    value: null,
                },
                min: {
                    title: 'На чилле',
                    description: 'Завершено меньше всего игр / фильмов.',
                    user: null,
                    value: null,
                },
            },
        },
        drops: {
            values: {
                max: {
                    title: 'Убегунчик',
                    description: 'Больше всего дропов.',
                    user: null,
                    value: null,
                },
                min: {
                    title: 'Терпила',
                    description: 'Меньше всего дропов.',
                    user: null,
                    value: null,
                },
            },
        },
        rerolls: {
            values: {
                max: {
                    title: 'Любитель мошны',
                    description: 'Реролил больше всех.',
                    user: null,
                    value: null,
                },
                min: {
                    title: 'Честный',
                    description: 'Реролил меньше всех.',
                    user: null,
                    value: null,
                },
            },
        },
        wasInJail: {
            values: {
                max: {
                    title: 'Матёрый зек',
                    description: 'Чаще всех был в тюрьме.',
                    user: null,
                    value: null,
                },
                min: {
                    title: 'Послушный гражданин',
                    description: 'Реже всех был в тюрьме.',
                    user: null,
                    value: null,
                },
            },
        },
        diceRolls: {
            values: {
                max: {
                    title: 'Азартный',
                    description: 'Больше всех кидал кубики.',
                    user: null,
                    value: null,
                },
                min: {
                    title: 'А что, надо было играть?',
                    description: 'Меньше всех кидал кубики.',
                    user: null,
                    value: null,
                },
            },
        },
        maxDiceRoll: {
            values: {
                max: {
                    title: 'Шарлатан',
                    description: 'Самый большой бросок кубиков.',
                    user: null,
                    value: null,
                },
                min: {
                    title: 'Микропенис',
                    description: 'Самый маленький бросок кубиков.',
                    user: null,
                    value: null,
                },
            },
        },
        itemsUsed: {
            values: {
                max: {
                    title: 'Главная гнида',
                    description: 'Использовал больше всего предметов.',
                    user: null,
                    value: null,
                },
                min: {
                    title: 'На харде',
                    description: 'Использовал меньше всего предметов.',
                    user: null,
                    value: null,
                },
            },
        },
        wheelRolled: {
            values: {
                max: {
                    title: 'Крутилкин',
                    description: 'Прокрутил больше всего колёс.',
                    user: null,
                    value: null,
                },
                min: {
                    title: 'Боится колёс',
                    description: 'Прокрутил меньше всего колёс.',
                    user: null,
                    value: null,
                },
            },
        },
    }

    const users = app.users.users;
    users.forEach(user => {
        if (stats.winner.values.max.user === null ||
            user.points > stats.winner.values.max.value) {
            stats.winner.values.max.user = user;
            stats.winner.values.max.value = user.points;
        }

        if (stats.cellsPassed.values.max.user === null ||
            user.cellsPassed > stats.cellsPassed.values.max.value) {
            stats.cellsPassed.values.max.user = user;
            stats.cellsPassed.values.max.value = user.cellsPassed;
        }

        if (stats.cellsPassed.values.min.user === null ||
            user.cellsPassed < stats.cellsPassed.values.min.value) {
            stats.cellsPassed.values.min.user = user;
            stats.cellsPassed.values.min.value = user.cellsPassed;
        }

        for (const statKey in user.stats) {
            const stat = user.stats[statKey];

            if (stats[statKey]?.values.max.user === null ||
                stat > stats[statKey]?.values.max.value) {
                stats[statKey].values.max.user = user;
                stats[statKey].values.max.value = stat;
            }

            if (stats[statKey]?.values.min.user === null ||
                stat < stats[statKey]?.values.min.value) {
                stats[statKey].values.min.user = user;
                stats[statKey].values.min.value = stat;
            }
        }
    });

    return stats;
}