package com.ielts.actual; // Papka nomiga qarab tekshirib oling

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class HomeController {
    @GetMapping("/")
    public String home() {
        return "<h1>Server ishlayapti, Otabek!</h1><p>IELTS loyihasiga xush kelibsiz.</p>";
    }
}