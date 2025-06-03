package DataBaseConnect

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func InitDb() {
	var err error
	Db, err = sql.Open("postgres", "host=postgres port=5432 user=postgres password=postgres dbname=practicum sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func MigrateDb() {
	query := `
	CREATE TABLE IF NOT EXISTS programs(
    	id INTEGER PRIMARY KEY NOT NULL,
    	code VARCHAR NOT NULL,
    	name VARCHAR NOT NULL,
    	description VARCHAR NOT NULL,
    	format_education VARCHAR NOT NULL
	);

	CREATE TABLE IF NOT EXISTS subjects(
	    id VARCHAR(255) PRIMARY KEY NOT NULL,
	    russian_name VARCHAR(32) NOT NULL,
	    english_name VARCHAR(32) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS program_subject(
	    program_id INTEGER,
	    subject_id VARCHAR,
	    FOREIGN KEY (program_id) REFERENCES programs (id),
	    FOREIGN KEY (subject_id) REFERENCES subjects (id),
	    PRIMARY KEY(program_id, subject_id)
	);

	CREATE TABLE IF NOT EXISTS passing_scores(
	    passing_score_id INTEGER PRIMARY KEY,
	    program_id INTEGER NOT NULL,
	    year INTEGER NOT NULL,
	    score_budget INTEGER NOT NULL,
	    score_paid INTEGER NOT NULL,
	    available_places INTEGER NOT NULL,
	    education_price INTEGER NOT NULL,
	    FOREIGN KEY(program_id) REFERENCES programs(id)
	);

	CREATE TABLE IF NOT EXISTS subject_groups(
	    id SERIAL PRIMARY KEY,
	    group_type VARCHAR NOT NULL
	);

	CREATE TABLE IF NOT EXISTS program_subject_groups(
	    program_id INTEGER,
	    group_id INTEGER,
	    PRIMARY KEY(program_id, group_id),
	    FOREIGN KEY (program_id) REFERENCES programs(id),
	    FOREIGN KEY (group_id) REFERENCES subject_groups(id)
	);

	CREATE TABLE IF NOT EXISTS subject_group_items(
	    group_id INTEGER,
	    subject_id VARCHAR(255),
	    PRIMARY KEY(group_id, subject_id),
	    FOREIGN KEY(subject_id) REFERENCES subjects(id),
	    FOREIGN KEY(group_id) REFERENCES subject_groups(id)
	);

	CREATE INDEX IF NOT EXISTS passing_scores_program_id_index ON 
	    passing_scores (program_id);

	CREATE INDEX IF NOT EXISTS programs_format_education_index ON
	    programs(format_education);
	`

	_, err := Db.Exec(query)
	if err != nil {
		log.Fatalf("Ошибка при выполнении миграции: %v", err)
	}

	log.Println("Миграция прошла успешно")
}

func DataFiling() {
	query := `
	INSERT INTO programs VALUES
		(1, '01.03.02', 'Прикладная математика и информатика', 'Программирование и искусственный интеллект', 'очная'),
		(2, '13.03.01', 'Теплоэнергетика и теплотехника', 'Информационные системы и технологии в топливно-энергетическом комплексе', 'очная'),
		(3, '15.03.06', 'Мехатроника и роботехника', 'Интеллектуальные робототехнические и мехатронные системы', 'очная'),
		(4, '37.03.01', 'Психология', 'Социальная психология', 'очная'),
		(5, '38.03.01', 'Экономика', 'Экономика и бизнес-аналитика', 'очная'),
		(6, '42.03.01', 'Реклама и связи с общественностью', 'Реклама и связи с общественностью в коммерческой сфере', 'очная'),
		(7, '09.03.02', 'Информационные системы и технологии', 'Информационные технологии и дизайн', 'очная'),
		(8, '39.03.01', 'Социология', 'Социология рекламы и связей с общественностью', 'очная'),
		(9, '54.05.02', 'Живопись', 'Художник-живописец', 'очная'),
		(10, '54.03.03', 'Искусство костюма и текстиля', 'Диджитал-арт и компьютерные технологии в современном искусстве', 'очная');

    INSERT INTO subjects VALUES
		('russian_language', 'Русский язык', 'Russian language'),
		('mathematics_is_specialized', 'Математика (профильная)', 'Mathematics is specialized'),
		('physics', 'Физика', 'Physics'),
		('chemistry', 'Химия', 'Chemistry'),
		('computer_science_and_ict', 'Информатика', 'Computer Science and ICT'),
		('history', 'История', 'History'),
		('social_studies', 'Обществознание', 'Social Studies'),
		('english_language', 'Английский язык', 'English language'),
		('biology', 'Биология', 'Biology'),
		('basic_mathematics', 'Математика (базовая)', 'Basic mathematics'),
		('geography', 'География', 'Geography'),
		('literature', 'Литература', 'Literature'),
		('composition', 'Композиция', 'Composition'),
		('drawing', 'Рисунок', 'Drawing');

	INSERT INTO passing_scores VALUES
		(1, 1, 2024, 244, 164, 16, 290000),
		(2, 2, 2024, 185, 184, 25, 320000),
		(3, 3, 2024, 230, 0, 14, 330000),
		(4, 4, 2024, 253, 146, 10, 290000),
		(5, 5, 2024, 269, 0, 5, 29000),
		(6, 6, 2024, 286, 145, 4, 320000),
		(7, 7, 2024, 246, 145, 92, 320000),
		(8, 8, 2024, 264, 144, 12, 290000),
		(9, 9, 2024, 308, 234, 3, 570000),
		(10, 10, 2024, 320, 191, 12, 570000);

	INSERT INTO program_subject VALUES 
		(1, 'russian_language'),
		(1, 'mathematics_is_specialized'),
		(1, 'physics'),
		(1, 'computer_science_and_ict'),
		(2, 'russian_language'),
		(2, 'mathematics_is_specialized'),
		(2, 'physics'),
		(2, 'computer_science_and_ict'),
		(2, 'chemistry'),	
		(3, 'russian_language'),
		(3, 'mathematics_is_specialized'),
		(3, 'physics'),
		(3, 'computer_science_and_ict'),
		(4, 'russian_language'),
		(4, 'biology'),
		(4, 'mathematics_is_specialized'),
		(4, 'social_studies'),
		(5, 'russian_language'),
		(5, 'mathematics_is_specialized'),
		(5, 'computer_science_and_ict'),
		(5, 'social_studies'),
		(5, 'history'),
		(5, 'geography'),
		(5, 'english_language'),
		(6, 'social_studies'),
		(6, 'history'),
		(6, 'english_language'),
		(6, 'russian_language'),
		(7, 'russian_language'),
		(7, 'mathematics_is_specialized'),
		(7, 'physics'),
		(7, 'computer_science_and_ict'),
		(8, 'russian_language'),
		(8, 'mathematics_is_specialized'),
		(8, 'social_studies'),
		(8, 'history'),
		(8, 'english_language'),
		(9, 'russian_language'),
		(9, 'literature'),
		(9, 'composition'),
		(9, 'drawing'),
		(10, 'russian_language'),
		(10, 'literature'),
		(10, 'composition'),
		(10, 'drawing');
		
		INSERT INTO subject_groups VALUES
		(1, 'required'),
		(2, 'required'),
		(3, 'optional'),
		(4, 'optional'),
		(5, 'required'),
		(6, 'optional'),
		(7, 'optional'),
		(8, 'required'),
		(9, 'optional'),
		(10, 'optional'),
		(11, 'required'),
		(12, 'optional');
		
		INSERT INTO subject_group_items VALUES
			(1, 'russian_language'),
			(2, 'mathematics_is_specialized'),
			(3, 'physics'),
			(3, 'computer_science_and_ict'),
			(4, 'physics'),
			(4, 'chemistry'),
			(4, 'computer_science_and_ict'),
			(5, 'biology'),
			(6, 'mathematics_is_specialized'),
			(7, 'english_language'),
			(7, 'geography'),
			(7, 'computer_science_and_ict'),
			(7, 'history'),
			(7, 'social_studies'),
			(8, 'social_studies'),
			(9, 'english_language'),
			(9, 'history'),
			(10, 'history'),
			(10, 'mathematics_is_specialized'),
			(10, 'english_language'),
			(11, 'literature'),
			(12, 'composition'),
			(12, 'drawing');

	INSERT INTO program_subject_groups VALUES
		(1, 1),
			(1, 2),
				(1, 3),
					(2, 1),
						(2, 2),
							(2, 4),
								(3, 1),
									(3, 2),
										(3, 3),
											(4, 1),
												(4, 5),
													(4, 6),
														(5, 1),
															(5, 2),
																(5, 7),
															(6, 1),
														(6, 8),
													(6, 9),
												(7, 1),
											(7, 2),
										(7, 3),
									(8, 1),
								(8, 8),
							(8, 10),
						(9, 1),
					(9, 11),
				(9, 12),
			(10, 1),
		(10, 11),
	(10, 12);
		 `
	_, err := Db.Exec(query)
	if err != nil {
		log.Fatalf("Ошибка при вставке данных: %v", err)
	}

	log.Println("Данные успешно внесены")
}
