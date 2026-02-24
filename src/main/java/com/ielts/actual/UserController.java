package com.ielts.actual;

import com.google.firebase.auth.FirebaseAuth;
import com.google.firebase.auth.FirebaseToken;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value; // Buni qo'shing
import org.springframework.web.bind.annotation.*;
import org.springframework.http.ResponseEntity;
import org.springframework.http.HttpStatus;

import java.util.Optional;
import java.util.Map; // Buni qo'shing

@RestController
@RequestMapping("/api") // /api/users dan /api ga o'zgartirdik, endpointlar aniqroq bo'lishi uchun
@CrossOrigin(origins = {"https://platform-v11.web.app", "https://platform-v11.firebaseapp.com", "http://localhost:5173"}) // Localhostni ham qo'shdik test uchun
public class UserController {

    @Autowired
    private UserRepository userRepository;

    // Render'dagi Environment Variables'dan oladi
    @Value("${ADMIN_ID}")
    private String adminId;

    @Value("${ADMIN_KEY}")
    private String adminKey;

    // --- ADMIN VERIFY ENDPOINT (SIZNING MAXFIY KIRISHINGIZ) ---
    @PostMapping("/admin/verify")
    public ResponseEntity<?> verifyAdmin(@RequestBody Map<String, String> credentials) {
        String inputId = credentials.get("id");
        String inputKey = credentials.get("key");

        // Solishtirish mantiqi
        if (adminId != null && adminId.equals(inputId) && 
            adminKey != null && adminKey.equals(inputKey)) {
            
            return ResponseEntity.ok().body(Map.of("status", "success", "message", "Access Granted"));
        } else {
            // Noto'g'ri bo'lsa 404 (Not Found) qaytarib, o'zini bildirmaydi
            return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
        }
    }

    // --- OLD USER SYNC LOGIC ---
    @PostMapping("/users/sync")
    public ResponseEntity<?> syncUser(@RequestHeader("Authorization") String idToken, @RequestBody User userDetails) {
        try {
            String token = idToken.replace("Bearer ", "");
            FirebaseToken decodedToken = FirebaseAuth.getInstance().verifyIdToken(token);
            String uid = decodedToken.getUid();

            return userRepository.findByFirebaseUid(uid)
                .map(ResponseEntity::ok)
                .orElseGet(() -> {
                    userDetails.setFirebaseUid(uid);
                    if (userDetails.getRole() == null) userDetails.setRole("USER");
                    User savedUser = userRepository.save(userDetails);
                    return ResponseEntity.ok(savedUser);
                });
        } catch (Exception e) {
            e.printStackTrace(); 
            return ResponseEntity.status(401).body("Xatolik: " + e.getMessage());
        }
    }

    @GetMapping("/users/me")
    public String test() {
        return "Backend ishlayapti, baza ulandi!";
    }
}
