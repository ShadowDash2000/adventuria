import {app} from "../app.js";
import Helper from "../helper.js";

const modal = document.getElementById('wheels-modal');
const modalContent = modal.querySelector('.wheels-modal__content');
const wheelsButton = document.getElementById('show-wheels');

wheelsButton.addEventListener('click', () => {
    app.modal.open('wheels', {
        speed: 100,
        animation: 'fadeInUp',
    });
});

document.addEventListener('modal.open.wheels', (e) => {
    putWheelsIntoModal();
}, {once: true});

function putWheelsIntoModal() {
    const wheelPresets = app.wheelItems.wheelItems;

    wheelPresets.forEach((preset, id) => {
        const presetDetail = app.wheelItems.getPresetById(id);
        const presetNode = document.createElement('div');
        presetNode.classList.add('wheels-modal__preset');

        const h2Node = document.createElement('h2');
        h2Node.innerText = presetDetail.name;
        presetNode.appendChild(h2Node);

        const itemsNode = document.createElement('div');
        itemsNode.classList.add('wheels-modal__items');

        preset.forEach(wheelItem => {
            const wheelItemNode = document.createElement('img');

            wheelItemNode.loading = 'lazy';
            wheelItemNode.src = wheelItem.icon ? Helper.getFile('icon', wheelItem) : 'img/undefined.jpg';
            wheelItemNode.dataset.description = wheelItem.name;

            itemsNode.appendChild(wheelItemNode);
        });

        presetNode.appendChild(itemsNode);

        modalContent.appendChild(presetNode);
    });
}
