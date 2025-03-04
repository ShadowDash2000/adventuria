export default class Helper {
    static getFile(key, item, params) {
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
            'color': '',
            'template': 'РЕРОЛЬНУЛ: {{VALUE}}',
        },
        'drop': {
            'name': 'DROP',
            'modal': '',
            'color': '',
            'template': 'ДРОПНУЛ: {{VALUE}}',
        },
        'chooseResult': {
            'name': 'ВЫБРАТЬ РЕЗУЛЬТАТ',
            'modal': 'game-result',
            'color': '',
        },
        'chooseGame': {
            'name': 'ВЫБРАТЬ ИГРУ',
            'modal': 'game-picker',
            'color': '',
        },
        'rollCell': {
            'name': 'КОЛЁСИКО',
            'modal': 'wheel',
            'color': '',
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
}