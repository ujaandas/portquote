INSERT OR IGNORE INTO users (username, password_hash, role) VALUES
  ('admin1',  '$2y$10$y4CUQzxBlt2yrTcQIz92BOBsRGhpvBXmZ8Tl4.B4no6tZiDXADBCW', 'admin'),
  ('agent1', '$2y$10$y4CUQzxBlt2yrTcQIz92BOBsRGhpvBXmZ8Tl4.B4no6tZiDXADBCW', 'agent'),
  ('agent2', '$2y$10$y4CUQzxBlt2yrTcQIz92BOBsRGhpvBXmZ8Tl4.B4no6tZiDXADBCW', 'agent'),
  ('crew1',  '$2y$10$y4CUQzxBlt2yrTcQIz92BOBsRGhpvBXmZ8Tl4.B4no6tZiDXADBCW', 'crew');

INSERT OR IGNORE INTO ports (name, country, city) VALUES
  ('Yangshan Deepwater Port', 'China', 'Shanghai'),
  ('Kwai Tsing Container Terminal', 'Hong Kong SAR', 'Hong Kong'),
  ('Euromax Terminal Rotterdam', 'Netherlands', 'Rotterdam'),
  ('Bell Street Cruise Terminal', 'USA', 'Seattle');

INSERT OR IGNORE INTO quotations (agent_id, port_id, rate, valid_until) VALUES
  (
    (SELECT id FROM users WHERE username='agent1'),
    (SELECT id FROM ports WHERE name='Yangshan Deepwater Port'),
    1200.00,
    '2025-12-31'
  ),
  (
    (SELECT id FROM users WHERE username='agent1'),
    (SELECT id FROM ports WHERE name='Yangshan Deepwater Port'),
    1150.50,
    '2025-11-30'
  ),
  (
    (SELECT id FROM users WHERE username='agent1'),
    (SELECT id FROM ports WHERE name='Yangshan Deepwater Port'),
    900,
    '2025-12-15'
  ),
  (
    (SELECT id FROM users WHERE username='agent2'),
    (SELECT id FROM ports WHERE name='Yangshan Deepwater Port'),
    1150.50,
    '2025-11-30'
  );