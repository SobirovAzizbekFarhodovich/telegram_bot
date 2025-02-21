import asyncio
import asyncpg
import logging
from aiogram import Bot, Dispatcher, types
from aiogram.filters import Command
from aiogram.types import (
    ReplyKeyboardMarkup, KeyboardButton,
    InlineKeyboardMarkup, InlineKeyboardButton, WebAppInfo
)

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

BOT_TOKEN = "7722907926:AAHe9pfBs74AbiC49nPpx8IcS9NpJ-vC-ew"
DATABASE_URL = "postgresql://azizbek:123@postgres_db/bot"

bot = Bot(token=BOT_TOKEN)
dp = Dispatcher()

async def get_db_connection():
    try:
        return await asyncpg.connect(DATABASE_URL)
    except Exception as e:
        logger.error(f"❌ PostgreSQL ulanishida xatolik: {e}")
        return None

def get_web_app_keyboard(web_app_url):
    return InlineKeyboardMarkup(
        inline_keyboard=[
            [InlineKeyboardButton(
                text="🔑 Web-App'ga kirish",
                web_app=WebAppInfo(url=web_app_url)
            )]
        ]
    )

@dp.message(Command("start"))
async def start_command(message: types.Message):
    user_id = message.from_user.id
    conn = await get_db_connection()
    
    if not conn:
        await message.answer("❌ Serverda xatolik yuz berdi. Keyinroq urinib ko‘ring.")
        return

    try:
        user = await conn.fetchrow("SELECT id FROM users WHERE telegram_id=$1", user_id)
    except Exception as e:
        logger.error(f"❌ Bazaga so‘rov yuborishda xatolik: {e}")
        await message.answer("❌ Serverda xatolik yuz berdi.")
        return
    finally:
        await conn.close()

    if user:
        web_app_url = f"https://password-manager.eslab.uz/?user_id={user['id']}"
        await message.answer("🔐 Web-App'ga kirishingiz mumkin!", reply_markup=get_web_app_keyboard(web_app_url))
    else:
        await message.answer(
            "📲 Iltimos, telefon raqamingizni yuboring:",
            reply_markup=ReplyKeyboardMarkup(
                resize_keyboard=True,
                keyboard=[[KeyboardButton(text="📱 Telefon raqamni yuborish", request_contact=True)]],
                one_time_keyboard=True
            )
        )

async def main():
    try:
        logger.info("🚀 Bot ishga tushdi!")
        await dp.start_polling(bot)
    finally:
        await bot.session.close()

if __name__ == "__main__":
    asyncio.run(main())
