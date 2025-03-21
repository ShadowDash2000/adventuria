export default class Helper {
    static getFile(key, item, params) {
        if (!item[key]) return '';

        let uri = "/api/files/" + item.collectionId + "/" + item.id + "/" + item[key];

        let i = 0;
        for (const param in params) {
            if (i === 0) uri += `?`;
            else uri += `&`;

            uri += `${param}=${params[param]}`;

            i++;
        }

        return uri;
    }

    static actions = {
        'roll': {
            'name': 'РОЛЛ',
            'modal': 'dice',
            'color': '',
            'template': 'БРОСИЛ КУБИКИ НА: {{VALUE}}',
        },
        'reroll': {
            'name': 'РЕРОЛЛ',
            'modal': '',
            'color': '#20c7bf',
            'template': 'РЕРОЛЬНУЛ: {{VALUE}}',
        },
        'drop': {
            'name': 'DROP',
            'modal': '',
            'color': '#c72020',
            'template': 'ДРОПНУЛ: {{VALUE}}',
        },
        'chooseResult': {
            'name': 'ВЫБРАТЬ РЕЗУЛЬТАТ',
            'modal': 'game-result',
            'color': '#20c723',
            'template': 'ЗАВЕРШИЛ: {{VALUE}}',
        },
        'chooseGame': {
            'name': 'ВЫБРАТЬ ИГРУ',
            'modal': 'game-picker',
            'color': '#20c7bf',
            'template': 'НАЧАЛ: {{VALUE}}',
        },
        'rollCell': {
            'name': 'КОЛЁСИКО',
            'modal': 'wheel',
            'color': '',
            'template': 'ВЫРОЛЯЛ НА КОЛЕСЕ КЛЕТКУ: {{VALUE}}',
        },
        'rollItem': {
            'name': 'КОЛЁСИКО',
            'modal': 'wheel',
            'color': '',
        },
        'rollWheelPreset': {
            'name': 'КОЛЁСИКО',
            'modal': 'wheel',
            'color': '',
            'template': 'ВЫРОЛЯЛ НА КОЛЕСЕ: {{VALUE}}',
        },
    };

    static formatDateLocalized(isoString) {
        return new Intl.DateTimeFormat("ru-RU", {
            day: "2-digit",
            month: "2-digit",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
            hour12: false,
        })
            .format(new Date(isoString))
            .replace(',', ' ');
    }

    static shuffleArray(array) {
        for (let i = array.length - 1; i >= 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [array[i], array[j]] = [array[j], array[i]];
        }
    }

    static formDataToJson(formData) {
        let o = {};
        formData.forEach((value, key) => o[key] = value);
        return JSON.stringify(o);
    }
}