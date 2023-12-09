CREATE TABLE birthdays(
    id INT NOT NULL AUTO_INCREMENT,
    chat_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    day INT NOT NULL,
    month INT NOT NULL,
    PRIMARY KEY (id)
)

SELECT * FROM birthdays