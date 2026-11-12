# W2: аудит схем update_*/create-* тулов vs InSales API

Три слоя: (а) inputSchema в `internal/mcp/tools.go`; (б) проводка аргументов в
`internal/mcp/handlers.go` (утекают ли роутинг-id из `req.GetArguments()`); (в) пэйлоад
клиента в `internal/integration/insales/client_impl.go` vs таблицы
`docs/api/Insales-api-docs.md` (источник истины). Классы: (i) нет поля в схеме,
(ii) несовпадение имени (молча теряется), (iii) лишний не-API ключ в пэйлоаде.

Статус-сноска: в рабочем дереве уже лежит (незакоммиченный) фикс блога W3
(`BuildBlogArticle`) — строки 1–2 отражают текущее состояние, seed-дефекты по ним
закрыты. Строки 3–4 (seed) остаются открытыми. Все update/create из `tools.go`
(~52 + ~40) пройдены; `update_product` и `create/update_page` — см. «Не дефекты».

Статус-сноска 2 (после W4–W7 и живых проб 2026-11-10): фиксы по таблице внедрены —
алиасы переименований (`title→name`, `event→topic`, `quantity→quantity_at_warehouse0`,
`handle→permalink` и др.), расширены схемы collection/page/blog/webhook/property_value/
accessory/discount_code/variant/external_discount, `CreateClient`/`CreateOrder` не теряют
`fields_values_attributes` (у заказа — follow-up PUT, т.к. InSales игнорирует их при
POST /orders.json), `CreatePropertyValue(args)`, роутинг-id вырезаются. `field_values`
добавлен как единый параметр (9 тулов) и убран у `create/update_category` — категории
InSales не хранят доп. поля. Живые пробы: testshop `myshop-dgi816` (блог дата/картинка,
field_values create+upsert на товарах/вариантах/заказах/клиентах/коллекциях).

## Таблица дефектов

