function getPerPage() {
    return parseInt(localStorage.getItem('perPage') || '20');
}

function setPerPage(n) {
    localStorage.setItem('perPage', String(n));
}

// handlerPath — строка вида "Admin.goUsersPage": inline onclick видит только window,
// а обработчики живут в объектах. Имя функции (fn.name) недостаточно — у AIKeys/Shop оба goPage.
function renderPagination(page, total, perPage, handlerPath) {
    const totalPages = Math.ceil(total / perPage);
    if (totalPages <= 1 && total <= perPage) return '';

    let html = '<div class="pagination">';

    html += '<div class="pagination__sizes">';
    html += '<span class="pagination__label">Показывать:</span>';
    [20, 50, 100].forEach(n => {
        const active = n === perPage ? ' pagination__size--active' : '';
        html += `<button class="pagination__size${active}" onclick="${handlerPath}(1, ${n})">${n}</button>`;
    });
    html += '</div>';

    html += '<div class="pagination__nav">';
    if (page > 1) {
        html += `<button class="pagination__btn" onclick="${handlerPath}(${page - 1}, ${perPage})">← Назад</button>`;
    }

    const maxVisible = 5;
    let start = Math.max(1, page - Math.floor(maxVisible / 2));
    let end = Math.min(totalPages, start + maxVisible - 1);
    if (end - start < maxVisible - 1) {
        start = Math.max(1, end - maxVisible + 1);
    }

    if (start > 1) {
        html += `<button class="pagination__btn" onclick="${handlerPath}(1, ${perPage})">1</button>`;
        if (start > 2) html += '<span class="pagination__dots">…</span>';
    }

    for (let i = start; i <= end; i++) {
        const cls = i === page ? ' pagination__btn--active' : '';
        html += `<button class="pagination__btn${cls}" onclick="${handlerPath}(${i}, ${perPage})">${i}</button>`;
    }

    if (end < totalPages) {
        if (end < totalPages - 1) html += '<span class="pagination__dots">…</span>';
        html += `<button class="pagination__btn" onclick="${handlerPath}(${totalPages}, ${perPage})">${totalPages}</button>`;
    }

    if (page < totalPages) {
        html += `<button class="pagination__btn" onclick="${handlerPath}(${page + 1}, ${perPage})">Далее →</button>`;
    }

    html += '</div>';
    html += `<div class="pagination__info">Всего: ${total} | Страница ${page} из ${totalPages}</div>`;
    html += '</div>';

    return html;
}
