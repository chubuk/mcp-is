const App = {
    routes: [
        { path: '/login', render: () => Auth.renderLogin(), auth: false, onMount: () => Auth.bindLogin() },
        { path: '/register', render: () => Auth.renderRegister(), auth: false, onMount: () => Auth.bindRegister() },
        { path: '/dashboard', render: () => Dashboard.render(), auth: true },
        { path: '/shops', render: () => Shops.render(), auth: true, onMount: () => Shops.bind() },
        { path: '/ai-keys', render: () => AIKeys.render(), auth: true, onMount: () => AIKeys.bind() },
        { path: '/settings', render: () => Settings.render(), auth: true },
        { path: '/forgot-password', render: () => Auth.renderForgotPassword(), auth: false, onMount: () => Auth.bindForgotPassword() },
        { path: '/reset-password', render: () => Auth.renderResetPassword(), auth: false, onMount: () => Auth.bindResetPassword() },
        { path: '/admin/users', render: () => Admin.renderUsers(), auth: true, role: ['admin', 'super_admin'] },
        { path: '/admin/users/:id', render: (params) => Admin.showUser(params.id), auth: true, role: ['admin', 'super_admin'] },
        { path: '/admin/stats', render: () => Admin.renderStats(), auth: true, role: ['admin', 'super_admin'] },
        { path: '/admin/server', render: () => Admin.renderServer(), auth: true, role: ['admin', 'super_admin'] },
        { path: '/admin/pdn', render: () => Admin.renderPDN(), auth: true, role: ['admin', 'super_admin'] },
        { path: '/force-change-password', render: () => Auth.renderForceChangePassword(), auth: true, onMount: () => Auth.bindForceChangePassword() },
    ],

    currentPage: null,

    init() {
        window.addEventListener('popstate', () => this.route());
        document.addEventListener('click', (e) => {
            const link = e.target.closest('a[href]');
            if (link) {
                e.preventDefault();
                this.navigate(link.getAttribute('href'));
            }
        });
        this.route();
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
        API.clearToken();
        localStorage.removeItem('nickname');
        localStorage.removeItem('role');
        localStorage.removeItem('user_id');
        this.navigate('/login');
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
                    <input type="email" id="settings-email" value="${email}" placeholder="your@email.com">
                </div>
                <button class="btn btn--primary btn--sm" onclick="Settings.saveEmail()">💾 Сохранить email</button>
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

    async saveEmail() {
        const email = document.getElementById('settings-email').value;
        if (!email) return;
        const res = await API.post('/api/set-email', { email });
        if (res && res.ok) {
            localStorage.setItem('email', email);
            alert('Email сохранён');
        } else {
            alert('Ошибка: ' + (res ? res.data.error : 'Сервер недоступен'));
        }
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
