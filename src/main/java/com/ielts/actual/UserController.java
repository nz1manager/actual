package com.ielts.actual;


import com.google.firebase.auth.FirebaseAuth;
import com.google.firebase.auth.FirebaseToken;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import org.springframework.http.ResponseEntity;

import java.util.Optional;

@RestController
@RequestMapping("/api/users")
@CrossOrigin(origins = {"https://platform-v11.web.app", "https://platform-v11.firebaseapp.com"})
public class UserController {

    @Autowired
    private UserRepository userRepository;

   @PostMapping("/sync")
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
        // Xatoni aniq ko'rish uchun terminalga chiqaramiz
        e.printStackTrace(); 
        return ResponseEntity.status(401).body("Xatolik: " + e.getMessage());
    }
}
    @GetMapping("/me")
    public String test() {
        return "Backend ishlayapti, baza ulandi!";
    }
}
