# Отчёт: повторная приёмка mcp-is на testshop (2026-11-02)

Сервер: `https://mcp-is.ru/mcp`, `serverInfo` = `mcp-is 0.9.4`, `protocolVersion` = `2025-06-18`.
Клиент: сырой JSON-RPC по Streamable HTTP (`probe.py` / `probe_c.py`), без моделей схем Hermes и без SDK `mcp` —
чтобы отличать дефект сервера от дефекта клиента (правило 1 брифа).
Магазин: `myshop-dgi816.myinsales.ru`, каталог 17 товаров (15 восстановленных + 2 тестовых).

## Сводная таблица

| # | Шаг | Ожидание | Факт | Источник |
|---|---|---|---|---|
| 1 | `tools/list` | `get_seo`/`update_seo` в core surface | ЕСТЬ, всего 32 инструмента, http 200 | tools/list |
| 2 | `update_product` схема | есть `clear_fields`, НЕТ `quantity` | Так и есть: whitelist `productUpdateFields` (+`clear_fields`), `quantity` отсутствует | код `client_impl.go` + поведение |
| 3 | Контрпример `get_seo` | ошибка, не тишина | `isError:true`, `unknown resource_type "zzz" (product\|page\|article\|blog\|collection)` | get_seo |
| 4 | `get_seo` bad id | ошибка | `isError:true`, `API error 404 Couldn't find Product with 'id'="xyz"` | get_seo |
| 5 | Пагинация `per_page:3` | ровно 3 | 3 товара, `total:3`, `page:1`, `per_page:3` | list_products |
| 6 | Пагинация `per_page:"3"` (строка) | 3 | 3, идентичный набор — строка приводится к числу | list_products |
| 7 | `page:5, per_page:3` | товары 13–15 | 1947395297, 1947395337, 1947395361 — позиции 13–15 полного списка | list_products |
| 8 | `page:6, per_page:3` | 2 последних | 2 товара (1947395409, 1947395465) | list_products |
| 9 | `page:7, per_page:3` | пусто | `items=[]`, http 200 | list_products |
| 10 | `per_page:250` | 17/17 | 17/17, ids совпадают с `get_products_count` | list_products |
| 11 | `per_page:300` | усечение до 250 | 17 товаров, конверт `per_page:250` | list_products |
| 12 | `total` = размер выдачи | да | при `per_page:3` вернул `total:3`, при 250 — `total:17` | list_products |
| 13 | `create_product` | товар создан | id `1947400777`, variant `2376378273`, price `123.0`, sku `ZZZ-HERMES-TEST`, qty 5; `get_products_count` 17→18 | create_product + get_products_count |
| 14 | `update_product` мусорный ключ | `ignored_keys`, объект не изменён | `ignored_keys:["zzz_probe_field"]`, `updated_at` не сдвинулся (`19:02:47`) | update_product + get_product |
| 15 | `update_product` `title:""` | Название не затёрлось | 200, `title` = `zzz-hermes-test`, `updated_at` тот же | update_product + get_product |
| 16 | `update_seo` | `get_seo` показывает ровно их | title/desc/keywords записаны и прочитаны байт-в-байт | update_seo → get_seo |
| 17 | `update_seo` мусорный ключ | `ignored_keys` | `ignored_keys:["zzz_probe_field"]`, значения целы | update_seo → get_seo |
| 18 | `clear_fields:["meta_keywords"]` | поле пустое (null) | `meta_keywords: null`, `html_title` сохранён | update_seo → get_seo |
| 19 | `clear_fields:["bogus_field"]` | ошибка | `isError:true`, `clear_fields: unknown field "bogus_field"` | update_seo |
| 20 | `update_variant` `quantity:0` | остаток реально 0 | qty 0 в ответе и в `get_variant`; возвращено 5 | update_variant → get_variant |
| 21 | `update_variant` `price:"321.5"` | цена строкой принята | `price:"321.5"` | update_variant |
| 22 | `add_product_to_collection` | привязка создана | collect `12015868569`, position 129; `get_collects_by_product` → total 1 | add_product_to_collection |
| 23 | `remove_product_from_collection` | привязка снята | с `collect_id` → `null`; `get_collects_by_product` → total 0 | remove_product_from_collection |
| 24 | `delete_product` | исчез из каталога | `get_products_count` 18→17, в списке нет, `get_product` → 404 | delete_product |
| 25 | Куда делся товар | корзина или навсегда | **в корзине**: `list_products deleted:true` содержит `1947400777` (всего 71 запись) | list_products deleted |
| 26 | Контрпример `get_product` | ошибка | `isError:true`, `API error 404 Couldn't find Product with 'id'="999999999999"` | get_product |
| 27 | Живые данные без кеша | изменение видно сразу | variant/seo/delete/collect читаются новыми значениями без паузы | все read-back |
| 28 | Витрина `/collection/all` | восстановленные товары с картинками | page1 = 12 карточек + page2 = 5, итого 17 slug'ов; cdn-картинок 78 и 38 соответственно | витрина (cookie) |
| 29 | SEO на витрине (`/product/testovyy-tovar-mcp`) | `<title>`/`meta description` из `update_seo` | `<title>Тестовый товар MCP — SEO заголовок</title>`, desc = `SEO-описание товара, записанное через MCP-сервер mcp-is.` | витрина |
| 30 | Админка «Все товары» = 17 | 17 | **не проверено мной**: браузер отказывает (`use_real_profile` требует Chromium по умолчанию), сессии InSales у меня нет. Требует ручного подтверждения | — |

