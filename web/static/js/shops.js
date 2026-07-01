const Shops = {
    async render(page, perPage) {
        page = page || 1;
        perPage = perPage || getPerPage();

        const res = await API.get(`/api/shops?page=${page}&per_page=${perPage}`);
        const data = (res && res.ok) ? res.data : { items: [], total: 0, page: 1, per_page: perPage };
        const shops = data.items || [];
        const total = data.total || 0;

        let rows = '';
        if (shops.length === 0) {
            rows = '<tr><td colspan="5" class="empty">Магазины не добавлены</td></tr>';
        } else {
            shops.forEach(s => {
                const roBadge = s.is_read_only
                    ? '<span style="color:#e67e22;font-weight:600">Только чтение</span>'
                    : '<span style="color:#27ae60">Чтение/Запись</span>';
                rows += `
                    <tr>
                        <td>${s.shop_url}</td>
                        <td>${s.api_key_id.substring(0, 12)}...</td>
                        <td>${s.is_default ? 'Да' : 'Нет'}</td>
                        <td>${roBadge}</td>
                        <td>
                            <button class="btn btn--danger btn--sm" onclick="Shops.delete('${s.id}')">Удалить</button>
                        </td>
                    </tr>`;
            });
        }

        const pagination = renderPagination(data.page, total, perPage, Shops.goPage);

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
                        <div id="test-connection-result" style="display:none;margin-bottom:16px"></div>
                        <button type="button" class="btn btn--secondary" onclick="Shops.testConnection()">Проверить соединение</button>
                        <button type="submit" class="btn btn--primary">Сохранить</button>
                        <button type="button" class="btn btn--ghost" onclick="Shops.hideForm()">Отмена</button>
                    </form>
                </div>
            </div>

            <div class="card">
                <div class="table-wrap">
                <table class="table">
                    <thead>
                        <tr>
                            <th>Домен</th>
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

    bind() {
        const form = document.getElementById('add-shop-form');
        if (!form) return;

        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            const res = await API.post('/api/shops', {
                shop_url: document.getElementById('shop-url').value,
                api_key_id: document.getElementById('shop-apikey').value,
                api_password: document.getElementById('shop-apipass').value,
                is_default: document.getElementById('shop-default').checked,
                is_read_only: document.getElementById('shop-readonly').checked,
            });

            if (res && res.ok) {
                App.navigate('/shops');
            } else {
                document.getElementById('shop-alert').innerHTML =
                    '<div class="alert alert--error">' + (res ? res.data.error : 'Ошибка') + '</div>';
            }
        });
    },

    async delete(id) {
        if (!confirm('Удалить магазин?')) return;
        await API.del('/api/shops/' + id);
        App.navigate('/shops');
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
                    Email: ${shop.email}<br>
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
