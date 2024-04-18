# xlsxsplit

## What Is This ?

xlsxsplit is a tool to split xlsx file with multiple sheets into multiple files.

## How To Use

Run such comand under directory where examples.xlsx exists.
```
$ xlsxsplit split -f examples.xlsx
```
then you'll get splitted files(Students.xlsx,ClassRoom.xlsx)


## TODO

- [ ] Log output to console (optional)
- [ ] Progress display。

## 

Thanks to [metafates/go-template](https://github.com/metafates/go-template) which make it easy to write a cli tool.
