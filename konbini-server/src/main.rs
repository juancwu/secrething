use axum::{
    routing::{get, post},
    Router,
};
use std::sync::Arc;
use tokio::sync::RwLock;
use tracing::info;

mod routes;
mod state;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();
    info!("Starting server...");

    let state = Arc::new(state::AppState {
        users: RwLock::new(Vec::new()),
    });

    let app = Router::new()
        .route("/users", get(routes::list_users))
        .route("/user", post(routes::create_user))
        .route("/aes", get(routes::aes_cipher))
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("127.0.0.1:3000")
        .await
        .unwrap();
    info!("Server running on http://127.0.0.1:3000");
    axum::serve(listener, app).await.unwrap();
}
