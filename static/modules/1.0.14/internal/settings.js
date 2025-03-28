import {app} from "../app.js";

export default class Settings {
    collectionName = 'settings';

    constructor(pb) {
        this.pb = pb;
        this.settings = null;

        const rulesModal = document.getElementById('rules-modal');
        this.rules = rulesModal.querySelector('.rules-modal');
        const showRulesButton = document.getElementById('show-rules');

        showRulesButton.addEventListener('click', () => {
            app.modal.open('rules', {
                speed: 100,
                animation: 'fadeInUp',
            });
        });
    }

    async fetch() {
        this.settings = await this.pb.collection(this.collectionName).getFirstListItem();

        this.rules.innerHTML = this.settings.rules;
    }
}