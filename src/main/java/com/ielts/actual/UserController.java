package com.ielts.actual;

import com.google.firebase.auth.FirebaseAuth;
import com.google.firebase.auth.FirebaseToken;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.Optional;

@RestController
@RequestMapping("/api/users")
@CrossOrigin(origins = "*") // Hamma joydan (Firebase hostingdan) keladigan so'rovlarga ruxsat berish
public class UserController {

    @Autowired
    private UserRepository userRepository;

    @PostMapping("/sync")
    public User syncUser(@RequestHeader("Authorization") String idToken, @RequestBody User userDetails) {
        try {
            // 1. Frontend yuborgan Bearer tokendan faqat kodni ajratib olish
            String token = idToken.replace("Bearer ", "");
            
            // 2. Firebase orqali tokenni tekshirish (haqiqiyligini aniqlash)
            FirebaseToken decodedToken = FirebaseAuth.getInstance().verifyIdToken(token);
            String uid = decodedToken.getUid();

            // 3. Bazadan ushbu foydalanuvchini qidirish
            Optional<User> existingUser = userRepository.findByFirebaseUid(uid);

            if (existingUser.isPresent()) {
                // Agar foydalanuvchi bazada bo'lsa, borini qaytaramiz
                return existingUser.get();
            } else {
                // Agar foydalanuvchi yangi bo'lsa, ma'lumotlarini to'ldirib saqlaymiz
                userDetails.setFirebaseUid(uid);
                if (userDetails.getRole() == null) userDetails.setRole("USER"); // Default rol
                return userRepository.save(userDetails);
            }
        } catch (Exception e) {
            throw new RuntimeException("Firebase tokenni tekshirishda xatolik: " + e.getMessage());
        }
    }

    @GetMapping("/me")
    public String test() {
        return "Backend ishlayapti, baza ulandi!";
    }
}
