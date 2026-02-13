package com.ielts.actual; // SIZNING LOYIHA NOMINGIZ SHU EDI

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

    private String firstName;
    private String lastName;
    
    @Column(unique = true)
    private String email;
    
    private String phone;
    private String role; // ADMIN yoki USER
    private Integer score; // IELTS balli
}