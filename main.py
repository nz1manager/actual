from fastapi import FastAPI, Depends, HTTPException, status
from sqlalchemy.orm import Session
from typing import List
import database, models, schemas
from fastapi.middleware.cors import CORSMiddleware

# Bazani yaratish (Jadvallar bo'lmasa, avtomatik yaratadi)
models.Base.metadata.create_all(bind=database.engine)

app = FastAPI(title="IELTS ACTUAL 2026 API")

# Frontend (Firebase/Web.app) ulanishi uchun
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"], 
    allow_methods=["*"],
    allow_headers=["*"],
)

# --- FOYDALANUVCHI QISMI (STUDENT) ---

@app.post("/auth/google", response_model=schemas.UserAuth)
def google_login(user_data: schemas.UserAuth, db: Session = Depends(database.get_db)):
    # Google orqali kirganda foydalanuvchini bazadan qidirish yoki yaratish
    db_user = db.query(models.User).filter(models.User.google_id == user_data.google_id).first()
    if not db_user:
        db_user = models.User(
            google_id=user_data.google_id, 
            email=user_data.email, 
            avatar_url=user_data.avatar_url
        )
        db.add(db_user)
        db.commit()
        db.refresh(db_user)
    return db_user

@app.post("/user/update")
def update_profile(user_data: schemas.UserUpdate, db: Session = Depends(database.get_db)):
    # Ism, Familiya va Tel raqamni saqlash (Save tugmasi)
    db_user = db.query(models.User).filter(models.User.email == user_data.email).first()
    if not db_user:
        raise HTTPException(status_code=404, detail="Foydalanuvchi topilmadi")
    
    db_user.first_name = user_data.first_name
    db_user.last_name = user_data.last_name
    db_user.phone_number = user_data.phone_number
    
    db.commit()
    return {"status": "success", "message": "Updated successfully!"}

# --- ADMIN PANEL QISMI ---

@app.post("/admin/login")
def admin_auth(creds: schemas.AdminLogin):
    # Siz aytgan maxfiy login va parol
    if creds.username == "drkz1manager" and creds.password == "4900728Tt":
        return {"access": "granted", "token": "admin_session_2026"}
    raise HTTPException(status_code=401, detail="Kirish taqiqlandi!")

@app.get("/admin/students")
def get_all_students(db: Session = Depends(database.get_db)):
    # Admin uchun hamma o'quvchilar ro'yxati
    return db.query(models.User).filter(models.User.role == "student").all()

@app.delete("/admin/students/{student_id}")
def delete_student(student_id: int, db: Session = Depends(database.get_db)):
    # O'quvchini bazadan butunlay o'chirish (Delete tugmasi)
    student = db.query(models.User).filter(models.User.id == student_id).first()
    if not student:
        raise HTTPException(status_code=404, detail="O'quvchi topilmadi")
    db.delete(student)
    db.commit()
    return {"message": "O'quvchi muvaffaqiyatli o'chirildi"}

# Ballar va Overview hisoblash mantiqi (Namuna)
def calculate_overview(score: float):
    if score >= 8.0: return "A'lochi (Expert)"
    if score >= 6.5: return "Yaxshi (Competent)"
    return "Qoniqarsiz (Improvement needed)"

@app.get("/user/scores/{user_id}")
def get_user_scores(user_id: int, db: Session = Depends(database.get_db)):
    # Dashboard uchun o'quvchi ballari
    return db.query(models.Score).filter(models.Score.user_id == user_id).all()