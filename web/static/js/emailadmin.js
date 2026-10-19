const EmailAdmin = {
    quill: null,
    currentKey: null,
    mode: 'visual', // visual | html | text

    render() {
        return `
            <div class="page">
                <h1>Почта</h1>
                <div class="card" id="smtp-settings-card">
                    <div class="card__title">SMTP-сервер</div>
                    <div id="smtp-settings-alert"></div>
                    <div id="smtp-settings-body">Загрузка…</div>
                </div>
                <div class="card">
                    <div class="card__title">DNS-записи домена (для доставки писем)</div>
                    <div class="alert alert--info">
                        Чтобы письма не попадали в спам, в DNS домена (например <code>mcp-is.ru</code>)
                        должны быть <strong>SPF</strong>, <strong>DKIM</strong> и <strong>DMARC</strong>.
                        Записи добавляются в панели регистратора домена (nameself, Reg.ru, Cloudflare и т.п.).
                    </div>
                    <table style="width:100%;font-size:13px">
                        <tr><th style="text-align:left;padding:4px">Тип</th><th style="text-align:left;padding:4px">Имя</th><th style="text-align:left;padding:4px">Значение</th><th style="text-align:left;padding:4px">Зачем</th></tr>
                        <tr>
                            <td style="padding:4px"><code>TXT</code></td>
                            <td style="padding:4px"><code>@</code></td>
                            <td style="padding:4px"><code>v=spf1 include:_spf.mail.ru ~all</code></td>
                            <td style="padding:4px">разрешить отправку через SMTP провайдера</td>
                        </tr>
                        <tr>
                            <td style="padding:4px"><code>TXT</code></td>
                            <td style="padding:4px"><code>_dmarc</code></td>
                            <td style="padding:4px"><code>v=DMARC1; p=none; rua=mailto:dmarc@ваш-домен</code></td>
                            <td style="padding:4px">политика и отчёты; позже <code>p=quarantine</code></td>
                        </tr>
                        <tr>
                            <td style="padding:4px"><code>TXT</code></td>
                            <td style="padding:4px"><code>&lt;selector&gt;._domainkey</code></td>
                            <td style="padding:4px"><code>v=DKIM1; k=rsa; p=…</code></td>
                            <td style="padding:4px">подпись писем — взять в панели почтового сервиса</td>
                        </tr>
                        <tr>
                            <td style="padding:4px"><code>MX</code></td>
                            <td style="padding:4px"><code>@</code></td>
                            <td style="padding:4px"><code>mxs.mail.ru.</code> (пример)</td>
                            <td style="padding:4px">если ящик в домене (не обязательно для только-отправки)</td>
                        </tr>
                    </table>
                    <p style="font-size:13px;color:#555;margin-top:10px">
                        <strong>Без DNS:</strong> можно отправлять с адреса провайдера
                        (<code>SMTP_FROM</code> вида <code>user@mail.ru</code>) — письма дойдут, но
                        «от кого» будет не ваш домен.
                        <br><strong>С доменом:</strong> <code>SMTP_FROM=MCP IS &lt;noreply@ваш-домен&gt;</code>
                        и все записи выше.
                        <br><strong>Mail.ru / biz.mail.ru:</strong> SMTP-пароль —
                        <strong>пароль приложений</strong> (Настройки → Безопасность → Пароли для внешних приложений),
                        не основной пароль ящика. Ошибка <code>535</code> = нужен app-password.
                        <br>Подробнее и скрипт проверки: <code>docs/setup/EMAIL-DNS.md</code>,
                        <code>./scripts/check-email-dns.sh</code>
                    </p>
                </div>
                <div class="card">
                    <div class="card__title">Шаблоны писем</div>
                    <div id="tpl-alert"></div>
                    <div class="form-group">
                        <label>Шаблон</label>
                        <select id="tpl-key" onchange="EmailAdmin.loadTemplate(this.value)"></select>
                    </div>
                    <div class="form-group">
                        <label>Тема письма</label>
                        <input type="text" id="tpl-subject" placeholder="Тема — MCP IS">
                    </div>
                    <div class="form-group">
                        <label>Тело письма</label>
                        <div style="margin-bottom:8px">
                            <button type="button" class="btn btn--sm btn--secondary" onclick="EmailAdmin.setMode('visual')">Визуально</button>
                            <button type="button" class="btn btn--sm btn--secondary" onclick="EmailAdmin.setMode('html')">HTML</button>
                            <button type="button" class="btn btn--sm btn--secondary" onclick="EmailAdmin.setMode('text')">Plain text</button>
                        </div>
                        <div style="margin:8px 0;font-size:13px">
                            <span style="color:#666;margin-right:6px">Вставить переменную:</span>
                            <button type="button" class="btn btn--sm"
                                onclick="EmailAdmin.insertVar('Name')"
                                title="Имя получателя&#10;Пример: Иван Петров">{{Name}}</button>
                            <button type="button" class="btn btn--sm"
                                onclick="EmailAdmin.insertVar('ResetURL')"
                                title="Ссылка сброса пароля (действует 1 час)&#10;Пример: https://mcp-is.ru/reset-password?token=abc123">{{ResetURL}}</button>
                            <button type="button" class="btn btn--sm"
                                onclick="EmailAdmin.insertVar('VerifyURL')"
                                title="Ссылка подтверждения email (24 часа)&#10;Пример: https://mcp-is.ru/verify-email?token=abc123">{{VerifyURL}}</button>
                            <button type="button" class="btn btn--sm"
                                onclick="EmailAdmin.insertVar('MCPURL')"
                                title="Ссылка на панель MCP IS&#10;Пример: https://mcp-is.ru/">{{MCPURL}}</button>
                        </div>
                        <div id="tpl-quill" style="background:#fff;min-height:220px"></div>
                        <textarea id="tpl-html" class="form-control" rows="12" style="display:none;font-family:monospace"></textarea>
                        <textarea id="tpl-text" class="form-control" rows="10" style="display:none;font-family:monospace"></textarea>
                    </div>
                    <div class="form-group">
                        <button class="btn btn--primary" onclick="EmailAdmin.save()">💾 Сохранить</button>
                        <button class="btn btn--secondary" onclick="EmailAdmin.preview()">👁 Предпросмотр</button>
                        <button class="btn btn--secondary" onclick="EmailAdmin.test()">📧 Тест-письмо</button>
                        <button class="btn btn--secondary" onclick="EmailAdmin.reset()">↩ Сбросить</button>
                    </div>
                    <div class="form-group">
                        <label>Email для теста</label>
                        <input type="email" id="tpl-test-email" placeholder="test@example.com">
                    </div>
                    <div id="tpl-preview-wrap" style="display:none">
                        <label>Предпросмотр</label>
                        <iframe id="tpl-preview" sandbox="" style="width:100%;height:360px;border:1px solid #ddd;background:#fff"></iframe>
                    </div>
                </div>
            </div>`;
    },

    async bind() {
        this.initQuill();
        await this.loadSMTPSettings();
        await this.loadList();
    },

    initQuill() {
        if (this.quill) return;
        this.quill = new Quill('#tpl-quill', {
            theme: 'snow',
            modules: {
                toolbar: [
                    [{ header: [2, 3, false] }],
                    ['bold', 'italic', 'underline'],
                    ['link'],
                    [{ list: 'ordered' }, { list: 'bullet' }],
                    ['clean'],
                ],
            },
            placeholder: 'Текст письма…',
        });
    },

    // Full HTML templates keep <head>/<style>; Quill edits only <body> and we restore chrome on save.
    extractBody(html) {
        if (!html) return '';
        const m = String(html).match(/<body[^>]*>([\s\S]*)<\/body>/i);
        return m ? m[1] : String(html);
    },

    wrapBody(html, body) {
        if (!html || !/<body[\s>]/i.test(html)) return body;
        return String(html).replace(/(<body[^>]*>)[\s\S]*(<\/body>)/i, (_all, open, close) => open + body + close);
    },

    setMode(mode) {
        const currentHTML = this.getHTML();
        if (currentHTML) this._htmlBody = currentHTML;
        this.mode = mode;
        document.getElementById('tpl-quill').style.display = mode === 'visual' ? '' : 'none';
        document.getElementById('tpl-html').style.display = mode === 'html' ? '' : 'none';
        document.getElementById('tpl-text').style.display = mode === 'text' ? '' : 'none';
        if (mode === 'html') {
            document.getElementById('tpl-html').value = this._htmlBody || '';
        }
        if (mode === 'visual') {
            this.loadIntoQuill(this.extractBody(this._htmlBody));
        }
        if (mode === 'text') {
            document.getElementById('tpl-text').value = this._textBody || '';
        }
    },

    // Quill ignores root.innerHTML for its Delta — must paste via clipboard.
    loadIntoQuill(bodyHTML) {
        if (!this.quill) return;
        this.quill.enable(false);
        this.quill.setContents([]);
        const clean = (bodyHTML || '').trim();
        if (clean && clean !== '<p><br></p>') {
            try {
                this.quill.clipboard.dangerouslyPasteHTML(clean);
            } catch (e) {
                this.quill.setText(clean.replace(/<[^>]+>/g, ' '));
            }
        }
        this.quill.enable(true);
    },

    getHTML() {
        if (this.mode === 'html') {
            this._htmlBody = document.getElementById('tpl-html').value;
            return this._htmlBody;
        }
        if (!this.quill) return this._htmlBody || '';
        const body = this.quill.root.innerHTML;
        const empty = !body || body === '<p><br></p>' || body === '<p></p>';
        // safety: never wipe a non-empty template because the editor is blank
        if (empty && this.extractBody(this._htmlBody).trim()) {
            return this._htmlBody;
        }
        this._htmlBody = this.wrapBody(this._htmlBody, body);
        return this._htmlBody;
    },

    setHTML(html) {
        this._htmlBody = html || '';
        this.loadIntoQuill(this.extractBody(this._htmlBody));
        const ta = document.getElementById('tpl-html');
        if (ta) ta.value = this._htmlBody;
    },

    insertVar(name) {
        const token = '{{.' + name + '}}';
        if (this.mode === 'html') {
            const ta = document.getElementById('tpl-html');
            const pos = ta.selectionStart || ta.value.length;
            ta.value = ta.value.slice(0, pos) + token + ta.value.slice(pos);
            return;
        }
        if (this.quill) {
            const range = this.quill.getSelection(true);
            this.quill.insertText(range ? range.index : 0, token, 'user');
        }
    },

    async loadList() {
        const res = await API.get('/api/admin/email/templates');
        const sel = document.getElementById('tpl-key');
        if (!res || !res.ok) {
            sel.innerHTML = '<option value="">Ошибка загрузки</option>';
            return;
        }
        sel.innerHTML = (res.data.items || [])
            .map((t) => `<option value="${t.key}">${t.name} (${t.key})</option>`)
            .join('');
        if (res.data.items && res.data.items.length) {
            await this.loadTemplate(res.data.items[0].key);
        }
    },

    async loadTemplate(key) {
        this.currentKey = key;
        const res = await API.get('/api/admin/email/templates/' + encodeURIComponent(key));
        if (!res || !res.ok) return;
        document.getElementById('tpl-subject').value = res.data.subject || '';
        this._textBody = res.data.text_body || '';
        this._htmlBody = res.data.html_body || '';
        if (!this.quill) this.initQuill();
        this.setHTML(this._htmlBody);
        this.mode = 'visual';
        document.getElementById('tpl-quill').style.display = '';
        document.getElementById('tpl-html').style.display = 'none';
        document.getElementById('tpl-text').style.display = 'none';
        document.getElementById('tpl-text').value = this._textBody;
        this.hideAlert();
    },

    collect() {
        return {
            subject: document.getElementById('tpl-subject').value,
            html_body: this.getHTML(),
            text_body: this.mode === 'text'
                ? document.getElementById('tpl-text').value
                : (this._textBody || ''),
        };
    },

    async save() {
        if (!this.currentKey) return;
        const body = this.collect();
        const res = await API.put('/api/admin/email/templates/' + encodeURIComponent(this.currentKey), body);
        if (res && res.ok) {
            this.alert('Шаблон сохранён', 'success');
        } else {
            this.alert((res && res.data && res.data.error) || 'Ошибка сохранения', 'error');
        }
    },

    async preview() {
        const body = this.collect();
        const res = await API.post('/api/admin/email/templates/' + encodeURIComponent(this.currentKey) + '/preview', {
            subject: body.subject,
            html_body: body.html_body,
            sample: { Name: 'Пользователь', ResetURL: 'https://mcp-is.ru/reset-password?token=sample', VerifyURL: 'https://mcp-is.ru/verify-email?token=sample', MCPURL: 'https://mcp-is.ru/' },
        });
        const wrap = document.getElementById('tpl-preview-wrap');
        const frame = document.getElementById('tpl-preview');
        wrap.style.display = 'block';
        if (res && res.ok) {
            frame.srcdoc = res.data.html || '';
        } else {
            frame.srcdoc = '<p style="color:#c00">Ошибка рендера шаблона</p>';
            this.alert((res && res.data && res.data.error) || 'Ошибка preview', 'error');
        }
    },

    async test() {
        const email = document.getElementById('tpl-test-email').value.trim();
        if (!email) {
            this.alert('Введите email для теста', 'error');
            return;
        }
        await this.save();
        const res = await API.post('/api/admin/email/templates/' + encodeURIComponent(this.currentKey) + '/test', { email });
        if (res && res.ok) {
            this.alert('Тестовое письмо отправлено на ' + email, 'success');
        } else {
            this.alert((res && res.data && res.data.error) || 'Ошибка отправки', 'error');
        }
    },

    async reset() {
        if (!this.currentKey) return;
        if (!confirm('Сбросить шаблон к значению по умолчанию?')) return;
        const res = await API.post('/api/admin/email/templates/' + encodeURIComponent(this.currentKey) + '/reset', {});
        if (res && res.ok) {
            await this.loadTemplate(this.currentKey);
            this.alert('Шаблон сброшен', 'success');
        } else {
            this.alert('Ошибка сброса', 'error');
        }
    },

    alert(msg, type) {
        document.getElementById('tpl-alert').innerHTML =
            `<div class="alert alert--${type === 'success' ? 'success' : 'error'}">${msg}</div>`;
    },
    hideAlert() {
        const el = document.getElementById('tpl-alert');
        if (el) el.innerHTML = '';
    },

    // --- SMTP settings ---
    async loadSMTPSettings() {
        const res = await API.get('/api/admin/smtp/settings');
        const body = document.getElementById('smtp-settings-body');
        if (!res || !res.ok) {
            body.innerHTML = '<div class="alert alert--error">Не удалось загрузить настройки SMTP</div>';
            return;
        }
        const s = res.data;
        body.innerHTML = `
            <div class="form-group" style="display:flex;align-items:center;gap:10px;flex-wrap:wrap;padding:10px 12px;background:#f7f9fc;border:1px solid #d9e2ef;border-radius:6px;margin-bottom:14px">
                <label style="margin:0;font-weight:600;cursor:pointer;display:flex;align-items:center;gap:8px">
                    <input type="checkbox" id="smtp-enabled" ${s.enabled ? 'checked' : ''} style="width:18px;height:18px;cursor:pointer">
                    Включить отправку почты
                </label>
                <span style="font-size:12px;color:#666">— письма сброса пароля и уведомления</span>
                <span id="smtp-enabled-status" style="font-size:12px;font-weight:600;color:${s.enabled ? '#1a7f37' : '#999'}">
                    ${s.enabled ? 'Отправка: включена' : 'Отправка: выключена'}
                </span>
            </div>
            <div class="form-row" style="display:flex;gap:12px;flex-wrap:wrap">
                <div class="form-group" style="flex:2">
                    <label for="smtp-host">Хост</label>
                    <input type="text" id="smtp-host" value="${s.host || ''}" placeholder="smtp.mail.ru">
                </div>
                <div class="form-group" style="flex:1">
                    <label for="smtp-port">Порт</label>
                    <input type="text" id="smtp-port" value="${s.port || ''}" placeholder="587">
                    <small style="color:#888;font-size:11px">587 = STARTTLS, 465 = SSL</small>
                </div>
            </div>
            <div class="form-row" style="display:flex;gap:12px;flex-wrap:wrap">
                <div class="form-group" style="flex:1">
                    <label for="smtp-user">Пользователь</label>
                    <input type="text" id="smtp-user" value="${s.user || ''}" placeholder="user@mail.ru">
                </div>
                <div class="form-group" style="flex:1">
                    <label for="smtp-pass">Пароль${s.has_password ? ' <small style="color:#1a7f37">(сохранён ✓ — пустое = не менять)</small>' : ''}</label>
                    <input type="password" id="smtp-pass" data-pw-no-reveal value="" autocomplete="new-password" placeholder="${s.has_password ? '•••••• (сохранён)' : 'пароль SMTP'}">
                </div>
            </div>
            <div class="form-row" style="display:flex;gap:12px;flex-wrap:wrap">
                <div class="form-group" style="flex:2">
                    <label for="smtp-from">От кого</label>
                    <input type="text" id="smtp-from" value="${s.from || ''}" placeholder="MCP IS <user@mail.ru>">
                </div>
                <div class="form-group" style="flex:1">
                    <label for="smtp-tls">TLS</label>
                    <select id="smtp-tls">
                        <option value="starttls" ${s.tls_mode === 'starttls' || !s.tls_mode ? 'selected' : ''}>STARTTLS</option>
                        <option value="ssl" ${s.tls_mode === 'ssl' ? 'selected' : ''}>SSL</option>
                        <option value="none" ${s.tls_mode === 'none' ? 'selected' : ''}>Без TLS</option>
                    </select>
                </div>
            </div>
            <button class="btn btn--primary" onclick="EmailAdmin.saveSMTP()">💾 Сохранить SMTP</button>
        `;
    },

    async saveSMTP() {
        const body = {
            enabled: document.getElementById('smtp-enabled').checked,
            host: document.getElementById('smtp-host').value,
            port: document.getElementById('smtp-port').value,
            user: document.getElementById('smtp-user').value,
            from: document.getElementById('smtp-from').value,
            tls_mode: document.getElementById('smtp-tls').value,
        };
        const pass = document.getElementById('smtp-pass').value;
        if (pass) body.pass = pass;
        const res = await API.put('/api/admin/smtp/settings', body);
        const alertEl = document.getElementById('smtp-settings-alert');
        if (res && res.ok) {
            alertEl.innerHTML = '<div class="alert alert--success">Настройки SMTP сохранены</div>';
            await this.loadSMTPSettings();
        } else {
            alertEl.innerHTML = `<div class="alert alert--error">${(res && res.data && res.data.error) || 'Ошибка'}</div>`;
        }
    },
};
