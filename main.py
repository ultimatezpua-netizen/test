import asyncio
import datetime
import logging
import random
import string
from typing import Dict, Optional

from aiohttp import web, ClientSession
from aiogram import Bot, Dispatcher, types
from aiogram.utils import executor

from config import settings

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s | %(levelname)s | %(name)s | %(message)s",
)

bot = Bot(token=settings.bot_token)
dp = Dispatcher(bot)

pending_tokens: Dict[str, tuple[int, str]] = {}

def generate_payment_token(length: int = 8) -> str:
    return "".join(random.choices(string.ascii_uppercase + string.digits, k=length))

async def add_ticket(user_id: int, username: Optional[str], payment_id: str):
    ticket_num = f"T-{{random.randint(1000, 9999)}}"
    baserow_table = settings.baserow_table_id
    if settings.baserow_token and baserow_table:
        async with ClientSession() as session:
            url = f"https://api.baserow.io/api/database/rows/table/{{baserow_table}}/"
            payload = {
                "user_id": user_id,
                "username": username or "",
                "ticket_number": ticket_num,
                "payment_id": payment_id,
                "paid_at": datetime.datetime.utcnow().isoformat(),
                "rules_version": settings.rules_version,
                "status": "Paid"
            }
            headers = {
                "Authorization": f"Token {{settings.baserow_token}}",
                "Content-Type": "application/json",
            }
            resp = await session.post(url, headers=headers, json=payload)
            if resp.status >= 300:
                text = await resp.text()
                logging.error("Baserow error %s: %s", resp.status, text)
    else:
        logging.warning("Baserow creds not set; ticket only in memory.")
    return ticket_num

@dp.message_handler(commands=["start"])
async def start_cmd(msg: types.Message):
    kb = types.ReplyKeyboardMarkup(resize_keyboard=True)
    kb.add("🎟 Купить билет")
    await msg.answer(
        "Привет! Билет стоит 100 грн.\nНажми кнопку ниже, чтобы оплатить через Monobank.",
        reply_markup=kb
    )

@dp.message_handler(lambda m: m.text == "🎟 Купить билет")
async def buy_ticket(msg: types.Message):
    token = generate_payment_token()
    pending_tokens[token] = (msg.from_user.id, msg.from_user.username or "")
    pay_link = f"https://send.monobank.ua/jar/ВАШ_JAR_ID?amount=100&text={{token}}"
    await msg.answer(
        "1. Перейди по ссылке и оплати 100 грн.\n"
        "2. Комментарий к платежу (code) уже подставлен в ссылку: \n"
        f"`{{token}}`\n"
        "3. После успешной обработки придёт билет.",
        parse_mode="Markdown"
    )
    await msg.answer(f"Ссылка на оплату: {{pay_link}}")

async def mono_webhook(request: web.Request):
    if settings.mono_webhook_secret:
        provided = request.headers.get("X-Mono-Secret")
        if provided != settings.mono_webhook_secret:
            return web.Response(status=403, text="forbidden")

    try:
        data = await request.json()
    except Exception:
        return web.Response(status=400, text="invalid json")

    logging.info("Monobank webhook payload: %s", data)

    amount = data.get("amount")
    payment_id = data.get("id")
    comment = (data.get("comment") or "").strip().upper()

    if amount != settings.ticket_price:
        logging.warning("Amount mismatch: got %s expected %s", amount, settings.ticket_price)
        return web.Response(text="ignored")

    if not comment or comment not in pending_tokens:
        logging.warning("Unknown or missing token: %s", comment)
        return web.Response(text="pending_token_not_found")

    user_id, username = pending_tokens.pop(comment)
    ticket = await add_ticket(user_id, username, payment_id)
    try:
        await bot.send_message(user_id, f"✅ Оплата получена! Ваш билет: {{ticket}}")
    except Exception as e:
        logging.error("Failed to send ticket to user %s: %s", user_id, e)

    return web.Response(text="ok")

def create_app() -> web.Application:
    app = web.Application()
    app.router.add_post("/mono-webhook", mono_webhook)
    app.router.add_get("/health", lambda _: web.Response(text="ok"))
    return app

async def run_web_app(app: web.Application):
    runner = web.AppRunner(app)
    await runner.setup()
    site = web.TCPSite(runner, settings.host, settings.port)
    await site.start()
    logging.info("Webhook server started on %s:%s", settings.host, settings.port)

async def main():
    app = create_app()
    await run_web_app(app)
    loop = asyncio.get_event_loop()
    loop.create_task(executor._startup_polling(dp, timeout=20, relax=0.1, fast=True))
    while True:
        await asyncio.sleep(3600)

if __name__ == "__main__":
    try:
        asyncio.run(main())
    except (KeyboardInterrupt, SystemExit):
        logging.info("Stopped.")