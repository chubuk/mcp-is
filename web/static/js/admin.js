const Admin = {
    async renderUsers() {
        const res = await API.get('/api/admin/users');
        const users = (res && res.ok) ? res.data : [];

        let rows = '';
        users.forEach(u => {
            const roleBadge = u.role === 'super_admin' ? '<span style="color:#e74c3c;font-weight:600">Super Admin</span>'
                : u.role === 'admin' ? '<span style="color:#f39c12;font-weight:600">Admin</span>'
                : 'User';
            const statusBadge = u.is_active
                ? '<span style="color:#27ae60">Активен</span>'
                : '<span style="color:#e74c3c">Заблокирован</span>';

            rows += `
                <tr>
                    <td>${u.nickname}</td>
                    <td>${u.email || '—'}</td>
                    <td>${roleBadge}</td>
                    <td>${statusBadge}</td>
                    <td>${u.created_at ? u.created_at.split('T')[0] : '—'}</td>
                    <td>
                        <button class="btn btn--ghost btn--sm" onclick="Admin.showUser('${u.id}')">Подробнее</button>
                    </td>
                </tr>`;
        });

        return `
            <div class="page-header">
                <div>
                    <h1>Пользователи</h1>
                    <p class="desc">Управление пользователями системы</p>
                </div>
            </div>
            <div class="card">
                <table class="table">
                    <thead>
                        <tr>
                            <th>Никнейм</th>
                            <th>Email</th>
                            <th>Роль</th>
                            <th>Статус</th>
                            <th>Создан</th>
                            <th></th>
                        </tr>
                    </thead>
                    <tbody>${rows}</tbody>
                </table>
            </div>`;
    },

    async showUser(id) {
        const res = await API.get('/api/admin/users/' + id);
        if (!res || !res.ok) {
            alert('Пользователь не найден');
            return;
        }
        const u = res.data;
        const myId = localStorage.getItem('user_id');
        const myRole = localStorage.getItem('role');
        const isSuperAdmin = u.role === 'super_admin';
        const isMe = u.id === myId;

        let actions = '';
        if (!isSuperAdmin && !isMe) {
            actions += `
                <button class="btn btn--danger btn--sm" onclick="Admin.deleteUser('${u.id}')">Удалить</button>`;
            if (myRole === 'super_admin') {
                const newRole = u.role === 'admin' ? 'user' : 'admin';
                actions += `
                    <button class="btn btn--ghost btn--sm" onclick="Admin.changeRole('${u.id}', '${newRole}')">
                        ${u.role === 'admin' ? 'Понизить до User' : 'Повысить до Admin'}
                    </button>`;
            }
            if (u.is_active) {
                actions += `
                    <button class="btn btn--ghost btn--sm" onclick="Admin.blockUser('${u.id}')">Заблокировать</button>`;
            } else {
                actions += `
                    <button class="btn btn--ghost btn--sm" onclick="Admin.unblockUser('${u.id}')">Разблокировать</button>`;
            }
        } else if (isMe) {
            actions = '<span style="color:#888;font-size:13px">Вы не можете изменить свой профиль через админку</span>';
        } else {
            actions = '<span style="color:#888;font-size:13px">Super Admin защищён от изменений</span>';
        }

        return `
            <div class="page-header">
                <div>
                    <h1>${u.nickname}</h1>
                    <p class="desc">Просмотр пользователя</p>
                </div>
                <a href="/admin/users" class="btn btn--ghost">← Назад</a>
            </div>
            <div class="card">
                <div class="card__title">Основная информация</div>
                <table class="table">
                    <tr><td><strong>ID</strong></td><td><code>${u.id}</code></td></tr>
                    <tr><td><strong>Никнейм</strong></td><td>${u.nickname}</td></tr>
                    <tr><td><strong>Email</strong></td><td>${u.email || '—'}</td></tr>
                    <tr><td><strong>Телефон</strong></td><td>${u.phone || '—'}</td></tr>
                    <tr><td><strong>Имя</strong></td><td>${u.first_name || '—'}</td></tr>
                    <tr><td><strong>Фамилия</strong></td><td>${u.last_name || '—'}</td></tr>
                    <tr><td><strong>Роль</strong></td><td>${u.role}</td></tr>
                    <tr><td><strong>Статус</strong></td><td>${u.is_active ? 'Активен' : 'Заблокирован'}</td></tr>
                    <tr><td><strong>Создан</strong></td><td>${u.created_at || '—'}</td></tr>
                </table>
            </div>
            <div class="card">
                <div class="card__title">Действия</div>
                <div style="display:flex;gap:8px;flex-wrap:wrap">${actions}</div>
            </div>`;
    },

    async changeRole(id, newRole) {
        if (!confirm('Изменить роль пользователя?')) return;
        const res = await API.put('/api/admin/users/' + id + '/role', { role: newRole });
        if (res && res.ok) {
            App.navigate('/admin/users/' + id);
        } else {
            alert('Ошибка: ' + (res ? res.data.error : 'Сервер недоступен'));
        }
    },

    async blockUser(id) {
        if (!confirm('Заблокировать пользователя?')) return;
        const res = await API.put('/api/admin/users/' + id + '/status', { is_active: false });
        if (res && res.ok) {
            App.navigate('/admin/users/' + id);
        } else {
            alert('Ошибка: ' + (res ? res.data.error : 'Сервер недоступен'));
        }
    },

    async unblockUser(id) {
        if (!confirm('Разблокировать пользователя?')) return;
        const res = await API.put('/api/admin/users/' + id + '/status', { is_active: true });
        if (res && res.ok) {
            App.navigate('/admin/users/' + id);
        } else {
            alert('Ошибка: ' + (res ? res.data.error : 'Сервер недоступен'));
        }
    },

    async deleteUser(id) {
        if (!confirm('УДАЛИТЬ пользователя? Все его данные будут удалены безвозвратно!')) return;
        const res = await API.del('/api/admin/users/' + id);
        if (res && res.ok) {
            App.navigate('/admin/users');
        } else {
            alert('Ошибка: ' + (res ? res.data.error : 'Сервер недоступен'));
        }
    },

    async renderStats() {
        const res = await API.get('/api/admin/stats');
        const stats = (res && res.ok) ? res.data : {};

        return `
            <div class="page-header">
                <div>
                    <h1>Статистика</h1>
                    <p class="desc">Общая статистика системы</p>
                </div>
            </div>
            <div style="display:grid;grid-template-columns:1fr 1fr 1fr;gap:16px">
                <div class="card">
                    <div class="card__title">Пользователи</div>
                    <div style="font-size:32px;font-weight:700">${stats.total_users || 0}</div>
                </div>
                <div class="card">
                    <div class="card__title">Магазины</div>
                    <div style="font-size:32px;font-weight:700">${stats.total_shops || 0}</div>
                </div>
                <div class="card">
                    <div class="card__title">AI-ключи</div>
                    <div style="font-size:32px;font-weight:700">${stats.total_ai_keys || 0}</div>
                </div>
            </div>`;
    }
};
