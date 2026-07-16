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
            rows = '<tr><td colspan="5" class="empty">Ключи не созданы</td></tr>';
        } else {
            keys.forEach(k => {
                const pdnBadge = k.pdn_enabled
                    ? '<span style="color:#27ae60;font-weight:600">ПДн ✅</span>'
                    : '<span style="color:#888">ПДн —</span>';
                rows += `
                    <tr>
                        <td>${k.name || 'Без названия'}</td>
                        <td><code>${k.key_hash.substring(0, 16)}...</code></td>
                        <td>${pdnBadge}</td>
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
                <button class="btn btn--primary" onclick="AIKeys.showForm()">+ Создать ключ</button>
            </div>

            ${keyAlert}

            <div id="aikey-form" style="display:none">
                <div class="card">
                    <div class="card__title">Новый AI-ключ</div>
                    <div id="aikey-alert"></div>
                    <form id="add-aikey-form">
                        <div class="form-group">
                            <label>Название ключа *</label>
                            <input type="text" id="aikey-name" placeholder="Cursor MCP" required>
                        </div>
                        <div class="form-group">
                            <label>
                                <input type="checkbox" id="aikey-pdn"> Защита персональных данных (152-ФЗ)
                            </label>
                            <small style="color:#888">Маскировка email, телефона, ФИО для AI</small>
                        </div>
                        <button type="submit" class="btn btn--primary">Создать ключ</button>
                        <button type="button" class="btn btn--ghost" onclick="AIKeys.hideForm()">Отмена</button>
                    </form>
                </div>
            </div>

            <div class="card">
                <div class="table-wrap">
                <table class="table">
                    <thead>
                        <tr>
                            <th>Название</th>
                            <th>Хэш ключа</th>
                            <th>ПДн</th>
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

    showForm() {
        document.getElementById('aikey-form').style.display = 'block';
    },

    hideForm() {
        document.getElementById('aikey-form').style.display = 'none';
    },

    bind() {
        const form = document.getElementById('add-aikey-form');
        if (!form) return;

        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            const name = document.getElementById('aikey-name').value.trim();
            if (!name) {
                document.getElementById('aikey-alert').innerHTML =
                    '<div class="alert alert--error">Название ключа обязательно</div>';
                return;
            }

            const pdn = document.getElementById('aikey-pdn').checked;
            const res = await API.post('/api/ai-keys', { name, pdn_enabled: pdn });
            if (res && res.ok && res.data.key) {
                this.newKey = res.data.key;
                App.navigate('/ai-keys');
            } else {
                document.getElementById('aikey-alert').innerHTML =
                    '<div class="alert alert--error">' + (res ? res.data.error : 'Ошибка') + '</div>';
            }
        });
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
