# Concurrent Web Scraper

The goal of this project is to learn and practice concurrency in Go.

The aim will be to take a CSV in the form

```csv
suggestedName,Link
Jack Jones Portfolio,https://www.jackljones.com/
```

and return the following information:

- Title
- Status Code
- How Long the Request took to make

## Todo

- [x] Skeleton Project
- [ ] Structs to model data
- [ ] Using hard-coded values make HTTP request
- [ ] Format/Print output
- [ ] Read data from CSV file
- [ ] Make it work concurrently