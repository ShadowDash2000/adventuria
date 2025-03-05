export default class Wheel {
    constructor() {
        this.items = null;
    }

    createWheel(items) {
        this.wheel = document.getElementById('wheel');
        this.title = document.getElementById('wheel-title');
        this.wheel.innerHTML = '';
        this.items = items;

        items.forEach((item, index) => {
            const li = document.createElement('li');

            li.style.rotate = `${360 / items.length * index}deg`;
            li.style.background = `hsl(${360 / items.length * (index + 1)}deg, 100%, 75%)`;
            li.style.aspectRatio = `1 / ${(2 * Math.tan(180 * (Math.PI / 180) / items.length))}`;

            li.dataset.id = item.id;
            li.dataset.type = item.type;

            if (item.src) {
                li.style.background = `url(${item.src}) no-repeat`;
                li.style.backgroundSize = 'cover';
            }

            li.innerHTML = item.text;
            this.wheel.appendChild(li);
            this.items[index].element = li;
        });
    }

    clearWheel() {
        if (this.wheel) {
            this.wheel.innerHTML = '';
            this.wheel.setAttribute('style', '');
        }
        if (this.interval) {
            clearInterval(this.interval);
        }
    }

    startSpin(winnerId, duration) {
        if (!this.items) return;

        this.wheel.style.transition = '';
        this.wheel.style.transform = '';

        const winnerIndex = this.items.findIndex(item => item.id === winnerId);
        const segmentAngle = 360 / this.items.length;

        const halfOfSegmentAngle = segmentAngle / 2;
        const maxSegment = halfOfSegmentAngle - 5;
        const randomAddAngle = Math.floor(Math.random() * (maxSegment + 1)) - maxSegment;

        const finalAngle = 360 * duration - (segmentAngle * winnerIndex) + 90 + randomAddAngle;

        this.stopRotate();
        setTimeout(() => {
            this.wheel.style.transition = `transform ${duration}s cubic-bezier(0.4, 0.2, 0.3, 1)`;
            this.wheel.style.transform = `rotate(${finalAngle}deg)`;

            this.startTitleInterval(duration);
        }, 100);
    }

    rotate() {
        this.wheel.classList.add('rotate');
    }

    stopRotate() {
        this.wheel.classList.remove('rotate');
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

    startTitleInterval(duration) {
        this.interval = setInterval(() => {
            const currentWinner = this.getCurrentWinner();

            this.setTitle(currentWinner.text);
            this.updateHighlight(currentWinner);
        }, 100);

        setTimeout(() => {
            clearInterval(this.interval)
        }, duration * 1000);
    }

    setTitle(title) {
        this.title.innerText = title;
    }

    updateHighlight(itemToHighlight) {
        this.items.forEach(item => {
            if (item.id === itemToHighlight.id) {
                item.element.classList.remove('unactive');
            } else {
                item.element.classList.add('unactive');
            }
        });
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