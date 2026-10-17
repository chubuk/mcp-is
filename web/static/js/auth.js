const Auth = {
    renderLogin() {
        return `
            <div class="auth-page">
                <div class="auth-card">
                    <h1>MCP IS</h1>
                    <p class="subtitle">Вход в панель управления</p>
                    <div id="alert"></div>
                    <form id="login-form">
                        <div class="form-group">
                            <label>Email, телефон или никнейм</label>
                            <input type="text" id="login" required autocomplete="username">
                        </div>
                        <div class="form-group">
                            <label>Пароль</label>
                            <input type="password" id="password" required autocomplete="current-password">
                        </div>
                        <button type="submit" class="btn btn--primary" style="width:100%">🔑 Войти</button>
                    </form>
                    <p class="link"><a href="/forgot-password">Забыли пароль?</a></p>
                    <p class="link">Нет аккаунта? <a href="/register">Зарегистрироваться</a></p>
                </div>
            </div>`;
    },

    renderRegister() {
        return `
            <div class="auth-page">
                <div class="auth-card">
                    <h1>Регистрация</h1>
                    <p class="subtitle">Создайте аккаунт для управления магазином</p>
                    <div id="alert"></div>
                    <form id="register-form">
                        <div class="form-group">
                            <label>Никнейм *</label>
                            <input type="text" id="reg-nickname" required>
                        </div>
                        <div class="form-group">
                            <label>Email</label>
                            <input type="email" id="reg-email">
                        </div>
                        <div class="form-group">
                            <label>Телефон</label>
                            <input type="tel" id="reg-phone">
                        </div>
                        <div class="form-group">
                            <label>Имя</label>
                            <input type="text" id="reg-firstname">
                        </div>
                        <div class="form-group">
                            <label>Фамилия</label>
                            <input type="text" id="reg-lastname">
                        </div>
                        <div class="form-group">
                            <label>Пароль *</label>
                            <input type="password" id="reg-password" required minlength="6" autocomplete="new-password">
                        </div>
                        <button type="submit" class="btn btn--primary" style="width:100%">➕ Зарегистрироваться</button>
                    </form>
                    <p class="link">Уже есть аккаунт? <a href="/login">Войти</a></p>
                </div>
            </div>`;
    },

    bindLogin() {
        const pending = new URLSearchParams(location.search).get('pending_install');
        if (pending) {
            const alert = document.getElementById('alert');
            if (alert) {
                alert.innerHTML = '<div class="alert alert--info">Магазин InSales ожидает привязки. Войдите в существующий аккаунт, чтобы подключить его.</div>';
            }
        }
        document.getElementById('login-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            const login = document.getElementById('login').value;
            const password = document.getElementById('password').value;

            const res = await API.post('/api/login', { login, password });
            if (res && res.ok) {
                API.setToken(res.data.token);
                localStorage.setItem('nickname', res.data.nickname);
                localStorage.setItem('role', res.data.role);
                localStorage.setItem('user_id', res.data.user_id);

                const pending = new URLSearchParams(location.search).get('pending_install');
                if (pending) {
                    Auth.linkPendingInstall(pending);
                    return;
                }
                if (res.data.require_password_change) {
                    App.navigate('/force-change-password');
                } else if (res.data.require_email) {
                    App.navigate('/settings');
                } else {
                    App.navigate('/dashboard');
                }
            } else {
                Auth.showError(res ? res.data.error : 'Ошибка сервера');
            }
        });
    },

    async linkPendingInstall(insalesId) {
        const res = await API.post('/api/shops/link-pending', { insales_id: insalesId });
        if (res && res.ok) {
            if (res.data.api_key) {
                localStorage.setItem('fresh_api_key', res.data.api_key);
            }
            App.navigate('/onboarding');
        } else {
            Auth.showError(res ? res.data.error : 'Не удалось привязать магазин');
            App.navigate('/dashboard');
        }
    },

    renderOnboarding() {
        return `
            <div class="auth-page">
                <div class="auth-card">
                    <h1>Подключение MCP</h1>
                    <p class="subtitle">Магазин установлен из InSales</p>
                    <div id="alert"></div>
                    <div id="onboarding-body">
                        <p>Обмениваем код доступа…</p>
                    </div>
                    <p class="link"><a href="/dashboard">Перейти в панель</a></p>
                </div>
            </div>`;
    },

    bindOnboarding() {
        const params = new URLSearchParams(location.search);
        const code = params.get('code');
        const body = document.getElementById('onboarding-body');

        const showKey = (apiKey, mcpUrl) => {
            const url = mcpUrl || (location.origin + '/mcp');
            body.innerHTML = `
                <p><strong>MCP-ключ</strong> (сохраните — показывается один раз):</p>
                <div class="form-group">
                    <input type="text" id="shown-key" readonly value="${apiKey}">
                </div>
                <p><strong>URL сервера:</strong> <code>${url}</code></p>
                <p class="subtitle">Добавьте ключ в свой MCP-клиент (Claude Desktop, Cursor, MiMoCode и др.).</p>
                <button class="btn btn--primary" id="copy-key">📋 Скопировать ключ</button>`;
            const btn = document.getElementById('copy-key');
            if (btn) btn.addEventListener('click', () => {
                navigator.clipboard.writeText(apiKey);
                Auth.showSuccess('Ключ скопирован');
            });
        };

        const stored = localStorage.getItem('fresh_api_key');
        if (!code && stored) {
            showKey(stored, location.origin + '/mcp');
            localStorage.removeItem('fresh_api_key');
            return;
        }
        if (!code) {
            body.innerHTML = '<p>Нет кода подключения. Откройте приложение из бэк-офиса InSales.</p>';
            return;
        }
        API.post('/api/exchange-code', { code }).then((res) => {
            if (res && res.ok) {
                API.setToken(res.data.token);
                localStorage.setItem('nickname', res.data.nickname);
                localStorage.setItem('role', res.data.role);
                localStorage.setItem('user_id', res.data.user_id);
                showKey(res.data.api_key || '(ключ недоступен — создайте новый в разделе «Ключи»)');
            } else {
                body.innerHTML = '<p>Код недействителен или истёк. Откройте приложение из InSales ещё раз.</p>';
            }
        });
    },

    bindRegister() {
        document.getElementById('register-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            const body = {
                nickname: document.getElementById('reg-nickname').value,
                email: document.getElementById('reg-email').value,
                phone: document.getElementById('reg-phone').value,
                first_name: document.getElementById('reg-firstname').value,
                last_name: document.getElementById('reg-lastname').value,
                password: document.getElementById('reg-password').value,
            };

            const res = await API.post('/api/register', body);
            if (res && res.ok) {
                Auth.showSuccess('Регистрация успешна! Теперь войдите.');
                setTimeout(() => App.navigate('/login'), 1500);
            } else if (res && res.status === 403) {
                Auth.showError('Регистрация закрыта. Обратитесь к администратору.');
            } else {
                Auth.showError(res ? res.data.error : 'Ошибка сервера');
            }
        });
    },

    showError(msg) {
        const el = document.getElementById('alert');
        if (el) el.innerHTML = '<div class="alert alert--error">' + msg + '</div>';
    },

    showSuccess(msg) {
        const el = document.getElementById('alert');
        if (el) el.innerHTML = '<div class="alert alert--success">' + msg + '</div>';
    },

    renderForgotPassword() {
        return `
            <div class="auth-page">
                <div class="auth-card">
                    <h1>Восстановление пароля</h1>
                    <p class="subtitle">Введите email для получения ссылки сброса</p>
                    <div id="alert"></div>
                    <form id="forgot-form">
                        <div class="form-group">
                            <label>Email</label>
                            <input type="email" id="forgot-email" required>
                        </div>
                        <button type="submit" class="btn btn--primary" style="width:100%">📧 Отправить</button>
                    </form>
                    <p class="link"><a href="/login">← Вернуться к входу</a></p>
                </div>
            </div>`;
    },

    bindForgotPassword() {
        document.getElementById('forgot-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            const email = document.getElementById('forgot-email').value;
            const res = await API.post('/api/forgot-password', { email });
            if (res && res.ok) {
                Auth.showSuccess(res.data.message || 'Письмо отправлено (проверьте логи сервера)');
            } else {
                Auth.showError(res ? res.data.error : 'Ошибка сервера');
            }
        });
    },

    renderResetPassword() {
        const params = new URLSearchParams(location.search);
        const token = params.get('token') || '';
        return `
            <div class="auth-page">
                <div class="auth-card">
                    <h1>Новый пароль</h1>
                    <p class="subtitle">Введите новый пароль</p>
                    <div id="alert"></div>
                    <form id="reset-form">
                        <input type="hidden" id="reset-token" value="${token}">
                        <div class="form-group">
                            <label>Новый пароль</label>
                            <input type="password" id="reset-password" required minlength="6" autocomplete="new-password">
                        </div>
                        <div class="form-group">
                            <label>Повторите пароль</label>
                            <input type="password" id="reset-password2" required minlength="6" autocomplete="new-password">
                        </div>
                        <button type="submit" class="btn btn--primary" style="width:100%">💾 Сохранить</button>
                    </form>
                    <p class="link"><a href="/login">← Вернуться к входу</a></p>
                </div>
            </div>`;
    },

    bindResetPassword() {
        document.getElementById('reset-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            const token = document.getElementById('reset-token').value;
            const password = document.getElementById('reset-password').value;
            const password2 = document.getElementById('reset-password2').value;

            if (password !== password2) {
                Auth.showError('Пароли не совпадают');
                return;
            }

            const res = await API.post('/api/reset-password', { token, new_password: password });
            if (res && res.ok) {
                Auth.showSuccess('Пароль обновлён! Теперь войдите.');
                setTimeout(() => App.navigate('/login'), 2000);
            } else {
                Auth.showError(res ? res.data.error : 'Ошибка сервера');
            }
        });
    },

    renderForceChangePassword() {
        return `
            <div class="auth-page">
                <div class="auth-card">
                    <h1>Смена пароля</h1>
                    <p class="subtitle">Вы вошли с дефолтным паролем. Для безопасности необходимо задать новый пароль.</p>
                    <div id="alert"></div>
                    <form id="force-change-form">
                        <div class="form-group">
                            <label>Новый пароль</label>
                            <input type="password" id="force-new-password" required minlength="8" autocomplete="new-password">
                            <small style="color:#888">Минимум 8 символов, заглавная буква и цифра</small>
                        </div>
                        <div class="form-group">
                            <label>Повторите пароль</label>
                            <input type="password" id="force-new-password2" required minlength="6" autocomplete="new-password">
                        </div>
                        <button type="submit" class="btn btn--primary" style="width:100%">🔑 Сменить пароль</button>
                    </form>
                </div>
            </div>`;
    },

    bindForceChangePassword() {
        document.getElementById('force-change-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            const password = document.getElementById('force-new-password').value;
            const password2 = document.getElementById('force-new-password2').value;

            if (password !== password2) {
                Auth.showError('Пароли не совпадают');
                return;
            }

            if (password.length < 8) {
                Auth.showError('Пароль должен быть не менее 8 символов');
                return;
            }

            const res = await API.post('/api/force-change-password', {
                new_password: password,
            });

            if (res && res.ok) {
                Auth.showSuccess('Пароль обновлён! Теперь войдите с новым паролем.');
                setTimeout(() => App.navigate('/login'), 2000);
            } else {
                Auth.showError(res ? res.data.error : 'Ошибка сервера');
            }
        });
    }
};
