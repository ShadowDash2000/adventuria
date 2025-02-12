let spinning = true;
let angle = 0;
let wheel;

document.addEventListener('DOMContentLoaded', () => {
    wheel = document.getElementById('wheel');
});

function createWheel(items) {
    wheel.innerHTML = '';

    items.forEach((item, index) => {
        let li = document.createElement("li");

        li.style.rotate = `${360 / items.length * (index + 1) - 1}deg`;
        li.style.background = `hsl(${360 / items.length * (index + 1)}deg, 100%, 75%)`;
        li.style.aspectRatio = `1 / ${(2 * Math.tan(180 * (Math.PI / 180) / items.length))}`;

        if (item.src) {
            li.style.background = `url(${item.src}) no-repeat`;
            li.style.backgroundSize = 'cover';
        }

        li.innerHTML = item.text;
        wheel.appendChild(li);
    });

    autoSpin();
}

function autoSpin() {
    setInterval(() => {
        if (spinning) {
            angle += 1;
            wheel.style.transform = `rotate(${angle}deg)`;
        }
    }, 50);
}

function startSpin(items, winnerId, duration) {
    spinning = false;
    let segmentAngle = 360 / items.length;
    let winnerIndex = items.findIndex(item => item.id === winnerId);
    let finalAngle = 360 * 5 - (winnerIndex * segmentAngle);

    wheel.style.transition = `transform ${duration}s ease-in-out`;
    wheel.style.transform = `rotate(${finalAngle}deg)`;
}