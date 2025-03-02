import {app} from "/app.js";

document.addEventListener('DOMContentLoaded', () => {
    const tooltip = document.getElementById('tooltip');
    const body = document.body;

    document.addEventListener('mousemove', (e) => {
        requestAnimationFrame(() => {
            const bodyY = body.dataset.position ? parseInt(body.dataset.position) : 0;
            const y = e.pageY + bodyY;
            tooltip.style.transform = `translate(${e.clientX}px, ${y}px)`;
        });
    });
    document.addEventListener('mouseover', (e) => {
        const id = e.target.dataset.id;
        const type = e.target.dataset.type;

        if (!id || !type) return;

        let description;
        switch (type) {
            case 'item':
                description = app.items.get(id).description;
                break;
            case 'cell':
                const item = app.cellsList.find(item => item.id === id);
                description = item.description;
                break;
        }

        if (!description) return;

        tooltip.innerHTML = description;
        tooltip.classList.add('show');
    });
    document.addEventListener('mouseout', () => {
        tooltip.classList.remove('show');
    });
});