import {app} from "/app.js";
import Helper from "./helper.js";

const tooltip = document.getElementById('tooltip');
const tooltipImg = tooltip.querySelector('img');
const tooltipText = tooltip.querySelector('span');
const body = document.body;
let isActive = false;

document.addEventListener('mousemove', (e) => {
    if (!isActive) return;

    requestAnimationFrame(() => {
        const bodyY = body.dataset.position ? parseInt(body.dataset.position) : 0;
        const y = e.pageY + bodyY;
        let x = e.clientX;

        if (x > window.innerWidth / 2) {
            x -= tooltip.offsetWidth;
        }

        tooltip.style.transform = `translate(${x}px, ${y}px)`;
    });
});

document.addEventListener('mouseover', (e) => {
    const id = e.target.dataset.id;
    const type = e.target.dataset.type;

    if (!id || !type) return;

    isActive = true;

    // TODO This need to be changed. Instead of checking conditions here, we need to get data attrs with text and image.
    let description, src, item;
    switch (type) {
        case 'item':
            description = app.items.getById(id).description;
            break;
        case 'cell':
            item = app.cells.getAll().find(item => item.id === id);
            description = item.description;
            break;
        case 'wheelItem':
            item = app.wheelItems.getById(id);
            src = Helper.getFile('icon', item, {'thumb': '250x0'})
            description = item.name;
            break;
    }

    if (!description) return;

    tooltipText.innerHTML = description;
    if (src) {
        tooltipImg.src = src;
        tooltipImg.classList.remove('hidden');
    } else {
        tooltipImg.classList.add('hidden');
    }
    tooltip.classList.add('show');
});

document.addEventListener('mouseout', () => {
    if (!isActive) return;
    isActive = false;
    tooltip.classList.remove('show');
});