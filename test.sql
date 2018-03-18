.load go

.headers on
.mode column

CREATE TABLE TestData(
    Test TEXT,
    Val1 TEXT,
    Val2 TEXT,
    Val3 TEXT
);

INSERT INTO TestData VALUES
    ('Jaro', 'John j Doe', 'John j Doe', ''),
    ('Jaro', 'John j Doe', 'John John Doe', ''),
    ('Regex', 'V\d{2}', 'V0', ''),
    ('Regex', 'V\d{2}', 'V00', ''),
    ('Regex', 'V\d{2}', NULL, ''),
    ('ParseTime', '01/02/2006', 'date', '02/26/2018'),
    ('ParseTime', '01/02/2006', 'time', '02/26/2018'),
    ('ParseTime', '01/02/2006', 'datetime', '02/26/2018'),
    ('ParseTime', '01/02/2006', '2006-01-02', '02/26/2018'),
    ('ParseTime', '01/02/2006', 'datetime', '02/26/2018 654'),
    ('ParseTimeErr', '01/02/2006', 'date', '02/26/2018'),
    ('ParseTimeErr', '01/02/2006', 'time', '02/26/2018'),
    ('ParseTimeErr', '01/02/2006', 'datetime', '02/26/2018'),
    ('ParseTimeErr', '01/02/2006', '2006-01-02', '02/26/2018'),
    ('ParseTimeErr', '01/02/2006', 'datetime', '02/26/2018 654');

.mode list

SELECT
    Val1,
    Val2,
    jaro(Val1, Val2) AS Jaro
FROM TestData
WHERE Test = 'Jaro';

.print

SELECT
    Val1, Val2, regex(Val1, Val2) AS Regex
FROM TestData
WHERE Test = 'Regex';

.print

SELECT
    Val1,
    Val2,
    Val3,
    parsetime(Val1, Val2, Val3, 1) AS ParseTime
FROM TestData
WHERE Test = 'ParseTime';

.print

SELECT
    Val1,
    Val2,
    Val3,
    parsetime(Val1, Val2, Val3, 0) AS ParseTimeErr
FROM TestData
WHERE Test = 'ParseTimeErr';
