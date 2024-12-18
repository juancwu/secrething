use axum::{http::StatusCode, Json};
use konbini_core::crypto;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct AES {
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

pub async fn aes_cipher() -> Result<Json<AES>, StatusCode> {
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
