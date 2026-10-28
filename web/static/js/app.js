const App = {
    routes: [
        { path: '/login', render: () => Auth.renderLogin(), auth: false, onMount: () => Auth.bindLogin() },
        { path: '/register', render: () => Auth.renderRegister(), auth: false, onMount: () => Auth.bindRegister() },
        { path: '/onboarding', render: () => Auth.renderOnboarding(), auth: false, onMount: () => Auth.bindOnboarding() },
        { path: '/dashboard', render: () => Dashboard.render(), auth: true },
        { path: '/shops', render: () => Shops.render(), auth: true, onMount: () => Shops.bind() },
        { path: '/ai-keys', render: () => AIKeys.render(), auth: true, onMount: () => AIKeys.bind() },
        { path: '/settings', render: () => Settings.render(), auth: true },
        { path: '/forgot-password', render: () => Auth.renderForgotPassword(), auth: false, onMount: () => Auth.bindForgotPassword() },
        { path: '/reset-password', render: () => Auth.renderResetPassword(), auth: false, onMount: () => Auth.bindResetPassword() },
        { path: '/admin/users', render: () => Admin.renderUsers(), auth: true, role: ['admin', 'super_admin'] },
        { path: '/admin/users/:id', render: (params) => Admin.showUser(params.id), auth: true, role: ['admin', 'super_admin'] },
        { path: '/admin/stats', render: () => Admin.renderStats(), auth: true, role: ['admin', 'super_admin'] },
        { path: '/admin/server', render: () => Admin.renderServer(), auth: true, role: ['admin', 'super_admin'], onMount: () => Admin.loadSMTPStatus() },
        { path: '/admin/pdn', render: () => Admin.renderPDN(), auth: true, role: ['admin', 'super_admin'] },
        { path: '/admin/backup', render: () => Admin.renderBackup(), auth: true, role: ['admin', 'super_admin'] },
        { path: '/admin/email', render: () => EmailAdmin.render(), auth: true, role: ['admin', 'super_admin'], onMount: () => EmailAdmin.bind() },
        { path: '/force-change-password', render: () => Auth.renderForceChangePassword(), auth: true, onMount: () => Auth.bindForceChangePassword() },
    ],

    currentPage: null,

    // Единый разбор полей формы. schema: { key: [elementId, type] }, type: bool|text|trim|int
    readFields(schema) {
        const out = {};
        for (const [key, spec] of Object.entries(schema)) {
            const [id, type] = Array.isArray(spec) ? spec : [spec, 'text'];
            const el = document.getElementById(id);
            if (!el) {
                out[key] = type === 'bool' ? false : type === 'int' ? 0 : '';
                continue;
            }
            if (type === 'bool') out[key] = !!el.checked;
            else if (type === 'int') out[key] = parseInt(el.value, 10) || 0;
            else if (type === 'trim') out[key] = String(el.value || '').trim();
            else out[key] = el.value;
        }
        return out;
    },

    check(id) {
        const el = document.getElementById(id);
        return el ? !!el.checked : false;
    },

    // Баннер «период не оплачен/закончился» для app-магазинов. shops — items из /api/shops.
    billingBanner(shops) {
        const unpaid = (shops || []).find(s => {
            if (s.billing_status === 'unpaid') return true;
            return s.billing_until && new Date(s.billing_until) <= new Date();
        });
        if (!unpaid) return '';
        const expired = unpaid.billing_until && new Date(unpaid.billing_until) <= new Date();
        const text = expired
            ? 'Оплаченный период закончился.'
            : 'Текущий период не оплачен.';
        const payUrl = unpaid.billing_pay_url
            || (unpaid.shop_url ? `https://${unpaid.shop_url}/admin/applications` : '');
        const btn = payUrl
            ? `<a class="btn btn--primary btn--sm" href="${payUrl}" target="_blank" rel="noopener">Оплатить / Продлить</a>`
            : '';
        return `<div class="alert alert--error" style="margin:12px 0">
            <strong>${text}</strong> Данные магазина ${unpaid.name || unpaid.shop_url} недоступны для ИИ до оплаты — продлите приложение в админке InSales. ${btn}
        </div>`;
    },

    init() {
        window.addEventListener('popstate', () => this.route());
        window.addEventListener('pageshow', (e) => {
            if (!e.persisted) return;
            if (!API.isLoggedIn()) {
                this.navigate('/login');
                return;
            }
            this.route();
        });
        document.addEventListener('click', (e) => {
            const link = e.target.closest('a[href]');
            if (link) {
                e.preventDefault();
                this.navigate(link.getAttribute('href'));
            }
        });
        this.enhancePasswordInputs(document);
        new MutationObserver(() => this.enhancePasswordInputs(document))
            .observe(document.body, { childList: true, subtree: true });
        this.route();
    },

    // Adds a show/hide toggle to every password input (existing and dynamically rendered)
    enhancePasswordInputs(root) {
        root.querySelectorAll('input[type="password"]:not(.pw-done):not([data-pw-no-reveal])').forEach((input) => {
            input.classList.add('pw-done');
            const wrap = document.createElement('div');
            wrap.className = 'pw-wrap';
            input.parentNode.insertBefore(wrap, input);
            wrap.appendChild(input);
            const btn = document.createElement('button');
            btn.type = 'button';
            btn.className = 'pw-toggle';
            btn.title = 'Показать пароль';
            btn.textContent = '👁';
            btn.addEventListener('click', () => {
                const show = input.type === 'password';
                input.type = show ? 'text' : 'password';
                btn.textContent = show ? '🙈' : '👁';
                btn.title = show ? 'Скрыть пароль' : 'Показать пароль';
            });
            wrap.appendChild(btn);
        });
    },

    navigate(path) {
        history.pushState(null, '', path);
        this.closeNav();
        this.route();
    },

    matchRoute(path) {
        for (const route of this.routes) {
            const parts = route.path.split('/');
            const pathParts = path.split('/');
            if (parts.length !== pathParts.length) continue;
            const params = {};
            let match = true;
            for (let i = 0; i < parts.length; i++) {
                if (parts[i].startsWith(':')) {
                    params[parts[i].slice(1)] = pathParts[i];
                } else if (parts[i] !== pathParts[i]) {
                    match = false;
                    break;
                }
            }
            if (match) return { route, params };
        }
        return null;
    },

    async route() {
        const path = location.pathname;
        const result = this.matchRoute(path);

        if (!result) {
            this.navigate('/login');
            return;
        }

        const { route, params } = result;

        if (route.auth && !API.isLoggedIn()) {
            this.navigate('/login');
            return;
        }

        if (route.role) {
            const userRole = localStorage.getItem('role');
            if (!route.role.includes(userRole)) {
                this.navigate('/dashboard');
                return;
            }
        }

        if (!route.auth && API.isLoggedIn() && path === '/login') {
            this.navigate('/dashboard');
            return;
        }

        this.currentPage = path;
        this.updateNav();

        const header = document.getElementById('header');
        if (route.auth) {
            header.style.display = 'flex';
            document.getElementById('user-nickname').textContent = localStorage.getItem('nickname') || '';

            const role = localStorage.getItem('role');
            document.querySelectorAll('.admin-link').forEach(link => {
                link.style.display = (role === 'admin' || role === 'super_admin') ? 'inline' : 'none';
            });

            const banner = document.getElementById('impersonate-banner');
            if (API.isImpersonating()) {
                banner.style.display = 'block';
                document.getElementById('impersonate-text').textContent =
                    'Вы зашли как ' + localStorage.getItem('nickname');
                header.style.top = '32px';
            } else {
                banner.style.display = 'none';
                header.style.top = '0';
            }
        } else {
            header.style.display = 'none';
        }

        const app = document.getElementById('app');
        app.innerHTML = await route.render(params);

        if (route.onMount) route.onMount();
    },

    updateNav() {
        document.querySelectorAll('.nav-link').forEach(link => {
            link.classList.toggle('active', link.dataset.page === this.currentPage.replace('/', ''));
        });
    },

    logout() {
        if (API.isImpersonating()) {
            this.stopImpersonate();
            return;
        }
        API.clearToken();
        localStorage.removeItem('nickname');
        localStorage.removeItem('role');
        localStorage.removeItem('user_id');
        this.navigate('/login');
    },

    async stopImpersonate() {
        const res = await API.post('/api/admin/stop-impersonate');
        if (res && res.ok) {
            API.stopImpersonate();
            API.setToken(res.data.token);
            localStorage.setItem('user_id', res.data.user_id);
            localStorage.setItem('nickname', res.data.nickname);
            localStorage.setItem('role', res.data.role);
            this.navigate('/admin/users');
        } else {
            alert('Ошибка: ' + (res ? res.data.error : 'Сервер недоступен'));
        }
    },

    toggleNav() {
        document.getElementById('nav').classList.toggle('active');
    },

    closeNav() {
        document.getElementById('nav').classList.remove('active');
    }
};

