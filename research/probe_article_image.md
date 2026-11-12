# Проба: семантика картинки статьи блога (InSales API)

Дата: 2026-11-10. Магазин: `myshop-dgi816.myinsales.ru` (testshop). Метод: прямой REST
`/admin/blogs/8222025/articles.json`, Basic-auth. Тестовые статьи удалены.

## Вопрос

`article[image_attributes]` документирован только с полем `src` (создание). Семантика
**замены** и **удаления** в доках `docs/api/Insales-api-docs.md` отсутствует. Нужно знать
фактическое поведение для проектирования параметров тулов `create_blog_post`/`update_blog_post`.

## Результаты (по шагам)

| Шаг | Запрос | Результат |
|-----|--------|-----------|
| 1. Add | `POST article{..., image_attributes:{src:"https://…/no_image_compact.jpg"}}` | ✅ картинка создана: `image:{id:16532793, original_url:"/images/loading.gif"}` (асинхронная обработка) |
| 2. Replace | `PUT article{image_attributes:{src:"<другой url>"}}` (без id) | ✅ **замена**: новый id `16532801`, старая картинка ушла. Слот один — новая картинка вытесняет прежнюю |
| 3a. Delete `_destroy:1` (число) | `PUT article{image_attributes:{id, _destroy:1}}` | ❌ не удалил: id **сменился** (16532809) — запись пересоздана/заменилась, картинка осталась |
| 3b. Delete `src:""` | `PUT article{image_attributes:{id, src:""}}` | ❌ то же: новый id (16532817), картинка на месте |
| 3c. DELETE endpoint | `DELETE /blogs/:b/articles/:p/images/:img.json` | ❌ HTTP 404 |
| 3d. Delete `_destroy:"1"` (строка) | `PUT article{image_attributes:{id, _destroy:"1"}}` | ❌ новый id (16532857), не удалил |
| 3e. Пустые атрибуты | `PUT article{image_attributes:{}}` | ❌ новый id (16532865) — любой `image_attributes` пересоздаёт слот |
| 3f. DELETE singular | `DELETE /blogs/:b/articles/:p/image.json` | ❌ HTTP 404 |
| 3g. `image_attributes:null` | `PUT article{image_attributes:null}` | ➖ ничего не изменилось (id прежний) |

## Выводы

1. **Add/Replace — одно и то же действие**: `image_attributes:{src:"url"}` ставит единственную
   картинку статьи; повторный вызов с новым `src` **заменяет** прежнюю (слот один, id новый).
2. **Удаление картинки через API невозможно** ни одним проверенным способом (`_destroy` в обеих
   формах игнорируется, DELETE-эндпоинтов нет). Фактически `image_attributes` имеет только
   create-семантику.
3. **Практический эквивалент «удалить и поставить новую»** = замена через новый `src` (шаг 2).
   Полное удаление — только через админку InSales.

## Решение для тулов (принято по итогам пробы)

- Параметр `image_url` у `create_blog_post`/`update_blog_post` → `article[image_attributes]{src}`:
  добавляет **или заменяет** картинку. Описание тула: «картинка заменяется, не складывается».
- Параметра удаления **не делаем** (был бы фейковым). В доки — явная оговорка:
  «удаление картинки статьи API не поддерживается; замените картинку или используйте админку InSales».

---

# Проба-2: загрузка картинки статьи байтами (attachment)

Дата: 2026-11-12. Магазин: `myshop-dgi816.myinsales.ru` (testshop), blog `8222025`.
Метод: прямой REST, Basic-auth. Тестовые статьи удалены.

## Вопрос

Принимает ли `article[image_attributes]` base64-`attachment` (как `image[attachment]`
у товаров), или статьям доступен только `src`? Есть ли отдельный image-ресурс у статьи
(в админке — сущность `article_image`)?

## Результаты

| Шаг | Запрос | Результат |
|-----|--------|-----------|
| 1. Create | `article{..., image_attributes:{attachment:"<b64 png>", filename:"probe-a.png"}}` | ✅ **201**, `image:{id, original_url:"/images/loading.gif"}` — attachment принят |
| 2. Replace | `PUT article{image_attributes:{attachment:"<другой b64>", filename:"probe-b.png"}}` | ✅ замена: новый image id, add-or-replace как с `src` |
| 3. Без filename | `PUT article{image_attributes:{attachment:"<b64>"}}` | ⚠️ **200**, но картинка **сломана**: `original_url:"/images/no_image_original.jpg"` (заглушка) — имя файла обязательно, без него слот затирается заглушкой |
| 4. Settle | create с attachment → GET через ~8 с | ✅ `original_url:"https://cdn.insales-shop.ru/images/articles/1/8177/16547825/probe-settle.png"` — имя файла сохраняется |
| 5. Отдельный эндпоинт | `POST /blogs/:b/articles/:p/images.json` (attachment и src) | ❌ **404** — product-стиля image-эндпоинтов у статей нет |
| 6. Сущность админки | `GET /admin/article_images/:id.json` | ✅ **200** (read-only; id совпадает с id картинки) |
| 7. Create без `published_at`/`author` | `article{title, content, image_attributes}` | ❌ **422** `{"published_at":["Не заполнены обязательные поля"],"author":["Не заполнены обязательные поля"]}` — InSales **обязывает** оба поля при создании |

