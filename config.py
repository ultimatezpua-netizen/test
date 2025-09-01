import os
from dataclasses import dataclass
from dotenv import load_dotenv

load_dotenv()

@dataclass
class Settings:
    bot_token: str
    baserow_token: str | None
    baserow_table_id: str | None
    mono_webhook_secret: str | None
    rules_version: str
    ticket_price: int
    host: str
    port: int

settings = Settings(
    bot_token=os.getenv("BOT_TOKEN", ""),
    baserow_token=os.getenv("BASEROW_TOKEN"),
    baserow_table_id=os.getenv("BASEROW_TABLE_ID"),
    mono_webhook_secret=os.getenv("MONO_WEBHOOK_SECRET"),
    rules_version=os.getenv("RULES_VERSION", "2025-09-05"),
    ticket_price=int(os.getenv("TICKET_PRICE", "10000")),
    host=os.getenv("HOST", "0.0.0.0"),
    port=int(os.getenv("PORT", "8080")),
)

def ensure():
    if not settings.bot_token:
        raise RuntimeError("BOT_TOKEN not set")
ensure()