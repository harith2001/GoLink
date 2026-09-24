-- Curated demo rates for the assistant. Confirm against the Customs gazette before relying on them.
INSERT INTO import_policy (effective_date, duty_rate_percent, description) VALUES
('2018-09-01', 0.00, 'Demo: EV concession window. Curated figure, not a gazette extract.'),
('2020-01-15', 10.00, 'Demo: reduced concession still in place. Curated figure.'),
('2021-08-19', NULL, 'Demo: vehicle imports largely suspended; duty was not the binding constraint.'),
('2023-06-01', NULL, 'Demo: import restriction still the practical rule for most EVs.'),
('2025-02-01', 30.00, 'Demo: imports reopened with a duty component. Curated figure, verify before use.');
