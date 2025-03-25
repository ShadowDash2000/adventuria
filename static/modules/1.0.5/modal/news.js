import {app} from "../app.js";
import Helper from "../helper.js";

const modal = document.getElementById('news-modal');
const modalContent = modal.querySelector('.news-modal__content');
const newsSentinel = modal.querySelector('.sentinel');
const newsItemTemplate = document.getElementById('news-item-template');

let isLoading = false;
let page = 1;
let totalPages = 1;
const limit = 10;
let observer;

document.addEventListener('modal.open', (e) => {
    if (e.detail.modalName !== 'news') return;

    modalContent.innerHTML = '';

    observer = new IntersectionObserver(async (entries, observer) => {
        const entry = entries[0];

        if (entry.isIntersecting && !isLoading) {
            if (page > totalPages) {
                observer.unobserve(newsSentinel);
                return;
            }

            isLoading = true;
            const news = await fetchNews(page, limit);
            totalPages = news.totalPages;
            page++;
            isLoading = false;

            for (const newsItem of news.items) {
                modalContent.appendChild(createNewsNode(newsItem));
            }
        }
    });

    observer.observe(newsSentinel);
});

async function fetchNews(page, limit) {
    return await app.pb.collection('news').getList(page, limit, {
        sort: '-created',
    });
}

function createNewsNode(newsItem) {
    const newsItemNode = newsItemTemplate.content.cloneNode(true).firstElementChild;
    const dateTime = newsItemNode.querySelector('.news-item__date');
    const text = newsItemNode.querySelector('.news-item__text');

    dateTime.innerText = Helper.formatDateLocalized(newsItem.created);
    text.innerHTML = newsItem.text;

    return newsItemNode;
}

document.addEventListener('modal.close', (e) => {
    if (e.detail.modalName !== 'news') return;
    page = 1;
    totalPages = 1;
    observer.unobserve(newsSentinel);
});