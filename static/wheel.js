export default class Wheel {
    constructor() {
        this.spinning = true;
        this.autoSpinInterval = null;
        this.items = null;
    }

    createWheel(items) {
        this.angle = 0;
        this.wheel = document.getElementById('wheel');
        this.wheel.innerHTML = '';
        this.items = items;

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
            this.wheel.appendChild(li);
        });

        if (this.autoSpinInterval) clearInterval(this.autoSpinInterval);

        this.autoSpin();
    }

    clearWheel() {
        if (this.autoSpinInterval) clearInterval(this.autoSpinInterval);
        if (this.wheel) this.wheel.innerHTML = '';
    }

    autoSpin() {
        this.autoSpinInterval = setInterval(() => {
            if (this.spinning) {
                this.angle += 1;
                this.wheel.style.transform = `rotate(${this.angle}deg)`;
            }
        }, 50);
    }

    startSpin(winnerId, duration) {
        if (!this.items) return;

        this.spinning = false;
        let segmentAngle = 360 / this.items.length;
        let winnerIndex = this.items.findIndex(item => item.id === winnerId);
        let finalAngle = 360 * 10 - (winnerIndex * segmentAngle);

        this.wheel.style.transition = `transform ${duration}s ease-in-out`;
        this.wheel.style.transform = `rotate(${finalAngle}deg)`;
    }
}