import GraphModal from "/graph-modal/graph-modal.js";
import PocketBase from "/pocketbase/pocketbase.es.js";
import Submit from "./modal/submit.js";
import Timer from "./timer.js";
import Helper from "./helper.js";
import Users from "./internal/users.js";
import Cells from "./internal/cells.js";
import Items from "./internal/items.js";
import Inventories from "./internal/inventories.js";
import WheelItems from "./internal/wheel-items.js";
import Audios from "./internal/audios.js";
import Actions from "./internal/actions.js";
import Settings from "./internal/settings.js";

class App {
    constructor() {
        this.pb = new PocketBase('/');
        this.nextStepType = '';
        this.volume = localStorage.getItem('volume') ? localStorage.getItem('volume') : 30;
        this.isSlowPc = localStorage.getItem('slow_pc') === 'true';
        this.audioPlayer = new Audio();
        this.audioPlayer.volume = this.volume / 100;

        this.requestsWorker = new Worker('./internal/requests_worker.js');
        this.requestsWorker.postMessage({
            method: 'setToken',
            payload: this.getUserAuthToken(),
        });

        this.requestsWorker.onmessage = function (e) {
            document.dispatchEvent(new CustomEvent(`worker.${e.data.method}`, {
                detail: {
                    result: e.data.result,
                }
            }));
        }

        this.modal = new GraphModal({
            isOpen: (modal) => {
                const activeModal = modal.modalContainer;
                const modalName = activeModal.dataset.graphTarget;

                document.dispatchEvent(new CustomEvent("modal.open", {
                    detail: {
                        activeModal,
                        modalName,
                    },
                }));
            },
            isClose: (modal) => {
                const activeModal = modal.modalContainer;
                const modalName = activeModal.dataset.graphTarget;

                document.dispatchEvent(new CustomEvent("modal.close", {
                    detail: {
                        activeModal,
                        modalName,
                    },
                }));
            },
        });
        this.submit = new Submit(this.modal);

        document.addEventListener('record.actions.create', async (e) => {
            if (e.detail.record.user !== this.getUserId()) return;
            setTimeout(async () => {
                await this.updateInnerField();
            }, 1000);
        });
        document.addEventListener('record.actions.delete', async (e) => {
            if (e.detail.record.user !== this.getUserId()) return;
            setTimeout(async () => {
                await this.updateInnerField();
            }, 1000);
        });

        document.addEventListener('worker.getNextStepType', async (e) => {
            this.showActionButton(e.detail.result);
        });
    }

    async init() {
        const volumeSlider = document.getElementById('volume-slider');
        volumeSlider.value = this.volume;
        volumeSlider.addEventListener('change', () => {
            this.setVolume(volumeSlider.value);
        });

        document.addEventListener('modal.open', () => {
            volumeSlider.parentElement.classList.add('fixed');
        });
        document.addEventListener('modal.close', () => {
            volumeSlider.parentElement.classList.remove('fixed');
        });

        const gradientBg = document.querySelector('.gradient-bg');
        if (!this.isSlowPc) gradientBg.classList.remove('hidden');

        const slowPcCheckbox = document.getElementById('slow-pc-checkbox');
        slowPcCheckbox.checked = this.isSlowPc;
        slowPcCheckbox.addEventListener('change', () => {
            this.setIsSlowPc(slowPcCheckbox.checked);

            if (this.isSlowPc) {
                gradientBg.classList.add('hidden');
            } else {
                gradientBg.classList.remove('hidden');
            }
        });

        if (this.isUerAuthorized()) {
            const collections = [
                'users',
                'actions',
                'inventory',
                'audio',
                'items',
                'wheel_items',
            ];
            for (const collection of collections) {
                this.pb.collection(collection).subscribe('*', (e) => {
                    document.dispatchEvent(new CustomEvent(`record.${collection}.${e.action}`, {
                        detail: {
                            'record': e.record,
                        },
                    }));
                });
            }

            const avatar = Helper.getFile('avatar', this.getUserRecord());
            const profile = document.querySelector('.profile');
            const profileImg = profile.querySelector('img');
            const user = this.getUserRecord();

            profileImg.src = avatar;
            profileImg.style.borderColor = user.color;

            const timerBlock = document.getElementById('timer');
            const resetDate = document.querySelector('.timer .timer__next-reset');
            const timer = new Timer(this.getUserAuthToken(), timerBlock, resetDate);
            
            const timerStopButton = document.querySelector('.timer button.red');
            const timerStartButton = document.querySelector('.timer button.green');
            timerStopButton.addEventListener('click', () => {
                timer.stopTimer();
            });
            timerStartButton.addEventListener('click', () => {
                timer.startTimer();
            });

            const timerCopyBlock = document.getElementById('timer-copy');
            timerCopyBlock.addEventListener('click', () => {
                navigator.clipboard.writeText(`${window.location.origin}/timer.html?t=${this.getUserAuthToken()}`);
            });

            const hiddenBlocks = document.querySelectorAll('[data-authorized]');
            for (const hiddenBlock of hiddenBlocks) {
                hiddenBlock.classList.remove('hidden');
            }
        }

        this.settings = new Settings(this.pb);
        this.cells = new Cells(this.pb);
        this.users = new Users(this.pb, this.cells);
        this.items = new Items(this.pb);
        this.inventories = new Inventories(this.pb);
        this.wheelItems = new WheelItems(this.pb);
        this.audios = new Audios(this.pb);

        await this.settings.fetch();

        await this.cells.fetch();
        this.cells.refresh();

        await this.users.fetch();
        this.users.refreshCells();
        this.users.refreshTable();

        await this.items.fetch();
        await this.inventories.fetch();
        await this.wheelItems.fetch();
        await this.audios.fetch();

        this.actions = new Actions(this.pb, this.cells, this.users);

        this.updateInnerField();
    }

    setVolume(v) {
        this.volume = v;
        this.audioPlayer.volume = v / 100;
        localStorage.setItem('volume', v);
    }

    setIsSlowPc(b) {
        this.isSlowPc = b;
        localStorage.setItem('slow_pc', b);
    }

    setAudioSrc(src) {
        this.audioPlayer.src = src;
    }

    getUserId() {
        if (this.pb.authStore) {
            return this.pb.authStore.record.id;
        }
    }

    isUerAuthorized() {
        return !!this.getUserAuthToken();
    }

    getUserAuthToken() {
        if (this.pb.authStore) {
            return this.pb.authStore.token;
        }
    }

    getUserRecord() {
        if (this.pb.authStore) {
            return this.pb.authStore.record;
        }
    }

    requestNextStepType() {
        this.requestsWorker.postMessage({
            method: 'getNextStepType',
        });
    }

    showActionButton(nextStepType) {
        if (!this.isUerAuthorized()) return;

        this.nextStepType = nextStepType;

        const action = Helper.actions[this.nextStepType];
        if (!action) return;

        const actionButton = document.querySelector('.actions-buttons button');

        actionButton.dataset.graphPath = action.modal;
        actionButton.innerText = action.name;
        actionButton.style.background = action.color;

        actionButton.classList.remove('hidden');
    }

    updateInnerField() {
        this.requestNextStepType();
    }
}

const a = new App();
await a.init();
export const app = a;