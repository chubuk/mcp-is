const AIKeys = {
    newKey: null,
    newMCPUrl: null,

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
            const mcpUrl = this.newMCPUrl || 'https://localhost:8443/mcp';
            const configClaude = JSON.stringify({mcpServers:{"mcp-is":{url:mcpUrl,headers:{"Authorization":"Bearer "+keyToShow}}}}, null, 2);
            const configVSCode = JSON.stringify({servers:{"mcp-is":{type:"http",url:mcpUrl,headers:{"Authorization":"Bearer "+keyToShow}}}}, null, 2);
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
                    <details>
                        <summary style="cursor:pointer;color:#4a6cf7;font-size:13px">📦 Конфигурация для AI-клиента</summary>
                        <div style="margin-top:8px">
                            <div style="display:flex;gap:4px;margin-bottom:8px">
                                <button class="btn btn--ghost btn--sm" onclick="AIKeys.copyConfig('claude',this)">Claude Desktop</button>
                                <button class="btn btn--ghost btn--sm" onclick="AIKeys.copyConfig('vscode',this)">VS Code</button>
                            </div>
                            <pre style="background:#f8f9fa;padding:12px;border-radius:6px;font-size:12px;overflow-x:auto;white-space:pre-wrap">${configClaude}</pre>
                        </div>
                    </details>
                </div>`;
            this.newKey = null;
            this.newMCPUrl = null;
        }

        let rows = '';
        if (keys.length === 0) {
            rows = '<tr><td colspan="5" class="empty">Ключи не созданы</td></tr>';
        } else {
            keys.forEach(k => {
                const pdnBadge = k.pdn_enabled
                    ? '<span style="color:#27ae60;font-weight:600;cursor:pointer" onclick="AIKeys.togglePDN(\'' + k.id + '\', false)">ПДн ✅</span>'
                    : '<span style="color:#888;cursor:pointer" onclick="AIKeys.togglePDN(\'' + k.id + '\', true)">ПДн —</span>';
                rows += `
                    <tr>
                        <td data-label="Название">${k.name || 'Без названия'}</td>
                        <td data-label="Хэш"><code>${k.key_hash.substring(0, 16)}...</code></td>
                        <td data-label="ПДн">${pdnBadge}</td>
                        <td data-label="Создан">${k.created_at ? k.created_at.split('T')[0] : '—'}</td>
                        <td data-label="Использован">${k.last_used_at ? k.last_used_at.split('T')[0] : '—'}</td>
                        <td data-label="">
                            <button class="btn btn--danger btn--sm" onclick="AIKeys.delete('${k.id}')">🗑️ Отозвать</button>
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
                            <small style="color:#888">Маскировка email, телефона, ФИО для AI</small>
                        </div>
                        <button type="submit" class="btn btn--primary">💾 Создать ключ</button>
                        <button type="button" class="btn btn--ghost" onclick="AIKeys.hideForm()">✖️ Отмена</button>
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
                this.newMCPUrl = res.data.mcp_url;
                App.navigate('/ai-keys');
            } else {
                document.getElementById('aikey-alert').innerHTML =
                    '<div class="alert alert--error">' + (res ? res.data.error : 'Ошибка') + '</div>';
            }
        });
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
        let config;
        if (client === 'claude') {
            config = JSON.stringify({mcpServers:{"mcp-is":{url:url,headers:{"Authorization":"Bearer "+key}}}}, null, 2);
        } else if (client === 'vscode') {
            config = JSON.stringify({servers:{"mcp-is":{type:"http",url:url,headers:{"Authorization":"Bearer "+key}}}}, null, 2);
        }
        if (config) {
            navigator.clipboard.writeText(config).then(() => {
                showToast('Конфигурация скопирована!');
            }).catch(() => {});
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
