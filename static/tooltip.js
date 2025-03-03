import {app} from "/app.js";

document.addEventListener('DOMContentLoaded', () => {
    const tooltip = document.getElementById('tooltip');
    const body = document.body;
    let isActive = false;

    document.addEventListener('mousemove', (e) => {
        requestAnimationFrame(() => {
            if (!isActive) return;

            const bodyY = body.dataset.position ? parseInt(body.dataset.position) : 0;
            const y = e.pageY + bodyY;
            tooltip.style.transform = `translate(${e.clientX}px, ${y}px)`;
        });
    });

    document.addEventListener('mouseover', (e) => {
        const id = e.target.dataset.id;
        const type = e.target.dataset.type;

        if (!id || !type) return;

        isActive = true;

        let description;
        switch (type) {
            case 'item':
                description = app.items.getById(id).description;
                break;
            case 'cell':
                const item = app.cells.getAll().find(item => item.id === id);
                description = item.description;
                break;
        }

        if (!description) return;

        tooltip.innerHTML = description;
        tooltip.classList.add('show');
    });

    document.addEventListener('mouseout', () => {
        isActive = false;
        tooltip.classList.remove('show');
    });
});