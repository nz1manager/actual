use futures_util::{StreamExt, SinkExt};
use tokio_tungstenite::{connect_async, tungstenite::protocol::Message};
use tokio_tungstenite::tungstenite::client::IntoClientRequest;
use serde_json::json;
use reqwest::Client;
use std::time::Duration;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // 1. SOZLAMALAR
    let ws_url = "wss://crash-gateway-grm-cr.gamedev-tech.cc/websocket/lifecycle";
    let render_url = "https://rulsz.onrender.com/api/v1/update";
    
    let token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzU3MTc0MzgsImlhdCI6MTc3NTQ1ODIzOCwic3ViIjoiNDI5OTgwMTAiLCJjaGFubmVscyI6WyJsdWNreS1qZXQtOTQiXX0.C4t5galGLGYj1fzbjEoEkK7q9R8Xl_tn79GfD6zkPBk";
    let client_id = "6750d484-153d-45e9-af66-9c3dffd814d6";
    let user_agent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36 Edg/146.0.0.0";

    let http_client = Client::builder().tcp_nodelay(true).build()?;
    println!("Lucky Jet serveriga ulanishga harakat qilinmoqda...");

    // 2. HTTP Request tayyorlash (Headerlar bilan)
    let mut request = ws_url.into_client_request()?;
    let headers = request.headers_mut();
    headers.insert("User-Agent", user_agent.parse()?);
    headers.insert("Origin", "https://1w-24094.com".parse()?); // O'yin domeni

    // 3. WebSocket ulanish
    let (ws_stream, _) = connect_async(request).await.expect("Ulanib bo'lmadi!");
    let (mut write, mut read) = ws_stream.split();

    // 4. Centrifuge Connect (Birinchi handshake)
    let connect_payload = json!({
        "params": {
            "token": token,
            "name": "js"
        },
        "id": 1
    }).to_string();
    
    // Muhim: Centrifuge'da har bir xabar oxirida \n yoki maxsus separator bo'lishi mumkin
    write.send(Message::Text(connect_payload)).await?;

    // 5. Ping-Pong (Aloqani uzib qo'ymaslik uchun)
    let (tx, mut rx) = tokio::sync::mpsc::channel::<String>(32);
    let mut write_half = write;
    
    tokio::spawn(async move {
        loop {
            tokio::time::sleep(Duration::from_secs(20)).await;
            let ping = json!({}).to_string(); // Centrifuge bo'sh pingi
            if write_half.send(Message::Text(ping)).await.is_err() { break; }
        }
    });

    println!("Muvaffaqiyatli ulandik! Ma'lumotlar o'g'irlanib rulsz.onrender.com ga yuborilmoqda...");

    // 6. Ma'lumotlarni o'qish va Render'ga otish
    while let Some(Ok(msg)) = read.next().await {
        if let Message::Text(text) = msg {
            // Agar xabar o'yin koeffitsienti haqida bo'lsa
            if text.contains("changeCoefficient") || text.contains("stopCoefficient") {
                let client = http_client.clone();
                let url = render_url.to_string();
                let body = text.clone();

                tokio::spawn(async move {
                    let _ = client.post(&url)
                        .header("Content-Type", "application/json")
                        .body(body)
                        .send()
                        .await;
                });

                if text.contains("stopCoefficient") {
                    println!("--- STOP ---");
                }
            }
        }
    }

    Ok(())
}
