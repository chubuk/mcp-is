# Привязка аккаунта при установке приложения из маркетплейса: коллизия e-mail и механизмы первой привязки

> Generated 2026-10-29 · depth: standard · 31 sources · workspace: research/app-install-account-linking/

## Executive summary

- Крупные маркетплейсы приложений (Shopify, Bitrix24, Wix, Tilda) **не передают e-mail как ключ идентичности** — только платформенные ID (`sub`+`dest`, `member_id`, user/app-instance IDs); сценария «e-mail совпал с аккаунтом вендора» в их протоколах просто нет — это 100% политика вендора приложения [1][3][6][9].
- InSales — исключение по форме: `user_email` передаётся, но **неподписанным** параметром «URL входа»; официальные доки прямо требуют подтверждать пользователя через автологин (`token3`), а выбор механизма первой привязки к своему аккаунту оставляют интегратору [10][11].
- Авто-линк по совпадению e-mail без подтверждения владения = атака **Classic-Federated Merge** (USENIX Security 2022): ≥35 из 75 популярных сервисов были уязвимы хотя бы к одной pre-hijacking атаке [14].
- NIST SP 800-63C-4 (2025-07) нормативно требует: линковка — только в аутентифицированной сессии; нельзя ассоциировать аккаунт с идентификатором, не принадлежащим подписчику; о факте привязки — уведомлять пользователя [15].
- Индустрия принимает **два равнозначных proof-of-control**: пароль существующего аккаунта ИЛИ подтверждение доступа к почте (OTP/magic link); именно OTP/magic-link — общий механизм «proof-of-inbox» у всех провайдеров [14][18][27]. Требования к OTP: ≥6 цифр, одноразовый, срок 10 минут [16].
- Ни один крупный провайдер (Auth0, Firebase, Clerk, Better Auth) **не делает автоматический merge аккаунтов**; авто-линк разрешён только при верифицированной почте с обеих сторон и доверенном провайдере [18][19][22][23].
- UX-консенсус: dead-end «Аккаунт уже существует» нелюбим пользователем; предпочтительный экран — «подтвердите и привяжите» (вход+линк или код на почту) с фолбэком на пароль [29][30].
- **Для mcp-is**: текущий `pending_link`+пароль безопасен и соответствует практикам, но UX можно улучшить страницей «Введите код из письма» (OTP от нашего сервера) с фолбэком на пароль. Это индустриальный «middle path» (модель Clerk/WorkOS) [19][27].
- «Одна почта — много магазинов/ключей» — это **не коллизия**, а штатная мультимагазинная модель; «коллизия» касается только доказательства владения аккаунтом при первой привязке, после которого новый магазин спокойно линкуется к тому же аккаунту.

## Background & scope

Вопрос: как платформы с маркетплейсами приложений и identity-провайдеры решают привязку аккаунта при установке/первом входе из приложения маркетплейса — что делать при совпадении e-mail с существующим аккаунтом, какие passwordless-механизмы и anti-takeover-защиты существуют. Границы: задокументированные flows и security-модели (не общий курс OAuth, не наш код). Актуальность: практики 2022–2026. Полный контекст — `brief.md`; доказательная база — `findings/F1–F5.md` (62 находки).

## 1. Как маркетплейсы передают идентичность

**Shopify** выдаёт короткоживущий ID token (JWT) с claims `sub` (числовой ID сотрудника) + `dest` (магазин) — e-mail в токене отсутствует; онлайн-токены привязаны к конкретному staff-пользователю, давшему авторизацию [1][2]. **Bitrix24** OAuth-redirect несёт `code`, `domain`, `member_id`, `scope` — без e-mail; идентичность = (member_id, домен портала), вызовы идут от имени авторизованного пользователя [3]. Более того, Bitrix24 прямо проектирует OAuth для случая «пользователь уже авторизован в вашем сервисе» — привязка к аккаунту вендора считается зоной ответственности вендора [3], а для тиражных приложений может создать «системного пользователя», независимого от сотрудника-установщика [4]. **Wix** ведёт OAuth-выдачу прав целиком сам, модель идентичности — стабильные платформенные ID и роли [6][7][8]. **Tilda** маркетплейса приложений не имеет вовсе — интеграция через пару ключей аккаунта `publickey`/`secretkey` [9].

