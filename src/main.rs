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
    let (tx, _rx) = broadcast::channel(100);
    let app_state = Arc::new(AppState { tx });

    let app = Router::new()
        .route("/api/v1/update", post(update_handler))
        .route("/ws", get(ws_handler))
        .layer(CorsLayer::permissive())
        .with_state(app_state);

    let port = std::env::var("PORT").unwrap_or_else(|_| "10000".to_string());
    let addr = format!("0.0.0.0:{}", port);
    
    let listener = tokio::net::TcpListener::bind(&addr).await.unwrap();
    println!("Render Backend ishga tushdi: {}", addr);
    axum::serve(listener, app).await.unwrap();
}

async fn update_handler(
    axum::extract::State(state): axum::extract::State<Arc<AppState>>,
    Json(payload): Json<serde_json::Value>,
) -> impl IntoResponse {
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

    while let Ok(msg) = rx.recv().await {
        // Message turini aniq ko'rsatamiz
        let text_msg = Message::Text(msg.to_string().into());
        if sender.send(text_msg).await.is_err() { 
            break; 
        }
    }
}