## Подтверждение инцидента 02.11 (важно для разработчиков)

В корзине (`list_products deleted:true`, 71 запись) найдены **все 15 старых ID** из `tmp/restored_ids.json`
(`1943778241`, `1943778257`, … `1943778385`) — 15/15, ни одного «удалён навсегда».
То есть каталог не был уничтожен физически: он был удалён (в корзину) одной операцией, а восстановление
шло новыми ID `1947394761`–`1947395465`. Это согласуется с версией «импорт со стратегией, удаляющей
отсутствующие позиции», и не согласуется с версией «PHP/API-баг удаления по одному».

## Расхождения с брифом и «похоже на успех, но ничего не произошло»

1. **`get_products_count` нет в core `tools/list`** — доступен только как под-инструмент
   `call_insales_tool {name:"get_products_count"}`. Работает корректно (17→18→17), но в брифе он назван
   наравне с core-инструментами.
2. **`remove_product_from_collection` ждёт `collect_id`**, тогда как парный `add_product_to_collection` —
   `collection_id`+`product_id`. Несимметричная сигнатура; ошибка точная (`required argument "collect_id" not found`),
   но по аналогии её не угадать.
3. **Повторный `delete_product` возвращает `null` (успех)**, хотя товар уже удалён. Состояние не меняется,
   ответ выглядит успешным — это ровно тот класс «тихого успеха», который чеклист просит фиксировать.
4. **`per_page:0` → 1 товар** (а не дефолт 25 и не ошибка валидации). Мелочь, но конверт честно говорит `per_page:1`.
5. **`meta_keywords` не рендерится темой ShopBerry** (`count=0` в HTML) — запись подтверждена только через
   `get_seo`; на витрине проверить нельзя (ограничение темы, не сервера).
6. **`update_product` не пишет `price`/`quantity`** — только `variants_attributes` (см. `CreateProduct`).
   Для `zzz` это не проявилось, но при обновлении существующих товаров цена/остаток меняются только `update_variant`.

## Вердикт

- **Работает и подтверждено read-back:** защита payload'ов (`ignored_keys`, пустые значения не пишутся,
  `clear_fields` с валидацией имён), честная пагинация (число и строка, границы, усечение до 250),
  SEO-инструменты (запись/чтение/очистка), варианты (нулевой остаток — легитимная запись), коллекции
  (добавление/снятие), удаление с попаданием в корзину, отсутствие кеша в чтении.
- **Сломано:** ничего из заявленного в брифе не воспроизвелось как отказ.
- **Подозрительно / требует решения:** повторное удаление как «тихий успех»; асимметрия
  `collect_id` vs `collection_id`; отсутствие `get_products_count` в core surface; `per_page:0 → 1`.
- **Не проверено:** счётчик «Все товары = 17» в админке (нет доступа к сессии InSales).

Тестовые артефакты: `zzz-hermes-test` (`1947400777`) удалён и лежит в корзине — при необходимости
восстанавливается вручную из «Импорт/Экспорт → Корзина».
