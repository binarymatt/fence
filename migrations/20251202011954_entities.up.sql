CREATE TABLE entities (
  id TEXT NOT NULL,
  type TEXT NOT NULL,
  parents TEXT,
  attributes TEXT,
  tags TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMETAMP,
  PRIMARY KEY (`id`,`type`)
)
