# Concurrent Web Scraper

The goal of this project is to learn and practice concurrency in Go.

The aim will be to take a CSV in the form

```csv
suggestedName,Link
Jack Jones Portfolio,https://www.jackljones.com/
Books,https://books.toscrape.com/
"This is a test of getting, JSON",https://jsonplaceholder.typicode.com/todos/1
```

and return the following information:

- Title
- Status Code
- How Long the Request took to make

The program should automatically create a file to test against

## Todo

- [x] Skeleton Project
- [x] Structs to model data
- [x] Using hard-coded values make HTTP request
- [x] Format/Print output
- [x] Extract title
  - [x] `<title>` element
  - [x] the first `<h1>` element
- [x] Read data from CSV file
- [x] Make it work concurrently