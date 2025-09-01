# Promo Bot

Бот для выдачи “билетов” после оплаты через Monobank и логирования в Baserow.

## MVP
- /start
- Кнопка “Купить билет”
- Генерация токена (связь пользователь ↔ платёж)
- Webhook /mono-webhook получает уведомление
- Валидация суммы
- Запись в Baserow (опционально)
- Отправка билета пользователю

## Переменные окружения (.env)
BOT_TOKEN=123456:ABC...
BASEROW_TOKEN=...
BASEROW_TABLE_ID=...
MONO_WEBHOOK_SECRET=...
RULES_VERSION=2025-09-05
TICKET_PRICE=10000
HOST=0.0.0.0
PORT=8080

## Локально
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
python main.py

## Docker
docker build -t promo-bot .
docker run --rm -it -p 8080:8080 --env-file .env promo-bot

## Fly.io
fly launch --no-deploy
fly deploy
fly secrets set BOT_TOKEN=... MONO_WEBHOOK_SECRET=... TICKET_PRICE=10000
fly secrets set BASEROW_TOKEN=... BASEROW_TABLE_ID=...

Проверка:
curl https://promo-bot.fly.dev/health

## Baserow поля
user_id, username, ticket_number, payment_id, paid_at, rules_version, status

## Улучшения
- Хранить pending_tokens в Redis
- HMAC подпись webhook
- Админ-статистика
- Перевод на Telegram Webhook режим
- Повторная отправка билета по команде