## Выводы

1. **`image_attributes:{attachment, filename}` работает** на create и update — загрузка
   байтами без публичного URL возможна одним вызовом (ветка primary).
2. **`filename` обязателен** — без него InSales молча ставит заглушку
   `no_image_original.jpg` и затирает прежнюю картинку. Валидация в mcp-is: attachment
   без filename → ошибка, запрос не уходит.
3. **`published_at` и `author` обязательны при create** (InSales валидирует). Следствие
   для «без отсебятины»: mcp-is НЕ подставляет значения — оба поля становятся
   Required-параметрами `create_blog_post`; вызывающий обязан их прислать.
4. Отдельного create-канала для картинки статьи в API нет (`images.json` → 404);
   read-only `article_images` в API присутствует.
5. Имя файла сохраняется в CDN-пути (`.../16547825/probe-settle.png`) — случайное имя
   у обложки «Мужское здоровье…» = имя исходного файла на файлообменнике, что
   подтверждает ретроспективу «тогда шёл `src` с внешнего URL».

## Решение для тулов (обновление)

- `image_data` (base64) + `image_filename` у `create/update_blog_post` →
  `article[image_attributes]{attachment, filename}`: add-or-replace, как `image_url`.
- `image_filename` обязателен при `image_data` (инвариант InSales, п.2).
- `create_blog_post`: `published_at` и `author` — Required, без подстановок (п.3).
- Fallback-хостинг (`src` с временного URL) не нужен — primary-ветка работает.

---

## Лимит размера attachment (проба 2026-11-12, testshop)

Вопрос: какой максимальный размер картинки реально принимает InSales. Метод: scratch-статья
(`PROBE SIZE (DELETE ME)`), PUT `article[image_attributes]{attachment, filename:"probe-Nmb.png"}`,
PNG-шум заданного размера (файл ≈ N МБ), проверка осадки: GET статьи → `original_url` → HEAD на CDN.
Прямой REST (мимо нашего лимита 5 МБ). Уборка: обе scratch-статьи удалены.

| Размер файла | Тело запроса (base64) | PUT | Осадка на CDN |
|---|---|---|---|
| 5 МБ | 6.7 МБ | 200, 4.2 с | ✅ `probe-5mb.png`, content-length 5 237 415 |
| 10 МБ | 13.3 МБ | 200, 7.5 с | ✅ 10 483 243 |
| 15 МБ | 20.0 МБ | 200, 14.1 с | ✅ 15 723 663 |
| 20 МБ | 26.7 МБ | 200, 14.9 с | ✅ 20 962 721 |
| 22 МБ | 29.3 МБ | 200, 19.1 с | ✅ 23 075 471 |
| 24 МБ | 32.0 МБ | **413** (nginx), ×2 | — |
| 25 МБ | 33.3 МБ | **413** (nginx), ×2 | — |

### Выводы

- **InSales принимает картинки до ~22 МБ включительно** (проверено живьём), файл сохраняется
  на CDN как есть (без пережатия: content-length совпадает с размером исходного PNG).
- **Отлуп — `413 Request Entity Too Large` от nginx** (не валидация InSales): режется по размеру
  HTTP-тела, а не картинки. 29.3 МБ тела проходит, 32.0 МБ — нет; похоже на `client_max_body_size 32m`.
- Практический максимум картинки через JSON-attachment (base64) ≈ **23 МБ** (24 МБ → тело 32.0 МБ → 413).
  Ветки 25–50 МБ не прогонялись — за границей.
- У mcp-is сейчас лимит 5 МБ (раскодированных) / 8 МиБ тело — запас до потолка InSales большой;
  поднимать — отдельным решением (см. обсуждение).
- **add_image (товары) — тот же лимит** (проба 2026-11-12, scratch-товар): 20/22 МБ → 201 и файл как есть
  на CDN, 24 МБ → 413 nginx. Граница одна на весь хост (nginx `client_max_body_size`), канал тот же.
- **Решение 2026-11-12 (пользователь)**: лимит mcp-is поднят до 20 МБ раскодированных (тело 28 МиБ,
  nginx mcp-is 30m) — ступень ниже потолка InSales ~22 МБ.
