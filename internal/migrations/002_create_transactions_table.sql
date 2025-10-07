CREATE TABLE IF NOT EXISTS transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    service_id INT DEFAULT NULL,
    type ENUM('airtime','data','bill','fund') NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    status ENUM('pending','success','failed') DEFAULT 'pending',
    reference VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_transaction_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_transaction_service FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE SET NULL
);