**Вывод**: поведение при совпадении e-mail с существующим аккаунтом сервиса вендора у маркетплейсов не определено — это осознанная политика самого вендора приложения [3]. InSales тут отличается: он передаёт `user_email`/`user_id` в «URL входа» — но как неподписанные параметры, не как доказательство [10].

## 2. InSales-специфика

Колбэк установки несёт только `shop`, `insales_id`, `token` (`password = MD5(token + secret)`) — user-identity на установке нет ни в одной версии доков [10][12]. Идентичность появляется в «URL входа» (`user_id`, `user_email`, `user_name`, `email_confirmed`), и официальная рекомендация прямым текстом: «Самого получения этих параметров недостаточно для входа: если сессии приложения ещё нет, подтвердите пользователя через автологин InSales» [11]. Единственный официальный механизм верификации — `token3 = MD5(LOGIN_TOKEN + user_email + user_name + user_id + email_confirmed + password)`: сессию создавать только после успешной проверки, токен одноразовый [10][11]. `email_confirmed` — литеральная строка `"true"`/`"false"`: InSales прямо сообщает, что почта может быть неподтверждённой [11][13]. Назначение token3 — доказать, что пользователь действительно залогинен в бэк-офис InSales (proof-of-InSales-session), а **не** proof-of-email-ownership [12]. Модель сессии «первый вход — автологин, дальше сессия без автологина» уже заложена в доках [11]; выбор пароль/OTP/авто-линк для первой привязки к аккаунту приложения официальными доками **не регулируется** [10].

## 3. Anti-takeover: почему совпадения e-mail недостаточно

Класс «account pre-hijacking» (USENIX Security 2022): атакующий, зная только e-mail жертвы, выполняет подготовительные действия до создания аккаунта жертвой; описано 5 атак, ≥35 из 75 популярных сервисов уязвимы хотя бы к одной (56 уязвимостей на 252 попытки, тесты Jan–Jun 2021) [14]. **Classic-Federated Merge** — точный аналог нашего сценария: атакующий регистрирует «классический» аккаунт с e-mail жертвы и своим паролем; при входе жертвы через федерацию сервис линкует по совпадению e-mail — и у атакующего сохраняется доступ [14]. Вариант «Non-verifying IdP»: IdP без верификации почты (кейс OneLogin/Okta → Zoom) делает авто-линк по e-mail из утверждения IdP прямым вектором захвата [14]. Корень проблемы — отсутствие проверки владения идентификатором до начала использования аккаунта; обязательное подтверждение ДО любых действий митигирует все 5 атак [14]. Рекомендованная процедура первой привязки: при наличии «классического» аккаунта требовать «provide or reset the password for the classic account» до объединения [14].

Нормативно: NIST SP 800-63C-4 — все функции линковки обязаны выполняться в аутентифицированной сессии; RP обязан не ассоциировать аккаунт с федеративным идентификатором, не принадлежащим подписчику; о добавлении/удалении идентификатора — уведомлять подписчика [15]. NIST SP 800-63B-4 задаёт параметры OTP: ≥6 десятичных цифр, одноразовый, недействителен через 10 минут [16]. OWASP: владение e-mail подтверждать до активации аккаунта, токены одноразовые и с ограниченным временем жизни; смена e-mail = смена идентичности (re-auth + подтверждение нового адреса); единая нормализация e-mail для всех flows [17].

## 4. Что делают identity-провайдеры при совпадении e-mail

