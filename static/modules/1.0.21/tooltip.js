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
    const src = e.target.dataset.src;
    const description = e.target.dataset.description;

    if (!description) return;

    isActive = true;

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