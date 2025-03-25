import {app} from "../app.js";
import Helper from "../helper.js";

const modal = document.getElementById('items-modal');
const modalContent = modal.querySelector('.items-modal__content');

document.addEventListener('modal.open.items', (e) => {
    putItemsIntoModal();
}, {once: true});

function putItemsIntoModal() {
    const items = app.items.getAll();

    items.forEach(item => {
        const itemNode = document.createElement('div');
        itemNode.classList.add('items-modal__item');

        const img = document.createElement('img');
        img.loading = 'lazy';
        img.src = Helper.getFile('icon', item);
        img.dataset.description = item.description;
        itemNode.appendChild(img);

        const span = document.createElement('span');
        span.innerHTML = item.name;
        itemNode.appendChild(span);

        modalContent.appendChild(itemNode);
    });
}
