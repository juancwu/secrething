use konbini_core::User;
use tokio::sync::RwLock;

pub struct AppState {
    pub users: RwLock<Vec<User>>,
}
