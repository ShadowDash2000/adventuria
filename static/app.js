import GraphModal from "/graph-modal/graph-modal.js";
import PocketBase from "/pocketbase/pocketbase.es.js";
import Submit from "./modal/submit.js";
import Timer from './timer.js';
import Helper from "./helper.js";
import Users from "./internal/users.js";
import Cells from "./internal/cells.js";
import Items from "./internal/items.js";
import Inventories from "./internal/inventories.js";
import WheelItems from "./internal/wheel-items.js";
import Audios from "./internal/audios.js";
import Actions from "./internal/actions.js";

class App {
    constructor() {
        this.pb = new PocketBase('/');
        this.nextStepType = '';
        this.submit = new Submit();
        this.volume = localStorage.getItem('volume') ? localStorage.getItem('volume') : 30;
        this.audioPlayer = new Audio();
        this.audioPlayer.volume = this.volume / 100;

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

        document.addEventListener('DOMContentLoaded', async () => {
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
                const timer = new Timer(this.getUserAuthToken(), timerBlock);
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

            this.cells = new Cells(this.pb);
            this.users = new Users(this.pb, this.cells);
            this.items = new Items(this.pb);
            this.inventories = new Inventories(this.pb);
            this.wheelItems = new WheelItems(this.pb);
            this.audios = new Audios(this.pb);
            this.actions = new Actions(this.pb, this.cells, this.users);

            await this.cells.fetch();
            this.cells.refresh();

            await this.users.fetch();
            this.users.refreshCells();
            this.users.refreshTable();

            await this.items.fetch();
            await this.inventories.fetch();
            await this.wheelItems.fetch();
            await this.audios.fetch();

            await this.actions.fetch(1);
            this.actions.refresh();

            await this.updateInnerField();
        });

        document.addEventListener('record.actions.create', async (e) => {
            if (e.detail.record.user !== this.getUserId()) return;
            setTimeout(async () => {
                await this.showActionButtons();
            }, 1000);
        });
        document.addEventListener('record.actions.delete', async (e) => {
            if (e.detail.record.user !== this.getUserId()) return;
            setTimeout(async () => {
                await this.showActionButtons();
            }, 1000);
        });
    }

    setVolume(v) {
        this.volume = v;
        this.audioPlayer.volume = v / 100;
        localStorage.setItem('volume', v);
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

    async showActionButtons() {
        if (!this.isUerAuthorized()) return;

        const res = await fetch('/api/get-next-step-type', {
            method: "GET",
            headers: {
                "Authorization": this.getUserAuthToken(),
            },
        });

        if (!res.ok) return;

        const json = await res.json();
        this.nextStepType = json.nextStepType;

        const actionsButtons = document.querySelector('.actions-buttons');
        const buttons = actionsButtons.querySelectorAll('button');
        for (const button of buttons) {
            button.classList.add('hidden');
        }

        let button;

        switch (json.nextStepType) {
            case 'roll':
                button = actionsButtons.querySelector('button.game-roll');
                break;
            case 'chooseResult':
                button = actionsButtons.querySelector('button.game-result');
                break;
            case 'chooseMovieResult':
                button = actionsButtons.querySelector('button.movie-result');
                break;
            case 'chooseGame':
                button = actionsButtons.querySelector('button.game-picker');
                break;
            case 'rollCell':
            case 'rollWheelPreset':
            case 'rollItem':
                button = actionsButtons.querySelector('button.wheel');
                break;
        }

        button.classList.remove('hidden');
    }

    async updateInnerField() {
        await this.showActionButtons();
    }
}

export const app = new App();