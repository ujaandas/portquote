INSERT OR IGNORE INTO users (username, password_hash, role) VALUES
  ('admin',  '$2a$10$7EqJtq98hPqEX7fNZaFWoO6phm0NfkgR9SNTZ4o0nhHbFq8rJZ6Nm', 'admin'),
  ('agent1', '$2a$10$KbGkJDzhoI3um2yCMpQ/huXpr6sAQjFzsyqmZyBkvZjajkCpZg.3y', 'agent'),
  ('agent2', '$2a$10$4DOjQTyjEh0oCe9ZyMw5nuSJ/6bUHivSdLXGzj5EtYPrQ/NEIzXZ6', 'agent');

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
    (SELECT id FROM users WHERE username='agent2'),
    (SELECT id FROM ports WHERE name='Yangshan Deepwater Port'),
    1150.50,
    '2025-11-30'
  );