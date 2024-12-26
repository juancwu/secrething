use axum::{
    routing::{get, post},
    Router,
};
use libsql::Builder;
use std::sync::Arc;

mod config;
mod state;

#[tokio::main]
async fn main() {
    let config = config::Config::load();

    let db = Builder::new_remote(config.turso_database_url, config.turso_auth_token)
        .build()
        .await
        .unwrap();

    let app_state = Arc::new(state::AppState::new(db));

    let app = Router::new()
        .route("/", get(hello_world))
        .with_state(app_state);

    let listener = tokio::net::TcpListener::bind(format!("0.0.0.0:{}", config.port))
        .await
        .unwrap();
    axum::serve(listener, app).await.unwrap();
}

async fn hello_world() -> &'static str {
    "hello world"
}
