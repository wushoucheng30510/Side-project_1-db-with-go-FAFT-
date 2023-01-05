LOAD DATA INFILE '/var/lib/mysql-files/stainless.csv' 
INTO TABLE Stainless_Result
FIELDS TERMINATED BY ',' 
LINES TERMINATED BY '\n'
IGNORE 1 ROWS;