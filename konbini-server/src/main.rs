use axum::{routing::get, Router};
use dotenvy::dotenv;
use libsql::Builder;
use std::sync::Arc;

mod routes;
mod state;

#[tokio::main]
async fn main() {
    let app_env = match std::env::var("APP_ENV").ok() {
        Some(value) => value,
        None => String::from("production"),
    };
    if app_env == "development" {
        dotenv().expect("Failed to load .env file");
    }

    let url = std::env::var("TURSO_DATABASE_URL").expect("TURSO_DATABASE_URL must be set");
    let token = std::env::var("TURSO_AUTH_TOKEN").expect("TURSO_AUTH_TOKEN must be set");

    let db = Builder::new_remote(url, token).build().await.unwrap();

    let state = Arc::new(state::AppState::new(db));

    let app = Router::new()
        .route("/aes", get(routes::aes_cipher))
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("127.0.0.1:3000")
        .await
        .unwrap();
    axum::serve(listener, app).await.unwrap();
}
