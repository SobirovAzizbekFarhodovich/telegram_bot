import asyncio
import asyncpg
import logging
import re
from aiogram import Bot, Dispatcher, types, F
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
    return await asyncpg.connect(DATABASE_URL)

phone_keyboard = ReplyKeyboardMarkup(
    keyboard=[
        [KeyboardButton(text="📱 Telefon raqamni yuborish", request_contact=True)]
    ],
    resize_keyboard=True,
    one_time_keyboard=True
)

# Web-App tugmasi
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
    try:
        conn = await get_db_connection()
        user = await conn.fetchrow("SELECT id FROM users WHERE telegram_id=$1", user_id)
        await conn.close()

        if user:
            web_app_url = f"https://password-manager.eslab.uz?user_id={user['id']}"
            await message.answer("🔐 Web-App'ga kirishingiz mumkin!", reply_markup=get_web_app_keyboard(web_app_url))
        else:
            await message.answer("📲 Iltimos, telefon raqamingizni yuboring:", reply_markup=phone_keyboard)

    except Exception as e:
        logger.error(f"❌ Xatolik: {e}")
        await message.answer("❌ Serverda xatolik yuz berdi.")

@dp.message(F.contact)
async def save_user_data(message: types.Message):
    user_id = message.from_user.id
    first_name = message.from_user.first_name
    phone_number = message.contact.phone_number

    if not re.match(r"^\+998[0-9]{9}$", phone_number):
        await message.answer("❌ Telefon raqam noto‘g‘ri formatda.")
        return

    try:
        conn = await get_db_connection()
        user_id_uuid = await conn.fetchval(
            """
            INSERT INTO users (telegram_id, first_name, phone_number)
            VALUES ($1, $2, $3)
            ON CONFLICT (telegram_id) DO UPDATE SET first_name = $2, phone_number = $3
            RETURNING id
            """,
            user_id, first_name, phone_number
        )
        await conn.close()

        web_app_url = f"https://password-manager.eslab.uz?user_id={user_id_uuid}"
        await message.answer("✅ Ro'yxatdan o'tdingiz! Web-App'ga kiring:", reply_markup=get_web_app_keyboard(web_app_url))

    except Exception as e:
        logger.error(f"❌ Foydalanuvchi ma'lumotlarini saqlashda xatolik: {e}")
        await message.answer("❌ Xatolik yuz berdi.")

async def main():
    try:
        logger.info("🚀 Bot ishga tushdi!")
        await dp.start_polling(bot)
    finally:
        await bot.session.close()

if __name__ == "__main__":
    asyncio.run(main())
