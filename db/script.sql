-- Создание таблицы пользователей
DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    last_name VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);

-- Создание таблицы токенов
DROP TABLE IF EXISTS tokens;
CREATE TABLE tokens (
    id VARCHAR(36) PRIMARY KEY,
    token VARCHAR(255),
    user_id VARCHAR(36) REFERENCES users(id) ON DELETE CASCADE,
    is_regular BOOLEAN
);

-- Создание таблицы ролей
DROP TABLE IF EXISTS roles;
CREATE TABLE roles (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Создание таблицы для связи пользователей и ролей (многие ко многим)
DROP TABLE IF EXISTS user_roles;
CREATE TABLE user_roles (
    user_id VARCHAR(36) REFERENCES users(id) ON DELETE CASCADE,
    role_id VARCHAR(36) REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- Создание таблицы больниц
DROP TABLE IF EXISTS hospitals;
CREATE TABLE hospitals (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);

-- Создание таблицы кабинетов
DROP TABLE IF EXISTS rooms;
CREATE TABLE rooms (
    id VARCHAR(36) PRIMARY KEY,
    hospital_id VARCHAR(36) REFERENCES hospitals(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL
);

-- Создание таблицы врачей
DROP TABLE IF EXISTS doctors;
CREATE TABLE doctors (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) REFERENCES users(id) ON DELETE CASCADE,
    specialization VARCHAR(100)
);

-- Создание таблицы расписания
DROP TABLE IF EXISTS timetable;
CREATE TABLE timetable (
    id VARCHAR(36) PRIMARY KEY,
    hospital_id VARCHAR(36) REFERENCES hospitals(id) ON DELETE CASCADE,
    doctor_id VARCHAR(36) REFERENCES doctors(id) ON DELETE CASCADE,
    time_from TIMESTAMP NOT NULL,
    time_to TIMESTAMP NOT NULL,
    room VARCHAR(50) REFERENCES rooms(id) NULL,
    CHECK (time_to > time_from),
    CHECK ((EXTRACT(MINUTE FROM time_to - time_from) % 30 = 0) AND (EXTRACT(SECOND FROM time_to - time_from) = 0))
);

-- Создание таблицы для записей на прием
DROP TABLE IF EXISTS appointments;
CREATE TABLE appointments (
    id VARCHAR(36) PRIMARY KEY,
    timetable_id VARCHAR(36) REFERENCES timetable(id) ON DELETE CASCADE,
    user_id VARCHAR(36) REFERENCES users(id) ON DELETE CASCADE,
    appointment_time TIMESTAMP NOT NULL,
    UNIQUE (timetable_id, appointment_time)
);

-- Создание таблицы истории посещений
DROP TABLE IF EXISTS visit_history;
CREATE TABLE visit_history (
    id VARCHAR(36) PRIMARY KEY,
    patient_id VARCHAR(36) REFERENCES users(id) ON DELETE CASCADE,
    hospital_id VARCHAR(36) REFERENCES hospitals(id) ON DELETE CASCADE,
    doctor_id VARCHAR(36) REFERENCES doctors(id) ON DELETE CASCADE,
    room VARCHAR(50) REFERENCES rooms(id) NULL,
    visit_date TIMESTAMP NOT NULL,
    data TEXT,
    UNIQUE (patient_id, visit_date)
);

-- создание ролей
INSERT INTO roles (id, name) VALUES ('0', 'Admin');
INSERT INTO roles (id, name) VALUES ('1', 'Manager');
INSERT INTO roles (id, name) VALUES ('2', 'Doctor');
INSERT INTO roles (id, name) VALUES ('3', 'User');

-- создание пользователя по умолчанию
INSERT INTO users (id, last_name, first_name, username, password) VALUES ('ee7741a7-4301-36a8-fcfb-845fe752a35e', 'Ivanov', 'Egor', 'admin', 'jGl25bVBBBW96Qi9Te4V37Fnqchz/Eu4qB9vKrRIqRg=') RETURNING id;
INSERT INTO users (id, last_name, first_name, username, password) VALUES ('3f1c0b5e-972c-4c1c-819f-101fd5ffe30c', 'Ivanov', 'Egor', 'manager', 'buSkac1OkQU4R/XT/LYdvMkejw7xC+d0jaTEobo4LRc=') RETURNING id;
INSERT INTO users (id, last_name, first_name, username, password) VALUES ('d8b9f484-c3c7-4e14-a5fd-c1e9e2f73f5f', 'Ivanov', 'Egor', 'doctor', 'cvS+idbrqxSW4h44vNfIygpokorzCBrX3/h+dy6zUMI=') RETURNING id;
INSERT INTO users (id, last_name, first_name, username, password) VALUES ('469b2081-a77c-41d2-b45c-8c484350cd3c', 'Ivanov', 'Egor', 'user', 'BPiZbadjt6lpsQKO4wB1aerzpjVIbdqyEdUSyFud+Ps=') RETURNING id;

-- задание ролей
INSERT INTO user_roles (user_id, role_id) VALUES ('ee7741a7-4301-36a8-fcfb-845fe752a35e', '0');
INSERT INTO user_roles (user_id, role_id) VALUES ('3f1c0b5e-972c-4c1c-819f-101fd5ffe30c', '1');
INSERT INTO user_roles (user_id, role_id) VALUES ('d8b9f484-c3c7-4e14-a5fd-c1e9e2f73f5f', '2');
INSERT INTO user_roles (user_id, role_id) VALUES ('469b2081-a77c-41d2-b45c-8c484350cd3c', '3');

-- создание больниц
INSERT INTO hospitals (id, name, address, contact_phone) VALUES 
('a5d0a195-2a2c-4857-b2f0-4f7e1b552bd6', 'Городская больница', 'Улица Ленина, 10', '+7 (495) 000-0001'),
('b7dfad4c-336c-4da6-bc4a-1d8c15a7b85e', 'Центральная клиника', 'Улица Пушкина, 5', '+7 (495) 000-0002'),
('e5aef428-780f-4b49-842e-ee246908bc30', 'Детская больница', 'Улица Мира, 15', '+7 (495) 000-0003');

-- создание кабинетов
INSERT INTO rooms (id, hospital_id, name) VALUES 
('f3c77224-5a32-4e15-8ff9-9c464d3f5cd2', 'a5d0a195-2a2c-4857-b2f0-4f7e1b552bd6', 'Кабинет 101'),
('f50e8fa5-7bf6-4d8b-90be-60e1ad8b93a3', 'b7dfad4c-336c-4da6-bc4a-1d8c15a7b85e', 'Процедурная 201'),
('b01e3df4-19f4-470b-b488-be70f0c30f3a', 'e5aef428-780f-4b49-842e-ee246908bc30', 'Офис 301');

-- создание докторов
INSERT INTO doctors (id, user_id, specialization) VALUES 
('a3210ffa-d31b-4a98-8ba8-4557cd7b794c', 'd8b9f484-c3c7-4e14-a5fd-c1e9e2f73f5f', 'Терапевт');

-- создание расписаний
INSERT INTO timetable (id, hospital_id, doctor_id, time_from, time_to, room) VALUES 
('c7c6bca9-c382-4286-8536-cb714f0a3b91', 'a5d0a195-2a2c-4857-b2f0-4f7e1b552bd6', 'a3210ffa-d31b-4a98-8ba8-4557cd7b794c', '2024-10-01 09:00:00', '2024-10-01 12:00:00', 'f3c77224-5a32-4e15-8ff9-9c464d3f5cd2'),
('d295c222-3f45-49e3-af84-60b83679992c', 'b7dfad4c-336c-4da6-bc4a-1d8c15a7b85e', 'a3210ffa-d31b-4a98-8ba8-4557cd7b794c', '2024-10-02 10:00:00', '2024-10-02 14:00:00', 'f50e8fa5-7bf6-4d8b-90be-60e1ad8b93a3'),
('bfae49ad-5c5f-4829-884b-x0d98eaa0a2d', 'e5aef428-780f-4b49-842e-ee246908bc30', 'a3210ffa-d31b-4a98-8ba8-4557cd7b794c', '2024-10-03 08:00:00', '2024-10-03 12:00:00', 'b01e3df4-19f4-470b-b488-be70f0c30f3a');

-- создание записей на приемы
INSERT INTO appointments (id, timetable_id, user_id, appointment_time) VALUES 
('7f9bf888-571d-4ec8-8e7f-b4a7ac0d5b20', 'c7c6bca9-c382-4286-8536-cb714f0a3b91', '469b2081-a77c-41d2-b45c-8c484350cd3c', '2024-10-01 09:30:00'),
('e0f40453-431b-46e3-ba44-4d653ead7a17', 'd295c222-3f45-49e3-af84-60b83679992c', '469b2081-a77c-41d2-b45c-8c484350cd3c', '2024-10-02 11:00:00'),
('fd7c215b-4cd9-4ae1-bc14-a990c9455f28', 'bfae49ad-5c5f-4829-884b-x0d98eaa0a2d', '469b2081-a77c-41d2-b45c-8c484350cd3c', '2024-10-03 09:00:00');

-- создание истории посещений
INSERT INTO visit_history (id, patient_id, hospital_id, doctor_id, room, visit_date, data) VALUES 
('4da06b43-8d57-448c-95e1-750f283c89c3', '469b2081-a77c-41d2-b45c-8c484350cd3c', 'a5d0a195-2a2c-4857-b2f0-4f7e1b552bd6', 'a3210ffa-d31b-4a98-8ba8-4557cd7b794c', 'f3c77224-5a32-4e15-8ff9-9c464d3f5cd2', '2024-10-01 10:00:00', 'Общее состояние нормальное.'),
('2d97a6e0-eee5-4a56-9208-eaa959d2c830', '469b2081-a77c-41d2-b45c-8c484350cd3c', 'b7dfad4c-336c-4da6-bc4a-1d8c15a7b85e', 'a3210ffa-d31b-4a98-8ba8-4557cd7b794c', 'f50e8fa5-7bf6-4d8b-90be-60e1ad8b93a3', '2024-10-02 12:00:00', 'Посещение по поводу боли в животе.'),
('bd20f11f-ec63-4a3c-834b-e314227cb4f2', '469b2081-a77c-41d2-b45c-8c484350cd3c', 'e5aef428-780f-4b49-842e-ee246908bc30', 'a3210ffa-d31b-4a98-8ba8-4557cd7b794c', 'b01e3df4-19f4-470b-b488-be70f0c30f3a', '2024-10-03 11:00:00', 'Подписать справку.');

