package com.ielts.actual;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import java.util.Optional;

@Repository
public interface UserRepository extends JpaRepository<User, Long> {
    
    // Foydalanuvchini Firebase ID-si orqali bazadan qidirish uchun
    Optional<User> findByFirebaseUid(String firebaseUid);
    
    // Email orqali qidirish (agar kerak bo'lsa)
    Optional<User> findByEmail(String email);
}
