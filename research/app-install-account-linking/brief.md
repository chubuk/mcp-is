# Brief: привязка аккаунта при установке приложения из маркетплейса

Дата: 2026-10-29. Режим: standard (3-5 суб-агентов, 1 follow-up раунд, 15+ источников).

## Вопрос

Как платформы с маркетплейсами приложений (Shopify, Bitrix24, Wix, Tilda, amoCRM и др.)
и identity-провайдеры (Auth0, Clerk, Supabase, Firebase) решают привязку аккаунта при
установке/первом входе из приложения маркетплейса:

1. что происходит, когда e-mail из маркетплейса совпадает с существующим аккаунтом сервиса;
2. какие passwordless-механизмы привязки существуют (magic link, OTP на e-mail, passkey, SSO);
3. какие anti-takeover-защиты обязательны (атаки на e-mail-based linking, pre-hijacking);
4. какие UX-паттерны реально используют (pending link / merge / confirm-by-email).

Цель — выбрать схему для InSales-приложения mcp-is: оставить парольный pending_link,
перейти на OTP-код на e-mail, или авто-линк с ограничениями.

## Контекст

- InSales-приложение (mcp-is) при входе из админки получает `user_email` +
  `email_confirmed=true` + подписи (token/token2/token3).
- Сейчас: e-mail совпал с существующим аккаунтом → `pending_link` → вход паролем
  один раз → дальше автологин без пароля. Разработчик спрашивает про альтернативы.
- Аудитория отчёта — разработчик проекта; **REPORT.md писать на русском**.
- Актуальность: практики 2023–2026, маркетплейс-платформы и identity-провайдеры.

## Границы

- ВНЕ темы: общий курс по OAuth 2.1, наш код/схема БД, биллинг, маркетинг.
- В теме: задокументированные flows привязки, security-модели, UX-паттерны, кейсы
  уязвимостей e-mail-linking.

## Допущения

- Решение на стейкхолде: механизм первой привязки аккаунта в app-mode.
- Регион-провайдеры: глобальные + RU-платформы (Bitrix24, Tilda, amoCRM) — в приоритете
  сопоставимые с InSales (e-commerce/сайты).
- Язык поиска: en + ru.

## Angles

- **F1** — Как маркетплейсы приложений передают идентичность при установке (Shopify
  session tokens/JWT, Bitrix24 OAuth, Tilda, Wix) и что делают при совпадении e-mail
  с существующим аккаунтом (официальные доки).
- **F2** — Passwordless-привязка аккаунта: magic link / OTP на e-mail / passkeys —
  задокументированные link-identity flows у Auth0, Clerk, Supabase, Firebase, WorkOS
  при совпадении e-mail (link vs merge vs block).
- **F3** — Anti-takeover: известные атаки на e-mail-based account linking
  (pre-hijacking, federated identity takeover), рекомендации NIST/OWASP/исследований —
  какие подтверждения обязательны при первой привязке.
- **F4** — Практика UX: как решают «user already exists» при входе из приложения
  маркетплейса (форумы, guides, community) — pending link, merge accounts,
  confirm-by-email; что выбирают чаще и почему.
- **F5** — InSales-специфика: что официальные доки InSales говорят про user-identity
  в app-колбэках (`user_email`, `email_confirmed`, token3) и рекомендации по
  онбордингу/входу в приложение.
