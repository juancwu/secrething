use axum::{
    extract::State,
    http::StatusCode,
    routing::{get, post},
    Json, Router,
};
use konbini_core::User;
use std::sync::Arc;
use tokio::sync::RwLock;
use tracing::info;

struct AppState {
    users: RwLock<Vec<User>>,
}

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();
    info!("Starting server...");

    let state = Arc::new(AppState {
        users: RwLock::new(Vec::new()),
    });

    let app = Router::new()
        .route("/users", get(list_users))
        .route("/user", post(create_user))
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("127.0.0.1:3000")
        .await
        .unwrap();
    info!("Server running on http://127.0.0.1:3000");
    axum::serve(listener, app).await.unwrap();
}

async fn list_users(State(state): State<Arc<AppState>>) -> Json<Vec<User>> {
    let users = state.users.read().await;
    Json(users.clone())
}

#[derive(serde::Deserialize)]
struct CreateUserRequest {
    username: String,
    email: String,
}

async fn create_user(
    State(state): State<Arc<AppState>>,
    Json(payload): Json<CreateUserRequest>,
) -> Result<Json<User>, StatusCode> {
    let user = User::new(payload.username, payload.email);
    if user.validate() {
        state.users.write().await.push(user.clone());
        Ok(Json(user.clone()))
    } else {
        Err(StatusCode::BAD_REQUEST)
    }
}
