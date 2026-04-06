use axum::{
    extract::ws::{Message, WebSocket, WebSocketUpgrade},
    response::IntoResponse,
    routing::{get, post},
    Json, Router,
};
use futures_util::{sink::SinkExt, stream::StreamExt};
use std::sync::Arc;
use tokio::sync::broadcast;
use tower_http::cors::CorsLayer;

struct AppState {
    tx: broadcast::Sender<serde_json::Value>,
}

#[tokio::main]
async fn main() {
    // Kanal yaratamiz: 100 ta xabarni navbatda tuta oladi
    let (tx, _rx) = broadcast::channel(100);
    let app_state = Arc::new(AppState { tx });

    let app = Router::new()
        // 1. Uydagi Rust botdan ma'lumot qabul qilish
        .route("/api/v1/update", post(update_handler))
        // 2. Vanilla JS frontend uchun WebSocket ulanishi
        .route("/ws", get(ws_handler))
        .layer(CorsLayer::permissive())
        .with_state(app_state);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:10000").await.unwrap();
    println!("Render Backend ishga tushdi: port 10000");
    axum::serve(listener, app).await.unwrap();
}

async fn update_handler(
    axum::extract::State(state): axum::extract::State<Arc<AppState>>,
    Json(payload): Json<serde_json::Value>,
) -> impl IntoResponse {
    // Kelgan ma'lumotni barcha ulangan frontend larga tarqatamiz (Broadcast)
    let _ = state.tx.send(payload);
    "OK"
}

async fn ws_handler(
    ws: WebSocketUpgrade,
    axum::extract::State(state): axum::extract::State<Arc<AppState>>,
) -> impl IntoResponse {
    ws.on_upgrade(|socket| handle_socket(socket, state))
}

async fn handle_socket(socket: WebSocket, state: Arc<AppState>) {
    let (mut sender, _) = socket.split();
    let mut rx = state.tx.subscribe();

    // Yangi xabar kelishi bilan frontendga uchadi
    while let Ok(msg) = rx.recv().await {
        let text = msg.to_string();
        if sender.send(Message::Text(text)).await.is_err() {
            break;
        }
    }
}