| # | тул | класс | сегодня (file:line) | InSales ждёт (строка доков) | фикс |
|---|-----|-------|---------------------|------------------------------|------|
| 1 | update_blog_post (seed) | ii+iii | **исправлено в раб. дереве (W3)**: `BuildBlogArticle` client_impl.go:1476-1494 (body_html→content, image_url→image_attributes{src}, режет blog_id/post_id); схема tools.go:1127-1143 полная | `article[content]`, `article[image_attributes]`, без id — доки 2044-2058 | остаток (i): `related_products_attributes` (доки 1613) |
| 2 | create_blog_post (seed) | i | **исправлено (W3, уточнено v0.14)**: хардкода «2025-01-01»/«Admin» нет; **дефолтов тоже нет** — `published_at`/`author` обязательны (InSales 422 без них), + `image_data`/`image_filename` (attachment-загрузка обложки) | create: доки 1603-1619 (`published_at`,`author` обязательны) | остаток (i): `related_products_attributes` (доки 1618) |
| 3 | update_blog (seed) | iii (+i) | client_impl.go:1432 `blog{cleanArgs(args)}` + handlers.go:1638-1641 `args:=GetArguments()` → `blog_id` внутри `blog{}`; `handle` в схеме (tools.go:1066) — нет ни в одной таблице API | доки 2395-2412: `blog[title,html_title,meta_keywords,meta_description,notify_about_comment,captcha_enabled,disallow_links_in_comments,commentable]` | вырезать blog_id (паттерн `BuildBlogArticle`); +html_title/meta_* (частично дублируется update_seo); handle — убрать или задокументировать |
| 4 | класс `X: cleanArgs(args)` (seed) | iii | **подтверждён везде**, где хендлер отдаёт сырой args: клиент 637, категория 705, вебхук 938, аксессуар 1154, вариант 1199, accessory_value 1254, меню 1320, блог 1432, промокод 1586, property 1719, characteristic 1780, review 1843, custom_status 2074 + generic-класс (стр. 25) | id-entity в path, не в теле | общий фикс W4: вырезать роутинг-id в хендлере/клиенте (см. строки 3,5,7,8,10,14,19,21-26) |
| 5 | update_menu | ii+iii+i | client_impl.go:1320 шлёт `menu{title,menu_id}`; API не знает `title` → переименование молча не проходит | `menu[name]`, `menu[handle]` — доки 10292-10299 | title→name, вырезать menu_id, +handle |
| 6 | create_menu | ii | client_impl.go:1303-1304 `menu{title,handle}` | `menu[name]`, `menu[handle]` — доки 10164-10172 | title→name (как в `CreateAccessory`, client_impl.go:1136) |
| 7 | update_accessory | ii+iii+i | client_impl.go:1153-1154 `accessory{title,position,accessory_id}`; `create_accessory` при этом мапит title→name (1133-1140) — асимметрия | `accessory[name,max_count,min_count,permalink]` — доки 325-332 | title→name; `position` не параметр — убрать; вырезать accessory_id; +max_count/min_count/permalink |
| 8 | update_property_value | ii+iii+i | client_impl.go:1253-1254 `accessory_value{title,accessory_id,property_value_id}` | `accessory_value[name,price,accounting_variant_id]` — доки 872-882 | title→name; вырезать оба id; +price/accounting_variant_id |
| 9 | create_property_value | i | client_impl.go:1235-1242 хардкод `price:0`, в схеме только title (tools.go:880) | `accessory_value[name,price]` (обяз.) — доки 396-404 | +price, accounting_variant_id в схему |
| 10 | update_webhook | ii+iii+i | client_impl.go:936-938 `webhook{url,webhook_id}` | `webhook[address,topic,format_type,warehouse_id,sales_channel_id,batch_size]` — доки 27349-27360 | url→address; вырезать webhook_id; +остальные |
| 11 | create_webhook | ii | client_impl.go:918-923 `webhook{event,url,format}` | `webhook[address,topic,format_type]` — доки 27161-27175 (ref:1242-1247) | event→topic, url→address, format→format_type |
| 12 | create_bonus_system_transaction | ii+iii | client_impl.go:2539 `bonus_system_transaction{client_id,amount,comment}` (handlers.go:2792-2800 не мапит) | `bonus_system_transaction[bonus_points,description]` — доки 2474-2488 | amount→bonus_points, comment→description; вырезать client_id (он в path) |
| 13 | create_custom_status | ii | tools.go:1503 `status_type`; client_impl.go:2061 шлёт как есть | `custom_status[system_status,title]` — доки 5629-5645 | status_type→system_status |
| 14 | update_custom_status | iii | client_impl.go:2074 + genericUpdateGetHandler (handlers.go:2238) → `status_id` в теле; путь по id, доки — по permalink | `custom_status[title]` — доки 5860-5866 | вырезать status_id (см. стр. 25) |
| 15 | create/update_collection_filter | ii (+i) | ключ-обёртка `collection_filter` (client_impl.go:2283/2286, crudCreate/Update 2125/2138) и `filter_type` в схеме (tools.go:1804) | обёртка **`filter`**; `filter[id]` обяз. в теле update; поля title, collection_id, description, permalink, html_title, meta_*, *_attributes — доки 5218-5228, 5499-5522 | ключ `filter`; filter_type не существует в API; +description/permalink/SEO |
| 16 | create/update_external_discount | ii | схема title/percent/amount (tools.go:1860-1867) — ни одно не параметр API | `external_discount[description,url,order_lines_depend,active]` (+`observed_fields` create) — доки 7784-7795, 7939-7948 | пересобрать схему под API-имена |
| 17 | update_collection | i | tools.go:525 — только title (хендлер handlers.go:999-1010 передаёт лишь его) | +position, sort_type, description, html_title, seo_description, meta_description, meta_keywords, is_hidden, permalink, recursive, field_values_attributes — доки 4598-4612 | расширить схему (SEO частично дублирует update_seo) |
| 18 | create_collection | i | tools.go:516 title, parent_id | +position, sort_type, description, SEO, image_attributes, field_values_attributes — доки 4237-4255 | расширить схему |
| 19 | update_client | i+iii | tools.go:434 email,phone; client_impl.go:637 `client{cleanArgs}` → client_id в теле; `email` не параметр update | `client[name,phone,subscribe,client_group_id,ip_addr,juridical_address,inn,kpp,ogrn,okpo,bik,bank_name,correspondent_account,settlement_account,fields_values_attributes]` — доки 3594-3614 | вырезать client_id; +name/subscribe/group/fields; email убрать из update |
| 20 | create_client | i | client_impl.go:613-622 (name,surname,middlename,email,phone,password,registered) | +subscribe, client_group_id, type, ip_addr, fields_values_attributes, skip_send_mail — доки 3131-3160 | дописать поля |
| 21 | update_variant | ii+iii+i | client_impl.go:1198-1199 `variant{cleanArgs}` → product_id+variant_id в теле; `sku`/`quantity` не в таблице update | `variant[price,quantity_at_warehouse0,dimensions,variant_field_values_attributes]` — доки 26197-26210 | quantity→quantity_at_warehouse0; sku мимо (или задокументировать); вырезать id; +dimensions |
| 22 | update_review | iii | client_impl.go:1843 → review_id в теле; rating/published/manager_reply не параметры update | только `review[content]` — доки 24805-24812 (create полный: 24271-24282, product_id — легитимный параметр) | вырезать review_id; rating/published/manager_reply убрать из update |
| 23 | update_discount_code | iii | client_impl.go:1586 → discount_code_id в теле; code/description/disabled — create-only | update: только `discount_code[discount]` — доки 7383-7387 | вырезать id; create: +disabled (доки 7178-7195) |
| 24 | update_property, update_characteristic, create_characteristic | iii | client_impl.go:1719/1780/1767 → property_id[/characteristic_id] в теле (handlers.go:2023 не вырезает property_id) | update принимает только title — доки 23773-23777, 3075-3079; create characteristic: +permalink, position (доки 2934-2942) | вырезать id; is_navigational — create-only (доки 23465-23478) |
| 25 | genericUpdateGetHandler × ~23 | iii | handlers.go:2233-2250 → `crudUpdate` client_impl.go:2138-2139 `key{cleanArgs(args)}` — роутинг-id утекает: variant_id, gateway_id, domain_id, group_id, field_id, filter_id, tag_id, discount_id, warehouse_id, channel_id, currency_id, kind_id, option_name_id, action_id, widget_id, app_id, task_id, outlet_id, marketplace_id и др. (tools.go:659-2349) | id — в path URL | вырезать param-ключ в genericUpdateGetHandler (как `nestedUpdateHandler`, handlers.go:2302-2310) |
| 26 | create/update_option_value | iii | handlers.go:2732-2758 → crudCreate/Update (2125/2138) с option_name_id/option_value_id в теле | `option_value[position,image...]` — доки 11406-11413 | вырезать оба id |
| 27 | update_page (i) / create_page | iii / OK | client_impl.go:1086 `page{cleanArgs}` → page_id в теле; **content нативно — переименование корректно** (проверено: tools.go:754/764, 1066/1085) | update: `page[title,content,permalink,html_title,meta_keywords,meta_description,created_at]` — доки 20621-20633 | вырезать page_id; +permalink/created_at (SEO — через update_seo) |
| 28 | update_order*, create_order | OK | `BuildOrderBody` (order_payload.go:6-41) — whitelist, order_id в тело не попадает; sub-хендлеры собирают *_attributes вручную (handlers.go:2499-2645) | order_lines_attributes/discounts_attributes/shipping_address_attributes — доки 19197-19213 и сл. | только (i): `marking_codes` в order line (доки 19212) |
| 29 | update_seo / get_seo | ii (alias) | маппинг полей **верный**: seo.go:22-45 product/page/article/blog/collection, `seo_description` только для collection (доки 4604-4608), пэйлоад фильтруется whitelist (seo.go:60-73) — утечки resource_id невозможны. **НО** алиас `keywords` мёртв: `resolveKeywordsAlias` (handlers.go:580-591) вызывается после `wrapArgsAudit` (handlers.go:610-624), который вырезает `keywords` (нет в inputSchema, tools.go:200-215) → ignored_keys, молча | `keywords` = legacy-имя meta_keywords | упомянуть keywords в схеме/фильтре или мапить до wrapArgsAudit |
| 30 | update_account | OK | client_impl.go:2755 `account{client_cookies_whitelist}`, роутинг-id нет | `account[client_cookies_whitelist]` — доки 996-1002 | — |

