import os
from sqlalchemy import create_engine, MetaData
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from dotenv import load_dotenv

load_dotenv()

# Neon.tech Connection String - Environment Variable-dan olinadi
DATABASE_URL = os.getenv("DATABASE_URL")

# Bazaga ulanish mexanizmi
engine = create_engine(DATABASE_URL)

# Ma'lumotlar bilan ishlash uchun sessiya (aloqa seansi)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

# Jadvallarni yaratish uchun asosiy klass
Base = declarative_base()

# Bazaga ulanishni tekshirish va yopish funksiyasi
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()