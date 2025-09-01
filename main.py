#!/usr/bin/env python3
"""
Simple Echo Telegram Bot

This bot echoes back any message sent to it and responds to /start command with a greeting.
The bot token is loaded from environment variables using python-dotenv.
"""

import logging
import os
from telegram import Update
from telegram.ext import Application, CommandHandler, MessageHandler, filters, ContextTypes
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Configure logging
logging.basicConfig(
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    level=logging.INFO
)
logger = logging.getLogger(__name__)

# Get bot token from environment variables
BOT_TOKEN = os.getenv('BOT_TOKEN')

if not BOT_TOKEN:
    raise ValueError("BOT_TOKEN environment variable is not set. Please check your .env file.")


async def start_command(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Handle the /start command with a greeting message."""
    user_name = update.effective_user.first_name
    greeting_message = f"Привет, {user_name}! 👋\n\nЯ простой эхо-бот. Отправь мне любое сообщение, и я повторю его обратно!"
    
    await update.message.reply_text(greeting_message)
    logger.info(f"User {update.effective_user.id} ({user_name}) started the bot")


async def echo_message(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Echo any text message back to the user."""
    user_message = update.message.text
    user_name = update.effective_user.first_name
    
    # Echo the message back
    await update.message.reply_text(user_message)
    
    logger.info(f"User {update.effective_user.id} ({user_name}) sent: {user_message}")


async def error_handler(update: object, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Log errors caused by Updates."""
    logger.error(msg="Exception while handling an update:", exc_info=context.error)


def main() -> None:
    """Main function to start the bot."""
    # Create the Application
    application = Application.builder().token(BOT_TOKEN).build()

    # Register handlers
    application.add_handler(CommandHandler("start", start_command))
    application.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, echo_message))
    
    # Register error handler
    application.add_error_handler(error_handler)

    # Start the bot
    logger.info("Starting Echo Bot...")
    application.run_polling(allowed_updates=Update.ALL_TYPES)


if __name__ == '__main__':
    main()