import {app} from "../app.js";
import Helper from "../helper.js";

const profileModal = document.getElementById('profile-modal');
const profileName = profileModal.querySelector('h2');
const profileAvatar = profileModal.querySelector('.profile-modal img');
const profileDescription = profileModal.querySelector('.profile-modal .profile-modal__description');

const profileCellImg = profileModal.querySelector('.current-cell img');
const profileCellName = profileModal.querySelector('.current-cell .profile-modal__name');
const profileCellDescription = profileModal.querySelector('.current-cell .profile-modal__description');

document.addEventListener('profile.open', (e) => {
    putUserInfoToProfile(e.detail.userId);
    app.modal.open('profile', {
        speed: 100,
        animation: 'fadeInUp',
    });
});

function putUserInfoToProfile(userId) {
    const user = app.users.getById(userId);
    const currentCell = app.users.getUserCurrentCell(userId);

    profileName.innerText = `ПРОФИЛЬ ${user.name}`;
    profileAvatar.src = Helper.getFile('avatar', user);
    profileAvatar.style.borderColor = user.color;
    profileDescription.innerHTML = user.description;

    profileCellImg.src = Helper.getFile('icon', currentCell, {'thumb': '250x0'});
    profileCellName.innerText = currentCell.name;
    profileCellDescription.innerHTML = currentCell.description;
}