const Settings = {
    render() {
        const nickname = localStorage.getItem('nickname') || '';
        const email = localStorage.getItem('email') || '';
        return `
            <div class="page-header">
                <div>
                    <h1>Настройки</h1>
                    <p class="desc">Управление профилем пользователя</p>
                </div>
            </div>
            <div class="card">
                <div class="card__title">Профиль</div>
                <div class="form-group">
                    <label>Никнейм</label>
                    <input type="text" value="${nickname}" disabled>
                </div>
                <div class="form-group">
                    <label>Email</label>
                    <input type="email" value="${email}" disabled>
                    <small style="color:#888">Смена e-mail — только через администратора</small>
                </div>
            </div>
            <div class="card">
                <div class="card__title">Смена пароля</div>
                <div id="pwd-alert"></div>
                <div class="form-group">
                    <label>Текущий пароль</label>
                    <input type="password" id="pwd-old">
                </div>
                <div class="form-group">
                    <label>Новый пароль</label>
                    <input type="password" id="pwd-new" minlength="6">
                </div>
                <div class="form-group">
                    <label>Повторите новый пароль</label>
                    <input type="password" id="pwd-new2" minlength="6">
                </div>
                <button class="btn btn--primary btn--sm" onclick="Settings.changePassword()">🔑 Изменить пароль</button>
            </div>`;
    },

    async changePassword() {
        const oldPwd = document.getElementById('pwd-old').value;
        const newPwd = document.getElementById('pwd-new').value;
        const newPwd2 = document.getElementById('pwd-new2').value;

        if (!oldPwd || !newPwd) {
            alert('Заполните все поля');
            return;
        }
        if (newPwd !== newPwd2) {
            alert('Пароли не совпадают');
            return;
        }
        if (newPwd.length < 6) {
            alert('Пароль должен быть не менее 6 символов');
            return;
        }

        const res = await API.post('/api/change-password', {
            old_password: oldPwd,
            new_password: newPwd,
        });

        if (res && res.ok) {
            alert('Пароль изменён');
            document.getElementById('pwd-old').value = '';
            document.getElementById('pwd-new').value = '';
            document.getElementById('pwd-new2').value = '';
        } else {
            alert('Ошибка: ' + (res ? res.data.error : 'Сервер недоступен'));
        }
    }
};

function showToast(message, type = 'success', duration = 3000) {
    document.querySelectorAll('.toast').forEach(t => t.remove());
    const toast = document.createElement('div');
    toast.className = 'toast toast--' + type;
    toast.textContent = message;
    document.body.appendChild(toast);
    setTimeout(() => toast.remove(), duration);
}

document.addEventListener('DOMContentLoaded', () => App.init());
