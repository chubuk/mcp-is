# Справка: заполнение карточки приложения в ЛК InSales

Раздел **Приложения → Разработчикам → Добавление нового расширения**.
Значения ниже — для приложения MCP IS (бесплатное, полный доступ к API).

---

## Описание для клиентов

| Поле | Значение |
|------|----------|
| **Название** | `MCP IS — AI-доступ к магазину` |
| **Категория** | ближайшая из списка (Интеграции / Продвижение / Прочее) |
| **Стоимость** | значение из настройки «Стоимость подписки» (Админ → Сервер); синхронизировать вручную |
| **Автор** | `Vladislav Chubuk` |
| **Контакты** | `chubuk@mail.ru` |
| **Image** | логотип 250×100 (PNG) |

### Описание (HTML)

```html
<p><strong>MCP IS</strong> — подключение AI-ассистентов (Claude, Cursor, MiMoCode и др.) к данным вашего магазина через MCP-протокол.</p>
<ul>
  <li>Товары, заказы, клиенты, коллекции — чтение и изменение</li>
  <li>Безопасные ключи доступа для каждого AI-клиента</li>
  <li>Маскирование персональных данных (152-ФЗ)</li>
  <li>Готовые инструкции подключения для 10+ MCP-клиентов</li>
</ul>
<p>После установки вы получите ключ и ссылку для подключения AI-ассистента к магазину.</p>
```

---

## Технические параметры

| Поле | Значение |
|------|----------|
| **Идентификатор** | латиница + цифры, напр. `mcpis` → станет `INSALES_APP_KEY` |
| **Триальный период** | `14` (дефолт формы) |
| **URL установки** | `https://mcp-is.ru/install` |
| **URL входа** | `https://mcp-is.ru/app` |
| **URL деинсталяции** | `https://mcp-is.ru/uninstall` |
| **URL карточки товара** | *пусто* |
| **URL карточки заказа** | *пусто* |
| **Отключить отправку вебхуков для изменений, вызванных самим расширением** | **включить** ✓ |

> Секрет (`secret_key`), который InSales сгенерирует рядом с идентификатором, — это `INSALES_SECRET_KEY`.
> Не сообщать его никому; он нужен только в `.env` сервера.

### После создания приложения

```bash
# /opt/mcp-is/.env на VDS
INSALES_APP_KEY=mcpis            # как в форме
INSALES_SECRET_KEY=<секрет>      # из карточки приложения

docker compose restart app
```

---

## Права и ограничения

| Право | Значение | Зачем |
|-------|----------|-------|
| Заказы → Доступ в раздел | **✓** | MCP-инструменты заказов |
| Заказы → Импортировать заказы старше двух месяцев | **✓** | полная выгрузка истории |
| Клиенты → Доступ в раздел | **✓** | MCP-инструменты клиентов |
| Клиенты → Доступ к персональным данным | **✓** | ФИО, бонусы и т.п. |
| Клиенты → Номер телефона | **✓** | PDN-маскирование (152-ФЗ) |
| Клиенты → Email | **✓** | PDN-маскирование (152-ФЗ) |
| Товары → Доступ в раздел | **✓** | MCP-инструменты товаров |
| **Недоступно в странах** | **ничего не выбирать** | без гео-ограничений |

Кнопка **Создать**.

---

## API-права (api_permissions)

Ниже «Права и ограничений» — поэндпоинтный список REST API. Галочка = эндпоинт
доступен ключу приложения (403 без неё). Отмечаем **только то, что реально
вызывает клиент** (`internal/integration/insales/client_impl.go`); остальное —
поверхность для аудита и модерации, не нужна.

Проверка 2026-11 на обычном платном магазине (myshop-dgi816): вызовы без
галочек (`list_variants`, `get_orders_count`) проходили с 200 — enforcement
не подтверждён. Галочки проставляем всё равно: это декларация прав для
публикации в каталоге InSales. Новые права подхватываются установленными
магазинами без переустановки (права — записи приложения, не инсталляции;
проверено: право, добавленное после установки, действует на старой).

