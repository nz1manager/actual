from sqlalchemy import Column, Integer, String, Boolean, ForeignKey, JSON, Float, DateTime
from sqlalchemy.sql import func
from database import Base

class User(Base):
    __tablename__ = "users"

    id = Column(Integer, PRIMARY KEY=True, index=True)
    google_id = Column(String, unique=True, index=True)
    email = Column(String, unique=True, index=True)
    first_name = Column(String, nullable=True)
    last_name = Column(String, nullable=True)
    phone_number = Column(String, nullable=True)
    avatar_url = Column(String, nullable=True)
    role = Column(String, default="student") # student yoki admin
    created_at = Column(DateTime(timezone=True), server_default=func.now())

class Test(Base):
    __tablename__ = "tests"

    id = Column(Integer, PRIMARY KEY=True, index=True)
    title = Column(String)
    skill_type = Column(String) # reading, listening, writing, speaking
    content = Column(JSON) # Savollar va javoblar shu yerda
    is_active = Column(Boolean, default=False) # Admin 'ptichka' qo'yishi uchun
    created_at = Column(DateTime(timezone=True), server_default=func.now())

class Score(Base):
    __tablename__ = "scores"

    id = Column(Integer, PRIMARY KEY=True, index=True)
    user_id = Column(Integer, ForeignKey("users.id"))
    test_id = Column(Integer, ForeignKey("tests.id"))
    correct_answers = Column(Integer)
    total_questions = Column(Integer)
    band_score = Column(Float) # Masalan: 7.5
    overview = Column(String) # Alochi, Qoniqarli va h.k.
    submitted_at = Column(DateTime(timezone=True), server_default=func.now())