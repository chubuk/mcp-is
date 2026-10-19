# InSales API Reference

> Полный справочник API InSales. Сгенерировано из официальной документации.
> Всего эндпоинтов: 144
>
> Примеры запросов и ответов: [Insales-api-docs.md](Insales-api-docs.md).

## Содержание

- [Accessories (Опции товаров)](#accessories)
- [Account (Аккаунт)](#account)
- [ApplicationActions (Действия приложения)](#application-actions)
- [ApplicationCharges (Счета приложения)](#application-charges)
- [ApplicationWidgets (Виджеты приложения)](#application-widgets)
- [Blogs (Блоги)](#blogs)
- [BonusSystemTransactions (Бонусные транзакции)](#bonus-system-transactions)
- [Categories (Категории)](#categories)
- [ClientGroups (Группы клиентов)](#client-groups)
- [Clients (Клиенты)](#clients)
- [CollectionFields (Поля коллекций)](#collection-fields)
- [CollectionFilters (Фильтры коллекций)](#collection-filters)
- [Collections (Коллекции)](#collections)
- [Collects (Связь товар-коллекция)](#collects)
- [CustomStatuses (Кастомные статусы)](#custom-statuses)
- [DeliveryVariants (Способы доставки)](#delivery-variants)
- [dialogs](#dialogs)
- [DiscountCodes (Промокоды)](#discount-codes)
- [DiscountCollectionsGroups (Группы коллекций скидок)](#discount-collections-groups)
- [Domains (Домены)](#domains)
- [ExternalDiscounts (Внешние скидки)](#external-discounts)
- [fields](#fields)
- [Files (Файлы)](#files)
- [JsTags (JS-теги)](#js-tags)
- [marketplaces](#marketplaces)
- [Menus (Меню)](#menus)
- [OAuthApps (OAuth приложения)](#o-auth-apps)
- [OptionNames (Имена опций)](#option-names)
- [OptionValues (Значения опций)](#option-values)
- [Orders (Заказы)](#orders)
- [outlets](#outlets)
- [Pages (Страницы)](#pages)
- [PaymentGateways (Способы оплаты)](#payment-gateways)
- [PriceKinds (Типы цен)](#price-kinds)
- [ProductFields (Поля товаров)](#product-fields)
- [Products (Товары)](#products)
- [Properties (Свойства)](#properties)
- [recurring_application_charge](#recurring-application-charge)
- [Redirects (Редиректы)](#redirects)
- [Reviews (Отзывы)](#reviews)
- [SalesChannels (Каналы продаж)](#sales-channels)
- [StockCurrencies (Валюты склада)](#stock-currencies)
- [UserTasks (Задачи пользователя)](#user-tasks)
- [VariantFields (Поля вариантов)](#variant-fields)
- [Variants (Варианты товаров)](#variants)
- [Warehouses (Склады)](#warehouses)
- [Webhooks (Вебхуки)](#webhooks)

---

## Accessories (Опции товаров)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/accessories.json` | Create accessory |
| GET | `/admin/accessories.json` | Get accessories |
| POST | `/admin/accessories/74/values.json` | Create accessory value |
| GET | `/admin/accessories/70/values.json` | Get accessory values |

### POST `/admin/accessories.json`

*Create accessory*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `accessory[name]` | да | Наименование |
| `accessory[max_count]` | да | Максимальное количество опций в заказе |
| `accessory[min_count]` | да | Минимальное количество опций в заказе |
| `accessory[permalink]` | да | Идентификатор |

---

## Account (Аккаунт)

| Метод | URL | Описание |
|-------|-----|----------|
| PUT | `/admin/account.json` | Update account |

### PUT `/admin/account.json`

*Update account*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `account[client_cookies_whitelist]` | нет | cookies whitelist |

---

## ApplicationActions (Действия приложения)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/application_actions.json` | Create application_action |
| GET | `/admin/application_actions.json` | Get application actions |

### POST `/admin/application_actions.json`

*Create application_action*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `application_action[title]` | да | title to show in backoffice |
| `application_action[handle]` | да | unique string identifier |
| `application_action[entity]` | да | string, defines entity that action belongs to: order, product, client, discount_code, file, page, article, collection_filter |
| `application_action[url]` | да | URL to which POST request with serialized entities will be made |

---

## ApplicationCharges (Счета приложения)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/application_charges.json` | Create application_charge |
| POST | `/admin/application_charges/10/decline.json` |  |
| GET | `/admin/application_charges.json` | Get application_charges |

### POST `/admin/application_charges.json`

*Create application_charge*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `application_charge[name]` | да | purpose of payment |
| `application_charge[price]` | да | bill amount |
| `application_charge[return_url]` | нет | url for notification about payment |
| `application_charge[test]` | нет | flag for debugging purpose; if true, bill can be confirmed without payment; false by default |

---

## ApplicationWidgets (Виджеты приложения)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/application_widgets.json` | Create application_widget |
| GET | `/admin/application_widgets.json` | Get application_widgets |

### POST `/admin/application_widgets.json`

*Create application_widget*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `application_widget[code]` | да | html or javascript code of widget |
| `application_widget[height]` | да | height of iframe block |
| `application_widget[page_type]` | нет | product or order |

---

## Blogs (Блоги)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/blogs/1/articles.json` |  |
| GET | `/admin/blogs/1/articles.json` | Get articles list |
| POST | `/admin/blogs.json` | Create blog |
| GET | `/admin/blogs.json` | Get blogs |

### POST `/admin/blogs.json`

*Create blog*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `blog[title]` | да | title |
| `blog[html_title]` | нет | html title tag content |
| `blog[meta_keywords]` | нет | html meta words on page |
| `blog[meta_description]` | нет | html meta description |
| `blog[notify_about_comment]` | нет | true by default |
| `blog[captcha_enabled]` | нет | false by default |
| `blog[disallow_links_in_comments]` | нет | false by default |
| `blog[commentable]` | нет | default: 'no', available: no, yes, moderable |

---

## BonusSystemTransactions (Бонусные транзакции)

| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/admin/bonus_system_transactions.json` | Get bonus_system_transactions |

### GET `/admin/bonus_system_transactions.json`

*Get bonus_system_transactions*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `updated_since` | нет | Set datetime to get only data updated after it. |
| `from_id` | нет | Set id to get only data starting from it. |
| `per_page` | нет | Quantity of objects per page. |

---

## Categories (Категории)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/categories.json` | Create category |
| GET | `/admin/categories.json` | Get categories |

### POST `/admin/categories.json`

*Create category*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `category[title]` | да | title |
| `category[parent_id]` | нет | parent category id |
| `category[position]` | нет | position in categories list |

---

## ClientGroups (Группы клиентов)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/client_groups.json` | Create client group |
| GET | `/admin/client_groups.json` | Get client groups |

### POST `/admin/client_groups.json`

*Create client group*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `client_group[title]` | нет | client group title |
| `client_group[discount]` | нет | discount |
| `client_group[discount_description]` | нет | discount description |
| `client_group[is_default]` | нет | set as default group |

---

## Clients (Клиенты)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/clients/2/bonus_system_transactions.json` | Create transaction for client |
| GET | `/admin/clients/2/bonus_system_transactions.json` | Get transactions for client |
| POST | `/admin/clients.json` |  |
| GET | `/admin/clients.json` | Get clients |

### POST `/admin/clients.json`

**

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `client[name]` | да | name of company or a sole trader |
| `client[registered]` | нет | if true, email and password are required |
| `client[email]` | нет | email |
| `client[password]` | нет | password at least 6 characters |
| `client[phone]` | нет | phone number |
| `client[subscribe]` | нет | set true to subscribe client for emails |
| `client[client_group_id]` | нет | client group id |
| `client[type]` | да | Client::Juridical |
| `client[ip_addr]` | нет | IP address |
| `client[juridical_address]` | нет | address of company or a sole trader |
| `client[inn]` | да | INN |
| `client[kpp]` | нет | KPP |
| `client[ogrn]` | нет | OGRN |
| `client[okpo]` | нет | OKPO |
| `client[bik]` | нет | BIK |
| `client[bank_name]` | нет | bank name |
| `client[correspondent_account]` | нет | Correspondent account |
| `client[settlement_account]` | нет | Settlement account |
| `client[fields_values_attributes]` | нет | fields values array |
| `client[fields_values_attributes][field_id]` | нет | field id |
| `client[fields_values_attributes][value]` | нет | value of field |
| `client[skip_send_mail]` | нет | set it true to skip sending email confirmation |

### GET `/admin/clients.json`

*Get clients*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `updated_since` | нет | Set datetime to get only data updated after it. |
| `from_id` | нет | Set id to get only data starting from it. |
| `per_page` | нет | Quantity of objects per page. |
| `q` | нет | Search queue. Matches name, surname, middlename, phone, or email. |
| `phone` | нет | Exact match search by phone number. |
| `email` | нет | Exact match search by email address. |
| `registered` | нет | Filter by registration status. Accepted values: 'yes' or 'no'. |
| `subscribe` | нет | Filter by subscription permission. Accepted values: 'yes' or 'no'. |
| `rfm_group` | нет | RFM group value or array. Allowed values: 'small_bills', 'not_frequent', 'not_recent', 'best', 'worst'. |

---

## CollectionFields (Поля коллекций)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/collection_fields.json` | Create collection_field |
| DELETE | `/admin/collection_fields/handle2.json` | Destroy collection_field by handle |
| GET | `/admin/collection_fields/handle2.json` | Get collection_field by handle |
| GET | `/admin/collection_fields.json` | Get collection_fields |
| PUT | `/admin/collection_fields/handle2.json` | Update collection_field by handle |

### POST `/admin/collection_fields.json`

*Create collection_field*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collection_field[title]` | да | title |
| `collection_field[handle]` | да | handle used in liquid |
| `collection_field[type]` | да | CollectionField::TextArea - textarea, CollectionField::TextField - textfield |
| `collection_field[position]` | нет | position in collection_fields list |
| `collection_field[is_hidden]` | нет | boolean |

---

## CollectionFilters (Фильтры коллекций)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/collection_filters.json` | Create collection filter |
| GET | `/admin/collection_filters.json` | Get collection filters |

### POST `/admin/collection_filters.json`

*Create collection filter*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `filter[title]` | да | filter title |
| `filter[collection_id]` | да | collection id |
| `filter[description]` | нет | description |
| `filter[permalink]` | нет | filter permalink |
| `filter[html_title]` | нет | tag title |
| `filter[meta_keywords]` | нет | tag keyword |
| `filter[meta_description]` | нет | tag description |
| `filter[properties_attributes]` | нет | properties |
| `filter[options_attributes]` | нет | options |
| `filter[field_values_attributes]` | нет | collection filter field values attributes |

### GET `/admin/collection_filters.json`

*Get collection filters*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collection_id` | нет | collection id |

---

## Collections (Коллекции)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/collections.json` | Create collection |
| GET | `/admin/collections.json` | Get collections |
| GET | `/admin/collections/1/collection_field_values/handle2.json` | Get collection field value by handle |
| GET | `/admin/collections/1/collection_field_values.json` | Get collection fields values |
| PUT | `/admin/collections/1/collection_field_values/handle2.json` | Update collection field value by handle |

### POST `/admin/collections.json`

*Create collection*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collection[title]` | да | title |
| `collection[parent_id]` | да | parent collection id |
| `collection[position]` | нет | position in collections list |
| `collection[sort_type]` | нет | Alphabetically: A - Z => 1, Alphabetically: Z - A => 2, By date: New - Old => 3, By date: Old - New => 4, By price: Expensive - Cheap => 5, By price: Cheap - Expensive => 6, Manual => 7, By discount: Large - Little => 8, By discount: Little - Large => 9, By sort rank: Low - High => 10, By sort rank: High - Low => 11, Sort by ascending popularity => 12, Sort by descending popularity => 13 |
| `collection[description]` | нет | description |
| `collection[html_title]` | нет | html title |
| `collection[seo_description]` | нет | seo description |
| `collection[meta_description]` | нет | meta description |
| `collection[meta_keywords]` | нет | meta keywords |
| `collection[is_hidden]` | нет | set true to hide collection |
| `collection[permalink]` | нет | generated |
| `collection[recursive]` | нет | recursive |
| `collection[image_attributes]` | нет | image for collection |
| `collection[image_attributes][src]` | нет | image url |
| `collection[field_values_attributes]` | нет | collection field values attributes |

### GET `/admin/collections.json`

*Get collections*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `updated_since` | нет | Set datetime to get only data updated after it. |
| `from_id` | нет | Set id to get only data starting from it. |
| `per_page` | нет | Quantity of objects per page. |

---

## Collects (Связь товар-коллекция)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/collects.json` | Add product to collection |
| GET | `/admin/collects.json` | Get collections by product |

### POST `/admin/collects.json`

*Add product to collection*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `collect[collection_id]` | да | collection id |
| `collect[product_id]` | да | product id |

---

## CustomStatuses (Кастомные статусы)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/custom_statuses.json` |  |
| DELETE | `/admin/custom_statuses/novyy.json` |  |
| GET | `/admin/custom_statuses/novyy.json` |  |
| GET | `/admin/custom_statuses.json` |  |
| PUT | `/admin/custom_statuses/novyy.json` |  |

### POST `/admin/custom_statuses.json`

**

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `custom_status[system_status]` | да | system status |
| `custom_status[title]` | да | description |

---

## DeliveryVariants (Способы доставки)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/delivery_variants.json` |  |
| GET | `/admin/delivery_variants.json` | Get delivery_variants |

### POST `/admin/delivery_variants.json`

**

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `delivery_variant[title]` | да | title |
| `delivery_variant[type]` | да | DeliveryVariant::External, DeliveryVariant::Fixed, DeliveryVariant::LocationDepend, DeliveryVariant::MapZoneDelivery, DeliveryVariant::NoDelivery, DeliveryVariant::None, DeliveryVariant::OfficialEmsShipping, DeliveryVariant::OfficialFirstClass, DeliveryVariant::OfficialRussianpost, DeliveryVariant::PickUp, DeliveryVariant::PriceDepend |
| `delivery_variant[description]` | нет | description |
| `delivery_variant[position]` | нет | position |
| `delivery_variant[add_payment_gateways]` | нет | set true to add all payment variants |
| `delivery_variant[delivery_locations_attributes]` | нет | array of delivery locations |
| `delivery_variant[delivery_locations_attributes][country]` | нет | if shop has many countries, fill this field;it should be enabled in checkout options to be binded to shipping method |
| `delivery_variant[delivery_locations_attributes][region]` | нет | field Region should be enabled in checkout options to be binded to shipping method |
| `delivery_variant[delivery_locations_attributes][city]` | нет | field Region should be filled |
| `delivery_variant[inverted]` | нет | false by default, if true this delivery variant will be available for all regions except passed in parameters 'country', 'region' and 'city' |
| `delivery_variant[charge_up_to]` | нет | free on orders from given sum |
| `delivery_variant[javascript]` | нет | if needed, paste javascript to be displayed on page of delivery variant selection |
| `delivery_variant[show_calendar_in_checkout]` | нет | show delivery date selection in checkout |
| `delivery_variant[delivery_date_required]` | да | delivery date selection is required in chekout |
| `delivery_variant[forbid_weekends]` | нет | weekend days are forbidden for delivery |
| `delivery_variant[estimated_delivery_period]` | нет | estimated delivery period |
| `delivery_variant[order_time_is_later]` | нет | increase estimated delivery time for one day if order time is later |
| `delivery_variant[forbidden_days]` | нет | forbidden days for delivery (23.02,08.03,..) |
| `delivery_variant[show_time_intervals_in_checkout]` | нет | show delivery time intervals in checkout |
| `delivery_variant[time_intervals_required]` | да | delivery time interval selection is required in checkout |
| `delivery_variant[time_intervals]` | нет | time intervals (10-12, 12-16, ...) |
| `delivery_variant[pick_up_sources_attributes]` | нет | array of pick up sources |
| `delivery_variant[pick_up_sources_attributes][title]` | да | title |
| `delivery_variant[pick_up_sources_attributes][http_method]` | да | http method (GET, POST) used for request |
| `delivery_variant[pick_up_sources_attributes][url]` | да | url for receiving list of delivery points |
| `delivery_variant[pick_up_sources_attributes][point_info_url]` | да | url for receiving delivery info to the selected point |

---

## dialogs

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/dialogs/provider_tokens.json` | Create or refresh messengers custom provider token |

### POST `/admin/dialogs/provider_tokens.json`

*Create or refresh messengers custom provider token*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `messenger_type` | да | Messenger type	["whats_app", "telegram", "viber", "sms", "email", "vk", "instagram"] |

---

## DiscountCodes (Промокоды)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/discount_codes.json` | Create discount code |
| GET | `/admin/discount_codes.json` | Get discount codes |

### POST `/admin/discount_codes.json`

*Create discount code*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `discount_code[code]` | да | code |
| `discount_code[description]` | нет | description |
| `discount_code[disabled]` | нет | disabled |
| `discount_code[act_once]` | нет | act once |
| `discount_code[act_once_for_client]` | нет | act once for client |
| `discount_code[expired_at]` | нет | expired at |
| `discount_code[type_id]` | да | type_id (1 - percent, 2 - money) |
| `discount_code[discount]` | да | discount |
| `discount_code[min_price]` | нет | min price |

### GET `/admin/discount_codes.json`

*Get discount codes*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `updated_since` | нет | Set datetime to get only data updated after it. |
| `from_id` | нет | Set id to get only data starting from it. |
| `q` | нет | Strict search by code or regexp for description. |
| `per_page` | нет | Quantity of objects per page. |
| `page` | нет | Page number. |

---

## DiscountCollectionsGroups (Группы коллекций скидок)

| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/admin/discount_collections_groups/220/discount_codes.json` | Get discount codes for discount collections group |
| GET | `/admin/discount_collections_groups.json` | Get discount collections groups |

---

## Domains (Домены)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/domains.json` | Create domain |
| GET | `/admin/domains.json` | Get domains |

### POST `/admin/domains.json`

*Create domain*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `domain[domain]` | да | domain |
| `domain[main]` | нет | boolean, could be false only if there are more than one domain |

---

## ExternalDiscounts (Внешние скидки)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/external_discounts.json` | Create external discount |
| GET | `/admin/external_discounts.json` | Get external discounts |

### POST `/admin/external_discounts.json`

*Create external discount*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `external_discount[description]` | нет | description |
| `external_discount[url]` | нет | external service url |
| `external_discount[order_lines_depend]` | нет | discount calculation depends on order lines |
| `external_discount[active]` | нет | active |
| `external_discount[observed_fields]` | нет | discount calculation depends on this order fields (array of field ids) |

---

## fields

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/fields.json` | Create select field |
| GET | `/admin/fields.json` | Get fields |

### POST `/admin/fields.json`

*Create select field*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `field[type]` | да | Field::Select |
| `field[office_title]` | да | title for backoffice |
| `field[destiny]` | да | integer, defines subject that field belongs to: 1 - ShippingAddress, 2 - Client::Individual, 5 - Client::Juridical, 3 - Order, 4 - OrderLine, 6 - Client |
| `field[handle]` | нет | You can create or update value by handle |
| `field[title]` | нет | title |
| `field[active]` | нет | true by default, defines ability of filling field by a customer |
| `field[for_buyer]` | нет | true by defaut, defines visibility for a customer |
| `field[obligatory]` | нет | false by default |
| `field[show_in_result]` | нет | true by default, field will be shown in result |
| `field[show_in_checkout]` | нет | true by default, field will be shown in checkout |
| `field[position]` | нет | position in fields list |
| `field[is_indexed]` | нет | boolean, true if searchable |
| `field[hide_in_backoffice]` | нет | boolean, true to hide in backoffice |
| `field[field_options_attributes]` | нет | array of field options objects |

---

## Files (Файлы)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/files.json` | Create file from attachment |
| GET | `/admin/files.json` | Get files |

### POST `/admin/files.json`

*Create file from attachment*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `file[attachment]` | да | file |
| `file[filename]` | да | file name |

---

## JsTags (JS-теги)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/js_tags.json` |  |
| GET | `/admin/js_tags.json` | Get js_tags |

### POST `/admin/js_tags.json`

**

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `js_tag[type]` | да | JsTag::TextTag |
| `js_tag[content]` | да | javascript code without tag SCRIPT |
| `js_tag[name]` | нет | name for js tag |

---

## marketplaces

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/marketplaces.json` | Create Model for Yandex marketplace |
| GET | `/admin/marketplaces.json` | Get marketplaces |

### POST `/admin/marketplaces.json`

*Create Model for Yandex marketplace*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `marketplace[name]` | да | Name |
| `marketplace[type]` | да | Marketplace type, e.g. Marketplace::ModelYandexMarket |
| `marketplace[shop_name]` | нет | Shop name |
| `marketplace[shop_company]` | нет | Company name |
| `marketplace[description_type]` | нет | 1 fo short description, 2 for full |
| `marketplace[vendor_id]` | да | id of corresponding property |
| `marketplace[use_sku_as_vendor_code]` | нет | Use sku as vendor code |
| `marketplace[pickup]` | нет | Pickup enable |
| `marketplace[store]` | нет | Offline store exists |
| `marketplace[delivery_style]` | нет | Delivery pricing type, 0 for old, 1 for new |
| `marketplace[local_delivery_cost]` | нет | Local delivery cost |
| `marketplace[adult]` | нет | 1 if products are adult-only, else 0 |
| `marketplace[utm_label]` | нет | UTM label |
| `marketplace[page_encoding]` | да | Encoding, windows-1251 or utf-8 |
| `marketplace[image_style]` | да | Image style: thumb, compact, medium, large or original |
| `marketplace[model_type]` | да | Model naming type: name or sku |
| `marketplace[collection_ids]` | нет | Collections to load |

---

## Menus (Меню)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/menus.json` | Create menu |
| GET | `/admin/menus.json` | Get menus list |
| POST | `/admin/menus/1/items.json` | Add blog to menu |
| GET | `/admin/menus/1/items.json` | Get menu items list |

### POST `/admin/menus.json`

*Create menu*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `menu[name]` | да | menu text |
| `menu[handle]` | да | handle used in liquid |

---

## OAuthApps (OAuth приложения)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/o_auth_apps.json` | Create o_auth_app |
| GET | `/admin/o_auth_apps.json` | Get o_auth_apps |

### POST `/admin/o_auth_apps.json`

*Create o_auth_app*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `o_auth_app[provider_handle]` | да | provider name |
| `o_auth_app[client_id]` | да | client id |
| `o_auth_app[client_secret]` | да | client_secret |
| `o_auth_app[issuer]` | да | issuer |

---

## OptionNames (Имена опций)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/option_names.json` | Create option_name |
| GET | `/admin/option_names.json` | Get option_names |
| POST | `/admin/option_names/1/option_values.json` | Create option_value |
| DELETE | `/admin/option_names/1/option_values/20/destroy_image.json` | Destroy option_value image |
| GET | `/admin/option_names/1/option_values/20/destroy_status.json` | Get destroy option_value status |
| GET | `/admin/option_names/1/option_values.json` | Get option_values for specific option |

### POST `/admin/option_names.json`

*Create option_name*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `option_name[title]` | да | title |
| `option_name[position]` | нет | position in option names list |
| `option_name[navigational]` | нет | true by default; show in products filter on category page |

---

## OptionValues (Значения опций)

| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/admin/option_values.json` | Get option_values for all options |

---

## Orders (Заказы)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/orders.json` | Create order with bundle |
| GET | `/admin/orders/count.json` | Get orders count |
| GET | `/admin/orders.json` | Get orders |

### POST `/admin/orders.json`

*Create order with bundle*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `order[order_lines_attributes]` | да | order lines attributes |
| `order[order_lines_attributes][][quantity]` | да | quantity |
| `order[order_lines_attributes][][variant_id]` | да | variant id |
| `order[order_lines_attributes][][vat]` | нет | do not set vat for bundle, because after order creation bundle will break into bundle components, which have their own VAT rates |
| `order[client]` | да | client info |
| `order[client][name]` | нет | client name |
| `order[client][email]` | нет | client email |
| `order[client][phone]` | да | client phone |
| `order[client][consent_to_personal_data]` | да | consent to personal data |
| `order[shipping_address_attributes]` | да | array of shipping address attributes |
| `order[shipping_address_attributes][address]` | нет | shipping address |
| `order[shipping_address_attributes][full_locality_name]` | да | full locality name |
| `order[delivery_variant_id]` | да | delivery variant id |
| `order[payment_gateway_id]` | да | payment gateway id |
| `order[coupon]` | нет | discount code |

### GET `/admin/orders/count.json`

*Get orders count*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `fulfillment_status` | нет | array of statuses |
| `delivery_variant` | нет | array of delivery variant ids |
| `payment_gateway_id` | нет | array of payment gateway ids |

### GET `/admin/orders.json`

*Get orders*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `fulfillment_status` | нет | array of statuses |
| `delivery_variant` | нет | array of delivery variant ids |
| `payment_gateway_id` | нет | array of payment gateway ids |
| `updated_since` | нет | Set datetime to get only data updated after it. |
| `from_id` | нет | Set id to get only data starting from it. |
| `per_page` | нет | Quantity of objects per page. |

---

## outlets

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/outlets.json` | Create outlet |
| GET | `/admin/outlets.json` | Get outlets |

### POST `/admin/outlets.json`

*Create outlet*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `outlet[title]` | да | title |
| `outlet[description]` | нет | description |
| `outlet[address]` | да | address in text way |
| `outlet[price]` | да | delivery price |
| `outlet[latitude]` | да | point latitude |
| `outlet[longitude]` | да | point longitude |

---

## Pages (Страницы)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/pages.json` | Create page |
| GET | `/admin/pages.json` | Get pages list |

### POST `/admin/pages.json`

*Create page*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `page[title]` | да | page title |
| `page[content]` | да | page content |
| `page[permalink]` | нет | permalink on /page/* |
| `page[html_title]` | нет | html title tag content |
| `page[meta_keywords]` | нет | html meta words on page |
| `page[meta_description]` | нет | html meta description |
| `page[created_at]` | нет | publication date |

---

## PaymentGateways (Способы оплаты)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/payment_gateways.json` | Create External payment_gateway |
| GET | `/admin/payment_gateways.json` | Get payment_gateways |

### POST `/admin/payment_gateways.json`

*Create External payment_gateway*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `payment_gateway[title]` | да | title |
| `payment_gateway[margin]` | нет | margin measured in percentage and added to order price with current payment method; it should fall in the range from -100 to 1000 |
| `payment_gateway[position]` | нет | position in payment gateways list |
| `payment_gateway[type]` | да | PaymentGateway::External |
| `payment_gateway[description]` | нет | description |
| `payment_gateway[payment_delivery_variants_attributes]` | нет | array of available delivery variants for this payment gateway; payment gateway is unavailable if list is empty |
| `payment_gateway[payment_delivery_variants_attributes][delivery_variant_id]` | нет | id of delivery variant |
| `payment_gateway[url]` | да | url of external payment gateway |
| `payment_gateway[shop_id]` | да | shop id - store identificator for external systems |
| `payment_gateway[send_order]` | нет | set true to send order_lines information to the "url"; parameter name - "order_json"; do not forget add "order_json" to signature validation to the end of the string. |
| `payment_gateway[add_delivery_variants]` | нет | set all delivery variants available for this payment gateway |
| `payment_gateway[redirect_to_payment]` | нет | true by default; if true automatically redirect to payment after successful order creation; |
| `payment_gateway[convert_currency]` | нет | iso code, if you would like to convert order amount from site currency to another currency |

---

## PriceKinds (Типы цен)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/price_kinds.json` | Create price kind |
| GET | `/admin/price_kinds.json` | Get price kinds |

### POST `/admin/price_kinds.json`

*Create price kind*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `price_kind[title]` | да | title |
| `price_kind[show_in_lists]` | нет | show price in products list in back-office |

---

## ProductFields (Поля товаров)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/product_fields.json` | Create product_field |
| DELETE | `/admin/product_fields/handle2.json` | Destroy product_field by handle |
| GET | `/admin/product_fields/handle2.json` | Get product_field by handle |
| GET | `/admin/product_fields.json` | Get product_fields |
| PUT | `/admin/product_fields/handle2.json` | Update product_field by handle |

### POST `/admin/product_fields.json`

*Create product_field*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `product_field[title]` | да | title |
| `product_field[handle]` | да | handle used in liquid |
| `product_field[type]` | да | ProductField::TextArea - textarea, ProductField::TextField - textfield |
| `product_field[position]` | нет | position in product_fields list |
| `product_field[is_hidden]` | нет | boolean |
| `product_field[is_indexed]` | нет | boolean |

---

## Products (Товары)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/products/293977313/product_accessory_links.json` | Create product accessory link |
| GET | `/admin/products/293977297/product_accessory_links.json` | Get product accessory links |
| POST | `/admin/products/1/images.json` |  |
| GET | `/admin/products/1/images.json` | Get images |
| POST | `/admin/products.json` | Create bundle |
| GET | `/admin/products/count.json` | Get products count |
| GET | `/admin/products.json` | Get products |
| POST | `/admin/products/1/product_field_values.json` | Create product_field_value by handle |
| DELETE | `/admin/products/1/product_field_values/handle.json` | Destroy product_field_value by handle |
| GET | `/admin/products/1/product_field_values/handle.json` | Get product_field_value by handle |
| GET | `/admin/products/1/product_field_values.json` | Get product_field_values |
| PUT | `/admin/products/1/product_field_values/handle.json` | Update product_field_value by handle |
| POST | `/admin/products/1/similars.json` | Create similar |
| GET | `/admin/products/1/similars.json` | Get similars |
| POST | `/admin/products/1/supplementaries.json` | Create supplementary |
| GET | `/admin/products/1/supplementaries.json` | Get supplementaries |
| POST | `/admin/products/1/variants.json` | Create variant |
| GET | `/admin/products/1/variants.json` | Get variants |
| PUT | `/admin/products/variants_group_update.json` | Update group of variants (multi warehouses case) |

### POST `/admin/products.json`

*Create bundle*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `product[category_id]` | да | category id |
| `product[title]` | да | title |
| `product[description]` | нет | description |
| `product[short_description]` | нет | short description |
| `product[sort_weight]` | нет | sort weight |
| `product[bundle]` | да | bundle |
| `product[variants_attributes]` | да | variants attributes |
| `product[vat]` | нет | do not set VAT for bundle because order lines will be created with its components, which have their own VAT |
| `product[product_bundle_components_attributes]` | да | product bundle components attributes |

### GET `/admin/products.json`

*Get products*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `updated_since` | нет | Set datetime to get only data updated after it. |
| `from_id` | нет | Set id to get only data starting from it. |
| `per_page` | нет | Quantity of objects per page. |
| `with_deleted` | нет | Set to true to get also deleted products. |
| `deleted` | нет | Set to true to get only deleted products. |
| `variant_fields` | нет | List of required fields separated with `,`. Possible values are: id,title,product_id,sku,barcode,dimensions,available. |

### PUT `/admin/products/variants_group_update.json`

*Update group of variants (multi warehouses case)*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `variants` | нет | Array with parameters for each variant	array |
| `variants[][id]` | да | Id of a variant	integer |
| `variants[][quantity]` | нет | Variant quantity	float |
| `variants[][cost_price]` | нет | Cost price of the variant	float |
| `variants[][old_price]` | нет | Old price of the variant	float |
| `variants[][cost_price]` | нет | Cost price of the variant	float |
| `variants[][price]` | нет | First price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price2]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price3]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price4]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price5]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price6]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price7]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price8]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price9]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price10]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price11]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price12]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price13]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price14]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price15]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price16]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price17]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price18]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price19]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][price20]` | нет | n-th price of the price kinds of the variant (if price kinds are enabled)	float |
| `variants[][quantity_at_warehouses]` | нет | Array of quantity at each of the warehouses. Could be replaced with quantity_at_warehousesN fields, which have higher priority	array	{"type" => "float"} |

---

## Properties (Свойства)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/properties/3/characteristics.json` | Create characteristic |
| GET | `/admin/properties/3/characteristics.json` | Get characteristics |
| POST | `/admin/properties.json` | Create property |
| GET | `/admin/properties.json` | Get properties |

### POST `/admin/properties.json`

*Create property*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `property[title]` | да | title |
| `property[permalink]` | нет | permalink |
| `property[position]` | нет | position in property list |
| `property[backoffice]` | нет | show in backoffice |
| `property[is_hidden]` | нет | hide in product page on site |
| `property[is_navigational]` | нет | show in products filter on category page |

---

## recurring_application_charge

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/recurring_application_charge.json` |  |
| DELETE | `/admin/recurring_application_charge.json` | Destroy recurring_application_charge |
| GET | `/admin/recurring_application_charge.json` | Get recurring_application_charge |
| PUT | `/admin/recurring_application_charge.json` | Update recurring_application_charge |
| PATCH | `/admin/recurring_application_charge/add_free_days.json` | add free days to recurring_application_charge |

### POST `/admin/recurring_application_charge.json`

**

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `recurring_application_charge[monthly]` | да | amount of monthly payment |
| `recurring_application_charge[trial_expired_at]` | нет | trial period ending date, by default it uses trial period from application |

### PUT `/admin/recurring_application_charge.json`

*Update recurring_application_charge*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `recurring_application_charge[monthly]` | нет | amount of monthly payment |

### PATCH `/admin/recurring_application_charge/add_free_days.json`

*add free days to recurring_application_charge*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `recurring_application_charge[days_count]` | нет | how many free days you want to add, it is possible that paid_till attribute would change but trial_expired_at would not, it is normal |

---

## Redirects (Редиректы)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/redirects.json` | Create redirect with removing redirect chains |
| GET | `/admin/redirects.json` | Get redirects |

### POST `/admin/redirects.json`

*Create redirect with removing redirect chains*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `redirect[old_url]` | да | old url |
| `redirect[new_url]` | да | new url |
| `redirect[remove_redirect_chains]` | нет | if set to 1 and given a redirect A=>B when creating a new redirect B=>C, it will destroy old redirect (A=>B) and create new one (A=>C) |

---

## Reviews (Отзывы)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/reviews.json` | Create review |
| GET | `/admin/reviews.json` | Get reviews |

### POST `/admin/reviews.json`

*Create review*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `review[product_id]` | да | product id |
| `review[author]` | да | author |
| `review[content]` | нет | content |
| `review[email]` | нет | email |
| `review[published]` | нет | should be published? |
| `review[rating]` | нет | product rating (1-10) |

---

## SalesChannels (Каналы продаж)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/sales_channels.json` | Create sales channel |
| GET | `/admin/sales_channels.json` | Get sales channels |

### POST `/admin/sales_channels.json`

*Create sales channel*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `sales_channel[title]` | да | sales channel title |

---

## StockCurrencies (Валюты склада)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/stock_currencies.json` |  |
| GET | `/admin/stock_currencies.json` | Get stock currencies |

### POST `/admin/stock_currencies.json`

**

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `stock_currency[code]` | да | Available codes:LVL, BYN, USD, EUR, EEK, LTL, AMD, KGS, MDL, EGP, CAD, JPY, TNB, TRY, AZN, PLN, ILS, CNY, BGN, GBP, RUR, UAH, AUD, GEL, AED, CHF, KZT, UZS, TJS, VND, TMT, MXN |
| `stock_currency[is_default]` | нет | false by default, set it true to make this currency default for account |
| `stock_currency[exchange_rate_use_cb]` | нет | true by default, set it false to use manual exchange rate |
| `stock_currency[exchange_rate_manual]` | нет | manual exchange rate |
| `stock_currency[exchange_rate_percent]` | нет | 0 by default, deviation as a percentage from defined exchange rate |

---

## UserTasks (Задачи пользователя)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/user_tasks.json` | Create user task |
| GET | `/admin/user_tasks.json` | Get user tasks list |

### POST `/admin/user_tasks.json`

*Create user task*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `user_task[title]` | да | user task title |
| `user_task[description]` | да | user task description |
| `user_task[assignee_id]` | да | user id |
| `user_task[deadline_at]` | нет | deadline |
| `user_task[complete]` | нет | is task completed? |
| `user_task[order_id]` | нет | order id |

---

## VariantFields (Поля вариантов)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/variant_fields.json` | Create variant_field |
| DELETE | `/admin/variant_fields/handle2.json` | Destroy variant_field by handle |
| GET | `/admin/variant_fields/handle2.json` | Get variant_field by handle |
| GET | `/admin/variant_fields.json` | Get variant_fields |
| PUT | `/admin/variant_fields/handle2.json` | Update variant_field by handle |

### POST `/admin/variant_fields.json`

*Create variant_field*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `variant_field[title]` | да | title |
| `variant_field[handle]` | да | handle used in liquid |
| `variant_field[type]` | да | VariantField::TextArea - textarea, VariantField::TextField - textfield |
| `variant_field[position]` | нет | position in variant_fields list |
| `variant_field[is_hidden]` | нет | boolean |

---

## Variants (Варианты товаров)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/variants/1/variant_field_values.json` | Create variant_field_value by handle |
| DELETE | `/admin/variants/1/variant_field_values/handle.json` | Destroy variant_field_value by handle |
| GET | `/admin/variants/1/variant_field_values/handle.json` | Get variant_field_value by handle |
| GET | `/admin/variants/1/variant_field_values.json` | Get variant_field_values |
| PUT | `/admin/variants/1/variant_field_values/handle.json` | Update variant_field_value by handle |

---

## Warehouses (Склады)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/warehouses.json` | Create warehouse |
| GET | `/admin/warehouses.json` | Get warehouses list |

### POST `/admin/warehouses.json`

*Create warehouse*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `warehouse[title]` | да | warehouse title |
| `warehouse[position]` | нет | warehouse position (1 - main priority) |

---

## Webhooks (Вебхуки)

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/admin/webhooks.json` | Create webhook |
| GET | `/admin/webhooks.json` | Get webhooks |

### POST `/admin/webhooks.json`

*Create webhook*

| Параметр | Обязательный | Описание |
|----------|:------------:|----------|
| `webhook[address]` | да | address |
| `webhook[topic]` | да | Topics: 'orders/create', 'orders/update', 'orders/destroy', 'products/create', 'products/update', 'client/create', 'client/update' |
| `webhook[format_type]` | нет | format of data sent by webhook: 'json' or 'xml' |
| `webhook[warehouse_id]` | нет | For 'orders/create', 'orders/update', 'orders/destroy' topics only, filters sent orders by warehouse |
| `webhook[sales_channel_id]` | нет | For 'products/update' topic only, filters sent products by sales channel |
| `webhook[batch_size]` | нет | For 'products/create', 'products/update' topics only, products are sent in batches, default value: 10 |

---
