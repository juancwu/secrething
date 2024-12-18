#[derive(PartialEq, Eq, Clone, Copy, Debug)]
pub enum AppEnv {
    DEVELOPMENT,
    STAGING,
    PRODUCTION,
}

pub struct Config {
    pub turso_database_url: String,
    pub turso_auth_token: String,
    pub app_env: AppEnv,
    pub backend_url: String,
    pub resend_api_key: String,
    pub noreply_email: String,
}

impl Config {
    pub fn load() -> Self {
        let app_env = match std::env::var("APP_ENV").ok() {
            Some(value) => match_app_env(value.as_str()),
            None => AppEnv::DEVELOPMENT,
        };
        if app_env == AppEnv::DEVELOPMENT {
            dotenvy::dotenv().expect("Failed to load .env file in development environment");
        }
        let turso_database_url =
            std::env::var("TURSO_DATABASE_URL").expect("TURSO_DATABASE_URL must be set");
        let turso_auth_token =
            std::env::var("TURSO_AUTH_TOKEN").expect("TURSO_AUTH_TOKEN must be set");
        let backend_url = std::env::var("BACKEND_URL").expect("BACKEND_URL must be set");
        let resend_api_key = std::env::var("RESEND_API_KEY").expect("RESEND_API_KEY must be set");
        let noreply_email = std::env::var("NOREPLY_EMAIL").expect("NOREPLY_EMAIL must be set");

        Self {
            turso_database_url,
            turso_auth_token,
            app_env,
            backend_url,
            resend_api_key,
            noreply_email,
        }
    }
}

/// Matches a string env value to the corresponding enum value.
fn match_app_env(app_env: &str) -> AppEnv {
    match app_env {
        "development" => AppEnv::DEVELOPMENT,
        "staging" => AppEnv::STAGING,
        "production" => AppEnv::PRODUCTION,
        _ => AppEnv::DEVELOPMENT,
    }
}
