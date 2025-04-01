export default class ChooseBetween {
    constructor(modal) {
        this.modal = modal;
        this.chooseBetweenModal = document.querySelector('.graph-modal__content.choose-between');
        this.chooseBetweenModalContent = document.querySelector('.graph-modal__content.choose-between');
        this.submit = this.chooseBetweenModal.querySelector('.button.accept');
        this.decline = this.chooseBetweenModal.querySelector('.button.decline');
        this.activeItemId = null;

        const eventHandler = (e) => {
            this.submitActions(e);
        }

        this.decline.addEventListener('click', eventHandler);
        this.submit.addEventListener('click', eventHandler);

        document.addEventListener('modal.close.choose-between', () => {
            this.activeItemId = null;
            this.chooseBetweenModalContent.innerHTML = '';
        }, {once: true});
    }

    open(options) {
        this.options = {
            onAccept: (itemId) => {},
            onDecline: () => {},
            text: '',
            backModal: '',
            content: null,
            ...options
        }

        this.chooseBetweenModal.querySelector('.text').innerText = this.options.text;
        this.chooseBetweenModalContent.innerHTML = this.options.content;

        const items = this.chooseBetweenModalContent.querySelectorAll('[data-id]');
        items.forEach(item => {
            item.addEventListener('click', (e) => {
                const prevActive = this.chooseBetweenModalContent.querySelector(`[data-id="${this.activeItemId}"]`);
                prevActive.classList.remove('selected');

                this.activeItemId = e.currentTarget.dataset.id;
                e.currentTarget.classList.add('selected');
            });
        });

        this.modal.close();
        this.modal.open('choose-between', {
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
                if (!this.activeItemId) return;

                this.modal.close();
                this.options.onAccept(this.activeItemId);
                break;
        }
    }
}