| Карточка (ресурс) | Отметить действия |
|---|---|
| Order | index, show, create, update, count |
| Client | index, show, create, update, destroy |
| ClientGroup | index, show, create, update, destroy |
| BonusSystemTransaction | index, show, create |
| CustomStatus | index, show, create, update, destroy |
| DiscountCode | index, show, create, update, destroy |
| ExternalDiscount | index, show, create, update, destroy |
| DiscountCollectionsGroup | index, show, discount_codes |
| Product | index, show, create, update, destroy, count |
| Image | index, create, destroy |
| Variant | index, show, update, destroy |
| ProductFieldValue | index, show, create, update, destroy |
| VariantFieldValue | index, show, create, update, destroy |
| CollectionFieldValue | index, show, update |
| ProductAccessoryLink | index |
| Similar | index, create, destroy |
| Accessory | index, show, create, update, destroy |
| AccessoryValue | index, show, create, update, destroy |
| ProductField | index, show, create, update, destroy |
| VariantField | index, show, create, update, destroy |
| CollectionField | index, show, create, update, destroy |
| CollectionFilter | index, show, create, update, destroy |
| ApiCategory | index, show, create, update, destroy |
| ApiCollection | index, show, create, update, destroy |
| Collect | index, show, create, update, destroy |
| OptionName | index, show, create, update, destroy |
| OptionValue | index, show, create, update, destroy |
| PriceKind | index, show, create, update, destroy |
| Property | index, show, create, update, destroy |
| ApiCharacteristic | index, show, create, update, destroy |
| Review | index, show, create, update, destroy |
| SalesChannel | index, show, create, destroy (update в форме не предлагается) |
| Account | show, update |
| Currency | show |
| DeliveryVariant | index, show, create, update, destroy |
| PaymentGateway | index, show, create, update, destroy |
| Domain | index, show, create, update, destroy |
| StockCurrency | index, show, create, update (destroy в форме не предлагается) |
| Unit | index |
| ApplicationAction | index, show, create, update, destroy |
| Page | index, show, create, update, destroy |
| Blog | index, show, create, update, destroy |
| Article | index, show, create, update, destroy |
| Menu | index, show, create, update, destroy |
| ApiMenuItem | index, show, create, update, destroy |
| Redirect | index, create, destroy |
| File | index, show, create, destroy |
| JsTag | index, show, create, update, destroy |
| Marketplace | index, show, create, update, destroy |
| Outlet | index, show, create, update, destroy |
| Warehouse | index, show, create, update, destroy |
| UserTask | index, show, create, update, destroy |
| Supplementary | index, create, destroy (show в форме нет) |
| OAuthApp | index, show, create, update, destroy |
| Webhook | index, show, create, update, destroy (обязательные, включены всегда) |

Сверено с каталогом `https://api.insales.ru` (2026-11). Форма повторяет его
ресурсы 1:1; расхождения каталог↔клиент:
- **Supplementary** — в каталоге только Get list/Create/Destroy, show у клиента
  (`GET /admin/supplementaries/:id.json`) и галочки `show` нет.
- **SalesChannel** — каталог без Update, форма без update; клиент шлёт
  `PUT /admin/sales_channels/:id.json` (вне каталога).
- **StockCurrency** — каталог без Destroy, форма без destroy; клиент шлёт
  `DELETE /admin/stock_currencies/:id.json` (вне каталога).
- **DiscountCollectionsGroup** — каталог только 3 GET (index/show/discount_codes);
  клиент дополнительно шлёт POST/PUT/DELETE (вне каталога и формы).
- **Currency** — карточка в форме есть (монтируется как `/admin/currency`),
  секции в каталоге нет; клиент шлёт `GET /admin/currencies.json`.
- **MetaFields** (`/admin/meta_fields.json`, key/value на Product/Variant/Order/
  Client/Category/Page) — есть в openapi-спеке, но нет ни секции в каталоге, ни
  карточки в форме — InSales их не разграничивает, отмечать нечего.
- **ProviderToken** (`POST /admin/dialogs/provider_tokens.json`) — в каталоге
  есть, карточки в форме нет.

Не отмечаем (клиент не вызывает): ClientTag, Horeca::Table, PickUpSource,
OrderDoc, AccountCountry, AccountLanguage, BonusSystem, Block/BlockList/Theme/
Asset/WidgetType, Options::Fields, ShippingCompany (обязательная).

Нюансы:
- Similar в форме монтируется как `/admin/products/:product_id/similars`, клиент
  шлёт top-level `/admin/similars.json`; Article — `/admin/articles` vs
  `/admin/blogs/:id/articles`. Контроллер один, галочки ставим по карточке.
- Галочки проставляются в консоли браузера на странице карточки, затем
  «Сохранить». Только докрутка вверх — ничего не снимаем:

