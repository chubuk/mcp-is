const AIKeys = {
    newKey: null,
    newKeyName: null,
    newMCPUrl: null,
    _mcpTypes: [],

    slugify(str) {
        if (!str || !str.trim()) return 'mcp-is';
        // Cyrillic transliteration map
        const cyr = {'а':'a','б':'b','в':'v','г':'g','д':'d','е':'e','ё':'e','ж':'zh','з':'z','и':'i','й':'y','к':'k','л':'l','м':'m','н':'n','о':'o','п':'p','р':'r','с':'s','т':'t','у':'u','ф':'f','х':'kh','ц':'ts','ч':'ch','ш':'sh','щ':'shch','ъ':'','ы':'y','ь':'','э':'e','ю':'yu','я':'ya'};
        let s = str.toLowerCase();
        // Transliterate Cyrillic
        let result = '';
        for (const ch of s) { result += cyr[ch] !== undefined ? cyr[ch] : ch; }
        // Keep only latin, digits, hyphens
        result = result.replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '');
        return result || 'mcp-is';
    },

    async render(page, perPage) {
        page = page || 1;
        perPage = perPage || getPerPage();

        const res = await API.get(`/api/ai-keys?page=${page}&per_page=${perPage}`);
        const data = (res && res.ok) ? res.data : { items: [], total: 0, page: 1, per_page: perPage };
        const keys = data.items || [];
        const total = data.total || 0;

        let keyAlert = '';
        if (this.newKey) {
            const keyToShow = this.newKey;
            const mcpUrl = this.newMCPUrl || (location.origin + '/mcp');
            keyAlert = `
                <div class="card" style="border:2px solid #27ae60">
                    <div class="card__title" style="color:#27ae60">Ключ создан!</div>
                    <div class="card__desc">Скопируйте данные СЕЙЧАС. Ключ больше не будет отображаться.</div>
                    <div style="margin-bottom:12px">
                        <label style="font-weight:600;font-size:13px;color:#555">🔑 API Key</label>
                        <div class="key-display">${keyToShow}</div>
                    </div>
                    <div style="margin-bottom:12px">
                        <label style="font-weight:600;font-size:13px;color:#555">🌐 MCP Server URL</label>
                        <div class="key-display">${mcpUrl}</div>
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
        }

        let rows = '';
        if (keys.length === 0) {
            rows = '<tr><td colspan="5" class="empty">Ключи не созданы</td></tr>';
        } else {
            keys.forEach(k => {
                const pdnBtn = k.pdn_enabled
                    ? '<button class="btn btn--ghost btn--sm" style="color:#27ae60;border-color:#27ae60" onclick="AIKeys.togglePDN(\'' + k.id + '\', false)">✅ ПДн</button>'
                    : '<button class="btn btn--ghost btn--sm" onclick="AIKeys.togglePDN(\'' + k.id + '\', true)">☐ ПДн</button>';
                rows += `
                    <tr>
                        <td data-label="Название">${k.name || 'Без названия'}</td>
                        <td data-label="Хэш"><code>${k.key_hash.substring(0, 16)}...</code></td>
                        <td data-label="ПДн">${pdnBtn}</td>
                        <td data-label="Создан">${k.created_at ? k.created_at.split('T')[0] : '—'}</td>
                        <td data-label="Использован">${k.last_used_at ? k.last_used_at.split('T')[0] : '—'}</td>
                        <td data-label="">
                            <div style="display:flex;gap:4px;flex-wrap:nowrap">
                                <button class="btn btn--ghost btn--sm" onclick="AIKeys.edit('${k.id}', '${k.name}', ${k.pdn_enabled})">✏️ Ред.</button>
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
                    <p class="desc">API-ключи для подключения AI-агентов к магазину</p>
                </div>
                <button class="btn btn--primary" onclick="AIKeys.showForm()">➕ Создать ключ</button>
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
                            <th>Хэш ключа</th>
                            <th>ПДн</th>
                            <th>Создан</th>
                            <th>Использован</th>
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

    edit(id, name, pdnEnabled) {
        document.getElementById('edit-aikey-id').value = id;
        document.getElementById('edit-aikey-name').value = name;
        document.getElementById('edit-aikey-pdn').checked = pdnEnabled;

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
                this.newKeyName = name;
                this.newMCPUrl = res.data.mcp_url;
                App.navigate('/ai-keys');
            } else {
                document.getElementById('aikey-alert').innerHTML =
                    '<div class="alert alert--error">' + (res ? res.data.error : 'Ошибка') + '</div>';
            }
        });

        // Edit form binding
        const editForm = document.getElementById('edit-aikey-form');
        if (editForm) {
            editForm.addEventListener('submit', async (e) => {
                e.preventDefault();
                const id = document.getElementById('edit-aikey-id').value;
                const name = document.getElementById('edit-aikey-name').value.trim();
                if (!name) {
                    document.getElementById('aikey-edit-alert').innerHTML =
                        '<div class="alert alert--error">Название обязательно</div>';
                    return;
                }
                const pdn = document.getElementById('edit-aikey-pdn').checked;
                const res = await API.put('/api/ai-keys/' + id, { name, pdn_enabled: pdn });
                if (res && res.ok) {
                    App.navigate('/ai-keys');
                } else {
                    document.getElementById('aikey-edit-alert').innerHTML =
                        '<div class="alert alert--error">' + (res ? res.data.error : 'Ошибка') + '</div>';
                }
            });
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
        }
    },

    copyConfig(client) {
        const keyEl = document.querySelectorAll('.key-display')[0];
        const urlEl = document.querySelectorAll('.key-display')[1];
        if (!keyEl || !urlEl) return;
        const key = keyEl.textContent.trim();
        const url = urlEl.textContent.trim();
        const keyName = AIKeys.slugify(AIKeys.newKeyName || '');
        let config;
        if (client === 'claude') {
            config = JSON.stringify({mcpServers:{[keyName]:{url:url,headers:{"Authorization":"Bearer "+key}}}}, null, 2);
        } else if (client === 'vscode') {
            config = JSON.stringify({servers:{[keyName]:{type:"http",url:url,headers:{"Authorization":"Bearer "+key}}}}, null, 2);
        }
        if (config) {
            navigator.clipboard.writeText(config).then(() => {
                showToast('Конфигурация скопирована!');
            }).catch(() => {});
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
                    select.innerHTML += `<option value="${t.id}">${t.name}</option>`;
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
        const keyName = AIKeys.slugify(AIKeys.newKeyName || '');
        const config = type.config_template
            .replace(/\{\{KEY_NAME\}\}/g, keyName)
            .replace(/\{\{URL\}\}/g, url)
            .replace(/\{\{API_KEY\}\}/g, apiKey);

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

        const blob = new Blob([code.textContent], { type: 'text/plain' });
        const a = document.createElement('a');
        a.href = URL.createObjectURL(blob);
        a.download = type.file_name;
        a.click();
        URL.revokeObjectURL(a.href);
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
