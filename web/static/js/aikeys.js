const AIKeys = {
    newKey: null,
    newKeyName: null,
    newMCPUrl: null,
    newConnName: null,
    _mcpTypes: [],
    _shops: [],
    _lastKeys: [],
    _configAlias: '',

    slugify(str, fallback = 'mcp-is') {
        if (!str || !str.trim()) return fallback;
        // Cyrillic transliteration map
        const cyr = {'а':'a','б':'b','в':'v','г':'g','д':'d','е':'e','ё':'e','ж':'zh','з':'z','и':'i','й':'y','к':'k','л':'l','м':'m','н':'n','о':'o','п':'p','р':'r','с':'s','т':'t','у':'u','ф':'f','х':'kh','ц':'ts','ч':'ch','ш':'sh','щ':'shch','ъ':'','ы':'y','ь':'','э':'e','ю':'yu','я':'ya'};
        let s = str.toLowerCase();
        // Transliterate Cyrillic
        let result = '';
        for (const ch of s) { result += cyr[ch] !== undefined ? cyr[ch] : ch; }
        // Keep only latin, digits, hyphens
        result = result.replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '');
        return result || fallback;
    },

    // Shared with the onboarding page (auth.js): fill a client config template.
    fillTemplate(template, alias, url, apiKey) {
        return template
            .replace(/\{\{KEY_NAME\}\}/g, alias)
            .replace(/\{\{URL\}\}/g, url)
            .replace(/\{\{API_KEY\}\}/g, apiKey);
    },

    downloadText(fileName, text) {
        const blob = new Blob([text], { type: 'text/plain' });
        const a = document.createElement('a');
        a.href = URL.createObjectURL(blob);
        a.download = fileName;
        a.click();
        URL.revokeObjectURL(a.href);
    },

    async loadShops() {
        const res = await API.get('/api/shops?per_page=100');
        AIKeys._shops = (res && res.ok && res.data.items) ? res.data.items : [];
    },

    shopOptions(selectedId) {
        return AIKeys._shops.map(s =>
            `<option value="${esc(s.id)}" ${s.id === selectedId ? 'selected' : ''}>${esc(s.name || s.shop_url)}</option>`
        ).join('');
    },

    async render(page, perPage) {
        page = page || 1;
        perPage = perPage || getPerPage();

        await AIKeys.loadShops();

        const res = await API.get(`/api/ai-keys?page=${page}&per_page=${perPage}`);
        const data = (res && res.ok) ? res.data : { items: [], total: 0, page: 1, per_page: perPage };
        const keys = data.items || [];
        AIKeys._lastKeys = keys;
        const total = data.total || 0;

        let keyAlert = '';
        if (this.newKey) {
            const keyToShow = this.newKey;
            const alias = this.newConnName || 'mcp-is';
            const mcpUrl = this.newMCPUrl || (location.origin + '/mcp');
            AIKeys._configAlias = alias;
            keyAlert = `
                <div class="card" style="border:2px solid #27ae60">
                    <div class="card__title" style="color:#27ae60">Ключ создан!</div>
                    <div class="card__desc">Скопируйте данные СЕЙЧАС. Ключ больше не будет отображаться. Потеряли — перевыпустите в редакторе ключа.</div>
                    <div style="margin-bottom:12px">
                        <label style="font-weight:600;font-size:13px;color:#555">🔑 API Key</label>
                        <div class="key-display">${keyToShow}</div>
                    </div>
                    <div style="margin-bottom:12px">
                        <label style="font-weight:600;font-size:13px;color:#555">🌐 MCP Server URL</label>
                        <div class="key-display">${mcpUrl}</div>
                    </div>
                    <div style="margin-bottom:12px">
                        <label style="font-weight:600;font-size:13px;color:#555">📛 Имя соединения (алиас)</label>
                        <div class="key-display">${alias}</div>
                    </div>
                    <div style="display:flex;gap:8px;flex-wrap:wrap;margin-bottom:16px">
                        <button class="btn btn--primary btn--sm" onclick="AIKeys.copyKey()">📋 Копировать ключ</button>
                        <button class="btn btn--secondary btn--sm" onclick="AIKeys.copyMCPUrl()">🌐 Копировать URL</button>
                        <button class="btn btn--ghost btn--sm" onclick="AIKeys.copyAll()">📋 Копировать всё</button>
                    </div>
                    <div style="margin-top:12px">
                        <div class="card__title">Подключение MCP-клиента</div>
                        <div class="form-group" style="margin-bottom:8px">
                            <label>Выберите систему</label>
                            <select id="mcp-client-type" onchange="AIKeys.generateMCPConfig()">
                                <option value="">— Выберите —</option>
                            </select>
                        </div>
                        <div id="mcp-config-output" style="display:none">
                            <p id="mcp-file-hint" style="color:#888;font-size:13px;margin-bottom:8px"></p>
                            <pre style="background:#f8f9fa;padding:12px;border-radius:6px;font-size:12px;overflow-x:auto;white-space:pre-wrap"><code id="mcp-config-code"></code></pre>
                            <div style="display:flex;gap:8px;margin-top:8px">
                                <button class="btn btn--ghost btn--sm" onclick="AIKeys.copyMCPConfig()">📋 Копировать</button>
                                <button class="btn btn--ghost btn--sm" id="mcp-download-btn" style="display:none" onclick="AIKeys.downloadMCPConfig()">💾 Скачать файл</button>
                            </div>
                        </div>
                    </div>
                </div>`;
            this.newKey = null;
            this.newKeyName = null;
            this.newMCPUrl = null;
            this.newConnName = null;
        }

        const noShops = AIKeys._shops.length === 0;
        const shopsNotice = noShops
            ? '<div class="alert alert--error">Сначала добавьте магазин — ИИ-ключ привязывается к магазину. <a href="/shops">Перейти к магазинам</a></div>'
            : '';

        let rows = '';
        if (keys.length === 0) {
            rows = '<tr><td colspan="6" class="empty">Ключи не созданы</td></tr>';
        } else {
            keys.forEach(k => {
                const pdnBtn = k.pdn_enabled
                    ? '<button class="btn btn--ghost btn--sm" style="color:#27ae60;border-color:#27ae60" onclick="AIKeys.togglePDN(\'' + k.id + '\', false)">✅ ПДн</button>'
                    : '<button class="btn btn--ghost btn--sm" onclick="AIKeys.togglePDN(\'' + k.id + '\', true)">☐ ПДн</button>';
                const shopCell = k.shop_id
                    ? `${k.shop_name || k.shop_id}<br><code style="font-size:11px;color:#888">${k.connection_name || ''}</code>`
                    : '<span style="color:#c0392b;font-weight:600">не привязан</span>';
                rows += `
                    <tr>
                        <td data-label="Название">${k.name || 'Без названия'}</td>
                        <td data-label="Магазин">${shopCell}</td>
                        <td data-label="Хэш"><code>${k.key_hash.substring(0, 16)}...</code></td>
                        <td data-label="ПДн">${pdnBtn}</td>
                        <td data-label="Создан">${k.created_at ? k.created_at.split('T')[0] : '—'}</td>
                        <td data-label="">
                            <div style="display:flex;gap:4px;flex-wrap:nowrap">
                                <button class="btn btn--ghost btn--sm" onclick="AIKeys.edit('${k.id}')">✏️ Ред.</button>
                                <button class="btn btn--danger btn--sm" onclick="AIKeys.delete('${k.id}')">🗑️ Отозвать</button>
                            </div>
                        </td>
                    </tr>`;
            });
        }

        const pagination = renderPagination(data.page, total, perPage, 'AIKeys.goPage');

        return `
            <div class="page-header">
                <div>
                    <h1>AI-ключи</h1>
                    <p class="desc">API-ключи для подключения AI-агентов к магазину. Каждый ключ привязан к одному магазину.</p>
                </div>
                <button class="btn btn--primary" onclick="AIKeys.showForm()" ${noShops ? 'disabled' : ''}>➕ Создать ключ</button>
            </div>

            ${keyAlert}
            ${shopsNotice}

            <div id="aikey-form" style="display:none">
                <div class="card">
                    <div class="card__title">Новый AI-ключ</div>
                    <div id="aikey-alert"></div>
                    <form id="add-aikey-form">
                        <div class="form-group">
                            <label>Магазин *</label>
                            <select id="aikey-shop" required onchange="AIKeys.onShopChange('aikey')">
                                <option value="">— Выберите магазин —</option>
                                ${AIKeys.shopOptions('')}
                            </select>
                        </div>
                        <div class="form-group">
                            <label>Имя соединения (алиас для конфига AI-клиента) *</label>
                            <input type="text" id="aikey-conn" placeholder="romashka" required>
                            <span class="form-check-hint">Только строчная латиница, цифры и дефис (1–64), например testshop-2. Заглавные и пробелы преобразуются автоматически. Подставляется в пример конфига.</span>
                        </div>
                        <div class="form-group">
                            <label>Название ключа *</label>
                            <input type="text" id="aikey-name" placeholder="Cursor MCP" required>
                        </div>
                        <div class="form-group">
                            <label>
                                <input type="checkbox" id="aikey-pdn"> Защита персональных данных (152-ФЗ)
                            </label>
                            <span class="form-check-hint">Маскировка email, телефона, ФИО для AI</span>
                        </div>
                        <button type="submit" class="btn btn--primary">💾 Создать ключ</button>
                        <button type="button" class="btn btn--ghost" onclick="AIKeys.hideForm()">✖️ Отмена</button>
                    </form>
                </div>
            </div>

            <div id="aikey-edit-form" style="display:none">
                <div class="card">
                    <div class="card__title">Редактирование AI-ключа</div>
                    <div id="aikey-edit-alert"></div>
                    <form id="edit-aikey-form">
                        <input type="hidden" id="edit-aikey-id">
                        <div class="form-group">
                            <label>Магазин *</label>
                            <select id="edit-aikey-shop" required onchange="AIKeys.onShopChange('edit-aikey')">
                                <option value="">— Выберите магазин —</option>
                                ${AIKeys.shopOptions('')}
                            </select>
                            <span class="form-check-hint">Ключ без привязки не работает — сохранение требует магазин.</span>
                        </div>
                        <div class="form-group">
                            <label>Имя соединения (алиас) *</label>
                            <input type="text" id="edit-aikey-conn" required>
                            <span class="form-check-hint">Только строчная латиница, цифры и дефис (1–64), например testshop-2. Заглавные и пробелы преобразуются автоматически.</span>
                        </div>
                        <div class="form-group">
                            <label>Название ключа</label>
                            <input type="text" id="edit-aikey-name" required>
                        </div>
                        <div class="form-group">
                            <label><input type="checkbox" id="edit-aikey-pdn"> Защита персональных данных (152-ФЗ)</label>
                            <span class="form-check-hint">Маскировка email, телефона, ФИО для AI</span>
                        </div>
                        <div style="margin-bottom:12px">
                            <label style="font-weight:600;font-size:13px;color:#555">🌐 MCP Server URL</label>
                            <div style="display:flex;gap:4px;align-items:center">
                                <div class="key-display" id="edit-mcp-url" style="flex:1;margin-bottom:0"></div>
                                <button type="button" class="btn btn--secondary btn--sm" onclick="AIKeys.copyEditMCPUrl()">📋</button>
                            </div>
                        </div>
                        <button type="submit" class="btn btn--primary">💾 Сохранить</button>
                        <button type="button" class="btn btn--secondary" onclick="AIKeys.rotate(document.getElementById('edit-aikey-id').value)">🔄 Перевыпустить ключ</button>
                        <button type="button" class="btn btn--ghost" onclick="AIKeys.hideEditForm()">✖️ Отмена</button>
                    </form>
                </div>
            </div>

            <div class="card">
                <div class="table-wrap">
                <table class="table">
                    <thead>
                        <tr>
                            <th>Название</th>
                            <th>Магазин</th>
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

    onShopChange(prefix) {
        const sel = document.getElementById(prefix + '-shop');
        const conn = document.getElementById(prefix + '-conn');
        if (!sel || !conn) return;
        const shop = AIKeys._shops.find(s => s.id === sel.value);
        if (shop && !conn.value) {
            conn.value = AIKeys.slugify(shop.name || shop.shop_url);
        }
    },

    showForm() {
        document.getElementById('aikey-form').style.display = 'block';
    },

    hideForm() {
        document.getElementById('aikey-form').style.display = 'none';
    },

    edit(id) {
        const k = AIKeys._lastKeys.find(x => x.id === id);
        if (!k) return;
        document.getElementById('edit-aikey-id').value = id;
        document.getElementById('edit-aikey-name').value = k.name || '';
        document.getElementById('edit-aikey-pdn').checked = !!k.pdn_enabled;
        // eslint-disable-next-line no-unsanitized/property -- HTML из shopOptions, данные экранированы внутри
        document.getElementById('edit-aikey-shop').innerHTML =
            '<option value="">— Выберите магазин —</option>' + AIKeys.shopOptions(k.shop_id || '');
        document.getElementById('edit-aikey-conn').value = k.connection_name || '';
        AIKeys._configAlias = k.connection_name || '';

        const scheme = location.protocol === 'https:' ? 'https' : 'http';
        const host = location.hostname + (location.port ? ':' + location.port : (scheme === 'https' ? ':443' : ':8080'));
        document.getElementById('edit-mcp-url').textContent = scheme + '://' + host + '/mcp';

        document.getElementById('aikey-edit-form').style.display = 'block';
        document.getElementById('aikey-form').style.display = 'none';
    },

    hideEditForm() {
        document.getElementById('aikey-edit-form').style.display = 'none';
    },

    bind() {
        if (document.getElementById('mcp-client-type')) {
            AIKeys.loadMCPClientTypes();
        }

        const form = document.getElementById('add-aikey-form');
        if (!form) return;

        for (const id of ['aikey-conn', 'edit-aikey-conn']) {
            const inp = document.getElementById(id);
            if (inp) inp.addEventListener('blur', () => {
                const v = inp.value.trim();
                if (v) inp.value = AIKeys.slugify(v, '');
            });
        }

        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            const f = {
                name: document.getElementById('aikey-name').value.trim(),
                shop_id: document.getElementById('aikey-shop').value,
                connection_name: document.getElementById('aikey-conn').value.trim(),
                pdn_enabled: document.getElementById('aikey-pdn').checked,
            };
            const alertEl = document.getElementById('aikey-alert');
            if (!f.shop_id) {
                alertEl.innerHTML = '<div class="alert alert--error">Выберите магазин — ИИ-ключ привязывается к магазину</div>';
                return;
            }
            if (!f.name) {
                alertEl.innerHTML = '<div class="alert alert--error">Название ключа обязательно</div>';
                return;
            }
            if (f.connection_name) f.connection_name = AIKeys.slugify(f.connection_name, '');
            if (!/^[a-z0-9-]{1,64}$/.test(f.connection_name)) {
                alertEl.innerHTML = '<div class="alert alert--error">Имя соединения: только латиница в нижнем регистре, цифры и дефис (1–64 символов)</div>';
                return;
            }

            const res = await API.post('/api/ai-keys', f);
            if (res && res.ok && res.data.key) {
                this.newKey = res.data.key;
                this.newKeyName = f.name;
                this.newMCPUrl = res.data.mcp_url;
                this.newConnName = res.data.connection_name;
                App.navigate('/ai-keys');
            } else {
                alertEl.innerHTML =
                    '<div class="alert alert--error">' + esc(res ? res.data.error : 'Ошибка') + '</div>';
            }
        });

        // Edit form binding
        const editForm = document.getElementById('edit-aikey-form');
        if (editForm) {
            editForm.addEventListener('submit', async (e) => {
                e.preventDefault();
                const id = document.getElementById('edit-aikey-id').value;
                const f = {
                    name: document.getElementById('edit-aikey-name').value.trim(),
                    shop_id: document.getElementById('edit-aikey-shop').value,
                    connection_name: document.getElementById('edit-aikey-conn').value.trim(),
                    pdn_enabled: document.getElementById('edit-aikey-pdn').checked,
                };
                const alertEl = document.getElementById('aikey-edit-alert');
                if (!f.shop_id) {
                    alertEl.innerHTML = '<div class="alert alert--error">Выберите магазин — без привязки ключ не работает</div>';
                    return;
                }
                if (!f.name) {
                    alertEl.innerHTML = '<div class="alert alert--error">Название обязательно</div>';
                    return;
                }
                if (f.connection_name) f.connection_name = AIKeys.slugify(f.connection_name, '');
                if (!/^[a-z0-9-]{1,64}$/.test(f.connection_name)) {
                    alertEl.innerHTML = '<div class="alert alert--error">Имя соединения: только латиница в нижнем регистре, цифры и дефис (1–64 символов)</div>';
                    return;
                }
                const res = await API.put('/api/ai-keys/' + id, f);
                if (res && res.ok) {
                    App.navigate('/ai-keys');
                } else {
                    alertEl.innerHTML =
                        '<div class="alert alert--error">' + esc(res ? res.data.error : 'Ошибка') + '</div>';
                }
            });
        }
    },

    async rotate(id) {
        if (!confirm('Перевыпустить ключ? Старый секрет перестанет работать у AI-агентов.')) return;
        const res = await API.post('/api/ai-keys/' + id + '/rotate', {});
        if (res && res.ok && res.data.key) {
            this.newKey = res.data.key;
            this.newKeyName = res.data.name;
            this.newMCPUrl = res.data.mcp_url;
            this.newConnName = res.data.connection_name;
            App.navigate('/ai-keys');
        } else {
            showToast(res && res.data && res.data.error ? res.data.error : 'Не удалось перевыпустить ключ');
        }
    },

    copyEditMCPUrl() {
        const el = document.getElementById('edit-mcp-url');
        if (el) {
            navigator.clipboard.writeText(el.textContent.trim()).then(() => {
                showToast('URL скопирован!');
            }).catch(() => {});
        }
    },

    copyKey() {
        const el = document.querySelectorAll('.key-display')[0];
        if (el) {
            navigator.clipboard.writeText(el.textContent.trim()).then(() => {
                showToast('Ключ скопирован!');
            }).catch(() => {});
        }
    },

    copyMCPUrl() {
        const el = document.querySelectorAll('.key-display')[1];
        if (el) {
            navigator.clipboard.writeText(el.textContent.trim()).then(() => {
                showToast('URL скопирован!');
            }).catch(() => {});
        }
    },

    copyAll() {
        const keyEl = document.querySelectorAll('.key-display')[0];
        const urlEl = document.querySelectorAll('.key-display')[1];
        if (keyEl && urlEl) {
            const text = 'MCP Server: ' + urlEl.textContent.trim() + '\nAPI Key: ' + keyEl.textContent.trim();
            navigator.clipboard.writeText(text).then(() => {
                showToast('Ключ и URL скопированы!');
            }).catch(() => {});
        }
    },

    async togglePDN(id, enable) {
        const res = await API.put('/api/ai-keys/' + id, { pdn_enabled: enable });
        if (res && res.ok) {
            App.navigate('/ai-keys');
        } else {
            showToast(res && res.data && res.data.error ? res.data.error : 'Не удалось изменить ПДн');
        }
    },

    async loadMCPClientTypes() {
        const res = await API.get('/api/mcp-client-types');
        if (res && res.ok) {
            AIKeys._mcpTypes = res.data || [];
            const select = document.getElementById('mcp-client-type');
            if (select) {
                select.innerHTML = '<option value="">— Выберите —</option>';
                AIKeys._mcpTypes.forEach(t => {
                    select.innerHTML += `<option value="${esc(t.id)}">${esc(t.name)}</option>`;
                });
                const saved = localStorage.getItem('mcp_client_type');
                if (saved) {
                    select.value = saved;
                    AIKeys.generateMCPConfig();
                }
            }
        }
    },

    generateMCPConfig() {
        const select = document.getElementById('mcp-client-type');
        if (!select) return;
        const type = AIKeys._mcpTypes.find(t => t.id === select.value);
        if (!type) {
            document.getElementById('mcp-config-output').style.display = 'none';
            localStorage.removeItem('mcp_client_type');
            return;
        }

        localStorage.setItem('mcp_client_type', type.id);

        const keyEl = document.querySelectorAll('.key-display')[0];
        const urlEl = document.querySelectorAll('.key-display')[1];
        if (!keyEl || !urlEl) return;

        const url = urlEl.textContent.trim();
        const apiKey = keyEl.textContent.trim();
        const config = AIKeys.fillTemplate(type.config_template, AIKeys._configAlias || 'mcp-is', url, apiKey);

        document.getElementById('mcp-config-code').textContent = config;
        document.getElementById('mcp-file-hint').textContent =
            type.file_description ? 'Файл: ' + (type.file_name || '—') + ' → ' + type.file_description : '';
        document.getElementById('mcp-download-btn').style.display = type.file_name ? 'inline-block' : 'none';
        document.getElementById('mcp-config-output').style.display = 'block';
    },

    copyMCPConfig() {
        const code = document.getElementById('mcp-config-code');
        if (!code) return;
        navigator.clipboard.writeText(code.textContent).then(() => {
            showToast('Конфигурация скопирована!');
        }).catch(() => {});
    },

    downloadMCPConfig() {
        const code = document.getElementById('mcp-config-code');
        const select = document.getElementById('mcp-client-type');
        if (!code || !select) return;
        const type = AIKeys._mcpTypes.find(t => t.id === select.value);
        if (!type || !type.file_name) return;
        AIKeys.downloadText(type.file_name, code.textContent);
    },

    async delete(id) {
        if (!confirm('Отозвать ключ? AI-агенты с этим ключом потеряют доступ.')) return;
        await API.del('/api/ai-keys/' + id);
        App.navigate('/ai-keys');
    },

    goPage(page, perPage) {
        return showPage((p, pp) => AIKeys.render(p, pp), page, perPage);
    }
};
