const API = {
    token: localStorage.getItem('token'),

    async request(method, path, body) {
        const headers = { 'Content-Type': 'application/json' };
        if (this.token) {
            headers['Authorization'] = 'Bearer ' + this.token;
        }

        const opts = { method, headers };
        if (body) opts.body = JSON.stringify(body);

        const res = await fetch(path, opts);
        const data = await res.json();

        if (res.status === 401 && this.token) {
            this.token = null;
            localStorage.removeItem('token');
            App.navigate('/login');
            return null;
        }

        return { ok: res.ok, status: res.status, data };
    },

    get(path) { return this.request('GET', path); },
    post(path, body) { return this.request('POST', path, body); },
    put(path, body) { return this.request('PUT', path, body); },
    del(path) { return this.request('DELETE', path); },

    setToken(token) {
        this.token = token;
        localStorage.setItem('token', token);
    },

    clearToken() {
        this.token = null;
        localStorage.removeItem('token');
    },

    isLoggedIn() {
        return !!this.token;
    },

    impersonate(token, userId, nickname, role) {
        localStorage.setItem('admin_token', this.token);
        localStorage.setItem('admin_user_id', localStorage.getItem('user_id'));
        localStorage.setItem('admin_nickname', localStorage.getItem('nickname'));
        localStorage.setItem('admin_role', localStorage.getItem('role'));

        this.setToken(token);
        localStorage.setItem('user_id', userId);
        localStorage.setItem('nickname', nickname);
        localStorage.setItem('role', role);
        localStorage.setItem('impersonating', 'true');
    },

    stopImpersonate() {
        const adminToken = localStorage.getItem('admin_token');
        this.setToken(adminToken);
        localStorage.setItem('user_id', localStorage.getItem('admin_user_id'));
        localStorage.setItem('nickname', localStorage.getItem('admin_nickname'));
        localStorage.setItem('role', localStorage.getItem('admin_role'));

        localStorage.removeItem('admin_token');
        localStorage.removeItem('admin_user_id');
        localStorage.removeItem('admin_nickname');
        localStorage.removeItem('admin_role');
        localStorage.removeItem('impersonating');
    },

    isImpersonating() {
        return localStorage.getItem('impersonating') === 'true';
    }
};