- **Supabase** — единственный с явным авто-merge по одинаковому e-mail, но только к уникальным верифицированным адресам и с принудительным снятием неподтверждённых идентичностей (анти pre-hijacking); второй аккаунт с занятым e-mail не создаётся (анти-enumeration) [18].
- **Clerk** — авто-линк при верифицированной почте с обеих сторон, причём OAuth-линк обходит пароль; при неверифицированной почте — принудительная верификация/смена пароля перед линком [19]. Strict-режим не раскрывает существование аккаунта анониму [28].
- **Firebase** — блокирует `auth/account-exists-with-different-credential`, линк только после proof-of-control существующего провайдера; merge — обязанность приложения [20][21]. Авто-линк — только trusted+trusted (Google для @gmail.com и т.п.), смешанные пары дают ошибку «requiring account linking» осознанно — ради защиты от захвата [22].
- **Auth0** — авто-merge отсутствует; «Suggested Account Linking»: показать совпавшие аккаунты, запросить credentials целевого аккаунта, затем линк; «Every manual account link should prompt the user to enter credentials», линк — только к верифицированному аккаунту [23][24][25][26].
- **WorkOS AuthKit** — e-mail как уникальный ключ, но источник истины — доступ к инбоксу: новый кред прикрепляется к существующему юзеру только после верификации доступа к почте, иначе flow не завершается [27].

Сквозное правило всех пяти: **не линковать против неверифицированного e-mail**; верификация (OTP/magic link) — общий proof-of-inbox механизм, конкурирующий с паролем [18][19][27].

## 5. UX-практика «user already exists»

Никто не делает автоматический merge профилей: у Auth0 метаданные secondary-аккаунта отбрасываются, merge — только вручную [23]; Firebase merge вынесен в приложение [21]. Дефолты фреймворков безопасны, но UX-неудобны: Auth.js блокирует OAuth-вход на занятый e-mail (`OAuthAccountNotLinked`), авто-линк — за «dangerous»-флагом, community открыто жалуется на цену такого дефолта [29]; django-allauth шлёт письмо «Account Already Exists» с password-reset — issue называет это misleading и предлагает экран «Login and link this provider» + «Use password recovery» [30]. Better Auth предоставляет `disableImplicitLinking` именно для сценария «подтвердить линк на странице настроек, а не молча при входе» [31]. Подтверждение по e-mail (confirm-by-email) — индустриальный посерединный механизм между авто-линком и ручным (модель Clerk) [19].

## 6. Варианты для mcp-is

Контекст: InSales даёт `user_email` + `email_confirmed=true`, проверенный через `token3` (proof-of-InSales-session) [10][11][12]. Привязка к существующему аккаунту — разовое событие; дальше работает сессионная модель [11].

| Вариант | Безопасность | UX | Соответствие практикам | Источники |
|---|---|---|---|---|
| A. Пароль существующего аккаунта (текущий pending_link) | Высокая — proof-of-account | Средний: нужен пароль (может быть забыт) | Рекомендация USENIX («provide or reset the password»), модель Auth0 | [14][25] |
| **B. OTP-код на e-mail + фолбэк на пароль (рекомендуется)** | Высокая — proof-of-inbox; риск ≈ сбросу пароля (корень доверия — почта) | Высокий: не нужен пароль, один экран «введите код из письма» | Модель Clerk/WorkOS; параметры OTP по NIST 800-63B-4 | [16][19][27] |
| C. Авто-линк по совпадению e-mail (без подтверждения) | Низкая — Classic-Federated Merge | Максимальный | Прямо противоречит NIST 800-63C-4 и USENIX | [14][15] |
| D. Автоматический merge аккаунтов | Низкая; данные теряются у всех провайдеров | — | Не делает никто из крупных | [21][23] |

### Рекомендация: вариант B (с элементами A как фолбэка)

