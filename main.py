import os
import json
import firebase_admin
from firebase_admin import credentials, auth
from fastapi import FastAPI, Depends, HTTPException, status, Body
from sqlalchemy.orm import Session
from typing import List
import database, models, schemas
from fastapi.middleware.cors import CORSMiddleware

# Bazani yaratish
models.Base.metadata.create_all(bind=database.engine)

app = FastAPI(title="IELTS ACTUAL 2026 API")

# --- FIREBASE XAVFSIZ INITIALIZATION ---
firebase_creds_json = os.getenv("FIREBASE_SERVICE_ACCOUNT")
if firebase_creds_json:
    try:
        # JSON ichidagi ortiqcha qator tashlashlarni tozalaymiz
        creds_dict = json.loads(firebase_creds_json.replace('\n', '\\n'))
        if not firebase_admin._apps:
            cred = credentials.Certificate(creds_dict)
            firebase_admin.initialize_app(cred)
        print("Firebase Successfully Initialized!")
    except Exception as e:
        print(f"Firebase Init Error: {e}")

# CORS sozlamalari
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"], 
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
def home():
    return {"message": "IELTS ACTUAL 2026 API - Live and Safe!"}

# --- FOYDALANUVCHI QISMI (STUDENT) ---

@app.post("/auth/google")
def google_login(payload: dict = Body(...), db: Session = Depends(database.get_db)):
    id_token = payload.get("id_token")
    try:
        decoded_token = auth.verify_id_token(id_token)
        google_id = decoded_token['uid']
        email = decoded_token['email']
        avatar_url = decoded_token.get('picture')

        db_user = db.query(models.User).filter(models.User.google_id == google_id).first()
        if not db_user:
            db_user = models.User(
                google_id=google_id, 
                email=email, 
                avatar_url=avatar_url,
                role="student"
            )
            db.add(db_user)
            db.commit()
            db.refresh(db_user)
        
        return {
            "status": "success", 
            "user": {
                "id": db_user.id, 
                "email": db_user.email,
                "first_name": db_user.first_name,
                "last_name": db_user.last_name
            }
        }
    except Exception as e:
        print(f"Auth Error: {e}")
        raise HTTPException(status_code=401, detail="Token yaroqsiz!")

@app.post("/user/update")
def update_profile(user_data: schemas.UserUpdate, db: Session = Depends(database.get_db)):
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
    if creds.username == "drkz1manager" and creds.password == "4900728Tt":
        return {"access": "granted", "token": "admin_session_2026"}
    raise HTTPException(status_code=401, detail="Kirish taqiqlandi!")

@app.get("/admin/students")
def get_all_students(db: Session = Depends(database.get_db)):
    return db.query(models.User).filter(models.User.role == "student").all()

@app.delete("/admin/students/{student_id}")
def delete_student(student_id: int, db: Session = Depends(database.get_db)):
    student = db.query(models.User).filter(models.User.id == student_id).first()
    if not student:
        raise HTTPException(status_code=404, detail="O'quvchi topilmadi")
    db.delete(student)
    db.commit()
    return {"message": "O'quvchi muvaffaqiyatli o'chirildi"}

@app.get("/user/scores/{user_id}")
def get_user_scores(user_id: int, db: Session = Depends(database.get_db)):
    return db.query(models.Score).filter(models.Score.user_id == user_id).all()
