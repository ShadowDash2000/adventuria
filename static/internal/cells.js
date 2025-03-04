import Helper from "../helper.js";

export default class Cells {
    collectionName = 'cells';

    rows = 9;

    constructor(pb) {
        this.pb = pb;
        this.cells = null;

        this.cellTemplate = document.getElementById('cell-template');
        this.cellTemplateRight = document.getElementById('cell-template-right');
        this.specialCellTemplate = document.getElementById('special-cell-template');

        this.positions = [
            {
                container: document.getElementById('cell-bottom-left'),
                template: this.specialCellTemplate,
                type: 'corner',
            },
            {
                container: document.querySelector('.left-row'),
                template: this.cellTemplate,
                type: 'reverse',
            },
            {
                container: document.getElementById('cell-top-left'),
                template: this.specialCellTemplate,
                type: 'corner',
            },
            {
                container: document.querySelector('.top-row'),
                template: this.cellTemplate,
            },
            {
                container: document.getElementById('cell-top-right'),
                template: this.specialCellTemplate,
                type: 'corner',
            },
            {
                container: document.querySelector('.right-row'),
                template: this.cellTemplateRight,
            },
            {
                container: document.getElementById('cell-bottom-right'),
                template: this.specialCellTemplate,
                type: 'corner',
            },
            {
                container: document.querySelector('.bottom-row'),
                template: this.cellTemplateRight,
                type: 'reverse',
            },
        ];
    }

    async fetch() {
        this.cells = await this.pb.collection(this.collectionName).getFullList({
            sort: 'sort',
        });
    }

    refresh() {
        let posIndex = 0;
        let countInRow = 0;
        let cellContainer, cellNode, position;
        this.cells.forEach((cell, key) => {
            position = this.positions[posIndex];

            cellContainer = position.container;
            cellNode = position.template.content.cloneNode(true).firstElementChild;

            if (position.type === 'reverse') {
                cellNode.style.order = this.rows - countInRow;
            }

            countInRow++;
            if (position.type === 'corner' || countInRow >= this.rows) {
                countInRow = 0;
                posIndex = (posIndex + 1) % this.positions.length;
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

            this.cells[key]['cellElement'] = cellContainer.appendChild(cellNode);
        });
    }

    getById(cellId) {
        for (const cell of this.cells) {
            if (cell.id === cellId) return cell;
        }
    }

    getByPassed(cellsPassed) {
        return this.cells[cellsPassed % this.cells.length];
    }

    getAll() {
        return this.cells;
    }
}