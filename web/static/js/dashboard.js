const Dashboard = {
    async render() {
        const [shopsRes, keysRes] = await Promise.all([
            API.get('/api/shops'),
            API.get('/api/ai-keys')
        ]);

        const shops = (shopsRes && shopsRes.ok) ? shopsRes.data : [];
        const keys = (keysRes && keysRes.ok) ? keysRes.data : [];
        const nickname = localStorage.getItem('nickname') || 'пользователь';

        return `
            <div class="page-header">
                <div>
                    <h1>Панель управления</h1>
                    <p class="desc">Добро пожаловать, ${nickname}!</p>
                </div>
            </div>

            <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px">
                <div class="card">
                    <div class="card__title">Магазины</div>
                    <div class="card__desc">Подключённые магазины InSales</div>
                    <div style="font-size:32px;font-weight:700">${shops.length}</div>
                    <a href="/shops" class="btn btn--ghost btn--sm" style="margin-top:12px">Управление →</a>
                </div>
                <div class="card">
                    <div class="card__title">AI-ключи</div>
                    <div class="card__desc">Ключи для подключения AI-агентов</div>
                    <div style="font-size:32px;font-weight:700">${keys.length}</div>
                    <a href="/ai-keys" class="btn btn--ghost btn--sm" style="margin-top:12px">Управление →</a>
                </div>
            </div>

            <div class="card" style="margin-top:16px">
                <div class="card__title">Быстрый старт</div>
                <div class="card__desc">Как подключить AI-агента к вашему магазину</div>
                <ol style="font-size:14px;line-height:1.8;padding-left:20px">
                    <li>Добавьте магазин InSales на странице <a href="/shops">Магазины</a></li>
                    <li>Сгенерируйте API-ключ на странице <a href="/ai-keys">AI-ключи</a></li>
                    <li>Скопируйте ключ в настройки вашего AI-клиента (Cursor, Claude Desktop)</li>
                    <li>AI-агент готов работать с вашим магазином!</li>
                </ol>
            </div>`;
    }
};
