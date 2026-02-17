package com.ielts.actual;

import com.google.auth.oauth2.GoogleCredentials;
import com.google.firebase.FirebaseApp;
import com.google.firebase.FirebaseOptions;
import org.springframework.context.annotation.Configuration;

import jakarta.annotation.PostConstruct;
import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;

@Configuration
public class FirebaseConfig {

    @PostConstruct
    public void init() throws IOException {
        // 1. Render'dagi Environment Variable'dan qidiramiz
        String firebaseConfig = System.getenv("FIREBASE_CONFIG_JSON");
        InputStream serviceAccount;

        if (firebaseConfig != null && !firebaseConfig.isEmpty()) {
            // Render platformasida (Serverda) ishlayotgan bo'lsa
            serviceAccount = new ByteArrayInputStream(firebaseConfig.getBytes(StandardCharsets.UTF_8));
        } else {
            // Lokal kompyuterda fayldan o'qiymiz (agar fayl bo'lsa)
            serviceAccount = getClass().getClassLoader()
                    .getResourceAsStream("firebase-service-account.json");
        }

        if (serviceAccount == null) {
            // Agar na fayl, na Environment variable bo'lsa, xato beradi
            System.out.println("Ogohlantirish: Firebase xizmat hisobi kaliti topilmadi!");
            return;
        }

        FirebaseOptions options = FirebaseOptions.builder()
                .setCredentials(GoogleCredentials.fromStream(serviceAccount))
                .build();

        if (FirebaseApp.getApps().isEmpty()) {
            FirebaseApp.initializeApp(options);
        }
    }
}
