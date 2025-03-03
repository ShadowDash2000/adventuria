import Helper from "../helper.js";

export default class Cells {
    collectionName = 'cells';

    constructor(pb) {
        this.pb = pb;
        this.cells = null;

        this.positions = {
            'left': document.querySelector('.left-row'),
            'top': document.querySelector('.top-row'),
            'right': document.querySelector('.right-row'),
            'bottom': document.querySelector('.bottom-row'),
            'special': {
                'start': document.getElementById('start'),
                'jail': document.getElementById('jail'),
                'big-win': document.getElementById('big-win'),
                'preset': document.getElementById('preset'),
            },
        };
        this.cellTemplate = document.getElementById('cell-template');
        this.cellTemplateRight = document.getElementById('cell-template-right');
        this.specialCellTemplate = document.getElementById('special-cell-template');
    }

    async fetch() {
        this.cells = await this.pb.collection(this.collectionName).getFullList({
            sort: '-sort',
        });
    }

    refresh() {
        for (const key in this.cells) {
            const cell = this.cells[key];

            let cellContainer, cellNode;
            switch (cell.position) {
                case 'special':
                    cellContainer = this.positions[cell.position][cell.code];
                    cellNode = this.specialCellTemplate.content.cloneNode(true);
                    break;
                case 'right':
                case 'bottom':
                    cellContainer = this.positions[cell.position];
                    cellNode = this.cellTemplateRight.content.cloneNode(true);
                    break;
                default:
                    cellContainer = this.positions[cell.position];
                    cellNode = this.cellTemplate.content.cloneNode(true);
            }

            const colorBar = cellNode.querySelector('.color-bar');
            if (colorBar) {
                colorBar.style.background = cell.color;
            }

            cellNode.querySelector('img').src = Helper.getFile('icon', cell);
            const name = cellNode.querySelector('.name');
            name.innerHTML = cell.name;
            name.dataset.id = cell.id;
            name.dataset.type = 'cell';

            this.cells[key]['cellElement'] = cellContainer.appendChild(cellNode.firstElementChild);
        }
    }

    getCellById(cellId) {
        for (const cell of this.cells) {
            if (cell.id === cellId) return cell;
        }
    }

    getCellByPassed(cellsPassed) {
        return this.cells[this.cells.length - (cellsPassed % this.cells.length) - 1];
    }

    getAll() {
        return this.cells;
    }
}