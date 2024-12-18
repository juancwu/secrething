use serde::{Deserialize, Serialize};
use uuid::Uuid;

pub mod crypto;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct User {
    pub id: String,
    pub username: String,
    pub email: String,
}

impl User {
    pub fn new(username: String, email: String) -> Self {
        Self {
            id: Uuid::new_v4().to_string(),
            username,
            email,
        }
    }

    pub fn validate(&self) -> bool {
        !self.username.is_empty() && !self.email.is_empty() && self.email.contains('@')
    }
}