1. Экран при `pending_link` (вместо голого логина): «Магазин InSales ожидает привязки. Мы отправили одноразовый код на `<e-mail>`» + поле ввода + ссылка «Войти с паролем вместо этого».
2. Требования к коду (NIST 800-63B-4): ≥6 цифр, одноразовый, TTL 10 минут, лимит попыток (anti-brute-force), шлётся **только** на зарегистрированный e-mail аккаунта [16].
3. После подтверждения — штатный `link-pending` → provision → `/onboarding`; при успехе — **уведомление на почту** о привязке нового магазина (требование NIST 800-63C-4) [15].
4. Пароль остаётся фолбэком (недоставка письма), флоу не меняем — он уже соответствует практикам [14].
5. Реализация: новый шаблон письма «Код для входа/привязки» (в БД сейчас только «Подтверждение email», «Приветствие», «Сброс пароля», «Тестовое письмо» — шаблона кода нет), эндпоинт выдачи/проверки кода, страница в SPA. SMTP настроен и рабочий (`smtp.mail.ru:465`, `support@mcp-is.ru`) — канал доставки есть.

### Ответ на вопрос «много магазинов на одну почту — это коллизия?»

Нет. Мультимагазинность (один аккаунт — много магазинов и AI-ключей) — штатная модель mcp-is; совпадение e-mail при установке **второго** магазина не создаёт конфликта данных. «Коллизия» относится только к аутентификации: вход из InSales доказывает контроль над InSales-аккаунтом, но не над mcp-is-аккаунтом (у магазина могут быть менеджеры). Поэтому нужно одно доказательство владения (код/пароль), после которого новый магазин линкуется к тому же аккаунту — ровно так и работает текущий `provisionUserAndShop` (создаёт shop под существующим user). Т.е. механизм B не «разводит» аккаунты, а лишь подтверждает владельца.

## Open questions

- Количественных данных «что выбирают чаще» (survey) нет — выводы по дефолтам платформ и community-тредам [29][30].
- Passkey-инициированный link-identity flow ни у кого из проверенных провайдеров не задокументирован — вероятно, нишевый механизм.
- InSales `user_id`: глобальный или в пределах магазина — не выяснено; на текущую схему не влияет (ключ привязки — `insales_id` магазина + e-mail).
- Бенчмарки завершаемости OTP-vs-password для однократной привязки отсутствуют.
- Не сверен фактический состав install-колбэка InSales в проде с доками (доки обещают только `shop/token/insales_id`; наши тесты использовали `user_email` в install-параметрах) — см. `findings/F5.md`, follow-ups.

## Sources

