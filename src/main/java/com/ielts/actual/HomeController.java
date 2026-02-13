package com.ielts.actual;

import org.springframework.web.bind.annotation.CrossOrigin; // Mana buni qo'shish shart!
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@CrossOrigin(origins = "*") // 1. Mana shu qatorni qo'shing
@RestController
public class HomeController {

    @GetMapping("/")
    public String home() {
        return "<h1>Server ishlayapti, Otabek!</h1><p>IELTS loyihasiga xush kelibsiz.</p>";
    }
}