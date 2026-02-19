package com.ielts.actual;

import org.springframework.web.bind.annotation.*;
import java.util.List;

@RestController
@RequestMapping("/api/questions") // Bu Controller'ning manzili
@CrossOrigin(origins = {"https://platform-v11.web.app", "https://platform-v11.firebaseapp.com"})
public class QuestionController {

    // Hozircha bazasiz, shunchaki test uchun ro'yxat qaytaramiz
    @GetMapping
    public List<String> getQuestions() {
        return List.of("Describe a person you admire", "Talk about a traditional food");
    }
}
