export default class Timers {
    collectionName = 'timers';

    constructor(pb) {
        this.pb = pb;
        this.timers = new Map();

        document.addEventListener('record.timers.update', (e) => {
            this.timers.set(e.detail.record.user, e.detail.record);

            document.dispatchEvent(new Event('timer.update'));
        });
    }

    async fetch() {
        const timers = await this.pb.collection(this.collectionName).getFullList();

        for (const timer of timers) {
            this.timers.set(timer.user, timer);
        }
    }

    getByUserId(userId) {
        return this.timers.get(userId);
    }
}