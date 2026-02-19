package com.ielts.actual;

import jakarta.persistence.*;
import lombok.*;
import com.fasterxml.jackson.annotation.JsonProperty;

@Entity
@Table(name = "users")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class User {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    // Firebase UID - bazada firebase_uid bo'lib tushadi
    @Column(name = "firebase_uid", unique = true, nullable = false)
    private String firebaseUid; 

    @Column(name = "first_name")
    private String firstName;

    @Column(name = "last_name")
    private String lastName;
    
    @Column(unique = true)
    private String email;
    
    // Frontend-dan phoneNumber bo'lib kelsa ham Java-da phone-ga o'qiydi
    @JsonProperty("phoneNumber")
    private String phone;
    
    private String role; // ADMIN yoki USER
    
    private Integer score; // IELTS balli
}
