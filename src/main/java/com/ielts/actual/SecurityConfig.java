package com.ielts.actual;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.web.cors.CorsConfiguration;
import org.springframework.web.cors.CorsConfigurationSource;
import org.springframework.web.cors.UrlBasedCorsConfigurationSource;

import java.util.List;

@Configuration
public class SecurityConfig {

@Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        http
            .cors(Customizer.withDefaults()) // Sizda bor edi, zo'r!
            .csrf(csrf -> csrf.disable())    // Sizda bor edi, bu shart!
            .authorizeHttpRequests(auth -> auth
                // 1. Mana bu yo'lni hamma uchun ochamiz (Sync qilish uchun)
                .requestMatchers("/api/users/sync").permitAll() 
                
                // 2. Mana bu yo'lni faqat login qilganlar ko'radi
                .requestMatchers("/api/users/me").authenticated() 
                
                // 3. Qolgan hamma narsaga ruxsat beramiz
                .anyRequest().permitAll() 
            )
            .httpBasic(Customizer.withDefaults());
        
        return http.build();
    }

    @Bean
    public CorsConfigurationSource corsConfigurationSource() {
        CorsConfiguration configuration = new CorsConfiguration();
        configuration.setAllowedOrigins(List.of("*")); // Hamma saytlarga ruxsat
        configuration.setAllowedMethods(List.of("GET", "POST", "PUT", "DELETE", "OPTIONS"));
        configuration.setAllowedHeaders(List.of("*"));
        
        UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource();
        source.registerCorsConfiguration("/**", configuration);
        return source;
    }
}
