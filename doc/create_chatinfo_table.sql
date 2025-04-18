-- 创建对话信息表chatinfo
CREATE TABLE chat_info (
	id serial NOT NULL,
	user_id varchar(10) NOT NULL,
	session_id varchar(32) NOT NULL,
	chat_name varchar(255) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP(0) NOT NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP(0) NOT NULL,
	CONSTRAINT chatinfo_pkey PRIMARY KEY (id)
);

ALTER TABLE chat_info ADD COLUMN is_delete boolean DEFAULT false NOT NULL;