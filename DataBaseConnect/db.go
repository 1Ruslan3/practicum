package DataBaseConnect

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func InitDb() {
	var err error
	Db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=practicum sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	CREATE TABLE programs(
    id INTEGER NOT NULL,
    code CHAR(8) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    format_education VARCHAR(12) NOT NULL
);


CREATE TABLE subjects(
    id VARCHAR(255) NOT NULL,
    russian_name VARCHAR(32) NOT NULL,
    English_name VARCHAR(32) NOT NULL
);

CREATE TABLE program_subject(
    program_id INTEGER NOT NULL,
    subject_id VARCHAR(255) NOT NULL
);


CREATE TABLE passing_scores(
    passing_score_id SMALLINT NOT NULL,
    program_id INTEGER NOT NULL,
    year SMALLINT NOT NULL,
    score_budget SMALLINT NOT NULL,
    score_paid SMALLINT NOT NULL,
    available_places SMALLINT NOT NULL,
    education_price INTEGER NOT NULL
);

CREATE TABLE subject_scores(
    id INTEGER NOT NULL,
    subject_id VARCHAR(255) NOT NULL,
    year SMALLINT NOT NULL,
    score SMALLINT NOT NULL
);

ALTER TABLE
    subjects ADD PRIMARY KEY(id);

ALTER TABLE
    programs ADD PRIMARY KEY(id);

ALTER TABLE
    subject_scores ADD PRIMARY KEY(id);

ALTER TABLE
    subject_scores ADD CONSTRAINT subject_scores_subject_id_foreign FOREIGN KEY(subject_id) REFERENCES subjects(id);

ALTER TABLE
    passing_scores ADD PRIMARY KEY(passing_score_id);

ALTER TABLE
    passing_scores ADD CONSTRAINT passing_scores_program_id_foreign FOREIGN KEY(program_id) REFERENCES programs(id);

ALTER TABLE 
    program_subject ADD CONSTRAINT program_subject_program_id_foreign FOREIGN KEY (program_id) REFERENCES programs (id);

ALTER TABLE 
    program_subject ADD CONSTRAINT program_subject_subject_id_foreign FOREIGN KEY (subject_id) REFERENCES subjects (id);

ALTER TABLE
    program_subject ADD PRIMARY KEY(program_id, subject_id);

CREATE INDEX passing_scores_program_id_index ON 
    passing_scores (program_id);

CREATE INDEX subject_scores_subject_id_index ON 
    subject_scores (subject_id);

CREATE INDEX programs_format_education_index ON
    programs(format_education);
	
	INSERT INTO subjects VALUES
 (1, 'математика', 'math'),
 (2, 'русский', 'russian');
	`)
	if err != nil {
		log.Fatalf("Error executing schema creation: %v", err)
	}

	log.Printf("IUDHWFGBDIWUDOWINDOWIDOIW")
}
