CREATE TABLE IF NOT EXISTS users(
	id varchar(255) PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	project_v2 TEXT
);

CREATE TABLE IF NOT EXISTS repositories(
	id varchar(255) PRIMARY KEY NOT NULL,
	owner varchar(255) NOT NULL,
	name TEXT NOT NULL,
	created_at DATETIME NOT NULL DEFAULT NOW(),
	FOREIGN KEY (owner) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS issues(
	id varchar(255) PRIMARY KEY NOT NULL,
	url TEXT NOT NULL,
	title TEXT NOT NULL,
	closed INTEGER NOT NULL DEFAULT 0,
	number INTEGER NOT NULL,
	repository varchar(255) NOT NULL,
	CHECK (closed IN (0, 1)),
	FOREIGN KEY (repository) REFERENCES repositories(id)
);

CREATE TABLE IF NOT EXISTS projects(
	id varchar(255) PRIMARY KEY NOT NULL,
	title TEXT NOT NULL,
	url TEXT NOT NULL,
	owner varchar(255) NOT NULL,
	FOREIGN KEY (owner) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS pullrequests(
	id varchar(255) PRIMARY KEY NOT NULL,
	base_ref_name TEXT NOT NULL,
	closed INTEGER NOT NULL DEFAULT 0,
	head_ref_name TEXT NOT NULL,
	url TEXT NOT NULL,
	number INTEGER NOT NULL,
	repository varchar(255) NOT NULL,
	CHECK (closed IN (0, 1)),
	FOREIGN KEY (repository) REFERENCES repositories(id)
);

CREATE TABLE IF NOT EXISTS projectcards(
	id varchar(255) PRIMARY KEY NOT NULL,
	project varchar(255) NOT NULL,
	issue varchar(255),
	pullrequest varchar(255),
	FOREIGN KEY (project) REFERENCES projects(id),
	FOREIGN KEY (issue) REFERENCES issues(id),
	FOREIGN KEY (pullrequest) REFERENCES pullrequests(id),
	CHECK (issue IS NOT NULL OR pullrequest IS NOT NULL)
);

INSERT INTO users(id, name) VALUES
	('U_1', 'user')
;

INSERT INTO repositories(id, owner, name) VALUES
	('REPO_1', 'U_1', 'repo1')
;

INSERT INTO issues(id, url, title, closed, number, repository) VALUES
	('ISSUE_1', 'http://example.com/repo1/issue/1', 'First Issue', 1, 1, 'REPO_1'),
	('ISSUE_2', 'http://example.com/repo1/issue/2', 'Second Issue', 0, 2, 'REPO_1'),
	('ISSUE_3', 'http://example.com/repo1/issue/3', 'Third Issue', 0, 3, 'REPO_1')
;

INSERT INTO projects(id, title, url, owner) VALUES
	('PJ_1', 'My Project', 'http://example.com/project/1', 'U_1')
;

INSERT INTO pullrequests(id, base_ref_name, closed, head_ref_name, url, number, repository) VALUES
	('PR_1', 'main', 1, 'feature/kinou1', 'http://example.com/repo1/pr/1', 1, 'REPO_1'),
	('PR_2', 'main', 0, 'feature/kinou2', 'http://example.com/repo1/pr/2', 2, 'REPO_1')
;
