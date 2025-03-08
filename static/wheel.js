export default class Wheel {
    constructor() {
        this.items = null;
        this.spinning = false;
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

            const div = document.createElement('div');
            if (item.src) {
                div.style.background = `url(${item.src}) no-repeat`;
                div.style.backgroundSize = 'cover';
                div.style.backgroundPosition = 'center';
            }
            li.appendChild(div);

            const span = document.createElement('span');
            span.innerHTML = item.text;

            li.dataset.id = item.id;
            li.dataset.description = item.description;
            li.dataset.src = item.src;

            div.dataset.id = item.id;
            div.dataset.description = item.description;
            div.dataset.src = item.src;

            span.dataset.id = item.id;
            span.dataset.description = item.description;
            span.dataset.src = item.src;

            li.appendChild(span);
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

        this.spinning = true;
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

        setTimeout(() => {this.spinning = false;}, duration);
    }

    rotate() {
        this.wheel.classList.add('rotate');
    }

    stopRotate() {
        this.wheel.classList.remove('rotate');
    }

    isSpinning() {
        return this.spinning;
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