Секции в `Insales-api-docs.md` найдены для всех проверенных тулсов (пропусков
«нет секции» нет); `insales-api-reference.md` согласен по webhook/blog.

## Сводка

- (i) нет поля в схеме: **~12** (стр. 1-2 остатки, 3, 7-9, 17, 18, 20, 21, 27, 28 + мелочи в 15/24).
- (ii) несовпадение имени/обёртки: **12** (стр. 5-13, 15, 16, 21, 29; стр. 1-2 — были, исправлены W3).
- (iii) лишний не-API ключ (роутинг-id и др.): **~15 групп** (стр. 3-4, 5, 7, 8, 10, 12, 14, 19, 21-27).

Критичные (молча ломают запись): 5-8, 10-12, 15-16 — класс «имя не то» (как B2
`body_html→content`). Утечки id (класс iii) InSales молча игнорирует — грязь и риск.

## Не дефекты

- **`productUpdateFields`/`productArgFields` в `UpdateProduct`** (client_impl.go:408-415,
  handlers.go:664-671) — осознанный whitelist, не трогать.
- **`cleanArgs`** (client_impl.go:390-405) — отбрасывает только nil/«»; `0`/`false`
  проходят (quantity:0 легитимен), явная очистка через `clear_fields`.
- **`nestedUpdateHandler`** (handlers.go:2302-2331) сам вырезает parent/id —
  update_product_field_value / variant_field_value / collection_field_value чистые.
- **update_collect** (handlers.go:1862-1882): collect_id вырезается, payload
  {collection_id, position} — оба легальны (доки 4112-4118, 4182-4186).
- **add_product_to_collection/create_collect**: `collect{collection_id,product_id}`
  позиционно (client_impl.go:1639) = доки 3903-3910.
- **create_accessory** мапит title→name (client_impl.go:1133-1140) — эталонный пример.
- **create_product_field**: field_type→type маппинг (handlers.go:2353-2358);
  update_product_field вырезает field_id (handlers.go:2379); `type` на update
  не-API (доки 23076-23084) — мелочь.
- **create_menu_item** вырезает menu_id (handlers.go:2658); update_menu_item через
  nestedUpdateHandler — чист; (i): `url` в update (доки 10635-10645).
- **update_recurring_application_charge** — singleton-путь, `monthly` = доки 23982-23988.
- **update_product seo-поля инлайн + update_seo** — оба пути рабочие (как задумано).

## Попутно (вне классов, чтение)

`InSalesAccessoryResponse`/`InSalesPropertyValueResponse`/`InSalesMenuResponse`
(types.go:212-236) парсят `json:"title"`, а API возвращает `name` (доки 355, 886,
10318) — GET-ответы этих сущностей вернут пустой title. Проверить/починить в W4.
