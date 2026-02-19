package com.ielts.actual;

import jakarta.persistence.*;
import lombok.*;

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

    // BU JUDA MUHIM: Firebase va Neon DB ni bog'lovchi zanjir
    @Column(unique = true, nullable = false)
    private String firebaseUid; 

    private String firstName;
    private String lastName;
    
    @Column(unique = true)
    private String email;
  @JsonProperty("phoneNumber") // Frontend'dan "phoneNumber" bo'lib kelsa, Java buni "phone"ga o'qiydi
    private String phone;
    private String role; // ADMIN yoki USER
    private Integer score; // IELTS balli
}
