use headless_chrome::Browser;
use std::sync::Arc;
use reqwest::header::{HeaderMap, HeaderValue, USER_AGENT};
use std::time::Duration;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let remote_url = "https://sizning-backend-manzilingiz.onrender.com/api/v1/update";
    
    // User-Agent va Token ma'lumotlarini tayyorlash
    let user_agent_str = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36 Edg/146.0.0.0";
    let auth_token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzU3MTYwOTMsImlhdCI6MTc3NTQ1Njg5Mywic3ViIjoiNDI5OTgwMTAiLCJjaGFubmVscyI6WyJsdWNreS1qZXQtOTQiXX0.P1bu7QFrG3ec5bC-yLig6rPUJz3ScYfsU6BLNk7e6Y8";
    let client_id = "bb245762-d1c0-4050-bd23-3306f4fa5bda";

    let mut headers = HeaderMap::new();
    headers.insert(USER_AGENT, HeaderValue::from_str(user_agent_str)?);
    headers.insert("Authorization", HeaderValue::from_str(&format!("Bearer {}", auth_token))?);
    headers.insert("X-Client-ID", HeaderValue::from_str(client_id)?);

    let http_client = reqwest::Client::builder()
        .default_headers(headers)
        .tcp_nodelay(true)
        .build()?;

    println!("Edge ulanmoqda...");
    let browser = Browser::connect("http://127.0.0.1:9222".to_string())?;
    let tabs = browser.get_tabs().unwrap();
    let tab = tabs.lock().unwrap().get(0).cloned().expect("Tab topilmadi");

    tab.enable_network_events()?;
    let client_arc = Arc::new(http_client);

    // --- PING-PONG MEXANIZMI ---
    let ping_client = Arc::clone(&client_arc);
    let ping_url = remote_url.to_string();
    tokio::spawn(async move {
        loop {
            let _ = ping_client.post(&ping_url)
                .json(&serde_json::json!({"type": "ping", "client": "rust_bot"}))
                .send()
                .await;
            tokio::time::sleep(Duration::from_secs(25)).await; // Har 25 sekundda ping
        }
    });

    // --- WS MONITORING ---
    tab.add_event_listener(Arc::new(move |event: &headless_chrome::browser::tab::TabEvent| {
        if let headless_chrome::browser::tab::TabEvent::NetworkWebSocketFrameReceived(frame) = event {
            let payload = frame.params.response.payload_data.clone();

            if payload.contains("changeCoefficient") || payload.contains("stopCoefficient") {
                let client = Arc::clone(&client_arc);
                let url = remote_url.to_string();
                
                tokio::spawn(async move {
                    let mut data: serde_json::Value = serde_json::from_str(&payload).unwrap_or_default();
                    
                    // Stop bo'lganda qo'shimcha rang flagini qo'shamiz
                    if payload.contains("stopCoefficient") {
                        data["ui_color"] = serde_json::json!("red");
                    } else {
                        data["ui_color"] = serde_json::json!("white");
                    }

                    let _ = client.post(&url).json(&data).send().await;
                });
            }
        }
    }))?;

    println!("Bridge ishlamoqda. Ma'lumotlar Render'ga uchmoqda...");
    loop { tokio::time::sleep(Duration::from_secs(1)).await; }
}