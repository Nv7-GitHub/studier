# studier
Command-line studying thing

## Install
To install and update `studier`, run the following command:
```sh
cargo install studier
```

## Syntax
The question file consists of a list of questions seperated by 2 line breaks:
```
question1

question2

question3
```
There are 3 types of questions:
- Default
- List
- Blanks

## Default Questions
A default question consists of the text followed by the answer:
```
What is 1 + 1?
2
```

## List Questions
A list question is similar to a default question, but there can be multiple answers. They are answered in any order:
```
What values, when multiplied by two, equal 2, 4, 6, or 8?
1
2
3
4
```

## Blanks Questions
A blanks question has the user fill in the blanks in the order the answers are given:
```
In the year `year`, `king` was the king of Uruk.
king: Gilgamesh
year: 2750 BCE
```

## Includes
A study set can be spread over multiple files by using the `include` statement:

**topic.txt**:
```
include
topic1.txt
topic2.txt
```