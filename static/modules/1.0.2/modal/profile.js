import {app} from "../app.js";
import Helper from "../helper.js";
import Timer from "../timer.js";

const profileModal = document.getElementById('profile-modal');
const profileName = profileModal.querySelector('h2');
const profileAvatar = profileModal.querySelector('.profile-modal img');
const profileDescription = profileModal.querySelector('.profile-modal .profile-modal__description');

const profileCellImg = profileModal.querySelector('.current-cell img');
const profileCellName = profileModal.querySelector('.current-cell .profile-modal__name');
const profileCellDescription = profileModal.querySelector('.current-cell .profile-modal__description');

const profileTimer = profileModal.querySelector('.profile-modal__timer div');

const drops = profileModal.querySelector('.profile-modal__stats .profile-modal__drops');
const rerolls = profileModal.querySelector('.profile-modal__stats .profile-modal__rerolls');
const finished = profileModal.querySelector('.profile-modal__stats .profile-modal__finished');
const wasInJail = profileModal.querySelector('.profile-modal__stats .profile-modal__was-in-jail');
const itemsUsed = profileModal.querySelector('.profile-modal__stats .profile-modal__items-used');
const diceRolls = profileModal.querySelector('.profile-modal__stats .profile-modal__dice-rolls');
const maxDiceRoll = profileModal.querySelector('.profile-modal__stats .profile-modal__max-dice-roll');

const actionContainer = profileModal.querySelector('.actions .container');
const actionsSentinel = profileModal.querySelector('.actions .sentinel');
let isLoading = false;
let page = 1;
let totalPages = 1;
const limit = 10;
let observer;

document.addEventListener('profile.open', (e) => {
    actionContainer.innerHTML = '';
    const userId = e.detail.userId;

    putUserInfoToProfile(userId);

    observer = new IntersectionObserver(async (entries, observer) => {
        const entry = entries[0];

        if (entry.isIntersecting && !isLoading) {
            if (page > totalPages) {
                observer.unobserve(actionsSentinel);
                return;
            }

            isLoading = true;
            const actions = await fetchUserActions(userId, page, limit);
            totalPages = actions.totalPages;
            page++;
            isLoading = false;

            for (const action of actions.items) {
                const actionNode = app.actions.createActionNode(action);
                actionContainer.appendChild(actionNode);
            }
        }
    });

    observer.observe(actionsSentinel);

    app.modal.open('profile', {
        speed: 100,
        animation: 'fadeInUp',
    });
});

document.addEventListener('modal.close', (e) => {
    if (e.detail.modalName !== 'profile') return;
    page = 1;
    totalPages = 1;
    observer.unobserve(actionsSentinel);
});

async function fetchUserActions(userId, page, limit) {
    return await app.pb.collection('actions').getList(page, limit, {
        filter: `user.id = "${userId}" && '["roll", "reroll", "drop", "chooseResult", "chooseGame", "rollCell", "rollWheelPreset"]' ~ type`,
        sort: '-created',
    });
}

function putUserInfoToProfile(userId) {
    const user = app.users.getById(userId);
    const currentCell = app.users.getUserCurrentCell(userId);

    profileName.innerText = `ПРОФИЛЬ ${user.name}`;
    profileAvatar.src = Helper.getFile('avatar', user);
    profileAvatar.style.borderColor = user.color;
    profileDescription.innerHTML = user.description;

    Timer.fetchUserTimeLeft(userId).then((time) => {
        profileTimer.innerText = time;
    });

    profileCellImg.src = Helper.getFile('icon', currentCell, {'thumb': '250x0'});
    profileCellName.innerText = currentCell.name;
    profileCellDescription.innerHTML = currentCell.description;

    drops.innerText = 'Дропнуто: ' + (user.stats?.drops || 0);
    rerolls.innerText = 'Рерольнуто: ' + (user.stats?.rerolls || 0);
    finished.innerText = 'Пройдено: ' + (user.stats?.finished || 0);
    wasInJail.innerText = 'Был в тюрьме: ' + (user.stats?.wasInJail || 0);
    itemsUsed.innerText = 'Использовано предметов: ' + (user.stats?.itemsUsed || 0);
    diceRolls.innerText = 'Сколько раз бросал кубики: ' + (user.stats?.diceRolls || 0);
    maxDiceRoll.innerText = 'Максимальный ролл кубиков: ' + (user.stats?.maxDiceRoll || 0);
}