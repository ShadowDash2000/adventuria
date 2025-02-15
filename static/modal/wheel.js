import Wheel from "../wheel.js";

const wheel = new Wheel();

document.addEventListener('DOMContentLoaded', () => {
    const wheelModal = document.querySelector('.wheel-modal');
    const startButton = wheelModal.querySelector('.start-btn');

    document.addEventListener('modal.open', async (e) => {
        const modalName = e.detail.modalName;

        if (modalName !== 'wheel') return;

        let items = [
            {id: 1, src: "img/kiryu.gif", text: 'TEXT 1'},
            {id: 2, src: "img/kiryu.gif", text: 'TEXT 2'},
            {id: 3, src: "img/kiryu.gif", text: 'TEXT 3'},
            {id: 4, src: "img/kiryu.gif", text: 'TEXT 4'}
        ];

        wheel.createWheel(items);
        startButton.addEventListener('click', getWheel);
    });

    document.addEventListener('modal.close', () => {
        startButton.removeEventListener('click', getWheel);
        wheel.clearWheel();
    });
});

function getWheel() {
    wheel.startSpin(2, 6);
}