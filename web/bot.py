import asyncio
import asyncpg
from aiogram import Bot, Dispatcher, F, types
from aiogram.types import ReplyKeyboardMarkup, KeyboardButton, InlineKeyboardMarkup, InlineKeyboardButton, WebAppInfo
from aiogram.filters import Command
import re
import logging

# Loggerni sozlash
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Bot tokeni
BOT_TOKEN = "7722907926:AAHe9pfBs74AbiC49nPpx8IcS9NpJ-vC-ew"
DATABASE_URL = "postgresql://azizbek:123@postgres_db/bot"

# Bot va Dispatcher obyektlarini yaratish
bot = Bot(token=BOT_TOKEN)
dp = Dispatcher()

# Telefon raqamini so‘rash uchun klaviatura
phone_keyboard = ReplyKeyboardMarkup(
    keyboard=[
        [KeyboardButton(text="Telefon raqamni yuborish", request_contact=True)]
    ],
    resize_keyboard=True,
    one_time_keyboard=True
)

@dp.message(Command("start"))
async def start_command(message: types.Message):
    user_id = message.from_user.id
    try:
        conn = await asyncpg.connect(DATABASE_URL)
        user_exists = await conn.fetchrow("SELECT id FROM users WHERE telegram_id=$1", user_id)
        await conn.close()

        if user_exists:
            web_app_url = f"https://password-manager.eslab.uz?user_id={user_exists['id']}"
            await message.answer("Web-Appga kirishingiz mumkin.", reply_markup=get_web_app_keyboard(web_app_url))
        else:
            await message.answer(
                "Salom! Iltimos, telefon raqamingizni yuboring:",
                reply_markup=phone_keyboard
            )

    except Exception as e:
        logger.error(f"Error: {e}")
        await message.answer("Xatolik yuz berdi.")

@dp.message(F.contact)
async def save_user_data(message: types.Message):
    user_id = message.from_user.id
    first_name = message.from_user.first_name
    phone_number = message.contact.phone_number

    if not re.match(r"^\+998[0-9]{9}$", phone_number):
        await message.answer("Telefon raqam noto‘g‘ri formatda.")
        return

    try:
        conn = await asyncpg.connect(DATABASE_URL)
        user_id_uuid = await conn.fetchval(
            """
            INSERT INTO users (telegram_id, first_name, phone_number)
            VALUES ($1, $2, $3)
            ON CONFLICT (telegram_id) DO NOTHING
            RETURNING id
            """,
            user_id, first_name, phone_number
        )
        await conn.close()

        web_app_url = f"https://password-manager.eslab.uz?user_id={user_id_uuid}"
        await message.answer("Ma'lumotlaringiz saqlandi! Web-Appga kirishingiz mumkin.", reply_markup=get_web_app_keyboard(web_app_url))

    except Exception as e:
        logger.error(f"Error saving user data: {e}")
        await message.answer(f"Xatolik yuz berdi: {e}")

def get_web_app_keyboard(web_app_url):
    return InlineKeyboardMarkup(
        inline_keyboard=[
            [
                InlineKeyboardButton(
                    text="Web-Appga kirish",
                    web_app=WebAppInfo(url=web_app_url)
                )
            ]
        ]
    )

async def main():
    try:
        logger.info("Starting bot")
        await dp.start_polling(bot)
    finally:
        await bot.session.close()

if __name__ == "__main__":
    asyncio.run(main())
