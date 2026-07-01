const AIKeys = {
    newKey: null,

    async render(page, perPage) {
        page = page || 1;
        perPage = perPage || getPerPage();

        const res = await API.get(`/api/ai-keys?page=${page}&per_page=${perPage}`);
        const data = (res && res.ok) ? res.data : { items: [], total: 0, page: 1, per_page: perPage };
        const keys = data.items || [];
        const total = data.total || 0;

        let keyAlert = '';
        if (this.newKey) {
            keyAlert = `
                <div class="card" style="border:2px solid #27ae60">
                    <div class="card__title" style="color:#27ae60">Ключ создан!</div>
                    <div class="card__desc">Скопируйте его СЕЙЧАС. Он больше не будет отображаться.</div>
                    <div class="key-display">${this.newKey}</div>
                    <button class="btn btn--primary btn--sm" onclick="AIKeys.copyKey()">Скопировать</button>
                </div>`;
            this.newKey = null;
        }

        let rows = '';
        if (keys.length === 0) {
            rows = '<tr><td colspan="4" class="empty">Ключи не созданы</td></tr>';
        } else {
            keys.forEach(k => {
                rows += `
                    <tr>
                        <td>${k.name || 'Без названия'}</td>
                        <td><code>${k.key_hash.substring(0, 16)}...</code></td>
                        <td>${k.created_at ? k.created_at.split('T')[0] : '—'}</td>
                        <td>
                            <button class="btn btn--danger btn--sm" onclick="AIKeys.delete('${k.id}')">Отозвать</button>
                        </td>
                    </tr>`;
            });
        }

        const pagination = renderPagination(data.page, total, perPage, AIKeys.goPage);

        return `
            <div class="page-header">
                <div>
                    <h1>AI-ключи</h1>
                    <p class="desc">API-ключи для подключения AI-агентов к магазину</p>
                </div>
                <button class="btn btn--primary" onclick="AIKeys.create()">+ Создать ключ</button>
            </div>

            ${keyAlert}

            <div class="card">
                <div class="table-wrap">
                <table class="table">
                    <thead>
                        <tr>
                            <th>Название</th>
                            <th>Хэш ключа</th>
                            <th>Создан</th>
                            <th></th>
                        </tr>
                    </thead>
                    <tbody>${rows}</tbody>
                </table>
                </div>
                ${pagination}
            </div>`;
    },

    async create() {
        const name = prompt('Название ключа (например, "Cursor MCP"):');
        if (name === null) return;

        const res = await API.post('/api/ai-keys', { name });
        if (res && res.ok && res.data.key) {
            this.newKey = res.data.key;
            App.navigate('/ai-keys');
        } else {
            alert('Ошибка: ' + (res ? res.data.error : 'Сервер недоступен'));
        }
    },

    copyKey() {
        if (this.newKey) {
            navigator.clipboard.writeText(this.newKey).then(() => {
                alert('Ключ скопирован!');
            });
        }
    },

    async delete(id) {
        if (!confirm('Отозвать ключ? AI-агенты с этим ключом потеряют доступ.')) return;
        await API.del('/api/ai-keys/' + id);
        App.navigate('/ai-keys');
    },

    goPage(page, perPage) {
        setPerPage(perPage);
        AIKeys.render(page, perPage).then(html => {
            document.getElementById('app').innerHTML = html;
        });
    }
};
