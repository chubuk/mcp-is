const Admin = {
    _usersPage: 1,
    _usersQuery: '',
    _usersDebounce: null,

    async renderUsers(page, perPage) {
        page = page || 1;
        perPage = perPage || getPerPage();
        const q = this._usersQuery || '';

        const res = await API.get(`/api/admin/users?page=${page}&per_page=${perPage}&q=${encodeURIComponent(q)}`);
        const data = (res && res.ok) ? res.data : { items: [], total: 0, page: 1, per_page: perPage };
        const users = data.items || [];
        const total = data.total || 0;

        let rows = '';
        if (users.length === 0) {
            rows = '<tr><td colspan="6" class="empty">Пользователи не найдены</td></tr>';
        } else {
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
        }

        const pagination = renderPagination(data.page, total, perPage, Admin.goUsersPage);

        return `
            <div class="page-header">
                <div>
                    <h1>Пользователи</h1>
                    <p class="desc">Управление пользователями системы</p>
                </div>
            </div>
            <div class="search-box">
                <input type="text" id="user-search" placeholder="Поиск по никнейму, email, телефону, ФИО..." 
                       value="${this._usersQuery}" oninput="Admin.onUserSearch(this.value)">
            </div>
            <div class="card">
                <div class="table-wrap">
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
                </div>
                ${pagination}
            </div>`;
    },

    onUserSearch(value) {
        clearTimeout(this._usersDebounce);
        this._usersDebounce = setTimeout(() => {
            this._usersQuery = value;
            this._usersPage = 1;
            this.renderUsers(1, getPerPage()).then(html => {
                document.getElementById('app').innerHTML = html;
                const input = document.getElementById('user-search');
                if (input) {
                    input.focus();
                    input.setSelectionRange(input.value.length, input.value.length);
                }
            });
        }, 700);
    },

    goUsersPage(page, perPage) {
        setPerPage(perPage);
        Admin._usersPage = page;
        Admin.renderUsers(page, perPage).then(html => {
            document.getElementById('app').innerHTML = html;
        });
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
                <div class="table-wrap">
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
            <div class="stats-grid stats-grid--3">
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
    },

    async renderServer() {
        const res = await API.get('/api/admin/server');
        const s = (res && res.ok) ? res.data : {};
        const myRole = localStorage.getItem('role');
        const isSuperAdmin = myRole === 'super_admin';

        let statusClass = 'alert--success';
        let statusText = '✅ Действителен';
        if (s.cert_needs_attention) {
            statusClass = 'alert--error';
            statusText = '⚠️ Обновление не удалось (' + (s.retry_count || 0) + '/10)';
        } else if (s.cert_issuer === 'Self-signed' || !s.cert_issuer) {
            statusClass = 'alert--info';
            statusText = 'ℹ️ Self-signed сертификат';
        }

        const domainField = isSuperAdmin
            ? `<div class="form-group">
                <label>Доменное имя</label>
                <input type="text" id="srv-domain" value="${s.domain || 'localhost'}" placeholder="example.com">
               </div>`
            : `<div class="form-group">
                <label>Доменное имя</label>
                <input type="text" value="${s.domain || 'localhost'}" disabled>
               </div>`;

        const httpsField = isSuperAdmin
            ? `<div class="form-group">
                <label>HTTPS (Let's Encrypt)</label>
                <label><input type="radio" name="srv-https" value="false" ${!s.https_enabled ? 'checked' : ''}> Выключен (self-signed)</label><br>
                <label><input type="radio" name="srv-https" value="true" ${s.https_enabled ? 'checked' : ''}> Включён (Let's Encrypt)</label>
               </div>`
            : `<div class="form-group">
                <label>HTTPS</label>
                <input type="text" value="${s.https_enabled ? 'Let\'s Encrypt' : 'Self-signed'}" disabled>
               </div>`;

        return `
            <div class="page-header">
                <div>
                    <h1>Настройки сервера</h1>
                    <p class="desc">Доменное имя и HTTPS-сертификат</p>
                </div>
            </div>

            <div class="card">
                <div class="card__title">Основные настройки</div>
                <div id="srv-alert"></div>
                ${domainField}
                ${httpsField}
                <div class="form-group">
                    <label>Email для Let's Encrypt</label>
                    <input type="email" id="srv-email" value="${s.letsencrypt_email || ''}" placeholder="admin@example.com" ${!isSuperAdmin ? 'disabled' : ''}>
                </div>
                <div class="form-group">
                    <label>
                        <input type="checkbox" id="srv-auto-renew" ${s.auto_renew !== false ? 'checked' : ''}> Автообновление сертификата
                    </label>
                </div>
                <div class="form-group">
                    <label>Обновлять за (дней до окончания)</label>
                    <input type="number" id="srv-renew-days" value="${s.renew_days_before || 30}" min="14" max="60" ${!isSuperAdmin ? 'disabled' : ''}>
                </div>
                ${isSuperAdmin ? '<button class="btn btn--primary" onclick="Admin.saveServer()">Сохранить</button>' : ''}
            </div>

            <div class="card">
                <div class="card__title">Статус сертификата</div>
                <div class="alert ${statusClass}">${statusText}</div>
                <div class="table-wrap">
                <table class="table">
                    <tr><td><strong>Тип</strong></td><td>${s.cert_issuer || 'Self-signed'}</td></tr>
                    <tr><td><strong>Домен</strong></td><td>${s.domain || 'localhost'}</td></tr>
                    <tr><td><strong>Истекает</strong></td><td>${s.cert_expires_at ? new Date(s.cert_expires_at).toLocaleDateString('ru-RU') : '—'}</td></tr>
                    ${s.last_renew_error ? `<tr><td><strong>Ошибка</strong></td><td style="color:#e74c3c">${s.last_renew_error.substring(0, 100)}</td></tr>` : ''}
                    ${s.retry_count > 0 ? `<tr><td><strong>Попыток</strong></td><td>${s.retry_count} / 10</td></tr>` : ''}
                </table>
                </div>
                <div style="margin-top:12px;display:flex;gap:8px;flex-wrap:wrap">
                    <button class="btn btn--secondary" onclick="Admin.issueCert()">Выпустить сертификат</button>
                    <button class="btn btn--ghost" onclick="Admin.renewCert()">Продлить сертификат</button>
                </div>
            </div>`;
    },

    async saveServer() {
        const domain = document.getElementById('srv-domain').value.trim();
        const https = document.querySelector('input[name="srv-https"]:checked').value === 'true';
        const email = document.getElementById('srv-email').value.trim();
        const autoRenew = document.getElementById('srv-auto-renew').checked;
        const renewDays = parseInt(document.getElementById('srv-renew-days').value) || 30;

        const res = await API.put('/api/admin/server', {
            domain: domain,
            https_enabled: https,
            letsencrypt_email: email,
            auto_renew: autoRenew,
            renew_days_before: renewDays,
        });

        if (res && res.ok) {
            document.getElementById('srv-alert').innerHTML =
                '<div class="alert alert--success">Настройки сохранены</div>';
            setTimeout(() => Admin.renderServer().then(html => {
                document.getElementById('app').innerHTML = html;
            }), 1000);
        } else {
            document.getElementById('srv-alert').innerHTML =
                '<div class="alert alert--error">' + (res ? res.data.error : 'Ошибка') + '</div>';
        }
    },

    async issueCert() {
        if (!confirm('Выпустить сертификат Let\'s Encrypt? Это может занять несколько минут.')) return;
        const res = await API.post('/api/admin/server/cert-issue');
        if (res && res.ok) {
            alert('Запрос на выпуск сертификата отправлен');
            Admin.renderServer().then(html => document.getElementById('app').innerHTML = html);
        } else {
            alert('Ошибка: ' + (res ? res.data.error : 'Сервер недоступен'));
        }
    },

    async renewCert() {
        if (!confirm('Продлить сертификат?')) return;
        const res = await API.post('/api/admin/server/cert-renew');
        if (res && res.ok) {
            alert('Запрос на продление отправлен');
        } else {
            alert('Ошибка: ' + (res ? res.data.error : 'Сервер недоступен'));
        }
    },

    async renderPDN() {
        const [settingsRes, logRes, statsRes] = await Promise.all([
            API.get('/api/admin/pdn'),
            API.get('/api/admin/pdn/log?limit=20'),
            API.get('/api/admin/pdn/stats'),
        ]);

        const s = (settingsRes && settingsRes.ok) ? settingsRes.data : {};
        const log = (logRes && logRes.ok) ? logRes.data : { entries: [] };
        const stats = (statsRes && statsRes.ok) ? statsRes.data : {};

        let logRows = '';
        if (log.entries && log.entries.length > 0) {
            log.entries.forEach(e => {
                logRows += `
                    <tr>
                        <td>${e.timestamp ? new Date(e.timestamp).toLocaleString('ru-RU') : '—'}</td>
                        <td><code>${e.tool}</code></td>
                        <td>${(e.data_types || []).join(', ')}</td>
                        <td>${e.jurisdiction || '—'}</td>
                        <td>${e.llm_provider || '—'}</td>
                        <td>${e.masked ? '<span style="color:#27ae60">Да</span>' : 'Нет'}</td>
                    </tr>`;
            });
        } else {
            logRows = '<tr><td colspan="6" class="empty">Нет записей за сегодня</td></tr>';
        }

        return `
            <div class="page-header">
                <div>
                    <h1>Персональные данные (152-ФЗ)</h1>
                    <p class="desc">Настройки и журнал доступа к персональным данным</p>
                </div>
            </div>

            <div class="stats-grid">
                <div class="card">
                    <div class="card__title">Всего запросов</div>
                    <div style="font-size:32px;font-weight:700">${stats.total_today || 0}</div>
                </div>
                <div class="card">
                    <div class="card__title">Замаскировано</div>
                    <div style="font-size:32px;font-weight:700;color:#27ae60">${stats.masked_today || 0}</div>
                </div>
            </div>

            <div class="card">
                <div class="card__title">Настройки</div>
                <div id="pdn-alert"></div>
                <div class="form-group">
                    <label>
                        <input type="checkbox" id="pdn-enabled" ${s.enabled ? 'checked' : ''}> Контроль персональных данных
                    </label>
                </div>
                <div class="form-group">
                    <label>
                        <input type="checkbox" id="pdn-log" ${s.log ? 'checked' : ''}> Логирование доступа
                    </label>
                </div>
                <div class="form-group">
                    <label>
                        <input type="checkbox" id="pdn-mask" ${s.mask_pii ? 'checked' : ''}> Маскировка ПДн
                    </label>
                </div>
                <button class="btn btn--primary" onclick="Admin.savePDN()">Сохранить</button>
            </div>

            <div class="card">
                <div class="card__title">Журнал доступа к ПДн (сегодня)</div>
                <div class="table-wrap">
                <table class="table">
                    <thead>
                        <tr>
                            <th>Время</th>
                            <th>Команда</th>
                            <th>Данные</th>
                            <th>Юрисдикция</th>
                            <th>LLM</th>
                            <th>Маска</th>
                        </tr>
                    </thead>
                    <tbody>${logRows}</tbody>
                </table>
                </div>
            </div>`;
    },

    async savePDN() {
        const enabled = document.getElementById('pdn-enabled').checked;
        const log = document.getElementById('pdn-log').checked;
        const mask = document.getElementById('pdn-mask').checked;

        const res = await API.put('/api/admin/pdn', {
            enabled: enabled,
            log: log,
            mask_pii: mask,
        });

        if (res && res.ok) {
            document.getElementById('pdn-alert').innerHTML =
                '<div class="alert alert--success">Настройки сохранены</div>';
            setTimeout(() => Admin.renderPDN().then(html => {
                document.getElementById('app').innerHTML = html;
            }), 1000);
        } else {
            document.getElementById('pdn-alert').innerHTML =
                '<div class="alert alert--error">' + (res ? res.data.error : 'Ошибка') + '</div>';
        }
    }
};
