import PocketBase, {BaseAuthStore} from "/pocketbase/pocketbase.es.js";
import Helper from "./helper.js";

export default class Timer {
    constructor(token, timer, resetDate) {
        const store = new BaseAuthStore();
        store.save(token);
        this.pb = new PocketBase('/', store);

        this.pb.collection('timers').subscribe('*', (e) => {
            document.dispatchEvent(new CustomEvent(`record.timers.${e.action}`, {
                detail: {
                    'record': e.record,
                },
            }));
        });

        this.timer = timer;
        this.resetDate = resetDate;
        this.interval = null;
        this.remainingTime = 0;
        this.nextTimerResetDate = null;
        this.isNegative = false;

        document.addEventListener('record.timers.update', (e) => {
            this.fetchTimeLeft();
        });

        this.fetchTimeLeft();
    }

    async fetchTimeLeft() {
        const res = await fetch('/api/timer/left', {
            method: "GET",
            headers: {
                "Authorization": this.pb.authStore.token,
            },
        });

        const data = await res.json();

        this.isNegative = data.time < 0;
        this.remainingTime = Math.abs(data.time);
        this.nextTimerResetDate = data.nextTimerResetDate;

        this.updateDisplay();

        if (data.isActive) {
            this.startInterval();
        } else {
            this.stopInterval();
        }
    }

    updateDisplay() {
        let hours = String(Math.floor(this.remainingTime / 3600)).padStart(2, '0');
        let minutes = String(Math.floor((this.remainingTime % 3600) / 60)).padStart(2, '0');
        let seconds = String(this.remainingTime % 60).padStart(2, '0');
        this.timer.textContent = (this.isNegative ? '-' : '') + `${hours}:${minutes}:${seconds}`;

        if (this.resetDate) {
            this.resetDate.innerText = Helper.formatDateLocalized(this.nextTimerResetDate);
        }
    }

    async startTimer() {
        const res = await fetch('/api/timer/start', {
            method: "POST",
            headers: {
                "Authorization": this.pb.authStore.token,
            },
        });

        if (!res.ok) return;

        this.startInterval();
    }

    startInterval() {
        if (!this.interval) {
            this.interval = setInterval(() => {
                if (this.isNegative) {
                    this.remainingTime++;
                } else {
                    this.remainingTime--;
                    this.isNegative = this.remainingTime <= 0;
                }
                this.updateDisplay();
            }, 1000);
        }
    }

    stopInterval() {
        clearInterval(this.interval);
        this.interval = null;
    }

    async stopTimer() {
        const res = await fetch('/api/timer/stop', {
            method: "POST",
            headers: {
                "Authorization": this.pb.authStore.token,
            },
        });

        if (!res.ok) return;

        clearInterval(this.interval);
        this.interval = null;
    }
}