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
                        <td data-label="Никнейм">${u.nickname}</td>
                        <td data-label="Email">${u.email || '—'}</td>
                        <td data-label="Роль">${roleBadge}</td>
                        <td data-label="Статус">${statusBadge}</td>
                        <td data-label="Создан">${u.created_at ? u.created_at.split('T')[0] : '—'}</td>
                        <td data-label="">
                            <button class="btn btn--ghost btn--sm" onclick="App.navigate('/admin/users/${u.id}')">👁️ Подробнее</button>
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
                <button class="btn btn--danger btn--sm" onclick="Admin.deleteUser('${u.id}')">🗑️ Удалить</button>`;
            if (myRole === 'super_admin' || myRole === 'admin') {
                const canImpersonate = (myRole === 'super_admin') || (myRole === 'admin' && u.role === 'user');
                if (canImpersonate) {
                    actions += `
                        <button class="btn btn--ghost btn--sm" onclick="Admin.impersonate('${u.id}', '${u.nickname}')">
                            👤 Войти как ${u.nickname}
                        </button>`;
                }
            }
            if (myRole === 'super_admin') {
                const newRole = u.role === 'admin' ? 'user' : 'admin';
                actions += `
                    <button class="btn btn--ghost btn--sm" onclick="Admin.changeRole('${u.id}', '${newRole}')">
                        ${u.role === 'admin' ? '⬇️ Понизить' : '⬆️ Повысить'}
                    </button>`;
            }
            if (u.is_active) {
                actions += `
                    <button class="btn btn--ghost btn--sm" onclick="Admin.blockUser('${u.id}')">🔒 Заблокировать</button>`;
            } else {
                actions += `
                    <button class="btn btn--ghost btn--sm" onclick="Admin.unblockUser('${u.id}')">🔓 Разблокировать</button>`;
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

    async impersonate(userId, nickname) {
        if (!confirm('Войти как ' + nickname + '? Вы сможете вернуться кнопкой в баннере.')) return;
        const res = await API.post('/api/admin/users/' + userId + '/impersonate');
        if (res && res.ok) {
            API.impersonate(res.data.token, res.data.user_id, res.data.nickname, res.data.role);
            App.navigate('/dashboard');
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
                ${isSuperAdmin ? '<button class="btn btn--primary" onclick="Admin.saveServer()">💾 Сохранить</button>' : ''}
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
                    <button class="btn btn--secondary" onclick="Admin.issueCert()">📜 Выпустить сертификат</button>
                    <button class="btn btn--ghost" onclick="Admin.renewCert()">🔐 Продлить сертификат</button>
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
        const [settingsRes, logRes, statsRes, jurisRes] = await Promise.all([
            API.get('/api/admin/pdn'),
            API.get('/api/admin/pdn/log?limit=20'),
            API.get('/api/admin/pdn/stats'),
            API.get('/api/admin/pdn/jurisdictions'),
        ]);

        const s = (settingsRes && settingsRes.ok) ? settingsRes.data : {};
        const log = (logRes && logRes.ok) ? logRes.data : { entries: [] };
        const stats = (statsRes && statsRes.ok) ? statsRes.data : {};
        const juris = (jurisRes && jurisRes.ok) ? jurisRes.data : {};

        let logRows = '';
        if (log.entries && log.entries.length > 0) {
            log.entries.forEach(e => {
                logRows += `
                    <tr>
                        <td data-label="Время">${e.timestamp ? new Date(e.timestamp).toLocaleString('ru-RU') : '—'}</td>
                        <td data-label="Команда"><code>${e.tool}</code></td>
                        <td data-label="Данные">${(e.data_types || []).join(', ')}</td>
                        <td data-label="Юрисдикция">${e.jurisdiction || '—'}</td>
                        <td data-label="LLM">${e.llm_provider || '—'}</td>
                        <td data-label="Маска">${e.masked ? '<span style="color:#27ae60">Да</span>' : 'Нет'}</td>
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
                <div class="card__title">🛡️ АвтоПДн — Юрисдикция и автоматическая защита</div>
                <p style="color:#888;font-size:13px;margin-bottom:12px">Система автоматически включает защиту ПДн если магазин в юрисдикции РФ, а LLM-провайдер иностранный (152-ФЗ).</p>
                <div style="display:flex;gap:12px;flex-wrap:wrap;margin-bottom:12px">
                    <div class="card" style="flex:1;min-width:120px;text-align:center">
                        <div style="font-size:24px">🇷🇺</div>
                        <div style="font-size:20px;font-weight:700">${juris.jurisdictions ? (juris.jurisdictions.RU || 0) : 0}</div>
                        <div style="font-size:12px;color:#888">Россия</div>
                    </div>
                    <div class="card" style="flex:1;min-width:120px;text-align:center">
                        <div style="font-size:24px">🇰🇿</div>
                        <div style="font-size:20px;font-weight:700">${juris.jurisdictions ? (juris.jurisdictions.KZ || 0) : 0}</div>
                        <div style="font-size:12px;color:#888">Казахстан</div>
                    </div>
                    <div class="card" style="flex:1;min-width:120px;text-align:center">
                        <div style="font-size:24px">🇧🇾</div>
                        <div style="font-size:20px;font-weight:700">${juris.jurisdictions ? (juris.jurisdictions.BY || 0) : 0}</div>
                        <div style="font-size:12px;color:#888">Беларусь</div>
                    </div>
                    <div class="card" style="flex:1;min-width:120px;text-align:center">
                        <div style="font-size:24px">🌍</div>
                        <div style="font-size:20px;font-weight:700">${juris.jurisdictions ? (juris.jurisdictions.INTL || 0) : 0}</div>
                        <div style="font-size:12px;color:#888">Международные</div>
                    </div>
                </div>
                <div style="display:flex;gap:16px;font-size:13px;color:#555">
                    <span>🛡️ АвтоПДн: <strong>${juris.auto_protected || 0}</strong> запросов</span>
                    <span>✋ Ручная: <strong>${juris.manual_protected || 0}</strong> запросов</span>
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
                <div class="form-group">
                    <label>
                        <input type="checkbox" id="pdn-auto" ${s.auto_protect !== false ? 'checked' : ''}> АвтоПДн (автоматическая защита для РФ + иностранных LLM)
                    </label>
                    <small style="color:#888">Отключите для тестирования. В .env: PDN_AUTO_PROTECT=false</small>
                </div>
                <button class="btn btn--primary" onclick="Admin.savePDN()">💾 Сохранить</button>
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
        const autoPdn = document.getElementById('pdn-auto').checked;

        const res = await API.put('/api/admin/pdn', {
            enabled: enabled,
            log: log,
            mask_pii: mask,
            auto_protect: autoPdn,
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
    },

    // ==================== BACKUP ====================

    async renderBackup() {
        const [statusRes, logRes, incidentsRes, settingsRes] = await Promise.all([
            API.get('/api/admin/backup/status'),
            API.get('/api/admin/backup/log'),
            API.get('/api/admin/backup/incidents'),
            API.get('/api/admin/backup'),
        ]);

        const status = (statusRes && statusRes.ok) ? statusRes.data : {};
        const log = (logRes && logRes.ok) ? logRes.data : { entries: [] };
        const incidents = (incidentsRes && incidentsRes.ok) ? incidentsRes.data : { incidents: [] };
        const settings = (settingsRes && settingsRes.ok) ? settingsRes.data : {};

        // Stats
        const lastBackup = status.last_backup;
        const lastBackupTime = lastBackup ? new Date(lastBackup.timestamp).toLocaleString('ru-RU') : 'Нет данных';
        const lastBackupStatus = lastBackup ? (lastBackup.status === 'ok' ? '<span style="color:#27ae60">OK</span>' : '<span style="color:#e74c3c">Ошибка</span>') : '—';
        const totalBackups = status.total_backups || 0;
        const load = status.load_5m ? status.load_5m.toFixed(2) : '0.00';

        // Log rows
        let logRows = '';
        if (!log.entries || log.entries.length === 0) {
            logRows = '<tr><td colspan="6" class="empty">Нет записей</td></tr>';
        } else {
            log.entries.forEach(e => {
                const statusClass = e.status === 'ok' ? 'color:#27ae60' : e.status === 'error' ? 'color:#e74c3c' : 'color:#f39c12';
                const size = e.size_bytes ? (e.size_bytes / 1024 / 1024).toFixed(1) + ' MB' : '—';
                const time = new Date(e.timestamp).toLocaleString('ru-RU');
                logRows += `
                    <tr>
                        <td data-label="Время">${time}</td>
                        <td data-label="Статус"><span style="${statusClass}">${e.status}</span></td>
                        <td data-label="Размер">${size}</td>
                        <td data-label="Нагрузка">${e.load_5m || '—'}</td>
                        <td data-label="Длительность">${e.duration_ms ? e.duration_ms + 'мс' : '—'}</td>
                        <td data-label="Ошибка">${e.error_msg || '—'}</td>
                    </tr>`;
            });
        }

        // Incidents rows
        let incidentRows = '';
        if (!incidents.incidents || incidents.incidents.length === 0) {
            incidentRows = '<tr><td colspan="5" class="empty">Нет инцидентов</td></tr>';
        } else {
            incidents.incidents.forEach(i => {
                const sevClass = i.severity === 'critical' ? 'color:#e74c3c' : i.severity === 'warning' ? 'color:#f39c12' : 'color:#3498db';
                const time = new Date(i.timestamp).toLocaleString('ru-RU');
                incidentRows += `
                    <tr>
                        <td data-label="Время">${time}</td>
                        <td data-label="Тип">${i.type}</td>
                        <td data-label="Серьёзность"><span style="${sevClass}">${i.severity}</span></td>
                        <td data-label="Сообщение">${i.message}</td>
                        <td data-label="">${i.resolved ? '✅' : ''}</td>
                    </tr>`;
            });
        }

        return `
            <div class="page-header">
                <div>
                    <h1>Бэкапы</h1>
                    <p class="desc">Мониторинг и управление резервными копиями</p>
                </div>
            </div>

            <div class="stats-grid">
                <div class="card">
                    <div class="card__title">Последний бэкап</div>
                    <div style="font-size:24px;font-weight:700">${lastBackupTime}</div>
                    <div style="margin-top:8px">${lastBackupStatus}</div>
                </div>
                <div class="card">
                    <div class="card__title">Всего бэкапов</div>
                    <div style="font-size:32px;font-weight:700">${totalBackups}</div>
                </div>
                <div class="card">
                    <div class="card__title">Нагрузка (5min)</div>
                    <div style="font-size:32px;font-weight:700">${load}</div>
                </div>
            </div>

            <div class="card" style="margin-top:24px">
                <div class="card__title">Настройки</div>
                <div id="backup-alert"></div>
                <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px">
                    <div class="form-group">
                        <label>Метод</label>
                        <select id="backup-method" class="input">
                            <option value="pg_basebackup" ${settings.method === 'pg_basebackup' ? 'selected' : ''}>pg_basebackup (файловый снэпшот)</option>
                            <option value="pg_dump" ${settings.method === 'pg_dump' ? 'selected' : ''}>pg_dump (SQL-дамп)</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label>Интервал (часов)</label>
                        <select id="backup-interval" class="input">
                            <option value="6" ${settings.interval_hours === 6 ? 'selected' : ''}>6 часов</option>
                            <option value="12" ${settings.interval_hours === 12 ? 'selected' : ''}>12 часов</option>
                            <option value="24" ${settings.interval_hours === 24 ? 'selected' : ''}>24 часа</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label>Хранение (дней)</label>
                        <input type="number" id="backup-retention" class="input" value="${settings.retention_days || 14}" min="1" max="365">
                    </div>
                    <div class="form-group">
                        <label>Макс. нагрузка CPU (%)</label>
                        <input type="number" id="backup-load" class="input" value="${settings.max_load_pct || 50}" min="10" max="100">
                    </div>
                </div>
                <div style="margin-top:16px;display:flex;gap:8px">
                    <button class="btn btn--primary" onclick="Admin.saveBackup()">💾 Сохранить</button>
                    <button class="btn btn--ghost" onclick="Admin.runBackup()">▶️ Бэкап сейчас</button>
                </div>
            </div>

            <div class="card" style="margin-top:24px">
                <div class="card__title">Журнал бэкапов</div>
                <div class="table-wrap">
                    <table class="table">
                        <thead>
                            <tr>
                                <th>Время</th>
                                <th>Статус</th>
                                <th>Размер</th>
                                <th>Нагрузка</th>
                                <th>Длительность</th>
                                <th>Ошибка</th>
                            </tr>
                        </thead>
                        <tbody>${logRows}</tbody>
                    </table>
                </div>
            </div>

            <div class="card" style="margin-top:24px">
                <div class="card__title">Инциденты</div>
                <div class="table-wrap">
                    <table class="table">
                        <thead>
                            <tr>
                                <th>Время</th>
                                <th>Тип</th>
                                <th>Серьёзность</th>
                                <th>Сообщение</th>
                                <th></th>
                            </tr>
                        </thead>
                        <tbody>${incidentRows}</tbody>
                    </table>
                </div>
            </div>`;
    },

    async saveBackup() {
        const method = document.getElementById('backup-method').value;
        const interval = parseInt(document.getElementById('backup-interval').value);
        const retention = parseInt(document.getElementById('backup-retention').value);
        const load = parseInt(document.getElementById('backup-load').value);

        const res = await API.put('/api/admin/backup', {
            method: method,
            interval_hours: interval,
            retention_days: retention,
            max_load_pct: load,
        });

        if (res && res.ok) {
            document.getElementById('backup-alert').innerHTML =
                '<div class="alert alert--success">Настройки сохранены</div>';
            setTimeout(() => Admin.renderBackup().then(html => {
                document.getElementById('app').innerHTML = html;
            }), 1000);
        } else {
            document.getElementById('backup-alert').innerHTML =
                '<div class="alert alert--error">' + (res ? res.data.error : 'Ошибка') + '</div>';
        }
    },

    async runBackup() {
        if (!confirm('Запустить бэкап сейчас?')) return;

        const res = await API.post('/api/admin/backup/run');
        if (res && res.ok) {
            document.getElementById('backup-alert').innerHTML =
                '<div class="alert alert--success">Бэкап запущен</div>';
        } else {
            document.getElementById('backup-alert').innerHTML =
                '<div class="alert alert--error">' + (res ? res.data.error : 'Ошибка') + '</div>';
        }
    }
};
