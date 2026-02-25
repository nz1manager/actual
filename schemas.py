from pydantic import BaseModel, EmailStr
from typing import Optional, List, Any

# Foydalanuvchi bazaga birinchi marta kirganda (Google orqali)
class UserAuth(BaseModel):
    google_id: str
    email: EmailStr
    avatar_url: Optional[str] = None

# Profilni yangilash uchun (Ism, Familiya, Tel)
class UserUpdate(BaseModel):
    email: EmailStr
    first_name: str
    last_name: str
    phone_number: str

    class Config:
        # Bo'sh yozuvlarni oldini olish uchun
        min_anystr_length = 1

# Admin login uchun format
class AdminLogin(BaseModel):
    username: str
    password: str

# Testlarni ko'rish va natijalar uchun
class ScoreResponse(BaseModel):
    test_title: str
    correct_answers: int
    total_questions: int
    band_score: float
    overview: str

    class Config:
        from_attributes = True