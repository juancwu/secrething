use aes_gcm::{
    aead::{Aead, OsRng},
    AeadCore, Aes256Gcm, Key, KeyInit, Nonce,
};
use anyhow::{anyhow, Ok, Result};

/// Generates a random new key for AES-256-GCM
pub fn generate_key() -> Result<[u8; 32]> {
    let key = Aes256Gcm::generate_key(OsRng);
    let key: [u8; 32] = key
        .try_into()
        .map_err(|e| anyhow!("Failed to convert key to byte array: {}", e))?;
    Ok(key)
}

/// Encrypts the given byte array using the provided byte array key.
pub fn encrypt(key: &[u8; 32], data: &[u8]) -> Result<Vec<u8>> {
    let key: &Key<Aes256Gcm> = key.into();
    let cipher = Aes256Gcm::new(&key);
    let nonce = Aes256Gcm::generate_nonce(&mut OsRng);
    let ciphertext = cipher
        .encrypt(&nonce, data)
        .map_err(|e| anyhow!("Failed to encrypt data: {}", e))?;
    let mut data = nonce.to_vec();
    data.extend(&ciphertext);
    Ok(data)
}

/// Decrypts the given byte array using the provided byte array key.
pub fn decrypt(key: &[u8], data: &[u8]) -> Result<Vec<u8>> {
    if data.len() < 12 {
        return Err(anyhow!("Encrypted data too short"));
    }
    let key: &Key<Aes256Gcm> = key.into();
    let cipher = Aes256Gcm::new(&key);
    let (nonce, ciphertext) = data.split_at(12);
    let nonce = Nonce::from_slice(nonce);
    let plaintext = cipher
        .decrypt(nonce, ciphertext)
        .map_err(|e| anyhow!("Failed to decrypt data: {}", e))?;
    Ok(plaintext)
}

pub fn encode_to_hex(data: &[u8]) -> String {
    hex::encode(data)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_complete_integration() {
        let key = generate_key().unwrap();
        assert_eq!(key.len(), 32);
        let data = b"plaintext";
        let encrypted = encrypt(&key, data).unwrap();
        let decrypted = decrypt(&key, &encrypted).unwrap();
        assert_eq!(decrypted, data.to_vec());
    }
}
