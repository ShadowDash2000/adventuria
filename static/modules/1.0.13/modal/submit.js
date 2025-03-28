export default class Submit {
    constructor(modal) {
        this.modal = modal;
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

        this.modal.close();
        this.modal.open('submit', {
            speed: 100,
            animation: 'fadeInUp',
        });
    }

    submitActions(e) {
        e.preventDefault();

        const action = e.currentTarget.dataset.action;

        switch (action) {
            case 'decline':
                this.modal.close();
                if (this.options.backModal) {
                    this.modal.open(this.options.backModal);
                }
                this.options.onDecline();
                break;
            case 'accept':
                this.modal.close();
                this.options.onAccept();
                break;
        }
    }
}