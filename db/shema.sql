-- db/schema.sql

--divisiones dentro de la liga 
CREATE TABLE IF NOT EXISTS divisiones (
    id          SERIAL PRIMARY KEY,
    nombre      VARCHAR(100) NOT NULL,
    temporada   VARCHAR(50)
);
--equipos dentro de cada division
CREATE TABLE IF NOT EXISTS equipos (
    id          SERIAL PRIMARY KEY,
    division_id INT NOT NULL REFERENCES divisiones(id) ON DELETE CASCADE,
    nombre      VARCHAR(100) NOT NULL,
    logo_url    VARCHAR(255)
);
--jugadores dentro de cada equipo
CREATE TABLE IF NOT EXISTS jugadores (
    id          SERIAL PRIMARY KEY,
    equipo_id   INT NOT NULL REFERENCES equipos(id) ON DELETE CASCADE,
    nombre      VARCHAR(100) NOT NULL,
    numero      INT,
    posicion    VARCHAR(50)
);
--partidos dentro de cada division
CREATE TABLE IF NOT EXISTS partidos (
    id              SERIAL PRIMARY KEY,
    division_id     INT NOT NULL REFERENCES divisiones(id) ON DELETE CASCADE,
    equipo_local_id INT NOT NULL REFERENCES equipos(id),
    equipo_visita_id INT NOT NULL REFERENCES equipos(id),
    carreras_local  INT,
    carreras_visita INT,
    campo            VARCHAR(100),
    fecha           DATE,
    estado          VARCHAR(20) DEFAULT 'programado'
);