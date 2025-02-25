import {app} from "../app.js";

export default class Submit {
    constructor(options) {
        const defaultOptions = {
            onAccept: () => {
            },
            text: '',
            backModal: '',
        }
        this.options = Object.assign(defaultOptions, options);

        const submitModal = document.querySelector('.graph-modal__content.submit');
        const submitDeclineButton = submitModal.querySelector('.button.decline');
        const submitAcceptButton = submitModal.querySelector('.button.accept');

        submitModal.querySelector('.text').innerText = options.text;

        const eventHandler = (e) => {
            this.submitActions(e);
        }

        submitDeclineButton.addEventListener('click', eventHandler);
        submitAcceptButton.addEventListener('click', eventHandler);

        document.addEventListener('modal.close', (e) => {
            if (e.detail.modalName !== 'submit') return;

            if (this.options.backModal) {
                app.modal.open(this.options.backModal);
            }

            submitDeclineButton.removeEventListener('click', eventHandler);
            submitAcceptButton.removeEventListener('click', eventHandler);
        });
    }

    open() {
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
                break;
            case 'accept':
                this.options.onAccept();
                break;
        }
    }
}