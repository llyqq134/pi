-- ============================================
-- 5 INSERT-запросов для каждой таблицы
-- Порядок: departments → workers → equipment → records
-- ============================================

-- ---------- DEPARTMENTS ----------
INSERT INTO departments (name) VALUES ('Бухгалтерия');
INSERT INTO departments (name) VALUES ('IT-отдел');
INSERT INTO departments (name) VALUES ('Отдел кадров');
INSERT INTO departments (name) VALUES ('Маркетинг');
INSERT INTO departments (name) VALUES ('Продажи');

-- ---------- WORKERS ----------
INSERT INTO workers (name, jobtitle, department_id, department_name, password, accesslevel)
VALUES ('qwe', 'admin', (SELECT id FROM departments WHERE name = 'evil'), 'evil', '123', 3);

INSERT INTO workers (name, jobtitle, department_id, department_name, password, accesslevel)
VALUES ('Петрова Анна', 'admin', (SELECT id FROM departments WHERE name = 'IT-отдел'), 'IT-отдел', 'admin456', 3);

INSERT INTO workers (name, jobtitle, department_id, department_name, password, accesslevel)
VALUES ('Сидоров Пётр', 'employee', (SELECT id FROM departments WHERE name = 'Отдел кадров'), 'Отдел кадров', 'qwerty', 1);

INSERT INTO workers (name, jobtitle, department_id, department_name, password, accesslevel)
VALUES ('Козлова Мария', 'manager', (SELECT id FROM departments WHERE name = 'Маркетинг'), 'Маркетинг', 'mkt789', 2);

INSERT INTO workers (name, jobtitle, department_id, department_name, password, accesslevel)
VALUES ('Новиков Алексей', 'employee', (SELECT id FROM departments WHERE name = 'Продажи'), 'Продажи', 'sales01', 1);

-- ---------- EQUIPMENT ----------
INSERT INTO equipment (name, type, serial_number, inventory_number, status, location)
VALUES ('Ноутбук Lenovo', 'laptop', 'SN-LAP-001', 'INV-001', 'available', 'Склад');

INSERT INTO equipment (name, type, serial_number, inventory_number, status, location)
VALUES ('Монитор Dell 24"', 'monitor', 'SN-MON-002', 'INV-002', 'available', 'Офис');

INSERT INTO equipment (name, type, serial_number, inventory_number, status, location)
VALUES ('Клавиатура Logitech', 'keyboard', 'SN-KBD-003', 'INV-003', 'available', 'Склад');

INSERT INTO equipment (name, type, serial_number, inventory_number, status, location)
VALUES ('Мышь беспроводная', 'mouse', 'SN-MOU-004', 'INV-004', 'available', 'Склад');

INSERT INTO equipment (name, type, serial_number, inventory_number, status, location)
VALUES ('Телефон Samsung', 'phone', 'SN-PHN-005', 'INV-005', 'available', 'Офис');

INSERT INTO equipment (name, type, serial_number, inventory_number, status, location)
VALUES ('kal', 'phone', 'SN-PHN-005', 'INV-005', 'available', 'Офис');

-- ---------- EQUIPMENT_RECORDS ----------
INSERT INTO equipment_records (equipment_id, worker_id, worker_name, department_id, department_name, issued_at, returned_at, expected_return_date, status)
VALUES (
  (SELECT id FROM equipment WHERE inventory_number = 'INV-001' LIMIT 1),
  (SELECT id FROM workers WHERE name = 'qwe' LIMIT 1),
  'qwe',
  (SELECT department_id FROM workers WHERE name = 'qwe' LIMIT 1),
  'evil',
  NOW(), NULL, '2025-03-01'::date, 'issued'
);

INSERT INTO equipment_records (equipment_id, worker_id, worker_name, department_id, department_name, issued_at, returned_at, expected_return_date, status)
VALUES (
  (SELECT id FROM equipment WHERE inventory_number = 'INV-002' LIMIT 1),
  (SELECT id FROM workers WHERE name = 'Петрова Анна' LIMIT 1),
  'Петрова Анна',
  (SELECT department_id FROM workers WHERE name = 'Петрова Анна' LIMIT 1),
  'IT-отдел',
  NOW(), NULL, '2025-03-15'::date, 'issued'
);

INSERT INTO equipment_records (equipment_id, worker_id, worker_name, department_id, department_name, issued_at, returned_at, expected_return_date, status)
VALUES (
  (SELECT id FROM equipment WHERE inventory_number = 'INV-003' LIMIT 1),
  (SELECT id FROM workers WHERE name = 'Сидоров Пётр' LIMIT 1),
  'Сидоров Пётр',
  (SELECT department_id FROM workers WHERE name = 'Сидоров Пётр' LIMIT 1),
  'Отдел кадров',
  NOW(), NULL, '2025-02-28'::date, 'issued'
);

INSERT INTO equipment_records (equipment_id, worker_id, worker_name, department_id, department_name, issued_at, returned_at, expected_return_date, status)
VALUES (
  (SELECT id FROM equipment WHERE inventory_number = 'INV-004' LIMIT 1),
  (SELECT id FROM workers WHERE name = 'Козлова Мария' LIMIT 1),
  'Козлова Мария',
  (SELECT department_id FROM workers WHERE name = 'Козлова Мария' LIMIT 1),
  'Маркетинг',
  NOW(), NULL, '2025-04-01'::date, 'issued'
);

INSERT INTO equipment_records (equipment_id, worker_id, worker_name, department_id, department_name, issued_at, returned_at, expected_return_date, status)
VALUES (
  (SELECT id FROM equipment WHERE inventory_number = 'INV-005' LIMIT 1),
  (SELECT id FROM workers WHERE name = 'Новиков Алексей' LIMIT 1),
  'Новиков Алексей',
  (SELECT department_id FROM workers WHERE name = 'Новиков Алексей' LIMIT 1),
  'Продажи',
  NOW(), NULL, '2025-03-20'::date, 'issued'
);
