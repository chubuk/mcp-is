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
