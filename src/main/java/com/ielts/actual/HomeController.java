package com.ielts.actual;

import org.springframework.web.bind.annotation.CrossOrigin; // Mana buni qo'shish shart!
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@CrossOrigin(origins = {"https://platform-v11.web.app", "https://platform-v11.firebaseapp.com"})
@RestController
public class HomeController {

    @GetMapping("/")
    public String home() {
        return "<h1>Server ishlayapti, Otabek!</h1><p>IELTS loyihasiga xush kelibsiz.</p>";
    }
}
