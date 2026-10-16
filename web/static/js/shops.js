const Shops = {
    async render(page, perPage) {
        page = page || 1;
        perPage = perPage || getPerPage();
        const myRole = localStorage.getItem('role');
        const isAdmin = myRole === 'admin' || myRole === 'super_admin';

        const res = await API.get(`/api/shops?page=${page}&per_page=${perPage}`);
        const data = (res && res.ok) ? res.data : { items: [], total: 0, page: 1, per_page: perPage };
        const shops = data.items || [];
        const total = data.total || 0;

        let rows = '';
        if (shops.length === 0) {
            rows = '<tr><td colspan="7" class="empty">Магазины не добавлены</td></tr>';
        } else {
            shops.forEach(s => {
                const roBadge = s.is_read_only
                    ? '<span style="color:#e67e22;font-weight:600">Только чтение</span>'
                    : '<span style="color:#27ae60">Чтение/Запись</span>';
                const pdnBadge = s.pdn_enabled
                    ? '<span style="color:#27ae60;font-weight:600">ПДн ✅</span>'
                    : '<span style="color:#888">ПДн —</span>';
                const jurBadge = isAdmin && s.jurisdiction ? `<br><small style="color:#888">Юрисдикция: ${s.jurisdiction}</small>` : '';
                rows += `
                    <tr>
                        <td data-label="Магазин"><strong>${s.name || 'Без названия'}</strong><br><small style="color:#888">${s.shop_url}</small>${jurBadge}</td>
                        <td data-label="API Key"><code>${s.api_key_id.substring(0, 12)}...</code></td>
                        <td data-label="По умолчанию">${s.is_default ? 'Да' : 'Нет'}</td>
                        <td data-label="Режим">${roBadge}<br>${pdnBadge}</td>
                        <td data-label="">
                            <div style="display:flex;gap:4px;flex-wrap:nowrap">
                                <button class="btn btn--ghost btn--sm" onclick="Shops.edit('${s.id}')">✏️ Ред.</button>
                                <button class="btn btn--ghost btn--sm" onclick="Shops.checkCommands('${s.id}')">⚙️ Команды</button>
                                <button class="btn btn--danger btn--sm" onclick="Shops.delete('${s.id}')">🗑️ Удалить</button>
                            </div>
                        </td>
                    </tr>`;
            });
        }

        const pagination = renderPagination(data.page, total, perPage, 'Shops.goPage');

        return `
            <div class="page-header">
                <div>
                    <h1>Магазины</h1>
                    <p class="desc">Управление подключёнными магазинами InSales</p>
                </div>
                <button class="btn btn--primary" onclick="Shops.showForm()">+ Добавить</button>
            </div>

            <div id="shop-form" style="display:none">
                <div class="card">
                    <div class="card__title">Новый магазин</div>
                    <div id="shop-alert"></div>
                    <form id="add-shop-form">
                        <div class="form-group">
                            <label>Быстрое заполнение из InSales</label>
                            <textarea id="shop-paste" rows="10" style="width:100%;min-height:180px;font-family:monospace;font-size:13px;padding:8px;cursor:text"
                                placeholder="Вставьте данные ключа из InSales:&#10;&#10;Название ключа&#10;Идентификатор&#10;c06a4060214a31ad...&#10;Пароль&#10;802a7fc14bafd094...&#10;Пример URL&#10;https://...@myshop.myinsales.ru/..."></textarea>
                            <div id="paste-status"></div>
                            <small style="color:#888">Скопируйте ключ в InSales и вставьте сюда — поля заполнятся автоматически</small>
                        </div>
                        <hr style="border:none;border-top:1px solid #eee;margin:12px 0">
                        <div class="form-group">
                            <label>Название магазина *</label>
                            <input type="text" id="shop-name" placeholder="Моя Аптека" required>
                        </div>
                        <div class="form-group">
                            <label>Домен магазина</label>
                            <input type="text" id="shop-url" placeholder="myshop.myinsales.ru" required>
                        </div>
                        <div class="form-group">
                            <label>ID ключа API</label>
                            <input type="text" id="shop-apikey" required>
                        </div>
                        <div class="form-group">
                            <label>Пароль ключа</label>
                            <input type="password" id="shop-apipass" required>
                        </div>
                        <div class="form-group">
                            <label>
                                <input type="checkbox" id="shop-default"> Магазин по умолчанию
                            </label>
                        </div>
                        <div class="form-group">
                            <label>
                                <input type="checkbox" id="shop-readonly"> Только чтение (запрет записи)
                            </label>
                        </div>
                        <div class="form-group">
                            <label>
                                <input type="checkbox" id="shop-pdn"> Защита персональных данных (152-ФЗ)
                            </label>
                            <small style="color:#888">Маскировка email, телефона, ФИО для AI</small>
                        </div>
                        <div id="test-connection-result" style="display:none;margin-bottom:16px"></div>
                        <button type="button" class="btn btn--secondary" onclick="Shops.testConnection()">🔍 Проверить соединение</button>
                        <button type="submit" class="btn btn--primary">💾 Сохранить</button>
                        <button type="button" class="btn btn--ghost" onclick="Shops.hideForm()">✖️ Отмена</button>
                    </form>
                </div>
            </div>

            <div id="edit-shop-form" style="display:none">
                <div class="card">
                    <div class="card__title">Редактирование магазина</div>
                    <div id="edit-shop-alert"></div>
                    <form id="edit-shop-frm">
                        <input type="hidden" id="edit-shop-id">
                        <div class="form-group">
                            <label>Название магазина *</label>
                            <input type="text" id="edit-shop-name" required>
                        </div>
                        <div class="form-group">
                            <label>Домен магазина</label>
                            <input type="text" id="edit-shop-url" required>
                        </div>
                        <div class="form-group">
                            <label>ID ключа API</label>
                            <input type="text" id="edit-shop-apikey" required>
                        </div>
                        <div class="form-group">
                            <label>Пароль ключа (оставьте пустым чтобы не менять)</label>
                            <input type="password" id="edit-shop-apipass">
                        </div>
                        <div class="form-group">
                            <label><input type="checkbox" id="edit-shop-default"> Магазин по умолчанию</label>
                        </div>
                        <div class="form-group">
                            <label><input type="checkbox" id="edit-shop-readonly"> Только чтение (запрет записи)</label>
                        </div>
                        <div class="form-group">
                            <label><input type="checkbox" id="edit-shop-pdn"> Защита персональных данных (152-ФЗ)</label>
                            <small style="color:#888">Маскировка email, телефона, ФИО для AI</small>
                        </div>
                        <button type="submit" class="btn btn--primary">💾 Сохранить</button>
                        <button type="button" class="btn btn--ghost" onclick="Shops.hideEditForm()">✖️ Отмена</button>
                    </form>
                </div>
            </div>

            <div id="commands-modal" style="display:none">
                <div class="card">
                    <div class="card__title">Доступные команды</div>
                    <div id="commands-result"></div>
                    <button class="btn btn--ghost" onclick="document.getElementById('commands-modal').style.display='none'">✖️ Закрыть</button>
                </div>
            </div>

            <div class="card">
                <div class="table-wrap">
                <table class="table">
                    <thead>
                        <tr>
                            <th>Магазин</th>
                            <th>API Key ID</th>
                            <th>По умолчанию</th>
                            <th>Режим</th>
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
        document.getElementById('shop-form').style.display = 'block';
    },

    hideForm() {
        document.getElementById('shop-form').style.display = 'none';
    },

    async edit(id) {
        const res = await API.get(`/api/shops?page=1&per_page=100`);
        if (!res || !res.ok) return;
        const shop = (res.data.items || []).find(s => s.id === id);
        if (!shop) return;

        document.getElementById('edit-shop-id').value = shop.id;
        document.getElementById('edit-shop-name').value = shop.name || '';
        document.getElementById('edit-shop-url').value = shop.shop_url || '';
        document.getElementById('edit-shop-apikey').value = shop.api_key_id || '';
        document.getElementById('edit-shop-apipass').value = '';
        document.getElementById('edit-shop-default').checked = shop.is_default;
        document.getElementById('edit-shop-readonly').checked = shop.is_read_only;
        document.getElementById('edit-shop-pdn').checked = shop.pdn_enabled;

        document.getElementById('edit-shop-form').style.display = 'block';
        document.getElementById('shop-form').style.display = 'none';
    },

    hideEditForm() {
        document.getElementById('edit-shop-form').style.display = 'none';
    },

    parseInsalesData(text) {
        if (!text || text.length < 20) return null;

        const lines = text.split('\n').map(l => l.trim()).filter(l => l);
        const result = {};
        const markers = ['Идентификатор', 'Пароль', 'Формат URL', 'Пример URL', 'Дата подключения'];

        for (let i = 0; i < lines.length; i++) {
            const line = lines[i];
            const next = lines[i + 1];

            if (line === 'Идентификатор' && next) {
                result.api_key_id = next;
            } else if (line === 'Пароль' && next) {
                result.api_password = next;
            } else if (line === 'Пример URL' && next) {
                const match = next.match(/@([^/]+)\//);
                if (match) result.shop_url = match[1];
            } else if (!result.name && !markers.includes(line) &&
                       !line.includes('http') && !line.includes(':') &&
                       line.length < 50 && line.length > 0) {
                result.name = line;
            }
        }

        return result.api_key_id && result.api_password ? result : null;
    },

    onPasteData() {
        const textarea = document.getElementById('shop-paste');
        if (!textarea) return;

        const self = this;
        const parseAndFill = () => {
            const text = textarea.value;
            const data = self.parseInsalesData(text);
            if (!data) return;

            if (data.name) document.getElementById('shop-name').value = data.name;
            if (data.shop_url) document.getElementById('shop-url').value = data.shop_url;
            if (data.api_key_id) document.getElementById('shop-apikey').value = data.api_key_id;
            if (data.api_password) document.getElementById('shop-apipass').value = data.api_password;

            const statusEl = document.getElementById('paste-status');
            if (statusEl) {
                statusEl.innerHTML = '<div class="alert alert--success">Данные заполнены автоматически! Проверьте и нажмите «Сохранить».</div>';
            }

            textarea.value = '';
        };

        textarea.addEventListener('input', () => {
            clearTimeout(self._parseTimer);
            self._parseTimer = setTimeout(parseAndFill, 800);
        });

        textarea.addEventListener('focus', () => {
            const statusEl = document.getElementById('paste-status');
            if (statusEl) statusEl.innerHTML = '';
        });
    },

    bind() {
        const form = document.getElementById('add-shop-form');
        if (!form) return;

        this.onPasteData();

        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            const name = document.getElementById('shop-name').value.trim();
            if (!name) {
                document.getElementById('shop-alert').innerHTML =
                    '<div class="alert alert--error">Название магазина обязательно</div>';
                return;
            }

            const res = await API.post('/api/shops', {
                name: name,
                shop_url: document.getElementById('shop-url').value,
                api_key_id: document.getElementById('shop-apikey').value,
                api_password: document.getElementById('shop-apipass').value,
                is_default: document.getElementById('shop-default').checked,
                is_read_only: document.getElementById('shop-readonly').checked,
                pdn_enabled: document.getElementById('shop-pdn').checked,
            });

            if (res && res.ok) {
                App.navigate('/shops');
            } else {
                document.getElementById('shop-alert').innerHTML =
                    '<div class="alert alert--error">' + (res ? res.data.error : 'Ошибка') + '</div>';
            }
        });

        // Edit form binding
        const editForm = document.getElementById('edit-shop-frm');
        if (editForm) {
            editForm.addEventListener('submit', async (e) => {
                e.preventDefault();
                const id = document.getElementById('edit-shop-id').value;
                const name = document.getElementById('edit-shop-name').value.trim();
                if (!name) {
                    document.getElementById('edit-shop-alert').innerHTML =
                        '<div class="alert alert--error">Название магазина обязательно</div>';
                    return;
                }

                const body = {
                    name: name,
                    shop_url: document.getElementById('edit-shop-url').value,
                    api_key_id: document.getElementById('edit-shop-apikey').value,
                    is_default: document.getElementById('edit-shop-default').checked,
                    is_read_only: document.getElementById('edit-shop-readonly').checked,
                    pdn_enabled: document.getElementById('edit-shop-pdn').checked,
                };

                const pass = document.getElementById('edit-shop-apipass').value;
                if (pass) body.api_password = pass;

                const res = await API.put('/api/shops/' + id, body);
                if (res && res.ok) {
                    App.navigate('/shops');
                } else {
                    document.getElementById('edit-shop-alert').innerHTML =
                        '<div class="alert alert--error">' + (res ? res.data.error : 'Ошибка') + '</div>';
                }
            });
        }
    },

    async delete(id) {
        if (!confirm('Удалить магазин?')) return;
        await API.del('/api/shops/' + id);
        App.navigate('/shops');
    },

    async checkCommands(id) {
        const modal = document.getElementById('commands-modal');
        const result = document.getElementById('commands-result');
        modal.style.display = 'block';
        result.innerHTML = '<div class="alert alert--info">Проверка доступных команд...</div>';

        const res = await API.post(`/api/shops/${id}/check-commands`);
        if (res && res.ok) {
            const cmds = res.data.commands || [];
            const updatedAt = res.data.updated_at || '';
            result.innerHTML = `
                <div class="alert alert--success">
                    Найдено команд: ${res.data.count}
                    ${updatedAt ? '<br><small>Актуально: ' + updatedAt + '</small>' : ''}
                </div>
                <div style="max-height:400px;overflow-y:auto">
                    <ul style="list-style:none;padding:0">
                        ${cmds.map(c => `<li style="padding:6px 0;border-bottom:1px solid #eee">
                            <code>${c.name}</code>
                            <span style="color:#666;margin-left:8px">${c.description}</span>
                        </li>`).join('')}
                    </ul>
                </div>`;
        } else {
            result.innerHTML = '<div class="alert alert--error">' + (res ? res.data.error : 'Ошибка') + '</div>';
        }
    },

    async testConnection() {
        const resultEl = document.getElementById('test-connection-result');
        const shopUrl = document.getElementById('shop-url').value.trim();
        const apiKey = document.getElementById('shop-apikey').value.trim();
        const apiPass = document.getElementById('shop-apipass').value.trim();

        if (!shopUrl || !apiKey || !apiPass) {
            resultEl.style.display = 'block';
            resultEl.innerHTML = '<div class="alert alert--error">Заполните домен, ID ключа и пароль</div>';
            return;
        }

        resultEl.style.display = 'block';
        resultEl.innerHTML = '<div class="alert alert--info">Проверка соединения...</div>';

        const res = await API.post('/api/shops/test-connection', {
            shop_url: shopUrl,
            api_key_id: apiKey,
            api_password: apiPass,
        });

        if (res && res.ok && res.data.ok) {
            const shop = res.data.shop;
            resultEl.innerHTML = `
                <div class="alert alert--success">
                    <strong>${res.data.message}</strong><br>
                    Магазин: ${shop.name}<br>
                    Домен: ${shop.domain}<br>
                    Валюта: ${shop.currency}
                </div>`;
        } else {
            const errMsg = (res && res.data && res.data.error) ? res.data.error : 'Неизвестная ошибка';
            resultEl.innerHTML = '<div class="alert alert--error"><strong>Ошибка соединения:</strong> ' + errMsg + '</div>';
        }
    },

    goPage(page, perPage) {
        setPerPage(perPage);
        Shops.render(page, perPage).then(html => {
            document.getElementById('app').innerHTML = html;
        });
    }
};
