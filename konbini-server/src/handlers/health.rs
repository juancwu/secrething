use std::sync::Arc;

use axum::extract::State;

use crate::state;

pub async fn get_health(_: State<Arc<state::AppState>>) -> &'static str {
    "Healthy!"
}
