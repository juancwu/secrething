use axum::{
    extract::State,
    http::StatusCode,
    routing::{get, post},
    Json, Router,
};
use konbini_core::crypto;
use konbini_core::User;
use serde::{Deserialize, Serialize};
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
        .route("/aes", get(aes_cipher))
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

#[derive(Debug, Serialize, Deserialize, Clone)]
struct AES {
    pub key: String,
    pub plaintext: String,
    pub ciphertext: String,
}

impl AES {
    pub fn new(key: String, plaintext: String, ciphertext: String) -> Self {
        Self {
            plaintext,
            ciphertext,
            key,
        }
    }
}

async fn aes_cipher() -> Result<Json<AES>, StatusCode> {
    let key = match crypto::aes::generate_key() {
        Ok(k) => k,
        Err(_) => {
            return Err(StatusCode::BAD_REQUEST);
        }
    };
    let data = b"some text";
    let ciphertext = match crypto::aes::encrypt(&key, data) {
        Ok(k) => k,
        Err(_) => {
            return Err(StatusCode::BAD_REQUEST);
        }
    };
    let plaintext = match crypto::aes::decrypt(&key, &ciphertext) {
        Ok(k) => k,
        Err(_) => {
            return Err(StatusCode::BAD_REQUEST);
        }
    };
    let plaintext = match String::from_utf8(plaintext) {
        Ok(k) => k,
        Err(_) => {
            return Err(StatusCode::BAD_REQUEST);
        }
    };
    let key = crypto::aes::encode_to_hex(&key);
    let ciphertext = crypto::aes::encode_to_hex(&ciphertext);
    let aes = AES::new(key, plaintext, ciphertext);
    Ok(Json(aes.clone()))
}
