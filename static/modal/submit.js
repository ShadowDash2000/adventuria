import {app} from "../app.js";

export default class Submit {
    constructor() {
        this.submitModal = document.querySelector('.graph-modal__content.submit');
        this.submitDeclineButton = this.submitModal.querySelector('.button.decline');
        this.submitAcceptButton = this.submitModal.querySelector('.button.accept');

        const eventHandler = (e) => {
            this.submitActions(e);
        }

        this.submitDeclineButton.addEventListener('click', eventHandler);
        this.submitAcceptButton.addEventListener('click', eventHandler);
    }

    open(options) {
        this.options = {
            onAccept: () => {},
            onDecline: () => {},
            text: '',
            backModal: '',
            ...options
        }

        this.submitModal.querySelector('.text').innerText = options.text;

        app.modal.close();
        app.modal.open('submit', {
            speed: 100,
            animation: 'fadeInUp',
        });
    }

    submitActions(e) {
        e.preventDefault();

        const action = e.currentTarget.dataset.action;

        switch (action) {
            case 'decline':
                app.modal.close();
                if (this.options.backModal) {
                    app.modal.open(this.options.backModal);
                }
                this.options.onDecline();
                break;
            case 'accept':
                app.modal.close();
                this.options.onAccept();
                break;
        }
    }
}