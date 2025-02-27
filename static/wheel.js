export default class Wheel {
    constructor() {
        this.items = null;
    }

    createWheel(items) {
        this.wheel = document.getElementById('wheel');
        this.wheel.innerHTML = '';
        this.items = items;

        items.forEach((item, index) => {
            let li = document.createElement("li");

            li.style.rotate = `${360 / items.length * index}deg`;
            li.style.background = `hsl(${360 / items.length * (index + 1)}deg, 100%, 75%)`;
            li.style.aspectRatio = `1 / ${(2 * Math.tan(180 * (Math.PI / 180) / items.length))}`;

            li.dataset.id = item.id;

            if (item.src) {
                li.style.background = `url(${item.src}) no-repeat`;
                li.style.backgroundSize = 'cover';
            }

            li.innerHTML = item.text;
            this.wheel.appendChild(li);
        });
    }

    clearWheel() {
        if (this.wheel) {
            this.wheel.innerHTML = '';
            this.wheel.setAttribute('style', '');
        }
    }

    startSpin(winnerId, duration) {
        if (!this.items) return;

        const winnerIndex = this.items.findIndex(item => item.id === winnerId);
        const segmentAngle = 360 / this.items.length;
        const finalAngle = 360 * 10 - (segmentAngle * winnerIndex) + 90;

        this.wheel.classList.remove('rotate');
        setTimeout(() => {
            this.wheel.style.transition = `transform ${duration}s ease-in-out`;
            this.wheel.style.transform = `rotate(${finalAngle}deg)`;
        }, 100);
    }

    rotate() {
        this.wheel.classList.add('rotate');
    }

    getCurrentAngle() {
        const style = window.getComputedStyle(this.wheel);
        const transform = style.getPropertyValue('transform');

        if (transform === 'none') {
            return 0;
        }

        const values = transform.match(/matrix\(([^)]+)\)/);
        if (!values) {
            return 0;
        }

        const [a, b] = values[1].split(",").map(parseFloat);
        return Math.round(Math.atan2(b, a) * (180 / Math.PI));

    }

    getCurrentWinner() {
        const currentAngle = this.getCurrentAngle() - 90;
        let segmentAngle = 360 / this.items.length;
        let normalizedAngle = ((currentAngle % 360) + 360) % 360;
        normalizedAngle = (360 - normalizedAngle) % 360;
        let winnerIndex = Math.round(normalizedAngle / segmentAngle);

        return this.items[winnerIndex % this.items.length];
    }
}