[1] Shopify — ID tokens — https://shopify.dev/docs/apps/build/authentication-authorization/id-tokens (accessed 2026-10-29)
[2] Shopify — Authentication and authorization — https://shopify.dev/docs/apps/build/authentication-authorization (accessed 2026-10-29)
[3] Bitrix24 — OAuth для внешних приложений — https://apidocs.bitrix24.ru/settings/oauth/index.html (accessed 2026-10-29)
[4] Bitrix24 — Системный пользователя — https://apidocs.bitrix24.ru/settings/system-user.html (accessed 2026-10-29)
[5] Bitrix24 — Installation callback (local apps) — https://apidocs.bitrix24.ru/settings/app-installation/local-apps/installation-callback.html (accessed 2026-10-29)
[6] Wix — About app installation — https://dev.wix.com/docs/build-apps/launch-your-app/app-distribution/install-your-app/about-app-installation (updated 2026-04-19, accessed 2026-10-29)
[7] Wix — Identities — https://dev.wix.com/docs/overview/auth-permissions/identities (updated 2026-08-02, accessed 2026-10-29)
[8] Wix — About authentication — https://dev.wix.com/docs/build-apps/develop-your-app/auth/about-authentication (updated 2026-08-31, accessed 2026-10-29)
[9] Tilda — API — https://help-ru.tilda.cc/api (accessed 2026-10-29)
[10] InSales — App integration readme (raw) — https://raw.githubusercontent.com/insales/insales-docs-source/master/docs/InSales%20API/01%20readme.md (updated 2026-10-22, accessed 2026-10-29)
[11] InSales — App integration readme (published) — https://docs.liquidhub.ru/InSales%20API/01%20readme/ (updated 2026-10-22, accessed 2026-10-29)
[12] InSales — Материалы для разработчиков приложений (локальное зеркало) — /home/chubuk/projects/mcp-is/docs/api/api_and_apps.md (accessed 2026-10-29)
[13] insales_app_go — protocol.go (community SDK) — https://raw.githubusercontent.com/VladimirIvanin/insales_app_go/master/internal/insales/protocol.go (accessed 2026-10-29)
[14] Sudhodanan, Paverd — Pre-hijacked Accounts (USENIX Security 2022) — https://www.usenix.org/system/files/sec22-sudhodanan.pdf (published 2022-08-10, accessed 2026-10-29)
[15] NIST SP 800-63C-4 (Federation and Assertions) — https://pages.nist.gov/800-63-4/sp800-63c.html (published 2025-07, accessed 2026-10-29)
[16] NIST SP 800-63B-4 (Authentication and Authenticator Management) — https://pages.nist.gov/800-63-4/sp800-63b.html (published 2025-07, accessed 2026-10-29)
[17] OWASP — Email Validation and Verification Cheat Sheet — https://cheatsheetseries.owasp.org/cheatsheets/Email_Validation_and_Verification_Cheat_Sheet.html (accessed 2026-10-29)
[18] Supabase — Auth identity linking — https://supabase.com/docs/guides/auth/auth-identity-linking (accessed 2026-10-29)
[19] Clerk — Account linking — https://clerk.com/docs/guides/configure/auth-strategies/social-connections/account-linking (updated 2026-10-27, accessed 2026-10-29)
[20] Firebase — Google sign-in (account-exists-with-different-credential) — https://firebase.google.com/docs/auth/web/google-signin (accessed 2026-10-29)
[21] Firebase — Account linking — https://firebase.google.com/docs/auth/web/account-linking (updated 2026-10-21, accessed 2026-10-29)
[22] Firebase — Verified email addresses (trusted providers) — https://firebase.google.com/docs/auth/users#verified_email_addresses (updated 2026-10-20, accessed 2026-10-29)
[23] Auth0 — User account linking (overview) — https://auth0.com/docs/manage-users/user-accounts/user-account-linking (accessed 2026-10-29)
[24] Auth0 — Suggested account linking (server-side) — https://auth0.com/docs/manage-users/user-accounts/user-account-linking/suggested-account-linking-server-side-implementation (accessed 2026-10-29)
[25] Auth0 — Link user accounts — https://auth0.com/docs/manage-users/user-accounts/user-account-linking/link-user-accounts (accessed 2026-10-29)
[26] Auth0 — User-initiated account linking (client-side) — https://auth0.com/docs/manage-users/user-accounts/user-account-linking/user-initiated-account-linking-client-side-implementation (accessed 2026-10-29)
[27] WorkOS — AuthKit identity linking — https://workos.com/docs/authkit/identity-linking (accessed 2026-10-29)
[28] Clerk — User enumeration protection — https://clerk.com/docs/protect/user-enumeration-protection (accessed 2026-10-29)
[29] NextAuth.js issue #9992 (allowDangerousEmailAccountLinking) — https://github.com/nextauthjs/next-auth/issues/9992 (published 2024-02-11, accessed 2026-10-29)
[30] wafer-space issue #22 (django-allauth «Account Already Exists» UX) — https://github.com/wafer-space/platform.wafer.space/issues/22 (published 2025-09-28, accessed 2026-10-29)
[31] Better Auth — Users & accounts (disableImplicitLinking) — https://better-auth.com/docs/concepts/users-accounts (accessed 2026-10-29)