```js
// Вставить в консоль на странице карточки приложения и нажать «Сохранить».
const need = {
  'back_office/orders/orders': 'index show create update count',
  'back_office/orders/clients': 'index show create update destroy',
  'back_office/orders/client_groups': 'index show create update destroy',
  'back_office/orders/bonus_system_transactions': 'index show create',
  'back_office/orders/custom_statuses': 'index show create update destroy',
  'back_office/orders/discount_codes': 'index show create update destroy',
  'back_office/orders/external_discounts': 'index show create update destroy',
  'back_office/orders/discount_collections_groups': 'index show discount_codes',
  'back_office/products/products': 'index show create update destroy count',
  'back_office/products/images': 'index create destroy',
  'back_office/products/variants': 'index show update destroy',
  'back_office/products/product_field_values': 'index show create update destroy',
  'back_office/products/variant_field_values': 'index show create update destroy',
  'back_office/products/collection_field_values': 'index show update',
  'back_office/products/product_accessory_links': 'index',
  'back_office/products/similars': 'index create destroy',
  'back_office/products/accessories': 'index show create update destroy',
  'back_office/products/accessory_values': 'index show create update destroy',
  'back_office/products/product_fields': 'index show create update destroy',
  'back_office/products/variant_fields': 'index show create update destroy',
  'back_office/products/collection_fields': 'index show create update destroy',
  'back_office/products/api_collection_filters': 'index show create update destroy',
  'back_office/products/api_categories': 'index show create update destroy',
  'back_office/products/api_collections': 'index show create update destroy',
  'back_office/products/collects': 'index show create update destroy',
  'back_office/products/option_names': 'index show create update destroy',
  'back_office/products/option_values': 'index show create update destroy',
  'back_office/products/price_kinds': 'index show create update destroy',
  'back_office/products/properties': 'index show create update destroy',
  'back_office/products/api_characteristics': 'index show create update destroy',
  'back_office/products/reviews': 'index show create update destroy',
  'back_office/products/sales_channels': 'index show create destroy',
  'back_office/products/supplementaries': 'index create destroy',
  'back_office/options/accounts': 'show update',
  'back_office/options/currencies': 'show',
  'back_office/options/delivery_variants': 'index show create update destroy',
  'back_office/options/payment_gateways': 'index show create update destroy',
  'back_office/options/domains': 'index show create update destroy',
  'back_office/options/stock_currencies': 'index show create update',
  'back_office/options/units': 'index',
  'back_office/options/application_actions': 'index show create update destroy',
  'back_office/shopwindow/pages': 'index show create update destroy',
  'back_office/shopwindow/blogs': 'index show create update destroy',
  'back_office/shopwindow/articles': 'index show create update destroy',
  'back_office/shopwindow/menus': 'index show create update destroy',
  'back_office/shopwindow/api_menu_items': 'index show create update destroy',
  'back_office/shopwindow/redirects': 'index create destroy',
  'back_office/shopwindow/files': 'index show create destroy',
  'back_office/shopwindow/js_tags': 'index show create update destroy',
  'back_office/promotion/marketplaces': 'index show create update destroy',
  'back_office/outlets': 'index show create update destroy',
  'back_office/warehouses': 'index show create update destroy',
  'back_office/user_tasks': 'index show create update destroy',
  'back_office/applications/o_auth_apps': 'index show create update destroy',
};
let n = 0;
document.querySelectorAll('input[type=hidden][name$="[controller]"]').forEach(h => {
  const idx = h.name.match(/\[(\d+)\]/)[1];
  const act = document.querySelector(
    `input[type=hidden][name="application[api_permissions][${idx}][action]"]`)?.value;
  const cb = document.querySelector(
    `input[type=checkbox][name="application[api_permissions][${idx}][_destroy]"]`);
  const want = (need[h.value] || '').split(' ');
  if (cb && !cb.disabled && want.includes(act) && !cb.checked) { cb.click(); n++; }
});
console.log('отмечено новых галочек:', n);
```

---

## Проверка после публикации

1. Тестовый магазин → установить приложение.
2. `GET https://mcp-is.ru/install?...` должен ответить `200` (в логах `app installed`).
3. Кнопка «Войти в приложение» → автологин → экран с MCP-ключом.
4. Удалить приложение → в логах `app uninstalled`.

Подробности протокола: `docs/setup/INSALES-APP.md`, `docs/api/api_and_apps.md`.
