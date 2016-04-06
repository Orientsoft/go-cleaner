#go-cleaner
How to Deploy
-------------
1. git clone https://github.com/Orientsoft/go-cleaner  
2. cd go-cleaner  
3. ./install  
4. cd bin  
5. edit config file  
6. ./go-cleaner  

Parameter
---------
Usage: go-cleaner [flags]  
  
  -config="./default.conf": Config file path  
		Default: ./default.conf  
  -input="./input.csv": Input file path  
		Default: ./input.csv  
		
CSV Charset
-----------
Since golang uses UTF-8 strings internally, input CSV should be converted to **UTF-8** before processing.  

Config File
-----------
TOML, Tom's Obvious, Minimal Language, is a de facto standard as GO config file description language.  
go-cleaner uses github.com/BurntSushi/toml library to parse its config file into internal structs.  
Refer to [TOML v0.2.0 spec](https://github.com/toml-lang/toml/blob/master/versions/en/toml-v0.2.0.md) for detailed info about this language.  
For go-cleaner, you should define column position to clean and corresponding cleansing type. For example:  
```toml
[columns]
  [columns.name1]
  columnNo = 1
  columnType = "allx"
  [columns.name2]
  columnNo = 3
  columnType = "last4x"
  [columns.name3]
  columnNo = 6
  columnType = "allhash"
  [columns.name4]
  columnNo = 10
  columnType = "last4z"
```
There are four cleansing types supported:  
1. allx - replace the whole column with 'X's  
2. allhash - replace the whole column with its SHA-1 sum  
3. last4x - replace the last 4 charactors with 'X's  
4. last4z - replace the last 4 charactors with '0's  
