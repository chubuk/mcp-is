const Shops = {
    async render() {
        const res = await API.get('/api/shops');
        const shops = (res && res.ok) ? res.data : [];

        let rows = '';
        if (shops.length === 0) {
            rows = '<tr><td colspan="4" class="empty">Магазины не добавлены</td></tr>';
        } else {
            shops.forEach(s => {
                rows += `
                    <tr>
                        <td>${s.shop_url}</td>
                        <td>${s.api_key_id.substring(0, 12)}...</td>
                        <td>${s.is_default ? 'Да' : 'Нет'}</td>
                        <td>
                            <button class="btn btn--danger btn--sm" onclick="Shops.delete('${s.id}')">Удалить</button>
                        </td>
                    </tr>`;
            });
        }

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
                        <button type="submit" class="btn btn--primary">Сохранить</button>
                        <button type="button" class="btn btn--ghost" onclick="Shops.hideForm()">Отмена</button>
                    </form>
                </div>
            </div>

            <div class="card">
                <table class="table">
                    <thead>
                        <tr>
                            <th>Домен</th>
                            <th>API Key ID</th>
                            <th>По умолчанию</th>
                            <th></th>
                        </tr>
                    </thead>
                    <tbody>${rows}</tbody>
                </table>
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
    }
};
